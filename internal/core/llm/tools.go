package llm

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

func (t ToolInputSchema) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	if t.Type == "" {
		m["type"] = "object"
	} else {
		m["type"] = t.Type
	}

	// Marshal Properties to '{}' rather than `nil` when its length equals zero
	if t.Properties != nil {
		m["properties"] = t.Properties
	}

	if len(t.Required) > 0 {
		m["required"] = t.Required
	}

	return json.Marshal(m)
}

type ToolResult struct {
	ID        string // block id - function identifier returned in the ModelResponse->FunctionCall
	Name      string
	Result    []byte
	Arguments json.RawMessage
	Error     error
}

func NewToolResponseFrom(funcCall FunctionCall, toolResult []byte) *ToolResult {
	return &ToolResult{
		ID:        funcCall.ID,
		Name:      funcCall.Name,
		Arguments: funcCall.Arguments,
		Result:    toolResult,
		Error:     nil,
	}
}

func NewToolResponseError(funcCall FunctionCall, err error) *ToolResult {
	return &ToolResult{
		ID:     funcCall.ID,
		Name:   funcCall.Name,
		Result: nil,
		Error:  err,
	}
}

func (t *ToolResult) IsError() bool {
	return t.Error != nil
}
