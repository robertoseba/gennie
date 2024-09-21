package models

import (
	"github.com/robertoseba/gennie/internal/chat"
	"github.com/robertoseba/gennie/internal/models/profile"
)

type IModel interface {
	Ask(question string, profile *profile.Profile, history *chat.ChatHistory) (*chat.Chat, error)
}
