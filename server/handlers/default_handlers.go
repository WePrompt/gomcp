package handlers

import (
	"context"

	"github.com/WePrompt/gomcp/mcp"
)

// Default implementation of ResourceHandler
type DefaultResourceHandler struct{}

func NewDefaultResourceHandler() ResourceHandler {
	return &DefaultResourceHandler{}
}

func (h *DefaultResourceHandler) List(ctx context.Context, cursor *string) (*mcp.ListResourcesResult, error) {
	return &mcp.ListResourcesResult{
		Resources: []mcp.Resource{},
	}, nil
}

func (h *DefaultResourceHandler) Read(ctx context.Context, uri string) (*mcp.ReadResourceResult, error) {
	result := &mcp.ReadResourceResult{}
	result.AddTextContent(mcp.TextResourceContents{})
	return result, nil
}

func (h *DefaultResourceHandler) Subscribe(ctx context.Context, uri string) error {
	return nil
}

func (h *DefaultResourceHandler) Unsubscribe(ctx context.Context, uri string) error {
	return nil
}

// Default implementation of PromptHandler
type DefaultPromptHandler struct{}

func NewDefaultPromptHandler() PromptHandler {
	return &DefaultPromptHandler{}
}

func (h *DefaultPromptHandler) List(ctx context.Context, cursor *string) (*mcp.ListPromptsResult, error) {
	return &mcp.ListPromptsResult{
		Prompts: []mcp.Prompt{},
	}, nil
}

func (h *DefaultPromptHandler) Get(ctx context.Context, name string, arguments map[string]string) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Messages: []mcp.PromptMessage{},
	}, nil
}

// Default implementation of ToolHandler
type DefaultToolHandler struct{}

func NewDefaultToolHandler() ToolHandler {
	return &DefaultToolHandler{}
}

func (h *DefaultToolHandler) List(ctx context.Context, cursor *string) (*mcp.ListToolsResult, error) {
	return &mcp.ListToolsResult{
		Tools: []mcp.Tool{},
	}, nil
}

func (h *DefaultToolHandler) Call(ctx context.Context, name string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	result := &mcp.CallToolResult{}
	result.AddTextContent(mcp.TextContent{})
	return result, nil
}

// Default implementation of SystemHandler
type DefaultSystemHandler struct{}

func NewDefaultSystemHandler() SystemHandler {
	return &DefaultSystemHandler{}
}

func (h *DefaultSystemHandler) Initialize(ctx context.Context, capabilities mcp.ClientCapabilities, clientInfo mcp.Implementation, protocolVersion string) (*mcp.InitializeResult, error) {
	return &mcp.InitializeResult{
		ServerInfo: mcp.Implementation{
			Name:    clientInfo.Name,
			Version: clientInfo.Version,
		},
		ProtocolVersion: protocolVersion,
		Capabilities: mcp.ServerCapabilities{
			Resources: &mcp.ServerCapabilitiesResources{
				ListChanged: true,
				Subscribe:   true,
			},
		},
	}, nil
}

func (h *DefaultSystemHandler) Ping(ctx context.Context) error {
	return nil
}

func (h *DefaultSystemHandler) SetLevel(ctx context.Context, level mcp.LoggingLevel) error {
	return nil
}

func (h *DefaultSystemHandler) Complete(ctx context.Context, ref interface{}, argument mcp.CompleteRequest) (*mcp.CompleteResult, error) {
	return &mcp.CompleteResult{}, nil
}

// Default implementation of NotificationHandler
type DefaultNotificationHandler struct{}

func NewDefaultNotificationHandler() NotificationHandler {
	return &DefaultNotificationHandler{}
}

func (h *DefaultNotificationHandler) Handle(ctx context.Context, notification mcp.Notification) error {
	return nil
}
