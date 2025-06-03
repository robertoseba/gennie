package mcp

import (
	"context"
	"fmt"
	"slices"

	"github.com/robertoseba/gennie/internal/core/llm"
)

type (
	Group struct {
		tools map[string]toolDetails // tools indexed by name
	}

	toolDetails struct {
		llm.Tool
		mcpClient *mcpClient
	}
)

func NewGroup() *Group {
	return &Group{
		tools: make(map[string]toolDetails),
	}
}

func (g *Group) Add(ctx context.Context, server Server) error {
	// TODO: Add support for SSE Mcp servers
	mcpClient, err := newStdioClient(ctx, server)
	if err != nil {
		return fmt.Errorf("failed to create MCP client: %w", err)
	}

	err = g.retrieveToolsFrom(ctx, mcpClient, server.AllowedTools)
	if err != nil {
		return fmt.Errorf("failed to retrieve tools from MCP server %s: %w", mcpClient.server.Cmd, err)
	}

	// TODO: curretly we leave all mcp tools enabled.Not sure if it might be better to close and open as needed
	// because we have to load all of the them to get the tool list.
	// mcpClient.Close()

	return nil
}

func (g *Group) Shutdown() {
	for _, tool := range g.tools {
		tool.mcpClient.close()
	}
}

func (g *Group) ListTools() []llm.Tool {
	result := make([]llm.Tool, 0, len(g.tools))
	for _, tool := range g.tools {
		result = append(result, tool.Tool)
	}

	return result
}

func (g *Group) ExecTool(ctx context.Context, name string, args map[string]any) (string, error) {
	tool, ok := g.tools[name]
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
	tool, ok := g.tools[toolName]
	if !ok {
		return false
	}
	return tool.mcpClient.requiresApproval()
}

func (g *Group) retrieveToolsFrom(ctx context.Context, client *mcpClient, allowedTools []string) error {
	mcpToolsResponse, err := client.ListTools(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tools from MCP server %s: %w", client.server.Cmd, err)
	}

	for name, tool := range mcpToolsResponse {
		if !slices.Contains(allowedTools, name) && len(allowedTools) > 0 {
			continue
		}

		g.tools[name] = toolDetails{
			Tool:      tool,
			mcpClient: client,
		}
	}

	return nil
}
