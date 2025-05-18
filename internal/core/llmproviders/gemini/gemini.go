package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmcore"
	"google.golang.org/genai"
)

const (
	ExportedModelSlug        = "gemini"
	ExportedModelDescription = "Gemini Flash 2.0"
)

type provider struct {
	client       *genai.Client
	model        string
	systemPrompt string
	tools        []*genai.FunctionDeclaration
}

func NewProvider(apiKey string, model string, httpClient *http.Client) *provider {
	if model == "" || model == ExportedModelSlug {
		model = "gemini-2.0-flash"
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:     apiKey,
		Backend:    genai.BackendGeminiAPI,
		HTTPClient: httpClient,
	})
	if err != nil {
		panic("Failed to create GenAI client: " + err.Error())
	}

	return &provider{
		client: client,
		model:  model,
	}
}

func (p *provider) SetSystemPrompt(systemPrompt string) {
	p.systemPrompt = systemPrompt
}

func (p *provider) SetTools(tools []llmcore.Tool) {
	// p.tools = convertToolsToProvider(tools)
}

func (p *provider) Complete(ctx context.Context, conversation *conversation.Conversation, toolResults []llmcore.ToolResult) <-chan llmcore.LlmResponse {
	output := make(chan llmcore.LlmResponse, 10)

	go func() {
		defer close(output)

		messages := make([]*genai.Content, 0, conversation.Len()+len(toolResults))

		for _, qa := range conversation.QAs {
			question := genai.NewContentFromText(qa.GetQuestion(), genai.RoleUser)
			messages = append(messages, question)

			if qa.HasAnswer() {
				answer := genai.NewContentFromText(qa.GetAnswer(), genai.RoleModel)
				messages = append(messages, answer)
			}
		}

		if len(toolResults) > 0 {
			messages = addToolResultsToMessages(messages, toolResults)
		}

		chat, err := p.client.Chats.Create(
			ctx,
			"gemini-2.0-flash",
			&genai.GenerateContentConfig{
				SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: p.systemPrompt}}},
				Tools:             []*genai.Tool{
					// {FunctionDeclarations: p.tools},
				},
			},
			messages,
		)
		if err != nil {
			output <- llmcore.LlmResponse{
				Error:      err,
				Text:       "something went wrong",
				StopReason: llmcore.StopReasonError,
			}
		}

		for result, err := range chat.SendMessageStream(ctx, genai.Part{Text: conversation.LastQuestion()}) {
			if err != nil {
				output <- llmcore.LlmResponse{
					Error:      err,
					Text:       "something went wrong",
					StopReason: llmcore.StopReasonError,
				}
			}
			debugPrint(result)
			// output <- llmcore.LlmResponse{
			// 	Text:       result.Text(),
			// 	StopReason: llmcore.StopReasonNone,
			// }
		}
	}()

	return output
}

func debugPrint[T any](r *T) {
	// Marshal the result to JSON.
	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	// Log the output.
	fmt.Println(string(response))
}

func convertToolsToProvider(tools []llmcore.Tool) []*genai.FunctionDeclaration {
	functions := make([]*genai.FunctionDeclaration, 0, len(tools))

	// for i, tool := range tools {
	//    genai.NewPartFromFunctionCall(tool.Name, args map[string]any)
	// 	functions[i] = &genai.FunctionDeclaration{
	// 		Name:        tool.Name,
	// 		Description: tool.Description,
	// 		Parameters: &genai.Schema{
	// 			Type:       tool.InputSchema.Type,
	// 			Properties: tool.InputSchema.Properties,
	// 			Required:   tool.InputSchema.Required,
	// 		},
	// 	}
	// }
	return functions
}

func addToolResultsToMessages(messages []*genai.Content, toolResults []llmcore.ToolResult) []*genai.Content {
	// for _, toolResult := range toolResults {
	// 	toolMessage := genai.Newfun(toolResult.Result, genai.RoleUser)
	// 	messages = append(messages, toolMessage)
	// }
	return messages
}
