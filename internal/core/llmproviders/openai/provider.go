package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmcore"
)

const (
	ExportedModelSlug        = "gpt-4o"
	ExportedModelDescription = "OpenAI GPT-4o"

	ExportedModelSlug_Mini        = "gpt-4o-mini"
	ExportedModelDescription_Mini = "OpenAI GPT-4o Mini"

	DefaultModelSlug        = ExportedModelSlug_Mini
	DefaultModelDescription = ExportedModelDescription_Mini
)

type provider struct {
	client       *openai.Client
	model        string
	systemPrompt string
	tools        []openai.ChatCompletionToolParam
	options      []option.RequestOption
	logger       *slog.Logger
}

type opts func(p *provider)

func NewProvider(apiKey string, opts ...opts) *provider {
	p := &provider{
		model:   DefaultModelSlug,
		tools:   make([]openai.ChatCompletionToolParam, 0),
		options: make([]option.RequestOption, 0, len(opts)+1),
	}
	p.options = append(p.options, option.WithAPIKey(apiKey))

	for _, opt := range opts {
		opt(p)
	}

	client := openai.NewClient(p.options...)
	p.client = &client

	return p
}

func (p *provider) SetSystemPrompt(systemPrompt string) {
	p.systemPrompt = systemPrompt
}

func (p *provider) SetTools(tools []llmcore.Tool) {
	p.tools = convertToolsToProvider(tools)
}

func (p *provider) Complete(ctx context.Context, conversation *conversation.Conversation, toolResults []llmcore.ToolResult) <-chan llmcore.LlmResponse {
	output := make(chan llmcore.LlmResponse, 10)

	go func() {
		defer close(output)

		messages := make([]openai.ChatCompletionMessageParamUnion, 0, conversation.Len()+len(toolResults))

		for _, message := range conversation.QAs {
			messages = append(messages, openai.UserMessage(message.GetQuestion()))
			if message.HasAnswer() {
				messages = append(messages, openai.AssistantMessage(message.GetAnswer()))
			}
		}

		if len(toolResults) > 0 {
			messages = addToolResultsToMessages(messages, toolResults)
		}

		stream := p.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Messages: messages,
			Seed:     openai.Int(0),
			Model:    p.model,
			Tools:    p.tools,
		})
		defer stream.Close()

		acc := openai.ChatCompletionAccumulator{}

		for stream.Next() {
			chunk := stream.Current()
			acc.AddChunk(chunk)

			if tool, ok := acc.JustFinishedToolCall(); ok {
				msg := acc.ChatCompletion.Choices[0].Message.ToParam()
				b, _ := json.Marshal(msg)
				fmt.Println("Tool call finished:", string(b))
				output <- llmcore.LlmResponse{
					Text:       chunk.Choices[0].Delta.Content,
					Error:      nil,
					StopReason: llmcore.StopReasonTools,
					FunctionCall: llmcore.FunctionCall{
						ID:        tool.ID,
						Name:      tool.Name,
						Arguments: []byte(tool.Arguments),
					},
				}
			}

			if refusal, ok := acc.JustFinishedRefusal(); ok {
				println("Refusal stream finished:", refusal)
			}

			if len(chunk.Choices) > 0 {
				output <- llmcore.LlmResponse{
					Text:       chunk.Choices[0].Delta.Content,
					Error:      nil,
					StopReason: llmcore.StopReasonNone,
				}
			}
		}

		if stream.Err() != nil {
			output <- llmcore.LlmResponse{
				Error:      stream.Err(),
				Text:       "something went wrong",
				StopReason: llmcore.StopReasonError,
			}
		}
	}()

	return output
}

func convertToolsToProvider(tools []llmcore.Tool) []openai.ChatCompletionToolParam {
	converted := make([]openai.ChatCompletionToolParam, len(tools))
	for i, tool := range tools {
		byteSchema, err := json.Marshal(tool.InputSchema)
		if err != nil {
			log.Printf("Error marshaling tool input schema: %v", err)
			continue
		}

		var schema shared.FunctionParameters
		err = json.Unmarshal(byteSchema, &schema)
		if err != nil {
			log.Printf("Error unmarshaling tool input schema: %v", err)
			continue
		}

		converted[i] = openai.ChatCompletionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: param.NewOpt(tool.Description),
				Parameters:  schema,
			},
		}
	}
	return converted
}

func addToolResultsToMessages(messages []openai.ChatCompletionMessageParamUnion, toolResults []llmcore.ToolResult) []openai.ChatCompletionMessageParamUnion {
	assistantCall := openai.ChatCompletionAssistantMessageParam{
		ToolCalls: make([]openai.ChatCompletionMessageToolCallParam, 0, len(toolResults)),
	}

	for _, toolResult := range toolResults {
		toolCall := openai.ChatCompletionMessageToolCallParam{
			ID: toolResult.ID,
			Function: openai.ChatCompletionMessageToolCallFunctionParam{
				Name:      toolResult.Name,
				Arguments: string(toolResult.Arguments),
			},
		}
		assistantCall.ToolCalls = append(assistantCall.ToolCalls, toolCall)
		msg := openai.ChatCompletionMessageParamUnion{
			OfAssistant: &assistantCall,
		}
		messages = append(messages, msg)

		message := openai.ToolMessage(string(toolResult.Result), toolResult.ID)
		messages = append(messages, message)
	}
	return messages
}
