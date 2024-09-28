package models

import (
	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/models/profile"
)

type IModel interface {
	//TODO: remove profile dependency and use only the profile data needed. Call it system prompt
	//TODO: should model receive a chat history and question or a chat history prepare for next question
	//TODO: Ask would then be replaced by CompleteChat where the model would receive
	// the chat history and fillout the next answer.
	// In this scenario, chatHistory should also have profile.Data instead of being sent to the model.
	Ask(question string, profile *profile.Profile, history *chat.ChatHistory) (*chat.Chat, error)
	Model() string
}
