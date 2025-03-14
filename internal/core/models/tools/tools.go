package tools

import (
	"encoding/json"

	"github.com/robertoseba/gennie/internal/core/models/response"
)

// Example of a tool marshalled
// "tools": [
//
//	  {
//	    "name": "get_weather",
//	    "description": "Get the current weather in a given location",
//	    "input_schema": {
//	      "type": "object",
//	      "properties": {
//	        "location": {
//	          "type": "string",
//	          "description": "The city and state, e.g. San Francisco, CA"
//	        }
//	      },
//	      "required": ["location"]
//	    }
//	  }
//	],

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

func NewToolResponseFrom(modelResponse *response.ModelResponse, toolResult []byte) *ToolResult {
	return &ToolResult{
		ID:        modelResponse.FunctionCall.ID,
		Name:      modelResponse.FunctionCall.Name,
		Arguments: modelResponse.FunctionCall.Arguments,
		Result:    toolResult,
		Error:     nil,
	}
}

func NewToolResponseError(modelResponse *response.ModelResponse, err error) *ToolResult {
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
