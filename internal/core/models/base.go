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

	payload, err := m.modelProvider.PreparePayload(conversation, systemPrompt)
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
