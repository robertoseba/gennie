package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/robertoseba/gennie/internal/core/llmcore/entities"
)

func parseToolUseFrom(message *anthropic.Message) []entities.LlmResponse {
	var toolResponses []entities.LlmResponse

	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.ToolUseBlock:
			response := entities.LlmResponse{
				StopReason: entities.StopReasonTools,
				FunctionCall: entities.FunctionCall{
					ID:        variant.ID,
					Name:      variant.Name,
					Arguments: variant.Input,
				},
			}
			toolResponses = append(toolResponses, response)
		}
	}
	return toolResponses
}

func convertToolsToProvider(tools []entities.Tool) []anthropic.ToolUnionParam {
	if len(tools) == 0 {
		return nil
	}

	toolReturn := make([]anthropic.ToolUnionParam, 0)
	for _, tool := range tools {
		toolParam := anthropic.ToolParam{
			Name:        tool.Name,
			Description: anthropic.String(tool.Description),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: tool.InputSchema.Properties,
			},
		}
		toolReturn = append(toolReturn, anthropic.ToolUnionParam{OfTool: &toolParam})
	}

	return toolReturn
}

func addToolResultsToMessages(
	messages []anthropic.MessageParam,
	toolResults []entities.ToolResult,
) []anthropic.MessageParam {
	messageBlock := []anthropic.ContentBlockParamUnion{}
	messageAssistantBlock := []anthropic.ContentBlockParamUnion{}

	for _, tr := range toolResults {
		// add back assistant block with tool use
		useBlock := anthropic.ToolUseBlockParam{
			ID:    tr.ID,
			Name:  tr.Name,
			Input: tr.Arguments,
		}
		toolUseBlock := anthropic.ContentBlockParamUnion{OfRequestToolUseBlock: &useBlock}
		messageAssistantBlock = append(messageAssistantBlock, toolUseBlock)

		if tr.IsError() {
			messageBlock = append(messageBlock, anthropic.NewToolResultBlock(tr.ID, tr.Error.Error(), true))
		} else {
			messageBlock = append(messageBlock, anthropic.NewToolResultBlock(tr.ID, string(tr.Result), false))
		}
	}
	messages = append(messages, anthropic.NewAssistantMessage(messageAssistantBlock...))
	messages = append(messages, anthropic.NewUserMessage(messageBlock...))

	return messages
}
