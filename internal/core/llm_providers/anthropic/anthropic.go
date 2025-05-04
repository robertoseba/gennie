package anthropic

import (
	"context"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llm_providers/base"
)

type provider struct {
	client       *anthropic.Client
	model        string
	tools        []anthropic.ToolUnionParam
	systemPrompt string
}

func NewProvider(apiKey string, model string, httpClient *http.Client) *provider {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if model == "" || model == base.ClaudeSonnet.Slug() {
		model = anthropic.ModelClaude3_7SonnetLatest
	}

	client := anthropic.NewClient(
		option.WithMiddleware(NewErrorMiddleware(os.Stdout)),
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
	)
	return &provider{
		client: &client,
		model:  model,
	}
}

func (p *provider) SetSystemPrompt(systemPrompt string) {
	p.systemPrompt = systemPrompt
}

func (p *provider) SetTools(tools []base.Tool) {
	p.tools = parseTools(tools)
}

func (p *provider) Complete(ctx context.Context, conversation *conversation.Conversation, toolResults []base.ToolResult) <-chan base.ModelResponse {
	output := make(chan base.ModelResponse, 10)

	go func() {
		defer close(output)
		messages := make([]anthropic.MessageParam, 0, conversation.Len())

		for _, qa := range conversation.QAs {
			messages = append(messages, createUserMessage(qa.GetQuestion()))
			if qa.HasAnswer() {
				messages = append(messages, createAssistantMessage(qa.GetAnswer()))
			}
		}

		// Adds tool Result to message before sending
		if len(toolResults) > 0 {
			messageBlock := []anthropic.ContentBlockParamUnion{}
			messageAssistantBlock := []anthropic.ContentBlockParamUnion{}
			for _, tr := range toolResults {
				// add back assistant block with tool use
				useBlock := anthropic.ToolUseBlockParam{
					ID:    tr.ID,
					Name:  tr.Name,
					Input: tr.Arguments,
				}
				toolUseBlock := anthropic.ContentBlockParamUnion{OfRequestToolUseBlock: &useBlock}
				messageAssistantBlock = append(messageAssistantBlock, toolUseBlock)

				if tr.IsError() {
					messageBlock = append(messageBlock, anthropic.NewToolResultBlock(tr.ID, tr.Error.Error(), true))
				} else {
					messageBlock = append(messageBlock, anthropic.NewToolResultBlock(tr.ID, string(tr.Result), false))
				}
			}
			messages = append(messages, anthropic.NewAssistantMessage(messageAssistantBlock...))
			messages = append(messages, anthropic.NewUserMessage(messageBlock...))
		}

		stream := p.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
			Model:     p.model,
			MaxTokens: 1024,
			Messages:  messages,
			System: []anthropic.TextBlockParam{
				{Text: p.systemPrompt},
			},
			Tools: p.tools,
		})

		message := anthropic.Message{}

		for stream.Next() {
			event := stream.Current()
			err := message.Accumulate(event)
			if err != nil {
				// TODO: figure out how to handle this error
				continue
				// panic(err)
			}

			switch eventVariant := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
				case anthropic.TextDelta:
					output <- base.ModelResponse{
						Text:       deltaVariant.Text,
						Error:      nil,
						StopReason: base.StopReasonNone,
					}
				}
			}
			if stream.Err() != nil {
				output <- base.ModelResponse{
					Text:       "something went wrong!!",
					Error:      stream.Err(),
					StopReason: base.StopReasonError,
				}
			}
		}

		if message.StopReason == anthropic.MessageStopReasonToolUse {
			for _, block := range message.Content {
				switch variant := block.AsAny().(type) {
				case anthropic.ToolUseBlock:
					output <- base.ModelResponse{
						Text:       "",
						Error:      nil,
						StopReason: base.StopReasonTools,
						FunctionCall: base.FunctionCall{
							ID:        variant.ID,
							Name:      variant.Name,
							Arguments: variant.Input,
						},
					}
				}
			}
		}
	}()

	return output
}

func createUserMessage(content string) anthropic.MessageParam {
	return anthropic.NewUserMessage(anthropic.NewTextBlock(content))
}

func createAssistantMessage(content string) anthropic.MessageParam {
	return anthropic.NewAssistantMessage(anthropic.NewTextBlock(content))
}

func parseTools(tools []base.Tool) []anthropic.ToolUnionParam {
	if len(tools) == 0 {
		return nil
	}

	toolReturn := make([]anthropic.ToolUnionParam, 0)
	for _, tool := range tools {
		toolParam := anthropic.ToolParam{
			Name:        tool.Name,
			Description: anthropic.String(tool.Description),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: tool.InputSchema.Properties,
			},
		}
		toolReturn = append(toolReturn, anthropic.ToolUnionParam{OfTool: &toolParam})
	}

	return toolReturn
}
