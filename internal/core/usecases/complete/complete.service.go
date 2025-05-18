package complete

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmcore"
	"github.com/robertoseba/gennie/internal/core/llmcore/entities"
	"github.com/robertoseba/gennie/internal/core/llmproviders/factory"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type toolDetails struct {
	mcpClient        *McpClient
	tool             entities.Tool
	requiresApproval bool
}

type CompleteService struct {
	conversationRepo conversation.IConversationRepository
	profileRepo      profile.IProfileRepository
	httpClient       *http.Client
	config           *config.Config
	tools            map[string]toolDetails // each toolName maps to a mcpClient so we can make a request
}

type InputDTO struct {
	Question    string
	ProfileSlug string
	ModelSlug   string
	IsFollowUp  bool
	AppendFile  string
}

func NewCompleteService(
	cr conversation.IConversationRepository,
	pr profile.IProfileRepository,
	httpClient *http.Client,
	config *config.Config,
) *CompleteService {
	return &CompleteService{
		conversationRepo: cr,
		profileRepo:      pr,
		httpClient:       httpClient,
		config:           config,
	}
}

func (s *CompleteService) Execute(input *InputDTO) (<-chan entities.CompleteResponse, error) {
	var conv *conversation.Conversation
	var err error

	conv, profile, model, err := s.processInput(input)
	if err != nil {
		return nil, err
	}

	model.SetSystemPrompt(profile.Data)

	modelResponseChan := make(chan entities.CompleteResponse)
	outputChan := s.pipeToSaveConversation(conv, modelResponseChan)

	go func() {
		defer close(modelResponseChan)

		outputChan <- entities.CompleteResponse{Data: conv.ModelSlug, Type: entities.ModelInfo, Err: nil}
		outputChan <- entities.CompleteResponse{Data: profile.Name, Type: entities.ProfileInfo, Err: nil}
		if len(profile.McpServers) > 0 {
			outputChan <- entities.CompleteResponse{Data: fmt.Sprintf("Loading MCP Servers..."), Type: entities.LoadingInfo, Err: nil}
			mcpTools, err := startMcpServers(profile)
			if err != nil {
				outputChan <- entities.CompleteResponse{Err: err}
			}

			s.tools = mcpTools
			var modelTools []entities.Tool
			for toolName := range s.tools {
				tool := entities.Tool{
					Name:        s.tools[toolName].tool.Name,
					Description: s.tools[toolName].tool.Description,
					InputSchema: entities.ToolInputSchema{
						Properties: s.tools[toolName].tool.InputSchema.Properties,
					},
				}
				modelTools = append(modelTools, tool)
			}
			model.SetTools(modelTools)
		}

		ctx := context.Background()
		toolResults := make([]entities.ToolResult, 0)

		// Keeps calling the model while it needs to return function calls
		outputChan <- entities.CompleteResponse{Data: "Asking the model...", Type: entities.LoadingInfo, Err: nil}
		for {
			resp := complete(ctx, model, conv, toolResults, modelResponseChan)

			// If does not need to send function call back to model than breaks out of the loop
			if resp.StopReason != entities.StopReasonTools {
				break
			}

			if s.tools[resp.FunctionCall.Name].requiresApproval {
				outputChan <- entities.CompleteResponse{Data: fmt.Sprintf("Can I run this tool: %s with parameters (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: entities.ApprovalRequest, Err: nil}
			}

			outputChan <- entities.CompleteResponse{Data: fmt.Sprintf("Using tool: %s -> (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: entities.LoadingInfo, Err: nil}
			result := s.callTool(&resp)
			toolResults = append(toolResults, *result)
		}
	}()

	return outputChan, nil
}

func (s *CompleteService) pipeToSaveConversation(conv *conversation.Conversation, inputChan <-chan entities.CompleteResponse) chan entities.CompleteResponse {
	outputChan := make(chan entities.CompleteResponse)

	go func() {
		defer close(outputChan)

		convBuffer := strings.Builder{}
		for msg := range inputChan {
			outputChan <- msg
			if msg.Err == nil {
				convBuffer.WriteString(msg.Data)
			}
		}
		err := conv.AnswerLastQuestion(convBuffer.String())
		if err != nil {
			outputChan <- entities.CompleteResponse{Err: err}
		}
		err = s.conversationRepo.SaveAsActive(conv)
		if err != nil {
			outputChan <- entities.CompleteResponse{Err: err}
		}
	}()

	return outputChan
}

func complete(ctx context.Context, model llmcore.LlmProvider, conv *conversation.Conversation, toolResults []entities.ToolResult, outputChan chan<- entities.CompleteResponse) entities.LlmResponse {
	respChan := model.Complete(ctx, conv, toolResults)

	var toolCallRequest entities.LlmResponse

	for modelResponse := range respChan {
		if modelResponse.StopReason == entities.StopReasonTools {
			toolCallRequest = modelResponse
			continue
		}

		outputChan <- entities.CompleteResponse{Data: modelResponse.Text, Err: modelResponse.Error}
	}

	return toolCallRequest
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return string(content), err
}

func (s *CompleteService) callTool(modelResponse *entities.LlmResponse) *entities.ToolResult {
	if _, ok := s.tools[modelResponse.FunctionCall.Name]; !ok {
		return entities.NewToolResponseError(modelResponse, fmt.Errorf("tool %s not found", modelResponse.FunctionCall.Name))
	}

	var args map[string]any
	if len(modelResponse.FunctionCall.Arguments) > 0 {
		args = make(map[string]any)
		err := json.Unmarshal([]byte(modelResponse.FunctionCall.Arguments), &args)
		if err != nil {
			return entities.NewToolResponseError(nil, err)
		}
	}

	toolResponse, err := s.tools[modelResponse.FunctionCall.Name].mcpClient.ExecTool(context.TODO(), modelResponse.FunctionCall.Name, args)
	if err != nil {
		return entities.NewToolResponseError(nil, err)
	}

	return entities.NewToolResponseFrom(modelResponse, toolResponse)
}

func startMcpServers(profile *profile.Profile) (map[string]toolDetails, error) {
	returnTools := make(map[string]toolDetails)
	for _, server := range profile.McpServers {

		// TODO: figure out context handling
		ctx := context.TODO()
		mcpServer, err := StartMcpServer(ctx, server.Command, server.Envs, server.Args)
		if err != nil {
			fmt.Printf("error %v", err)
			continue
		}

		mcpTools, err := mcpServer.ListTools(ctx)
		if err != nil {
			fmt.Printf("error %v", err)
			panic(err)
		}

		// Filter tools based on profile allowed tools
		for _, tool := range mcpTools {
			if slices.Contains(server.AllowedTools, tool.Name) || len(server.AllowedTools) == 0 {
				returnTools[tool.Name] = toolDetails{
					tool:             tool,
					mcpClient:        mcpServer,
					requiresApproval: server.RequiresApproval,
				}
			}
		}
	}

	return returnTools, nil
}

func (s *CompleteService) processInput(input *InputDTO) (*conversation.Conversation, *profile.Profile, llmcore.LlmProvider, error) {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return nil, nil, nil, err
	}

	profile, err := s.loadProfile(input.ProfileSlug, conv)
	if err != nil {
		return nil, nil, nil, err
	}
	conv.SetProfileTo(profile.Slug)

	model, modelEnum, err := s.loadModel(input.ModelSlug, conv)
	if err != nil {
		return nil, nil, nil, err
	}
	conv.SetModelTo(modelEnum.Slug())

	// Resets the conversation if not a follow up
	if !input.IsFollowUp {
		conv = conversation.NewConversation(profile.Slug, modelEnum.Slug())
	}

	err = s.setQuestion(conv, input.Question, input.AppendFile)
	if err != nil {
		return nil, nil, nil, err
	}

	return conv, profile, model, nil
}

func (s *CompleteService) loadProfile(profileSlug string, conv *conversation.Conversation) (*profile.Profile, error) {
	if profileSlug == "" {
		profileSlug = conv.ProfileSlug
	}
	return s.profileRepo.FindBySlug(profileSlug)
}

func (s *CompleteService) loadModel(modelSlug string, conv *conversation.Conversation) (llmcore.LlmProvider, factory.ModelEnum, error) {
	if modelSlug == "" {
		modelSlug = conv.ModelSlug
	}

	modelEnum, ok := factory.ParseFrom(modelSlug)
	if !ok {
		return nil, factory.DefaultModel, llmcore.ErrModelNotFound
	}

	return factory.NewModel(modelEnum, s.httpClient, *s.config), modelEnum, nil
}

func (s *CompleteService) setQuestion(conv *conversation.Conversation, question string, filename string) error {
	if filename != "" {
		content, err := readFile(filename)
		if err != nil {
			return err
		}

		question += "\n" + content
	}

	if err := conv.NewQuestion(question); err != nil {
		return err
	}
	return nil
}
