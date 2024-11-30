// Package server implements a JSON-RPC 2.0 server that manages notes over standard I/O.
// It provides a simple interface for creating and managing note-taking functionality
// through a stateful server instance that communicates using stdin/stdout.
package server

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "sync"
)

// NewServer creates and initializes a new Server instance with the specified name.
// It initializes an empty notes storage map and sets up the basic server configuration.
//
// Parameters:
//   - name: A string identifier for the server instance
//
// Returns:
//   - *Server: A pointer to the newly created Server instance
//
// Example:
//
//	server := NewServer("my-notes-server")
func NewServer(name string) *Server {
    return &Server{
        name:  name,
        notes: make(map[string]string),
    }
}

// Run starts the server and begins processing JSON-RPC 2.0 requests over stdin/stdout.
// It continues running until either the context is cancelled or EOF is received on stdin.
//
// The server handles JSON-RPC 2.0 protocol requirements including:
//   - Version validation ("2.0" only)
//   - Method presence verification
//   - Request parsing and error handling
//   - Response encoding
//
// Parameters:
//   - ctx: A context.Context for controlling server lifecycle
//
// Returns:
//   - error: An error if the server encounters a fatal condition, including:
//     * Context cancellation
//     * IO errors
//     * JSON encoding/decoding errors
//     * Protocol errors
//
// Error Handling:
//   - Returns nil on clean shutdown (EOF)
//   - Returns context.Canceled or context.DeadlineExceeded when context is done
//   - Returns encoding/decoding errors with appropriate JSON-RPC error responses
//
// Protocol Errors:
//   - ErrParse (-32700): Invalid JSON was received
//   - ErrInvalidReq (-32600): Invalid JSON-RPC request (version mismatch)
//
// Example:
//
//	ctx := context.Background()
//	if err := server.Run(ctx); err != nil {
//	    log.Fatal(err)
//	}
func (s *Server) Run(ctx context.Context) error {
    // Use stderr for logging
    fmt.Fprintf(os.Stderr, "Notes Server starting on stdio...\n")
    
    // Create encoders/decoders for stdin/stdout
    decoder := json.NewDecoder(os.Stdin)
    
    // Create a mutex for stdout to ensure thread-safe writing
    var stdoutMutex sync.Mutex
    encoder := json.NewEncoder(os.Stdout)

    for {
        select {
        case <-ctx.Done():
            fmt.Fprintf(os.Stderr, "Server shutting down: %v\n", ctx.Err())
            return ctx.Err()
            
        default:
            var req RPCRequest
            if err := decoder.Decode(&req); err != nil {
                if err == io.EOF {
                    fmt.Fprintf(os.Stderr, "Server stopped: EOF received\n")
                    return nil
                }
                fmt.Fprintf(os.Stderr, "Error decoding request: %v\n", err)
                
                // Lock stdout while writing error response
                stdoutMutex.Lock()
                encodeErr := encoder.Encode(&RPCResponse{
                    JSONRPC: "2.0",
                    Error: &RPCError{
                        Code:    ErrParse,
                        Message: "parse error",
                        Data:    err.Error(),
                    },
                })
                stdoutMutex.Unlock()
                
                if encodeErr != nil {
                    return fmt.Errorf("failed to encode error response: %w", encodeErr)
                }
                return fmt.Errorf("failed to decode request: %w", err)
            }

            if req.JSONRPC != "2.0" {
                stdoutMutex.Lock()
                encodeErr := encoder.Encode(&RPCResponse{
                    JSONRPC: "2.0",
                    ID:      req.ID,
                    Error: &RPCError{
                        Code:    ErrInvalidReq,
                        Message: "invalid JSON-RPC version",
                        Data:    "expected version 2.0",
                    },
                })
                stdoutMutex.Unlock()
                
                if encodeErr != nil {
                    return fmt.Errorf("failed to encode response: %w", encodeErr)
                }
                continue
            }

            if req.Method == "" {
                stdoutMutex.Lock()
                encodeErr := encoder.Encode(&RPCResponse{
                    JSONRPC: "2.0",
                    ID:      req.ID,
                    Error: &RPCError{
                        Code:    ErrInvalidReq,
                        Message: "method is required",
                        Data:    "empty method",
                    },
                })
                stdoutMutex.Unlock()
                
                if encodeErr != nil {
                    return fmt.Errorf("failed to encode response: %w", encodeErr)
                }
                continue
            }

            // Handle the request and get response
            response := s.handleRequest(&req)
            
            // Lock stdout while writing response
            stdoutMutex.Lock()
            err := encoder.Encode(response)
            stdoutMutex.Unlock()
            
            if err != nil {
                return fmt.Errorf("failed to encode response: %w", err)
            }
        }
    }
}