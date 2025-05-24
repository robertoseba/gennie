package gemini

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmcore"
	"google.golang.org/genai"
)

const (
	ExportedModelSlug        = "gemini"
	ExportedModelDescription = "Gemini Flash 2.5"
)

type provider struct {
	client       *genai.Client
	model        string
	systemPrompt string
	tools        []*genai.FunctionDeclaration
}

func NewProvider(apiKey string, model string, httpClient *http.Client) *provider {
	if model == "" || model == ExportedModelSlug {
		model = "gemini-2.5-flash-preview-05-20"
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
	p.tools = convertToolsToProvider(tools)
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

		config := &genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: p.systemPrompt}}},
			Tools: []*genai.Tool{
				{FunctionDeclarations: p.tools},
			},
		}
		streamRes := p.client.Models.GenerateContentStream(ctx, p.model, messages, config)

		for result, err := range streamRes {
			if err != nil {
				output <- llmcore.LlmResponse{
					Error:      err,
					Text:       "something went wrong",
					StopReason: llmcore.StopReasonError,
				}
				continue
			}

			response := llmcore.LlmResponse{
				Text:       result.Text(),
				StopReason: llmcore.StopReasonNone,
			}

			if result.Candidates[0].FinishReason == genai.FinishReasonStop {
				response.StopReason = llmcore.StopReasonEnd
			}

			if result.Candidates[0].Content.Parts[0].FunctionCall != nil {
				response.StopReason = llmcore.StopReasonTools

				args := result.Candidates[0].Content.Parts[0].FunctionCall.Args
				argsBytes, err := json.Marshal(args)
				if err != nil {
					output <- llmcore.LlmResponse{
						Error:      err,
						Text:       "something went wrong",
						StopReason: llmcore.StopReasonError,
					}
					continue
				}

				response.FunctionCall = llmcore.FunctionCall{
					ID:        result.Candidates[0].Content.Parts[0].FunctionCall.ID,
					Name:      result.Candidates[0].Content.Parts[0].FunctionCall.Name,
					Arguments: argsBytes,
				}
			}

			output <- response
		}
	}()

	return output
}

func addToolResultsToMessages(messages []*genai.Content, toolResults []llmcore.ToolResult) []*genai.Content {
	for _, toolResult := range toolResults {
		var mappedResult map[string]any
		if err := json.Unmarshal(toolResult.Result, &mappedResult); err != nil {
			log.Printf("Error marshaling tool result: %v", err)
			continue
		}

		toolMessage := genai.NewContentFromFunctionResponse(toolResult.Name, mappedResult, genai.RoleUser)
		messages = append(messages, toolMessage)
	}
	return messages
}

func convertToolsToProvider(tools []llmcore.Tool) []*genai.FunctionDeclaration {
	functions := make([]*genai.FunctionDeclaration, 0, len(tools))

	for _, tool := range tools {
		byteSchema, err := json.Marshal(tool.InputSchema)
		if err != nil {
			log.Printf("Error marshaling tool input schema: %v", err)
			continue
		}

		var schema *genai.Schema
		err = json.Unmarshal(byteSchema, &schema)
		if err != nil {
			log.Printf("Error marshaling tool input schema: %v", err)
			continue
		}
		schema.Type = genai.TypeObject

		functionDec := &genai.FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  schema,
		}
		functions = append(functions, functionDec)
	}
	return functions
}
