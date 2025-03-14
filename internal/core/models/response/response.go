package response

import "encoding/json"

type StopReason string

const (
	StopReasonTools StopReason = "toolCalling"
	StopReasonEnd   StopReason = "endResponse"
	StopReasonError StopReason = "error"
	StopReasonNone  StopReason = "none"
)

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
