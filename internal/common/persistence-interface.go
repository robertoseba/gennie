package common

import (
	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/profile"
)

type IPersistence interface {
	GetConfig() Config
	SetConfig(Config)
	GetProfile(string) (*profile.Profile, error)
	GetProfileSlugs() []string
	GetChatHistory() chat.ChatHistory
	SetChatHistory(chat.ChatHistory)
	GetCacheFilePath() string
}
