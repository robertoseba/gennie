package models

import (
	"errors"

	"github.com/robertoseba/gennie/internal/conversation"
	"github.com/robertoseba/gennie/internal/httpclient"
)

var ErrEmptyChatHistory = errors.New("Chat history is empty")
var ErrLastChatCompleted = errors.New("Last chat is already completed with answer")

type BaseModel struct {
	client        httpclient.IApiClient
	modelProvider IModelProvider
}

func NewBaseModel(client httpclient.IApiClient, modelProvider IModelProvider) *BaseModel {
	return &BaseModel{
		client:        client,
		modelProvider: modelProvider,
	}
}

func (m *BaseModel) CompleteChat(chatHistory *conversation.Conversation, systemPrompt string) error {

	if chatHistory.LastAnswer() != "" {
		return ErrLastChatCompleted
	}

	if chatHistory.Len() == 0 {
		return ErrEmptyChatHistory
	}

	payload, err := m.modelProvider.PreparePayload(chatHistory, systemPrompt)
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

	return chatHistory.AnswerLastQuestion(parsedResponse)
}
