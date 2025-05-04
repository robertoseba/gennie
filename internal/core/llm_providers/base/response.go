package base

import "encoding/json"

type (
	ResponseType string
	StopReason   string
)

const (
	LoadingInfo     ResponseType = "loading_info"
	ModelInfo       ResponseType = "model_info"
	ProfileInfo     ResponseType = "profile_info"
	ApprovalRequest ResponseType = "approval_request"
)

const (
	StopReasonTools StopReason = "tool_calling"
	StopReasonEnd   StopReason = "end_response"
	StopReasonError StopReason = "error"
	StopReasonNone  StopReason = "none"
)

type StreamResponse struct {
	Data string
	Err  error
	Type ResponseType
}

type ModelResponse struct {
	Text         string
	Error        error
	StopReason   StopReason
	FunctionCall FunctionCall
}

type FunctionCall struct {
	ID        string
	Name      string
	Arguments json.RawMessage
}
