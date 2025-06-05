package complete

import (
	"log/slog"

	"github.com/robertoseba/gennie/internal/core/conversation"
	"github.com/robertoseba/gennie/internal/core/llm/factory"
	"github.com/robertoseba/gennie/internal/core/profile"
)

type CompleteService struct {
	conversationRepo conversation.ConversationRepository
	profileRepo      profile.ProfileRepository
	logger           *slog.Logger
	providerFactory  factory.ProviderFactory
}

type Request struct {
	Question       string
	ProfileSlug    string
	ModelSlug      string
	IsFollowUp     bool
	AppendFilename string
}

type (
	ResponseType string
	Response     struct {
		Data string
		Err  error
		Type ResponseType
	}
)

const (
	RtLoading     ResponseType = "loading_info"
	RtModel       ResponseType = "model_info"
	RtProfile     ResponseType = "profile_info"
	RtApprovalReq ResponseType = "approval_request"
	RtLlmAnswer   ResponseType = "llm_answer"
)

func NewErrorResponse(err error) Response {
	return Response{Err: err}
}

func NewLoadingResponse(data string) Response {
	return Response{Data: data, Type: RtLoading}
}

func NewModelInfoResponse(data string) Response {
	return Response{Data: data, Type: RtModel}
}

func NewProfileInfoResponse(data string) Response {
	return Response{Data: data, Type: RtProfile}
}

func NewApprovalRequestResponse(data string) Response {
	return Response{Data: data, Type: RtApprovalReq}
}

func NewAnswerResponse(data string) Response {
	return Response{Data: data, Type: RtLlmAnswer}
}
