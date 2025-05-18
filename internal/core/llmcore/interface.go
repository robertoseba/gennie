package llmcore

import (
	"context"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llmcore/entities"
)

type LlmProvider interface {
	Complete(ctx context.Context, conversation *conversation.Conversation, toolResult []entities.ToolResult) <-chan entities.LlmResponse
	SetSystemPrompt(systemPrompt string)
	SetTools(tools []entities.Tool)
}
