package llmcore

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
	LlmAnswer       ResponseType = "llm_answer"
	ToolResultInfo  ResponseType = "tool_result_info"
)

const (
	StopReasonTools StopReason = "tool_calling"
	StopReasonEnd   StopReason = "end_response"
	StopReasonError StopReason = "error"
	StopReasonNone  StopReason = "none"
)

// This is the response the use case sends back to the client(console)
// Based on the response type, the client can decide how to handle it (loading bar, request users approval, etc)
type CompleteResponse struct {
	Data string
	Err  error
	Type ResponseType
}

// This is the response converted from a providers response
// Every provider should convert their response to this base response
type LlmResponse struct {
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
