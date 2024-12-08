package models

import (
	"errors"

	"github.com/robertoseba/gennie/internal/core/conversation"
)

var ErrEmptyConversation = errors.New("there are no questions to answer")
var ErrLastQuestionAlreadyAnswered = errors.New("last conversation has already been answered")

type BaseModel struct {
	model         ModelEnum
	apiClient     IApiClient
	modelProvider iModelProvider
}

type StreamResponse struct {
	Data string
	Err  error
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

	payload, err := m.modelProvider.PreparePayload(conversation, systemPrompt, false)
	if err != nil {
		return err
	}

	postRes, err := m.apiClient.Post(m.modelProvider.GetUrl(), payload, m.modelProvider.GetHeaders())

	if err != nil {
		return err
	}

	parsedResponse, err := m.modelProvider.ParseResponse(postRes)
	if err != nil {
		return err
	}

	return conversation.AnswerLastQuestion(parsedResponse)
}

func (m *BaseModel) CompleteStreamable(conversation *conversation.Conversation, systemPrompt string) (<-chan StreamResponse, error) {
	payload, err := m.modelProvider.PreparePayload(conversation, systemPrompt, true)
	if err != nil {
		return nil, err
	}

	outputChan := m.apiClient.PostWithStreaming(m.modelProvider.GetUrl(), payload, m.modelProvider.GetHeaders(), m.modelProvider.GetStreamParser())
	return outputChan, nil
}

func (m *BaseModel) CanStream() bool {
	return m.modelProvider.CanStream()
}
