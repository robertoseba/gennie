package models

import "github.com/robertoseba/gennie/internal/chat"

type IModel interface {
	Ask(question string, history *chat.ChatHistory) chat.Response
}
