package llmcore

import (
	"context"

	"github.com/robertoseba/gennie/internal/core/conversation"
)

type LlmProvider interface {
	Complete(ctx context.Context, conversation *conversation.Conversation, toolResult []ToolResult) <-chan LlmResponse
	SetSystemPrompt(systemPrompt string)
	SetTools(tools []Tool)
}
