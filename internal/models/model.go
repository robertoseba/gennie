package models

import (
	"errors"

	"github.com/robertoseba/gennie/internal/apiclient"
	"github.com/robertoseba/gennie/internal/conversation"
)

var ErrEmptyConversation = errors.New("Chat history is empty")
var ErrLastQuestionAlreadyAnswered = errors.New("Last chat is already completed with answer")

type BaseModel struct {
	client        apiclient.IApiClient
	modelProvider iModelProvider
}

func newBaseModel(client apiclient.IApiClient, modelProvider iModelProvider) *BaseModel {
	return &BaseModel{
		client:        client,
		modelProvider: modelProvider,
	}
}

func (m *BaseModel) CompleteChat(conversation *conversation.Conversation, systemPrompt string) error {

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

	postRes, err := m.client.Post(m.modelProvider.GetUrl(), payload, m.modelProvider.GetHeaders())

	if err != nil {
		return err
	}

	parsedResponse, err := m.modelProvider.ParseResponse(postRes)
	if err != nil {
		return err
	}

	return conversation.AnswerLastQuestion(parsedResponse)
}
