package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/robertoseba/gennie/internal/core/llm"
)

type Server struct {
	Cmd              string   `toml:"command"`
	Args             []string `toml:"args"`
	Envs             []string `toml:"envs"`
	AllowedTools     []string `toml:"allowed_tools"`
	RequiresApproval bool     `toml:"requires_approval"`
}

type mcpClient struct {
	server Server         // The server this client is connected to.
	client *client.Client // mcp client to communicate with the MCP server.

	// This can be used by clients to improve the LLM's understanding of
	// available tools, resources, etc. It can be thought of like a "hint" to the model.
	// For example, this information MAY be added to the system prompt.
	instructions string // TODO: currently not used, but we might use it in the future
}

func newStdioClient(ctx context.Context, server Server) (*mcpClient, error) {
	c, err := client.NewStdioMCPClient(server.Cmd, server.Envs, server.Args...)
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
		server:       server,
		instructions: result.Instructions,
		client:       c,
	}
	return mcpClient, nil
}

func (c *mcpClient) close() {
	c.client.Close()
}

func (c *mcpClient) requiresApproval() bool {
	return c.server.RequiresApproval
}

// Returns a list of tools available in the MCP server indexed by their name.
func (c *mcpClient) ListTools(ctx context.Context) (map[string]llm.Tool, error) {
	request := mcp.ListToolsRequest{}

	result, err := c.client.ListTools(ctx, request)
	if err != nil {
		return nil, err
	}

	tools := make(map[string]llm.Tool, len(result.Tools))
	for _, t := range result.Tools {
		tools[t.Name] = llm.Tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: llm.ToolInputSchema{
				Type:       t.InputSchema.Type,
				Properties: t.InputSchema.Properties,
				Required:   t.InputSchema.Required,
			},
		}
	}
	return tools, nil
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
