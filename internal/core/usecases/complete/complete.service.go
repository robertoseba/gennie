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
	"github.com/robertoseba/gennie/internal/core/llm"
	"github.com/robertoseba/gennie/internal/core/llm/factory"
	"github.com/robertoseba/gennie/internal/core/mcp"
	"github.com/robertoseba/gennie/internal/core/profile"
)

func NewCompleteService(cr conversation.ConversationRepository, pr profile.ProfileRepository, httpClient *http.Client, config *config.Config) *CompleteService {
	return &CompleteService{
		conversationRepo: cr,
		profileRepo:      pr,
		httpClient:       httpClient,
		config:           config,
		logger:           slog.Default(),
	}
}

func (s *CompleteService) Execute(ctx context.Context, req Request) (<-chan Response, error) {
	previousConversation, err := s.conversationRepo.LoadActive()
	if err != nil {
		return nil, err
	}

	currProfile, currConversation, err := s.processRequest(req, previousConversation)
	if err != nil {
		return nil, err
	}

	llmProvider, err := factory.NewProvider(currConversation.ModelSlug, s.httpClient, *s.config)
	if err != nil {
		return nil, err
	}

	llmProvider.SetSystemPrompt(currProfile.Data)

	outputChan := make(chan Response)

	go func() {
		defer close(outputChan)

		outputChan <- Response{Data: currConversation.ModelSlug, Type: RtModel}
		outputChan <- Response{Data: currProfile.Name, Type: RtProfile}

		// setup mcps
		mcpGroup := mcp.NewGroup()
		if len(currProfile.McpServers) > 0 {
			outputChan <- Response{Data: "Loading MCP Servers...", Type: RtLoading}
			for _, mcpServer := range currProfile.McpServers {
				err := mcpGroup.Add(ctx, mcpServer)
				if err != nil {
					outputChan <- Response{Err: fmt.Errorf("failed to add MCP server %s: %w", mcpServer.Cmd, err)}
				}
			}

			llmProvider.SetTools(mcpGroup.ListTools())
		}
		defer mcpGroup.Shutdown()

		outputChan <- Response{Data: "Asking the model...", Type: RtLoading, Err: nil}

		toolResults := make([]llm.ToolResult, 0)
		answer := strings.Builder{}
		for {
			llmCompleteChan := llmProvider.Complete(ctx, currConversation, toolResults)

			toolCallRequest := llm.Response{}
			for llmResponse := range llmCompleteChan {
				if llmResponse.StopReason == llm.StopReasonToolCall {
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

		s.answerConversation(ctx, currConversation, answer.String())
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

func (s *CompleteService) callTool(ctx context.Context, llmResponse *llm.Response, mcpGroup *mcp.Group) *llm.ToolResult {
	var args map[string]any
	if len(llmResponse.FunctionCall.Arguments) > 0 {
		args = make(map[string]any)
		err := json.Unmarshal([]byte(llmResponse.FunctionCall.Arguments), &args)
		if err != nil {
			return llm.NewToolResponseError(nil, err)
		}
	}

	toolResponse, err := mcpGroup.ExecTool(ctx, llmResponse.FunctionCall.Name, args)
	if err != nil {
		s.logger.Error("Failed to execute tool", "toolName", llmResponse.FunctionCall.Name, "error", err)
		return llm.NewToolResponseError(llmResponse, err)
	}
	s.logger.Debug("Tool response", "toolName", llmResponse.FunctionCall.Name, "response", toolResponse)

	return llm.NewToolResponseFrom(llmResponse, []byte(toolResponse))
}

func (s *CompleteService) processRequest(req Request, previousConversation *conversation.Conversation) (*profile.Profile, *conversation.Conversation, error) {
	currConversation := conversation.NewConversation(previousConversation.ProfileSlug, previousConversation.ModelSlug)

	currProfile, err := s.loadProfile(req.ProfileSlug, previousConversation)
	if err != nil {
		return currProfile, nil, err
	}
	currConversation.SetProfileTo(currProfile.Slug)

	if req.ModelSlug != "" {
		currConversation.ModelSlug = req.ModelSlug
		currConversation.SetModelTo(req.ModelSlug)
	}

	if req.IsFollowUp {
		currConversation.QAs = append(currConversation.QAs, previousConversation.QAs...)
	}

	if req.AppendFilename != "" {
		content, err := os.ReadFile(req.AppendFilename)
		if err != nil {
			return currProfile, nil, err
		}
		req.Question += "\n" + string(content)
	}
	currConversation.NewQuestion(req.Question)

	return currProfile, currConversation, nil
}

func (c *CompleteService) loadProfile(profileSlug string, currConversation *conversation.Conversation) (*profile.Profile, error) {
	if profileSlug == "" {
		profileSlug = currConversation.ProfileSlug
	}

	currProfile, err := c.profileRepo.FindBySlug(profileSlug)
	if err != nil {
		return nil, err
	}

	return currProfile, nil
}
