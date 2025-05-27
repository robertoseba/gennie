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
	"time"

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
	showToolOutput   bool
}

type CompleteService struct {
	conversationRepo conversation.IConversationRepository
	profileRepo      profile.IProfileRepository
	httpClient       *http.Client
	config           *config.Config
	logger           *slog.Logger
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
	filename := fmt.Sprintf("gennie-%s.log", time.Now().Format("20060102-150405"))
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	// TODO: shutdown file.Close() and mcpClient.Close() properly
	return &CompleteService{
		conversationRepo: cr,
		profileRepo:      pr,
		httpClient:       httpClient,
		config:           config,
		logger:           logger,
	}
}

func (s *CompleteService) Execute(input *InputDTO) (<-chan llmcore.CompleteResponse, error) {
	var conv *conversation.Conversation
	var err error

	conv, profile, llmProvider, err := s.processInput(input)
	if err != nil {
		return nil, err
	}

	llmProvider.SetSystemPrompt(profile.Data)

	modelResponseChan := make(chan llmcore.CompleteResponse)
	outputChan := s.pipeToSaveConversation(conv, modelResponseChan)

	go func() {
		defer close(modelResponseChan)

		outputChan <- llmcore.CompleteResponse{Data: conv.ModelSlug, Type: llmcore.ModelInfo, Err: nil}
		outputChan <- llmcore.CompleteResponse{Data: profile.Name, Type: llmcore.ProfileInfo, Err: nil}

		if len(profile.McpServers) > 0 {
			outputChan <- llmcore.CompleteResponse{Data: "Loading MCP Servers...", Type: llmcore.LoadingInfo, Err: nil}
			mcpTools, err := startMcpServers(profile)
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

		ctx := context.Background()
		toolResults := make([]llmcore.ToolResult, 0)

		// Keeps calling the model while it needs to return function calls
		outputChan <- llmcore.CompleteResponse{Data: "Asking the model...", Type: llmcore.LoadingInfo, Err: nil}
		for {
			resp := complete(ctx, llmProvider, conv, toolResults, modelResponseChan)

			// If does not need to send function call back to model than breaks out of the loop
			if resp.StopReason != llmcore.StopReasonTools {
				break
			}

			if s.tools[resp.FunctionCall.Name].requiresApproval {
				outputChan <- llmcore.CompleteResponse{Data: fmt.Sprintf("Can I run this tool: %s with parameters (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: llmcore.ApprovalRequest, Err: nil}
			}

			outputChan <- llmcore.CompleteResponse{Data: fmt.Sprintf("Using tool: %s -> (%s)", resp.FunctionCall.Name, resp.FunctionCall.Arguments), Type: llmcore.LoadingInfo, Err: nil}
			result := s.callTool(&resp)

			if s.tools[resp.FunctionCall.Name].showToolOutput {
				outputChan <- llmcore.CompleteResponse{Data: fmt.Sprintf("Tool %s returned: %s", resp.FunctionCall.Name, result.Result), Type: llmcore.ToolResultInfo, Err: nil}
			}

			toolResults = append(toolResults, *result)
		}
	}()

	return outputChan, nil
}

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

func complete(ctx context.Context, llmProvider llmcore.LlmProvider, conv *conversation.Conversation, toolResults []llmcore.ToolResult, outputChan chan<- llmcore.CompleteResponse) llmcore.LlmResponse {
	respChan := llmProvider.Complete(ctx, conv, toolResults)

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

func (s *CompleteService) callTool(llmResponse *llmcore.LlmResponse) *llmcore.ToolResult {
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

	toolResponse, err := s.tools[llmResponse.FunctionCall.Name].mcpClient.ExecTool(context.TODO(), llmResponse.FunctionCall.Name, args)
	if err != nil {
		return llmcore.NewToolResponseError(nil, err)
	}

	return llmcore.NewToolResponseFrom(llmResponse, toolResponse)
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
					showToolOutput:   server.ShowToolOutput,
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

	model, modelEnum, err := s.loadLlmProvider(input.ModelSlug, conv)
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

func (s *CompleteService) loadLlmProvider(modelSlug string, conv *conversation.Conversation) (llmcore.LlmProvider, factory.ModelEnum, error) {
	if modelSlug == "" {
		modelSlug = conv.ModelSlug
	}

	modelEnum, ok := factory.ParseFrom(modelSlug)
	if !ok {
		return nil, factory.DefaultModel, llmcore.ErrModelNotFound
	}

	return factory.NewProvider(modelEnum, s.logger, s.httpClient, *s.config), modelEnum, nil
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
