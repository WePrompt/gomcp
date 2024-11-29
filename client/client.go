package client

import (
	"context"

	"github.com/WePrompt/gomcp/mcp"
)

// MCPClient defines the interface for communicating with an MCP server
type MCPClient interface {
	// System operations
	Initialize(ctx context.Context, capabilities mcp.ClientCapabilities, clientInfo mcp.Implementation, protocolVersion string) (*mcp.InitializeResult, error)
	Ping(ctx context.Context) error
	SetLoggingLevel(ctx context.Context, level mcp.LoggingLevel) error
	Complete(ctx context.Context, ref interface{}, argument mcp.CompleteRequest) (*mcp.CompleteResult, error)

	// Resource operations
	ListResources(ctx context.Context, cursor *string) (*mcp.ListResourcesResult, error)
	ReadResource(ctx context.Context, uri string) (*mcp.ReadResourceResult, error)
	SubscribeResource(ctx context.Context, uri string) error
	UnsubscribeResource(ctx context.Context, uri string) error

	// Prompt operations
	ListPrompts(ctx context.Context, cursor *string) (*mcp.ListPromptsResult, error)
	GetPrompt(ctx context.Context, name string, arguments map[string]string) (*mcp.GetPromptResult, error)

	// Tool operations
	ListTools(ctx context.Context, cursor *string) (*mcp.ListToolsResult, error)
	CallTool(ctx context.Context, name string, arguments map[string]interface{}) (*mcp.CallToolResult, error)
}
