package anthropic

import (
	"context"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmcore/entities"
)

const (
	ExportedModelSlug        = "sonnet"
	ExportedModelDescription = "Claude Sonnet 3.7"
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

	if model == "" || model == ExportedModelSlug {
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

func (p *provider) SetTools(tools []entities.Tool) {
	p.tools = convertToolsToProvider(tools)
}

func (p *provider) Complete(ctx context.Context, conversation *conversation.Conversation, toolResults []entities.ToolResult) <-chan entities.LlmResponse {
	output := make(chan entities.LlmResponse, 10)

	go func() {
		defer close(output)
		messages := make([]anthropic.MessageParam, 0, conversation.Len()+len(toolResults))

		for _, qa := range conversation.QAs {
			messages = append(messages, createUserMessage(qa.GetQuestion()))
			if qa.HasAnswer() {
				messages = append(messages, createAssistantMessage(qa.GetAnswer()))
			}
		}

		if len(toolResults) > 0 {
			messages = addToolResultsToMessages(messages, toolResults)
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
					output <- entities.LlmResponse{
						Text:       deltaVariant.Text,
						Error:      nil,
						StopReason: entities.StopReasonNone,
					}
				}
			}
			if stream.Err() != nil {
				output <- entities.LlmResponse{
					Text:       "something went wrong!!",
					Error:      stream.Err(),
					StopReason: entities.StopReasonError,
				}
			}
		}

		if message.StopReason == anthropic.MessageStopReasonToolUse {
			toolResponses := parseToolUseFrom(&message)
			for idx := range toolResponses {
				output <- toolResponses[idx]
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
