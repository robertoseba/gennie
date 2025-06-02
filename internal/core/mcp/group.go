package mcp

import (
	"context"
	"fmt"
	"slices"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/robertoseba/gennie/internal/core/llmcore"
)

type (
	toolNameType string

	toolDetails struct {
		llmcore.Tool
		mcpClient *mcpClient
	}
)

type Group struct {
	tools map[toolNameType]toolDetails
}

func NewGroup() *Group {
	return &Group{
		tools: make(map[toolNameType]toolDetails),
	}
}

func (g *Group) Add(ctx context.Context, cmd string, env []string, args []string, requiresApproval bool, allowedTools []string) error {
	// TODO: Add support for SSE Mcp servers
	mcpClient, err := newStdioClient(ctx, cmd, env, args, requiresApproval)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}

	err = g.retrieveToolsFrom(ctx, mcpClient, allowedTools)
	if err != nil {
		return fmt.Errorf("failed to retrieve tools from MCP server %s: %w", mcpClient.name, err)
	}

	// TODO: curretly we leave all mcp tools enabled.But it might be better to close and open as needed
	// mcpClient.Close()

	return nil
}

func (g *Group) Shutdown() {
	for _, tool := range g.tools {
		tool.mcpClient.close()
	}
}

func (g *Group) ListTools() ([]llmcore.Tool, error) {
	result := make([]llmcore.Tool, 0, len(g.tools))
	for _, tool := range g.tools {
		result = append(result, tool.Tool)
	}

	return result, nil
}

func (g *Group) ExecTool(ctx context.Context, name string, args map[string]any) (string, error) {
	tool, ok := g.tools[toolNameType(name)]
	if !ok {
		return "", fmt.Errorf("tool %s not found", name)
	}

	result, err := tool.mcpClient.execTool(ctx, name, args)
	if err != nil {
		return "", fmt.Errorf("failed to execute tool %s: %w", name, err)
	}

	return result, nil
}

func (g *Group) RequiresApproval(toolName string) bool {
	tool, ok := g.tools[toolNameType(toolName)]
	if !ok {
		return false
	}
	return tool.RequiresApproval
}

func (g *Group) retrieveToolsFrom(ctx context.Context, client *mcpClient, allowedTools []string) error {
	mcpToolsResponse, err := client.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return fmt.Errorf("failed to list tools from MCP server %s: %w", client.name, err)
	}

	for _, t := range mcpToolsResponse.Tools {
		if !slices.Contains(allowedTools, t.Name) && len(allowedTools) > 0 {
			continue
		}

		toolItem := llmcore.Tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: llmcore.ToolInputSchema{
				Type:       t.InputSchema.Type,
				Properties: t.InputSchema.Properties,
				Required:   t.InputSchema.Required,
			},
			RequiresApproval: client.requiresApproval,
		}
		g.tools[toolNameType(t.Name)] = toolDetails{
			Tool:      toolItem,
			mcpClient: client,
		}
	}

	return nil
}
