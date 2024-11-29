package handlers

import (
	"context"

	"github.com/WePrompt/gomcp/mcp"
)

type ResourceHandler interface {
	List(ctx context.Context, cursor *string) (*mcp.ListResourcesResult, error)
	Read(ctx context.Context, uri string) (*mcp.ReadResourceResult, error)
	Subscribe(ctx context.Context, uri string) error
	Unsubscribe(ctx context.Context, uri string) error
}

type PromptHandler interface {
	List(ctx context.Context, cursor *string) (*mcp.ListPromptsResult, error)
	Get(ctx context.Context, name string, arguments map[string]string) (*mcp.GetPromptResult, error)
}

type ToolHandler interface {
	List(ctx context.Context, cursor *string) (*mcp.ListToolsResult, error)
	Call(ctx context.Context, name string, arguments map[string]interface{}) (*mcp.CallToolResult, error)
}

type SystemHandler interface {
	Initialize(ctx context.Context, capabilities mcp.ClientCapabilities, clientInfo mcp.Implementation, protocolVersion string) (*mcp.InitializeResult, error)
	Ping(ctx context.Context) error
	SetLevel(ctx context.Context, level mcp.LoggingLevel) error
	Complete(ctx context.Context, ref interface{}, argument mcp.CompleteRequest) (*mcp.CompleteResult, error)
}

type NotificationHandler interface {
	Handle(ctx context.Context, notification mcp.Notification) error
}
