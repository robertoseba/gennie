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
	"github.com/robertoseba/gennie/internal/core/mcp"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type CompleteService struct {
	conversationRepo conversation.IConversationRepository
	profileRepo      profile.IProfileRepository
	httpClient       *http.Client
	config           *config.Config
	logger           *slog.Logger
	tools            map[string]toolDetails // each toolName maps to a mcpClient so we can make a request
}

type toolDetails struct {
	mcpClient        mcp.McpClientInterface
	tool             llmcore.Tool
	requiresApproval bool
}

type (
	ResponseType string
	Response     struct {
		Data string
		Err  error
		Type ResponseType
	}
)

const (
	RtLoading     ResponseType = "loading_info"
	RtModel       ResponseType = "model_info"
	RtProfile     ResponseType = "profile_info"
	RtApprovalReq ResponseType = "approval_request"
	RtLlmAnswer   ResponseType = "llm_answer"
)

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

func (s *CompleteService) Execute(ctx context.Context, req Request) (<-chan Response, error) {
	var err error

	activeConversation, err := s.conversationRepo.LoadActive()
	if err != nil {
		return nil, err
	}

	llmProvider, activeProfile, err := s.processRequestToConversation(req, activeConversation)
	if err != nil {
		return nil, err
	}

	llmProvider.SetSystemPrompt(activeProfile.Data)

	modelResponseChan := make(chan Response)
	outputChan := s.pipeToSaveConversation(ctx, activeConversation, modelResponseChan)

	go func() {
		defer close(modelResponseChan)

		outputChan <- Response{Data: activeConversation.ModelSlug, Type: RtModel}
		outputChan <- Response{Data: activeProfile.Name, Type: RtProfile}

		if len(activeProfile.McpServers) > 0 {
			outputChan <- Response{Data: "Loading MCP Servers...", Type: RtLoading}

			mcpTools, err := startMcpServers(ctx, activeProfile)
			if err != nil {
				outputChan <- Response{Err: err}
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
		outputChan <- Response{Data: "Asking the model...", Type: RtLoading, Err: nil}
		for {
			resp := complete(ctx, activeConversation, llmProvider, toolResults, modelResponseChan)

			// If does not need to send function call back to model than breaks out of the loop
			if resp.StopReason != llmcore.StopReasonTools {
				break
			}

			if s.tools[resp.FunctionCall.Name].requiresApproval {
				outputChan <- Response{Data: fmt.Sprintf("Can I run this tool: %s with parameters (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: RtApprovalReq, Err: nil}
			}

			outputChan <- Response{Data: fmt.Sprintf("Using tool: %s -> (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: RtLoading, Err: nil}
			result := s.callTool(ctx, &resp)

			toolResults = append(toolResults, *result)
		}
	}()

	return outputChan, nil
}

func (s *CompleteService) pipeToSaveConversation(ctx context.Context, conv *conversation.Conversation, inputChan <-chan Response) chan Response {
	outputChan := make(chan Response)

	go func() {
		defer close(outputChan)

		convBuffer := strings.Builder{}
		for msg := range inputChan {
			select {
			case <-ctx.Done():
				s.logger.Debug("Context done, stopping conversation saving", "error", ctx.Err())
				return
			case outputChan <- msg:
				if msg.Err == nil {
					convBuffer.WriteString(msg.Data)
				}
			}
		}

		err := conv.AnswerLastQuestion(convBuffer.String())
		if err != nil {
			outputChan <- Response{Err: err}
		}
		err = s.conversationRepo.SaveAsActive(conv)
		if err != nil {
			outputChan <- Response{Err: err}
		}
	}()

	return outputChan
}

func complete(ctx context.Context, activeConversation *conversation.Conversation, llmProvider llmcore.LlmProvider, toolResults []llmcore.ToolResult, llmResponseChan chan<- Response) llmcore.LlmResponse {
	respChan := llmProvider.Complete(ctx, activeConversation, toolResults)

	var toolCallRequest llmcore.LlmResponse

	for llmResponse := range respChan {
		if llmResponse.StopReason == llmcore.StopReasonTools {
			toolCallRequest = llmResponse
			continue
		}

		llmResponseChan <- Response{Data: llmResponse.Text, Err: llmResponse.Error}
	}

	return toolCallRequest
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

func startMcpServers(ctx context.Context, activeProfile *profile.Profile) (map[string]toolDetails, error) {
	// TODO: we are not closing these clients
	returnTools := make(map[string]toolDetails)

	for _, server := range activeProfile.McpServers {

		mcpServer, err := mcp.NewStdioClient(ctx, server.Command, server.Envs, server.Args)
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

func (s *CompleteService) processRequestToConversation(req Request, activeConversation *conversation.Conversation) (llmcore.LlmProvider, *profile.Profile, error) {
	activeProfile, err := s.loadProfile(req.ProfileSlug, activeConversation)
	if err != nil {
		return nil, activeProfile, err
	}
	activeConversation.SetProfileTo(activeProfile.Slug)

	if req.ModelSlug != "" {
		activeConversation.ModelSlug = req.ModelSlug
		activeConversation.SetModelTo(req.ModelSlug)
	}

	llmProvider, err := factory.NewProvider(activeConversation.ModelSlug, s.httpClient, *s.config)
	if err != nil {
		return nil, activeProfile, err
	}

	if !req.IsFollowUp {
		activeConversation.Clear()
	}

	if req.AppendFilename != "" {
		content, err := os.ReadFile(req.AppendFilename)
		if err != nil {
			return nil, activeProfile, err
		}
		req.Question += "\n" + string(content)
	}
	activeConversation.NewQuestion(req.Question)

	return llmProvider, activeProfile, nil
}

func (c *CompleteService) loadProfile(profileSlug string, activeConversation *conversation.Conversation) (*profile.Profile, error) {
	if profileSlug == "" {
		profileSlug = activeConversation.ProfileSlug
	}

	activeProfile, err := c.profileRepo.FindBySlug(profileSlug)
	if err != nil {
		return nil, err
	}

	return activeProfile, nil
}
