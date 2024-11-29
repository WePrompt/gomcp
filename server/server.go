package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/WePrompt/gomcp/mcp"
	"github.com/WePrompt/gomcp/server/handlers"
)

type MCPServer struct {
	resourceHandler handlers.ResourceHandler
	promptHandler   handlers.PromptHandler
	toolHandler     handlers.ToolHandler
	systemHandler   handlers.SystemHandler
	notifyHandlers  map[string]handlers.NotificationHandler
	serverInfo      ServerInfo
}

type ServerInfo struct {
	name    string
	version string
}

type ServerOption func(*MCPServer)

func NewMCPServer(opts ...ServerOption) *MCPServer {
	s := &MCPServer{
		notifyHandlers: make(map[string]handlers.NotificationHandler),
		serverInfo: ServerInfo{
			name:    "default",
			version: "1.0.0",
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	// Set default handlers if not provided
	if s.resourceHandler == nil {
		s.resourceHandler = handlers.NewDefaultResourceHandler()
	}
	if s.promptHandler == nil {
		s.promptHandler = handlers.NewDefaultPromptHandler()
	}
	if s.toolHandler == nil {
		s.toolHandler = handlers.NewDefaultToolHandler()
	}
	if s.systemHandler == nil {
		s.systemHandler = handlers.NewDefaultSystemHandler()
	}

	return s
}

func WithServerInfo(name, version string) ServerOption {
	return func(s *MCPServer) {
		s.serverInfo = ServerInfo{name: name, version: version}
	}
}

func WithResourceHandler(h handlers.ResourceHandler) ServerOption {
	return func(s *MCPServer) {
		s.resourceHandler = h
	}
}

func WithPromptHandler(h handlers.PromptHandler) ServerOption {
	return func(s *MCPServer) {
		s.promptHandler = h
	}
}

func WithToolHandler(h handlers.ToolHandler) ServerOption {
	return func(s *MCPServer) {
		s.toolHandler = h
	}
}

func WithSystemHandler(h handlers.SystemHandler) ServerOption {
	return func(s *MCPServer) {
		s.systemHandler = h
	}
}

func WithNotificationHandler(method string, h handlers.NotificationHandler) ServerOption {
	return func(s *MCPServer) {
		s.notifyHandlers[method] = h
	}
}

func (s *MCPServer) Request(ctx context.Context, method string, params json.RawMessage) (interface{}, error) {
	if strings.HasPrefix(method, "notifications/") {
		var notification mcp.Notification
		if err := json.Unmarshal(params, &notification); err != nil {
			return nil, fmt.Errorf("failed to parse notification: %w", err)
		}
		err := s.notifyHandlers[method].Handle(ctx, notification)
		return struct{}{}, err
	}

	switch method {
	case "initialize":
		var p struct {
			Capabilities    *mcp.ClientCapabilities `json:"capabilities"`
			ClientInfo      *mcp.Implementation     `json:"clientInfo"`
			ProtocolVersion string                  `json:"protocolVersion"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.systemHandler.Initialize(ctx, *p.Capabilities, *p.ClientInfo, p.ProtocolVersion)

	case "ping":
		return struct{}{}, s.systemHandler.Ping(ctx)

	case "resources/list":
		var p struct {
			Cursor *string `json:"cursor,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.resourceHandler.List(ctx, p.Cursor)

	case "resources/read":
		var p struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.resourceHandler.Read(ctx, p.URI)

	case "resources/subscribe":
		var p struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return struct{}{}, s.resourceHandler.Subscribe(ctx, p.URI)

	case "resources/unsubscribe":
		var p struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return struct{}{}, s.resourceHandler.Unsubscribe(ctx, p.URI)

	case "prompts/list":
		var p struct {
			Cursor *string `json:"cursor,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.promptHandler.List(ctx, p.Cursor)

	case "prompts/get":
		var p struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.promptHandler.Get(ctx, p.Name, p.Arguments)

	case "tools/list":
		var p struct {
			Cursor *string `json:"cursor,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.toolHandler.List(ctx, p.Cursor)

	case "tools/call":
		var p struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.toolHandler.Call(ctx, p.Name, p.Arguments)

	case "logging/setLevel":
		var p struct {
			Level mcp.LoggingLevel `json:"level"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return struct{}{}, s.systemHandler.SetLevel(ctx, p.Level)

	case "completion/complete":
		var p struct {
			Ref      interface{}         `json:"ref"`
			Argument mcp.CompleteRequest `json:"argument"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return s.systemHandler.Complete(ctx, p.Ref, p.Argument)

	default:
		return nil, fmt.Errorf("method not found: %s", method)
	}
}
