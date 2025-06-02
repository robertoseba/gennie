package complete

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

	activeProfile, err := s.processRequestToConversation(req, activeConversation)
	if err != nil {
		return nil, err
	}

	llmProvider, err := factory.NewProvider(activeConversation.ModelSlug, s.httpClient, *s.config)
	if err != nil {
		return nil, err
	}

	llmProvider.SetSystemPrompt(activeProfile.Data)

	outputChan := make(chan Response)

	go func() {
		defer close(outputChan)

		outputChan <- Response{Data: activeConversation.ModelSlug, Type: RtModel}
		outputChan <- Response{Data: activeProfile.Name, Type: RtProfile}

		// setup mcps
		mcpGroup := mcp.NewGroup()
		if len(activeProfile.McpServers) > 0 {
			outputChan <- Response{Data: "Loading MCP Servers...", Type: RtLoading}
			for _, mcpProfile := range activeProfile.McpServers {
				err := mcpGroup.Add(ctx, mcpProfile.Command, mcpProfile.Envs, mcpProfile.Args, mcpProfile.RequiresApproval, mcpProfile.AllowedTools)
				if err != nil {
					outputChan <- Response{Err: fmt.Errorf("failed to add MCP server %s: %w", mcpProfile.Command, err)}
				}
			}
			tools, err := mcpGroup.ListTools()
			if err != nil {
				outputChan <- Response{Err: fmt.Errorf("failed to list tools from MCP servers: %w", err)}
			}
			llmProvider.SetTools(tools)
		}
		defer mcpGroup.Shutdown()

		outputChan <- Response{Data: "Asking the model...", Type: RtLoading, Err: nil}

		toolResults := make([]llmcore.ToolResult, 0)
		answer := strings.Builder{}
		for {
			llmCompleteChan := llmProvider.Complete(ctx, activeConversation, toolResults)

			toolCallRequest := llmcore.LlmResponse{}
			for llmResponse := range llmCompleteChan {
				if llmResponse.StopReason == llmcore.StopReasonToolCall {
					toolCallRequest = llmResponse
					continue
				}

				outputChan <- Response{Data: llmResponse.Text, Err: llmResponse.Error}
				answer.WriteString(llmResponse.Text)
			}

			if !toolCallRequest.IsToolCall() {
				break
			}

			if mcpGroup.RequiresApproval(toolCallRequest.FunctionCall.Name) {
				outputChan <- Response{Data: fmt.Sprintf("Can I run this tool: %s with parameters (%s)", toolCallRequest.FunctionCall.Name, toolCallRequest.FunctionCall.Arguments), Type: RtApprovalReq, Err: nil}
			}

			outputChan <- Response{Data: fmt.Sprintf("Using tool: %s -> (%s)", toolCallRequest.FunctionCall.Name, toolCallRequest.FunctionCall.Arguments), Type: RtLoading, Err: nil}
			result := s.callTool(ctx, &toolCallRequest, mcpGroup)

			toolResults = append(toolResults, *result)
		}

		s.answerConversation(ctx, activeConversation, answer.String())
	}()

	return outputChan, nil
}

func (s *CompleteService) answerConversation(ctx context.Context, conv *conversation.Conversation, answer string) error {
	err := conv.AnswerLastQuestion(answer)
	if err != nil {
		return fmt.Errorf("failed to answer the last question: %w", err)
	}

	err = s.conversationRepo.SaveAsActive(conv)
	if err != nil {
		return fmt.Errorf("failed to answer the last question: %w", err)
	}

	return nil
}

func (s *CompleteService) callTool(ctx context.Context, llmResponse *llmcore.LlmResponse, mcpGroup *mcp.Group) *llmcore.ToolResult {
	var args map[string]any
	if len(llmResponse.FunctionCall.Arguments) > 0 {
		args = make(map[string]any)
		err := json.Unmarshal([]byte(llmResponse.FunctionCall.Arguments), &args)
		if err != nil {
			return llmcore.NewToolResponseError(nil, err)
		}
	}

	toolResponse, err := mcpGroup.ExecTool(ctx, llmResponse.FunctionCall.Name, args)
	if err != nil {
		s.logger.Error("Failed to execute tool", "toolName", llmResponse.FunctionCall.Name, "error", err)
		return llmcore.NewToolResponseError(llmResponse, err)
	}
	s.logger.Debug("Tool response", "toolName", llmResponse.FunctionCall.Name, "response", toolResponse)

	return llmcore.NewToolResponseFrom(llmResponse, []byte(toolResponse))
}

func (s *CompleteService) processRequestToConversation(req Request, activeConversation *conversation.Conversation) (*profile.Profile, error) {
	activeProfile, err := s.loadProfile(req.ProfileSlug, activeConversation)
	if err != nil {
		return activeProfile, err
	}
	activeConversation.SetProfileTo(activeProfile.Slug)

	if req.ModelSlug != "" {
		activeConversation.ModelSlug = req.ModelSlug
		activeConversation.SetModelTo(req.ModelSlug)
	}

	if !req.IsFollowUp {
		activeConversation.Clear()
	}

	if req.AppendFilename != "" {
		content, err := os.ReadFile(req.AppendFilename)
		if err != nil {
			return activeProfile, err
		}
		req.Question += "\n" + string(content)
	}
	activeConversation.NewQuestion(req.Question)

	return activeProfile, nil
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
