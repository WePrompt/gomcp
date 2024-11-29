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

func (s *MCPServer) Request(ctx context.Context, method string, params json.RawMessage) (json.RawMessage, error) {
	if strings.HasPrefix(method, "notifications/") {
		var notification mcp.Notification
		if err := json.Unmarshal(params, &notification); err != nil {
			return nil, fmt.Errorf("failed to parse notification: %w", err)
		}
		err := s.notifyHandlers[method].Handle(ctx, notification)
		return nil, err
	}

	switch method {
	case mcp.MethodInitialize:
		var p struct {
			Capabilities    *mcp.ClientCapabilities `json:"capabilities"`
			ClientInfo      *mcp.Implementation     `json:"clientInfo"`
			ProtocolVersion string                  `json:"protocolVersion"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.systemHandler.Initialize(ctx, *p.Capabilities, *p.ClientInfo, p.ProtocolVersion)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	case mcp.MethodPing:
		return nil, s.systemHandler.Ping(ctx)

	case mcp.MethodResourcesList:
		var p struct {
			Cursor *string `json:"cursor,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.resourceHandler.List(ctx, p.Cursor)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	case mcp.MethodResourcesRead:
		var p struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.resourceHandler.Read(ctx, p.URI)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	case mcp.MethodResourcesSubscribe:
		var p struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return nil, s.resourceHandler.Subscribe(ctx, p.URI)

	case mcp.MethodResourcesUnsubscribe:
		var p struct {
			URI string `json:"uri"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return nil, s.resourceHandler.Unsubscribe(ctx, p.URI)

	case mcp.MethodPromptsList:
		var p struct {
			Cursor *string `json:"cursor,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.promptHandler.List(ctx, p.Cursor)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	case mcp.MethodPromptsGet:
		var p struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.promptHandler.Get(ctx, p.Name, p.Arguments)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	case mcp.MethodToolsList:
		var p struct {
			Cursor *string `json:"cursor,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.toolHandler.List(ctx, p.Cursor)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	case mcp.MethodToolsCall:
		var p struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments,omitempty"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.toolHandler.Call(ctx, p.Name, p.Arguments)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	case mcp.MethodLoggingSetLevel:
		var p struct {
			Level mcp.LoggingLevel `json:"level"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		return nil, s.systemHandler.SetLevel(ctx, p.Level)

	case mcp.MethodCompletionComplete:
		var p struct {
			Ref      interface{}         `json:"ref"`
			Argument mcp.CompleteRequest `json:"argument"`
		}
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		result, err := s.systemHandler.Complete(ctx, p.Ref, p.Argument)
		if err != nil {
			return nil, err
		}
		return result.ToJSON()

	default:
		return nil, fmt.Errorf("method not found: %s", method)
	}
}
