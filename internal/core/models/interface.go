package models

import (
	"context"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/models/response"
	"github.com/robertoseba/gennie/internal/core/models/tools"
)

type ProviderStreamParser func(b []byte) (string, error)

// The Provider is only responsible for preparing the payload,
// formatting it accordinly to the model's requirements
// and parsing the response back to the system.
type iModelProvider interface {
	CanStream() bool
	Complete(ctx context.Context, conversation *conversation.Conversation, toolResult []tools.ToolResult) <-chan response.ModelResponse
	SetSystemPrompt(systemPrompt string)
	SetTools(tools []tools.Tool)
}

type IApiClient interface {
	Post(url string, body string, headers map[string]string) ([]byte, error)
	PostWithStreaming(url string, body string, headers map[string]string, parser ProviderStreamParser) <-chan StreamResponse
}
