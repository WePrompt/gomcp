package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/WePrompt/gomcp/mcp"
)

// StdioServer wraps an MCPServer and handles stdio communication
type StdioServer struct {
	server    MCPServer
	stopChan  chan os.Signal
	done      chan struct{}
	errLogger *log.Logger
}

// NewStdioServer creates a stdio server wrapper around an existing MCPServer
func NewStdioServer(server MCPServer) *StdioServer {
	s := &StdioServer{
		server:    server,
		stopChan:  make(chan os.Signal, 1),
		errLogger: log.New(os.Stderr, "", log.LstdFlags),
		done:      make(chan struct{}),
	}

	// Listen for shutdown signals
	signal.Notify(s.stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-s.stopChan
		close(s.done)
	}()

	return s
}

func (s *StdioServer) WithLogger(logger *log.Logger) *StdioServer {
	s.errLogger = logger
	return s
}

func (s *StdioServer) Serve() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-s.done
		cancel()
	}()

	reader := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !reader.Scan() {
				if err := reader.Err(); err != nil {
					s.errLogger.Printf("Error reading input: %v", err)
					return err
				}
				// EOF reached
				return nil
			}

			if err := s.handleMessage(ctx, reader.Text()); err != nil {
				if err == io.EOF {
					return nil
				}
				s.errLogger.Printf("Error handling message: %v", err)
			}
		}
	}
}

func (s *StdioServer) handleMessage(ctx context.Context, line string) error {
	var request mcp.JSONRPCRequest
	if err := json.Unmarshal([]byte(line), &request); err != nil {
		s.writeError(nil, mcp.ErrorCodeParseError, "Failed to parse JSON-RPC request")
		return fmt.Errorf("failed to parse JSON-RPC request: %w", err)
	}

	if request.Jsonrpc != mcp.JSONRPCVersion {
		s.writeError(request.Id, mcp.ErrorCodeInvalidRequest, "Invalid JSON-RPC version")
		return fmt.Errorf("invalid JSON-RPC version")
	}

	result, err := s.server.Request(ctx, request.Method, request.Params)
	if err != nil {
		s.writeError(request.Id, mcp.ErrorCodeInternalError, "Internal server error")
		return fmt.Errorf("request handling error: %w", err)
	}

	response := mcp.JSONRPCResponse{
		Jsonrpc: mcp.JSONRPCVersion,
		Id:      request.Id,
		Result:  result,
	}

	if err := s.writeResponse(response); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func (s *StdioServer) writeError(id interface{}, code int, message string) {
	response := mcp.JSONRPCResponse{
		Jsonrpc: mcp.JSONRPCVersion,
		Id:      id,
		Error: &mcp.JSONRPCErrorData{
			Code:    code,
			Message: message,
		},
	}
	if err := s.writeResponse(response); err != nil {
		s.errLogger.Printf("Error writing response: %v", err)
	}
}

func (s *StdioServer) writeResponse(response mcp.JSONRPCResponse) error {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.errLogger.Printf("Error marshal response: %v", err)
		return err
	}

	if _, err := fmt.Fprintf(os.Stdout, "%s\n", responseBytes); err != nil {
		s.errLogger.Printf("Error writing response: %v", err)
		return err
	}

	return nil
}
