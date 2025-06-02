package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type mcpClient struct {
	name             string
	args             []string
	env              []string
	requiresApproval bool // If true, the tool requires approval before execution
	client           *client.Client

	// This can be used by clients to improve the LLM's understanding of
	// available tools, resources, etc. It can be thought of like a "hint" to the model.
	// For example, this information MAY be added to the system prompt.
	instructions string // TODO: currently not used, but we might use it in the future
}

func newStdioClient(ctx context.Context, cmd string, env []string, args []string, requiresApproval bool) (*mcpClient, error) {
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

	result, err := c.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}

	mcpClient := &mcpClient{
		name:             cmd,
		args:             args,
		env:              env,
		instructions:     result.Instructions,
		requiresApproval: requiresApproval,
		client:           c,
	}
	return mcpClient, nil
}

func (c *mcpClient) close() {
	c.client.Close()
}

func (c *mcpClient) execTool(ctx context.Context, toolName string, args map[string]any) (string, error) {
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
		return "", err
	}
	return parseToolResult(result)
}

func parseToolResult(toolResponse *mcp.CallToolResult) (string, error) {
	if len(toolResponse.Content) == 0 {
		return "", fmt.Errorf("no content in tool result")
	}

	result := strings.Builder{}

	for _, content := range toolResponse.Content {
		textContent, ok := content.(mcp.TextContent)
		if !ok {
			return "", fmt.Errorf("server return a content type not supported. We currently only support text responses from mcp servers")
		}
		result.WriteString(textContent.Text)
	}
	return result.String(), nil
}
