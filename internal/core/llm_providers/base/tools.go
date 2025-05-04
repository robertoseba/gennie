package base

import (
	"encoding/json"
)

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema ToolInputSchema `json:"input_schema"`
}

type ToolInputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

type ToolResult struct {
	ID        string // block id - function identifier returned in the ModelResponse->FunctionCall
	Name      string
	Result    []byte
	Arguments json.RawMessage
	Error     error
}

func NewToolResponseFrom(modelResponse *ModelResponse, toolResult []byte) *ToolResult {
	return &ToolResult{
		ID:        modelResponse.FunctionCall.ID,
		Name:      modelResponse.FunctionCall.Name,
		Arguments: modelResponse.FunctionCall.Arguments,
		Result:    toolResult,
		Error:     nil,
	}
}

func NewToolResponseError(modelResponse *ModelResponse, err error) *ToolResult {
	return &ToolResult{
		ID:     modelResponse.FunctionCall.ID,
		Name:   modelResponse.FunctionCall.Name,
		Result: nil,
		Error:  err,
	}
}

func (t *ToolResult) IsError() bool {
	return t.Error != nil
}
