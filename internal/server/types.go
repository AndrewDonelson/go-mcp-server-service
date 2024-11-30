// Package server provides types and structures for implementing a JSON-RPC 2.0
// compliant server with support for note management, prompts, and tools.
package server

import (
    "encoding/json"
    "sync"
    "fmt"
)

// JSON-RPC 2.0 error codes as defined by the specification.
// Custom error codes should be in the range -32000 to -32099.
const (
    // ErrParse indicates the server received invalid JSON.
    // Code -32700 is reserved for parse errors by the JSON-RPC 2.0 spec.
    ErrParse = -32700

    // ErrInvalidReq indicates the request object is not a valid JSON-RPC 2.0 request.
    // Code -32600 is reserved for invalid requests by the JSON-RPC 2.0 spec.
    ErrInvalidReq = -32600

    // ErrMethodNotFound indicates the requested method does not exist.
    // Code -32601 is reserved for method not found errors by the JSON-RPC 2.0 spec.
    ErrMethodNotFound = -32601

    // ErrInvalidParams indicates invalid method parameters.
    // Code -32602 is reserved for invalid parameters by the JSON-RPC 2.0 spec.
    ErrInvalidParams = -32602

    // ErrInternal indicates an internal JSON-RPC error.
    // Code -32603 is reserved for internal errors by the JSON-RPC 2.0 spec.
    ErrInternal = -32603

    // ErrNotFound is a custom error code indicating a resource was not found.
    // Custom code -32001.
    ErrNotFound = -32001

    // ErrUnsupported is a custom error code indicating an unsupported operation.
    // Custom code -32002.
    ErrUnsupported = -32002
)

// Server represents the main server instance that handles note management and RPC requests.
// It maintains thread-safe access to the notes storage through sync.RWMutex.
type Server struct {
    name     string              // Server instance identifier
    notes    map[string]string   // Storage for note content
    notesMap sync.RWMutex       // Mutex for thread-safe access to notes
}

// Resource represents a note resource in the system with its metadata.
// It provides information about the resource's location, name, and content type.
type Resource struct {
    URI         string `json:"uri"`          // Unique identifier for the resource
    Name        string `json:"name"`         // Display name of the resource
    Description string `json:"description"`   // Human-readable description
    MimeType    string `json:"mimeType"`     // MIME type of the resource content
}

// Prompt represents a command prompt that can be executed by the server.
// It includes metadata about the prompt and its required arguments.
type Prompt struct {
    Name        string           `json:"name"`        // Unique identifier for the prompt
    Description string           `json:"description"` // Human-readable description
    Arguments   []PromptArgument `json:"arguments,omitempty"` // Optional list of arguments
}

// PromptArgument defines an argument that can be passed to a Prompt.
// It includes metadata about the argument and whether it's required.
type PromptArgument struct {
    Name        string `json:"name"`        // Name of the argument
    Description string `json:"description"` // Human-readable description
    Required    bool   `json:"required"`    // Whether this argument must be provided
}

// Tool represents an executable tool in the system.
// It includes metadata and a JSON schema defining its input parameters.
type Tool struct {
    Name         string          `json:"name"`        // Unique identifier for the tool
    Description  string          `json:"description"` // Human-readable description
    InputSchema  json.RawMessage `json:"inputSchema"` // JSON Schema of valid inputs
}

// TextContent represents a text-based content item with its type.
// Used for representing various text-based data in the system.
type TextContent struct {
    Type string `json:"type"` // Content type identifier
    Text string `json:"text"` // The actual text content
}

// GetPromptResult represents the result of retrieving a prompt.
// It includes a description and a list of messages associated with the prompt.
type GetPromptResult struct {
    Description string          `json:"description"` // Human-readable description
    Messages    []PromptMessage `json:"messages"`    // List of prompt messages
}

// PromptMessage represents a single message in a prompt sequence.
// It includes the role of the message sender and the content of the message.
type PromptMessage struct {
    Role    string      `json:"role"`    // Role of the message sender
    Content TextContent `json:"content"` // Content of the message
}

// RPCRequest represents a JSON-RPC 2.0 request.
// It follows the JSON-RPC 2.0 specification for request structure.
type RPCRequest struct {
    JSONRPC string          `json:"jsonrpc"` // Must be "2.0"
    ID      interface{}     `json:"id"`      // Request identifier
    Method  string         `json:"method"`   // Name of the method to be invoked
    Params  json.RawMessage `json:"params"`  // Parameters for the method
}

// validate checks if the RPCRequest is valid according to the JSON-RPC 2.0 specification.
// Currently only checks for method presence, but can be extended for additional validation.
//
// Returns:
//   - error: nil if the request is valid, otherwise an error describing the validation failure
func (r *RPCRequest) validate() error {
    if r.Method == "" {
        return fmt.Errorf("method is required")
    }
    return nil
}

// RPCResponse represents a JSON-RPC 2.0 response.
// It follows the JSON-RPC 2.0 specification for response structure.
type RPCResponse struct {
    JSONRPC string          `json:"jsonrpc"` // Must be "2.0"
    ID      interface{}     `json:"id"`      // Same as the request ID
    Result  interface{}     `json:"result,omitempty"` // Method return value
    Error   *RPCError       `json:"error,omitempty"`  // Error object if an error occurred
}

// RPCError represents a JSON-RPC 2.0 error object.
// It includes an error code, message, and optional additional data.
type RPCError struct {
    Code    int         `json:"code"`    // Error code (see constants)
    Message string      `json:"message"` // Human-readable error message
    Data    interface{} `json:"data,omitempty"` // Additional error information
}