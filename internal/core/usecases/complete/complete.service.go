package complete

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/robertoseba/gennie/internal/core/config"
	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmcore"
	"github.com/robertoseba/gennie/internal/core/llmproviders/factory"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type toolDetails struct {
	mcpClient        *McpClient
	tool             llmcore.Tool
	requiresApproval bool
}

type CompleteService struct {
	conversationRepo conversation.IConversationRepository
	profileRepo      profile.IProfileRepository
	httpClient       *http.Client
	config           *config.Config
	logger           *slog.Logger
	tools            map[string]toolDetails // each toolName maps to a mcpClient so we can make a request
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
		logger:           slog.Default(),
	}
}

func (s *CompleteService) Execute(ctx context.Context) (<-chan llmcore.CompleteResponse, error) {
	var err error

	llmProvider, err := s.processInput(ctx)
	if err != nil {
		return nil, err
	}

	req, _ := GetRequestFromCtx(ctx)
	llmProvider.SetSystemPrompt(req.Profile.Data)

	modelResponseChan := make(chan llmcore.CompleteResponse)
	outputChan := s.pipeToSaveConversation(req.Conversation, modelResponseChan)

	go func() {
		defer close(modelResponseChan)

		outputChan <- llmcore.CompleteResponse{Data: req.Conversation.ModelSlug, Type: llmcore.ModelInfo, Err: nil}
		outputChan <- llmcore.CompleteResponse{Data: req.Profile.Name, Type: llmcore.ProfileInfo, Err: nil}

		if len(req.Profile.McpServers) > 0 {
			outputChan <- llmcore.CompleteResponse{Data: "Loading MCP Servers...", Type: llmcore.LoadingInfo, Err: nil}

			mcpTools, err := startMcpServers(ctx)
			if err != nil {
				outputChan <- llmcore.CompleteResponse{Err: err}
			}

			s.tools = mcpTools
			var modelTools []llmcore.Tool
			for toolName := range s.tools {
				tool := llmcore.Tool{
					Name:        s.tools[toolName].tool.Name,
					Description: s.tools[toolName].tool.Description,
					InputSchema: llmcore.ToolInputSchema{
						Properties: s.tools[toolName].tool.InputSchema.Properties,
					},
				}
				modelTools = append(modelTools, tool)
			}
			llmProvider.SetTools(modelTools)
		}

		toolResults := make([]llmcore.ToolResult, 0)

		// Keeps calling the model while it needs to return function calls
		outputChan <- llmcore.CompleteResponse{Data: "Asking the model...", Type: llmcore.LoadingInfo, Err: nil}
		for {
			resp := complete(ctx, llmProvider, toolResults, modelResponseChan)

			// If does not need to send function call back to model than breaks out of the loop
			if resp.StopReason != llmcore.StopReasonTools {
				break
			}

			if s.tools[resp.FunctionCall.Name].requiresApproval {
				outputChan <- llmcore.CompleteResponse{Data: fmt.Sprintf("Can I run this tool: %s with parameters (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: llmcore.ApprovalRequest, Err: nil}
			}

			outputChan <- llmcore.CompleteResponse{Data: fmt.Sprintf("Using tool: %s -> (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: llmcore.LoadingInfo, Err: nil}
			result := s.callTool(ctx, &resp)

			toolResults = append(toolResults, *result)
		}
	}()

	return outputChan, nil
}

// TODO: pass context and check for cancellation
func (s *CompleteService) pipeToSaveConversation(conv *conversation.Conversation, inputChan <-chan llmcore.CompleteResponse) chan llmcore.CompleteResponse {
	outputChan := make(chan llmcore.CompleteResponse)

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
			outputChan <- llmcore.CompleteResponse{Err: err}
		}
		err = s.conversationRepo.SaveAsActive(conv)
		if err != nil {
			outputChan <- llmcore.CompleteResponse{Err: err}
		}
	}()

	return outputChan
}

func complete(ctx context.Context, llmProvider llmcore.LlmProvider, toolResults []llmcore.ToolResult, outputChan chan<- llmcore.CompleteResponse) llmcore.LlmResponse {
	req, _ := GetRequestFromCtx(ctx)
	fmt.Printf("conversatoin %v\n", req.Conversation)
	respChan := llmProvider.Complete(ctx, req.Conversation, toolResults)

	var toolCallRequest llmcore.LlmResponse

	for llmResponse := range respChan {
		if llmResponse.StopReason == llmcore.StopReasonTools {
			toolCallRequest = llmResponse
			continue
		}

		outputChan <- llmcore.CompleteResponse{Data: llmResponse.Text, Err: llmResponse.Error}
	}

	return toolCallRequest
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return string(content), err
}

func (s *CompleteService) callTool(ctx context.Context, llmResponse *llmcore.LlmResponse) *llmcore.ToolResult {
	if _, ok := s.tools[llmResponse.FunctionCall.Name]; !ok {
		return llmcore.NewToolResponseError(llmResponse, fmt.Errorf("tool %s not found", llmResponse.FunctionCall.Name))
	}

	var args map[string]any
	if len(llmResponse.FunctionCall.Arguments) > 0 {
		args = make(map[string]any)
		err := json.Unmarshal([]byte(llmResponse.FunctionCall.Arguments), &args)
		if err != nil {
			return llmcore.NewToolResponseError(nil, err)
		}
	}

	toolResponse, err := s.tools[llmResponse.FunctionCall.Name].mcpClient.ExecTool(ctx, llmResponse.FunctionCall.Name, args)
	s.logger.Debug("Tool response", "toolName", llmResponse.FunctionCall.Name, "response", toolResponse)

	if err != nil {
		return llmcore.NewToolResponseError(nil, err)
	}

	return llmcore.NewToolResponseFrom(llmResponse, toolResponse)
}

func startMcpServers(ctx context.Context) (map[string]toolDetails, error) {
	returnTools := make(map[string]toolDetails)
	req, _ := GetRequestFromCtx(ctx)

	for _, server := range req.Profile.McpServers {

		mcpServer, err := StartMcpServer(ctx, server.Command, server.Envs, server.Args)
		if err != nil {
			fmt.Printf("error %v", err)
			continue
		}

		mcpTools, err := mcpServer.ListTools(ctx)
		if err != nil {
			fmt.Printf("error %v", err)
			return nil, fmt.Errorf("failed to list tools from MCP server %s: %w", server.Command, err)
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

func (s *CompleteService) processInput(ctx context.Context) (llmcore.LlmProvider, error) {
	conv, err := s.conversationRepo.LoadActive()
	if err != nil {
		return nil, err
	}
	req, _ := GetRequestFromCtx(ctx)

	profile, err := s.loadProfile(req.ProfileSlug, conv)
	if err != nil {
		return nil, err
	}
	conv.SetProfileTo(profile.Slug)
	req.Profile = profile

	llmProvider, modelEnum, err := s.loadLlmProvider(req.ModelSlug, conv)
	if err != nil {
		return nil, err
	}
	conv.SetModelTo(modelEnum.Slug())
	req.ModelSlug = modelEnum.Slug()

	// Resets the conversation if not a follow up
	if !req.IsFollowUp {
		conv = conversation.NewConversation(profile.Slug, modelEnum.Slug())
	}

	err = s.setQuestion(conv, req.Question, req.AppendFilename)
	if err != nil {
		return nil, err
	}

	// TODO: this is too messy
	req.Conversation = conv
	return llmProvider, nil
}

func (s *CompleteService) loadProfile(profileSlug string, conv *conversation.Conversation) (*profile.Profile, error) {
	s.logger.Debug("Loading profile: ", "profile flag", profileSlug, "prev. conversation profile", conv.ProfileSlug)

	if profileSlug == "" {
		profileSlug = conv.ProfileSlug
	}
	return s.profileRepo.FindBySlug(profileSlug)
}

func (s *CompleteService) loadLlmProvider(modelSlug string, conv *conversation.Conversation) (llmcore.LlmProvider, factory.ModelEnum, error) {
	s.logger.Debug("Loading model: ", "model flag", modelSlug, "prev. conversation model", conv.ModelSlug)

	if modelSlug == "" {
		modelSlug = conv.ModelSlug
	}

	modelEnum, ok := factory.ParseFrom(modelSlug)
	if !ok {
		return nil, factory.DefaultModel, llmcore.ErrModelNotFound
	}

	return factory.NewProvider(modelEnum, s.httpClient, *s.config), modelEnum, nil
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
