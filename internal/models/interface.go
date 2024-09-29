package models

import (
	"github.com/robertoseba/gennie/internal/chat"
)

type IModel interface {
	/**
	* Complete chat receives a chat history with the last chat being the one that needs to be completed.
	* It also receives a system prompt that can be used to generate the answer.
	* Once succeded the model will fill out the answer in the last chat.
	 */
	CompleteChat(chatHistory *chat.ChatHistory, systemPrompt string) error
}
