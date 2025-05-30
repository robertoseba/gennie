package llmcore

import "encoding/json"

type (
	ResponseType string
	StopReason   string
)

const (
	StopReasonTools StopReason = "tool_calling"
	StopReasonEnd   StopReason = "end_response"
	StopReasonError StopReason = "error"
	StopReasonNone  StopReason = ""
)

// This is the response converted from a providers response
// Every provider should convert their response to this base response
type LlmResponse struct {
	Text         string
	Error        error
	StopReason   StopReason
	FunctionCall FunctionCall // TODO: maybe this should be an slice of function calls
}

type FunctionCall struct {
	ID        string
	Name      string
	Arguments json.RawMessage
}
