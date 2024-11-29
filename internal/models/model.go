package models

import (
	"errors"

	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/httpclient"
)

var ErrEmptyChatHistory = errors.New("Chat history is empty")
var ErrLastChatCompleted = errors.New("Last chat is already completed with answer")

type BaseModel struct {
	client        httpclient.IHttpClient
	modelProvider IModelProvider
}

func NewBaseModel(client httpclient.IHttpClient, modelProvider IModelProvider) *BaseModel {
	return &BaseModel{
		client:        client,
		modelProvider: modelProvider,
	}
}

func (m *BaseModel) CompleteChat(chatHistory *chat.ChatHistory, systemPrompt string) error {
	lastChat, ok := chatHistory.LastChat()
	if !ok {
		return ErrLastChatCompleted
	}

	if lastChat.GetAnswer() != "" {
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

	return chatHistory.SetNewAnswerToLastChat(parsedResponse)
}
