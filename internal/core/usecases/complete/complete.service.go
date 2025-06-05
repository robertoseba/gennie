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
		logger:           slog.Default(),
		providerFactory:  factory.NewProviderFactory(httpClient, *config),
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

	llmProvider, err := s.providerFactory.CreateProvider(currConversation.ModelSlug)
	if err != nil {
		return nil, err
	}

	llmProvider.SetSystemPrompt(currProfile.Data)

	outputChan := make(chan Response)

	go func() {
		defer close(outputChan)

		outputChan <- NewModelInfoResponse(currConversation.ModelSlug)
		outputChan <- NewProfileInfoResponse(currProfile.Slug)

		// setup mcps
		mcpGroup := mcp.NewGroup()
		if len(currProfile.McpServers) > 0 {
			outputChan <- NewLoadingResponse("Setting up MCP servers...")
			for name, mcpServer := range currProfile.McpServers {
				outputChan <- NewLoadingResponse(fmt.Sprintf("Adding MCP server: %s", name))
				err := mcpGroup.Add(ctx, mcpServer)
				if err != nil {
					outputChan <- NewErrorResponse(fmt.Errorf("failed to add MCP server %s: %w", name, err))
				}
			}

			llmProvider.SetTools(mcpGroup.ListTools())
		}
		defer mcpGroup.Shutdown()

		outputChan <- NewLoadingResponse("Asking question...")
		toolResults := make([]llm.ToolResult, 0)
		answerAcc := strings.Builder{}

		for {
			llmCompleteChan := llmProvider.Complete(ctx, currConversation, toolResults)

			var toolCallRequest *llm.Response

			for llmResponse := range llmCompleteChan {
				if llmResponse.IsError() {
					outputChan <- NewErrorResponse(llmResponse.Error)
					continue
				}

				if llmResponse.IsToolCall() {
					toolCallRequest = &llmResponse
				}

				outputChan <- NewAnswerResponse(llmResponse.Text)

				answerAcc.WriteString(llmResponse.Text)
			}

			if toolCallRequest == nil {
				break
			}

			for _, toolReq := range toolCallRequest.FunctionCalls {
				if mcpGroup.RequiresApproval(toolReq.Name) {
					outputChan <- NewApprovalRequestResponse(fmt.Sprintf("Tool call request for '%s' with params (%s) requires approval.",
						toolReq.Name, toolReq.Arguments))
				}

				outputChan <- NewLoadingResponse(fmt.Sprintf("Using tool: %s -> (%s)", toolReq.Name, toolReq.Arguments))

				result := s.callTool(ctx, toolReq, mcpGroup)
				toolResults = append(toolResults, *result)
			}

			outputChan <- NewLoadingResponse("Sending tool response to model...")
		}

		s.saveConversation(currConversation, answerAcc.String())
	}()

	return outputChan, nil
}

func (s *CompleteService) saveConversation(conv *conversation.Conversation, answer string) error {
	err := conv.AnswerLastQuestion(answer)
	if err != nil {
		return fmt.Errorf("failed to answer the last question: %w", err)
	}

	err = s.conversationRepo.SaveAsActive(conv)
	if err != nil {
		return fmt.Errorf("failed to save the conversation: %w", err)
	}

	return nil
}

func (s *CompleteService) callTool(ctx context.Context, funcCall llm.FunctionCall, mcpGroup *mcp.Group) *llm.ToolResult {
	var args map[string]any

	if len(funcCall.Arguments) > 0 {
		args = make(map[string]any)
		err := json.Unmarshal(funcCall.Arguments, &args)
		if err != nil {
			return llm.NewToolResponseError(funcCall, err)
		}
	}

	toolResponse, err := mcpGroup.ExecTool(ctx, funcCall.Name, args)
	if err != nil {
		s.logger.Error("Failed to execute tool", "toolName", funcCall.Name, "error", err)
		return llm.NewToolResponseError(funcCall, err)
	}
	s.logger.Debug("Tool response", "toolCallId", funcCall.ID, "toolName", funcCall.Name, "response", toolResponse)

	return llm.NewToolResponseFrom(funcCall, []byte(toolResponse))
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
