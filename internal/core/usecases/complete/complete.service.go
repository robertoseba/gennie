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

func NewCompleteService(cr conversation.IConversationRepository, pr profile.IProfileRepository, httpClient *http.Client, config *config.Config) *CompleteService {
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

			mcpTools, shutdown, err := startMcpServers(ctx, activeProfile)
			defer shutdown()

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

		outputChan <- Response{Data: "Asking the model...", Type: RtLoading, Err: nil}
		for {
			resp := askLlm(ctx, activeConversation, llmProvider, toolResults, modelResponseChan)

			if !resp.IsToolCall() {
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

func askLlm(ctx context.Context, activeConversation *conversation.Conversation, llmProvider llmcore.LlmProvider, toolResults []llmcore.ToolResult, llmResponseChan chan<- Response) llmcore.LlmResponse {
	respChan := llmProvider.Complete(ctx, activeConversation, toolResults)

	var toolCallRequest llmcore.LlmResponse

	for llmResponse := range respChan {
		if llmResponse.StopReason == llmcore.StopReasonToolCall {
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

type shutdownFunc func()

func startMcpServers(ctx context.Context, activeProfile *profile.Profile) (map[string]toolDetails, shutdownFunc, error) {
	result := make(map[string]toolDetails)

	var cleanupFuncs []shutdownFunc

	for _, server := range activeProfile.McpServers {
		mcpClient, err := mcp.NewStdioClient(ctx, server.Command, server.Envs, server.Args)
		if err != nil {
			fmt.Printf("error %v", err)
			continue
		}

		mcpTools, err := mcpClient.ListTools(ctx)
		if err != nil {
			fmt.Printf("error %v", err)
			return nil, nil, fmt.Errorf("failed to list tools from MCP server %s: %w", server.Command, err)
		}

		// Filter tools based on profile allowed tools
		for _, tool := range mcpTools {
			if len(server.AllowedTools) == 0 || slices.Contains(server.AllowedTools, tool.Name) {
				result[tool.Name] = toolDetails{
					tool:             tool,
					mcpClient:        mcpClient,
					requiresApproval: server.RequiresApproval,
				}
				cleanupFuncs = append(cleanupFuncs, mcpClient.Close)
			}
		}
	}

	var shutdown shutdownFunc = func() {
		for _, cleanup := range cleanupFuncs {
			cleanup()
		}
	}

	return result, shutdown, nil
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
