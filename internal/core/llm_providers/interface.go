package llmproviders

import (
	"context"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llm_providers/base"
)

type LlmProvider interface {
	Complete(ctx context.Context, conversation *conversation.Conversation, toolResult []base.ToolResult) <-chan base.ModelResponse
	SetSystemPrompt(systemPrompt string)
	SetTools(tools []base.Tool)
}
