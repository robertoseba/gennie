package complete

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models"
	"github.com/robertoseba/gennie/internal/core/models/response"
	"github.com/robertoseba/gennie/internal/core/models/tools"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type toolDetails struct {
	mcpClient        *McpClient
	tool             tools.Tool
	requiresApproval bool
}

type CompleteService struct {
	conversationRepo conversation.IConversationRepository
	profileRepo      profile.IProfileRepository
	apiClient        models.IApiClient
	config           *config.Config
	tools            map[string]toolDetails // each toolName maps to a mcpClient so we can make a request
}

type InputDTO struct {
	Question    string
	ProfileSlug string
	Model       string
	IsFollowUp  bool
	AppendFile  string
}

func NewCompleteService(
	cr conversation.IConversationRepository,
	pr profile.IProfileRepository,
	apiClient models.IApiClient,
	config *config.Config,
) *CompleteService {
	return &CompleteService{
		conversationRepo: cr,
		profileRepo:      pr,
		apiClient:        apiClient,
		config:           config,
	}
}

func (s *CompleteService) Execute(input *InputDTO) (<-chan models.StreamResponse, error) {
	var conv *conversation.Conversation
	var err error

	conv, profile, model, err := s.processInput(input)
	if err != nil {
		return nil, err
	}

	model.SetSystemPrompt(profile.Data)

	modelResponseChan := make(chan models.StreamResponse)
	outputChan := s.pipeToSaveConversation(conv, modelResponseChan)

	go func() {
		defer close(modelResponseChan)

		outputChan <- models.StreamResponse{Data: model.Model().String(), Type: models.ModelInfo, Err: nil}
		outputChan <- models.StreamResponse{Data: profile.Name, Type: models.ProfileInfo, Err: nil}
		if len(profile.McpServers) > 0 {
			outputChan <- models.StreamResponse{Data: fmt.Sprintf("Loading MCP Servers..."), Type: models.LoadingInfo, Err: nil}
			mcpTools, err := startMcpServers(profile)
			if err != nil {
				outputChan <- models.StreamResponse{Err: err}
			}
			s.tools = mcpTools
			var modelTools []tools.Tool
			for toolName := range s.tools {
				tool := tools.Tool{
					Name:        s.tools[toolName].tool.Name,
					Description: s.tools[toolName].tool.Description,
					InputSchema: tools.ToolInputSchema{
						Properties: s.tools[toolName].tool.InputSchema.Properties,
					},
				}
				modelTools = append(modelTools, tool)
			}
			model.SetTools(modelTools)
		}

		ctx := context.Background()
		toolResults := make([]tools.ToolResult, 0)

		// Keeps calling the model while it needs to return function calls
		outputChan <- models.StreamResponse{Data: "Asking the model...", Type: models.LoadingInfo, Err: nil}
		for {
			resp := complete(ctx, model, conv, toolResults, profile.Data, modelResponseChan)

			// If does not need to send function call back to model than breaks out of the loop
			if resp.StopReason != response.StopReasonTools {
				break
			}

			if s.tools[resp.FunctionCall.Name].requiresApproval {
				outputChan <- models.StreamResponse{Data: fmt.Sprintf("Can I run this tool: %s with parameters (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: models.ApprovalRequest, Err: nil}
			}

			outputChan <- models.StreamResponse{Data: fmt.Sprintf("Using tool: %s -> (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: models.LoadingInfo, Err: nil}
			result := s.callTool(&resp)
			toolResults = append(toolResults, *result)
		}
	}()

	return outputChan, nil
}

func (s *CompleteService) pipeToSaveConversation(conv *conversation.Conversation, inputChan <-chan models.StreamResponse) chan models.StreamResponse {
	outputChan := make(chan models.StreamResponse)

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
			outputChan <- models.StreamResponse{Err: err}
		}
		err = s.conversationRepo.SaveAsActive(conv)
		if err != nil {
			outputChan <- models.StreamResponse{Err: err}
		}
	}()

	return outputChan
}

// TODO: clean up and refactor this function
func complete(ctx context.Context, model *models.BaseModel, conv *conversation.Conversation, toolResults []tools.ToolResult, profileData string, outputChan chan<- models.StreamResponse) response.ModelResponse {
	respChan, err := model.CompleteStreamable(ctx, conv, profileData, toolResults)
	if err != nil {
		outputChan <- models.StreamResponse{Err: err}
	}

	var toolCallRequest response.ModelResponse

	for modelResponse := range respChan {
		if modelResponse.StopReason == response.StopReasonTools {
			toolCallRequest = modelResponse
			continue
		}

		outputChan <- models.StreamResponse{Data: modelResponse.Text, Err: modelResponse.Error}
	}

	return toolCallRequest
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return string(content), err
}

func (s *CompleteService) callTool(modelResponse *response.ModelResponse) *tools.ToolResult {
	if _, ok := s.tools[modelResponse.FunctionCall.Name]; !ok {
		return tools.NewToolResponseError(modelResponse, fmt.Errorf("tool %s not found", modelResponse.FunctionCall.Name))
	}

	var args map[string]any
	if len(modelResponse.FunctionCall.Arguments) > 0 {
		args = make(map[string]any)
		err := json.Unmarshal([]byte(modelResponse.FunctionCall.Arguments), &args)
		if err != nil {
			return tools.NewToolResponseError(nil, err)
		}
	}

	toolResponse, err := s.tools[modelResponse.FunctionCall.Name].mcpClient.ExecTool(context.TODO(), modelResponse.FunctionCall.Name, args)
	if err != nil {
		return tools.NewToolResponseError(nil, err)
	}

	return tools.NewToolResponseFrom(modelResponse, toolResponse)
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

func (s *CompleteService) processInput(input *InputDTO) (*conversation.Conversation, *profile.Profile, *models.BaseModel, error) {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return nil, nil, nil, err
	}

	profile, err := s.loadProfile(input.ProfileSlug, conv)
	if err != nil {
		return nil, nil, nil, err
	}
	conv.SetProfileTo(profile.Slug)

	model, err := s.loadModel(input.Model, conv)
	if err != nil {
		return nil, nil, nil, err
	}
	conv.SetModelTo(model.Model().Slug())

	// Resets the conversation if not a follow up
	if !input.IsFollowUp {
		conv = conversation.NewConversation(profile.Slug, model.Model().Slug())
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

func (s *CompleteService) loadModel(modelSlug string, conv *conversation.Conversation) (*models.BaseModel, error) {
	if modelSlug == "" {
		modelSlug = conv.ModelSlug
	}

	return models.NewModel(modelSlug, s.apiClient, *s.config)
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
