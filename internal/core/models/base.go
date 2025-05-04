package models

import (
	"context"
	"errors"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models/response"
	"github.com/robertoseba/gennie/internal/core/models/tools"
)

var (
	ErrEmptyConversation           = errors.New("there are no questions to answer")
	ErrLastQuestionAlreadyAnswered = errors.New("last conversation has already been answered")
)

type BaseModel struct {
	model         ModelEnum
	apiClient     IApiClient
	modelProvider iModelProvider
}

type ResponseType string

const (
	LoadingInfo     ResponseType = "loading_info"
	ModelInfo       ResponseType = "model_info"
	ProfileInfo     ResponseType = "profile_info"
	ApprovalRequest ResponseType = "approval_request"
)

type StreamResponse struct {
	Data string
	Err  error
	Type ResponseType
}

func newBaseModel(model ModelEnum, client IApiClient, modelProvider iModelProvider) *BaseModel {
	return &BaseModel{
		model:         model,
		apiClient:     client,
		modelProvider: modelProvider,
	}
}

func (m *BaseModel) Model() ModelEnum {
	return m.model
}

func (m *BaseModel) Complete(conversation *conversation.Conversation, systemPrompt string) error {
	if conversation.LastAnswer() != "" {
		return ErrLastQuestionAlreadyAnswered
	}

	if conversation.Len() == 0 {
		return ErrEmptyConversation
	}

	response := m.modelProvider.Complete(context.TODO(), conversation, nil)

	responseText := ""
	for r := range response {
		responseText += r.Text
	}

	conversation.AnswerLastQuestion(responseText)
	return nil
}

func (m *BaseModel) CompleteStreamable(ctx context.Context, conversation *conversation.Conversation, systemPrompt string, toolResult []tools.ToolResult) (<-chan response.ModelResponse, error) {
	return m.modelProvider.Complete(ctx, conversation, toolResult), nil
}

func (m *BaseModel) SetSystemPrompt(systemPrompt string) {
	m.modelProvider.SetSystemPrompt(systemPrompt)
}

func (m *BaseModel) SetTools(tools []tools.Tool) {
	m.modelProvider.SetTools(tools)
}

func (m *BaseModel) CanStream() bool {
	return m.modelProvider.CanStream()
}
