package common

import (
	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/profile"
)

type IPersistence interface {
	Save() error
	Clear() error
	GetConfig() Config
	SetConfig(Config)
	GetProfile(string) (profile.Profile, error)
	GetChatHistory() chat.ChatHistory
	SetChatHistory(chat.ChatHistory)
}
