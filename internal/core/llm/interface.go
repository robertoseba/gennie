package llm

import (
	"context"

	"github.com/robertoseba/gennie/internal/core/conversation"
)

type Provider interface {
	Complete(ctx context.Context, conversation *conversation.Conversation, toolResult []ToolResult) <-chan Response
	SetSystemPrompt(systemPrompt string)
	SetTools(tools []Tool)
}
