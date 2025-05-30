package complete

import (
	"context"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type ctxKey string

var ctxRequestKey ctxKey = "request"

type Request struct {
	Question       string
	ProfileSlug    string
	ModelSlug      string
	IsFollowUp     bool
	AppendFilename string
	Conversation   *conversation.Conversation
	Profile        *profile.Profile
}

func AddRequestToCtx(ctx context.Context, req *Request) context.Context {
	return context.WithValue(ctx, ctxRequestKey, req)
}

func GetRequestFromCtx(ctx context.Context) (*Request, bool) {
	req, ok := ctx.Value(ctxRequestKey).(*Request)
	if !ok {
		return nil, false
	}
	return req, true
}
