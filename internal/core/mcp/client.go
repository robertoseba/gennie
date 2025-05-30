package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/robertoseba/gennie/internal/core/llmcore"
)

type McpClientInterface interface {
	Close()
	ListTools(ctx context.Context) ([]llmcore.Tool, error)
	ExecTool(ctx context.Context, toolName string, args map[string]any) ([]byte, error)
}

var _ McpClientInterface = &McpClient{}

type McpClient struct {
	client *client.Client
}

func NewStdioClient(ctx context.Context, cmd string, env []string, args []string) (*McpClient, error) {
	c, err := client.NewStdioMCPClient(cmd, env, args...)
	if err != nil {
		return nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "gennie-client",
		Version: "1.0.0",
	}

	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}

	mcpClient := &McpClient{
		client: c,
	}
	return mcpClient, nil
}

func (c *McpClient) Close() {
	c.client.Close()
}

func (c *McpClient) ListTools(ctx context.Context) ([]llmcore.Tool, error) {
	toolsRequest := mcp.ListToolsRequest{}
	mcpToolsResponse, err := c.client.ListTools(ctx, toolsRequest)
	if err != nil {
		return nil, err
	}

	var toolItems []llmcore.Tool
	for _, t := range mcpToolsResponse.Tools {
		toolItem := llmcore.Tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: llmcore.ToolInputSchema{
				Type:       t.InputSchema.Type,
				Properties: t.InputSchema.Properties,
				Required:   t.InputSchema.Required,
			},
		}
		toolItems = append(toolItems, toolItem)
	}
	return toolItems, nil
}

func (c *McpClient) ExecTool(ctx context.Context, toolName string, args map[string]any) ([]byte, error) {
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	request.Params.Name = toolName
	if len(args) > 0 {
		request.Params.Arguments = args
	}

	result, err := c.client.CallTool(ctx, request)
	if err != nil {
		return nil, err
	}
	return parseToolResult(result)
}

func parseToolResult(toolResponse *mcp.CallToolResult) ([]byte, error) {
	if len(toolResponse.Content) == 0 {
		return nil, fmt.Errorf("no content in tool result")
	}

	result := make([]byte, 0)

	for _, content := range toolResponse.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			result = append(result, []byte(textContent.Text)...)
		} else {
			jsonBytes, err := json.Marshal(content)
			if err != nil {
				return nil, fmt.Errorf("error marshalling content: %v", err)
			}
			result = append(result, jsonBytes...)
		}
	}
	return result, nil
}
