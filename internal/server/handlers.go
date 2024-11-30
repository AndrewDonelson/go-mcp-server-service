// Package server provides JSON-RPC 2.0 request handlers for the notes server.
// It implements methods for resource management, prompt handling, and tool execution.
//
// The handlers support the following JSON-RPC 2.0 methods:
//   - list_resources: Lists all available resources
//   - read_resource: Reads content of a specific resource by URI
//   - list_prompts: Lists all available prompts
//   - get_prompt: Retrieves and processes a specific prompt with arguments
//   - list_tools: Lists all available tools
//   - call_tool: Executes a specific tool with provided arguments
//
// Error Handling:
// All handlers follow JSON-RPC 2.0 error specifications with the following error codes:
//   - ErrInvalidReq (-32600): Invalid request format
//   - ErrMethodNotFound (-32601): Requested method not found
//   - ErrInvalidParams (-32602): Invalid method parameters
//   - ErrInternal (-32603): Internal server error
//   - ErrNotFound (404): Resource or item not found
//   - ErrUnsupported (400): Unsupported operation
package server

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

// handleListResources processes the list_resources RPC method.
// It returns a list of all available resources in the server.
//
// The response contains:
//   - JSONRPC: Version string (always "2.0")
//   - ID: Request ID from the original request
//   - Result: Array of available resources
func (s *Server) handleListResources(req *RPCRequest) *RPCResponse {
    fmt.Fprintf(os.Stderr, "Handling list_resources request\n")
    resources := s.ListResources()
    return &RPCResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  resources,
    }
}

// handleReadResource processes the read_resource RPC method.
// It retrieves the content of a specific resource identified by its URI.
//
// Parameters:
//   - uri: String identifying the resource to read
//
// Returns a response with the resource content or an error if:
//   - URI parameter is missing or invalid
//   - Resource is not found
//   - URI scheme is unsupported
//   - Internal error occurs during reading
func (s *Server) handleReadResource(req *RPCRequest) *RPCResponse {
    if req.Params == nil {
        return newErrorResponse(req.ID, ErrInvalidParams, "params required", nil)
    }

    var params struct {
        URI string `json:"uri"` // Resource URI to read
    }
    if err := json.Unmarshal(req.Params, &params); err != nil {
        fmt.Fprintf(os.Stderr, "Error unmarshaling read_resource params: %v\n", err)
        return newErrorResponse(req.ID, ErrInvalidParams, "invalid URI parameter", err)
    }

    if params.URI == "" {
        return newErrorResponse(req.ID, ErrInvalidParams, "URI is required", nil)
    }

    fmt.Fprintf(os.Stderr, "Reading resource: %s\n", params.URI)
    content, err := s.ReadResource(params.URI)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading resource: %v\n", err)
        switch {
        case strings.Contains(err.Error(), "note not found"):
            return newErrorResponse(req.ID, ErrNotFound, "note not found", err)
        case strings.Contains(err.Error(), "unsupported URI scheme"):
            return newErrorResponse(req.ID, ErrUnsupported, "unsupported URI scheme", err)
        default:
            return newErrorResponse(req.ID, ErrInternal, "internal error", err)
        }
    }

    return &RPCResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  content,
    }
}

// handleListPrompts processes the list_prompts RPC method.
// It returns a list of all available prompt templates.
//
// The response contains:
//   - JSONRPC: Version string (always "2.0")
//   - ID: Request ID from the original request
//   - Result: Array of available prompts
func (s *Server) handleListPrompts(req *RPCRequest) *RPCResponse {
    fmt.Fprintf(os.Stderr, "Handling list_prompts request\n")
    prompts := s.ListPrompts()
    return &RPCResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  prompts,
    }
}

// handleGetPrompt processes the get_prompt RPC method.
// It retrieves and processes a specific prompt template with provided arguments.
//
// Parameters:
//   - name: String identifying the prompt template
//   - arguments: Optional map of key-value pairs for template processing
//
// Returns a response with the processed prompt or an error if:
//   - Name parameter is missing or invalid
//   - Prompt template is not found
//   - Internal error occurs during processing
func (s *Server) handleGetPrompt(req *RPCRequest) *RPCResponse {
    if req.Params == nil {
        return newErrorResponse(req.ID, ErrInvalidParams, "params required", nil)
    }

    var params struct {
        Name      string            `json:"name"`      // Name of the prompt template
        Arguments map[string]string `json:"arguments"` // Template arguments
    }
    if err := json.Unmarshal(req.Params, &params); err != nil {
        fmt.Fprintf(os.Stderr, "Error unmarshaling get_prompt params: %v\n", err)
        return newErrorResponse(req.ID, ErrInvalidParams, "invalid prompt parameters", err)
    }

    if params.Name == "" {
        return newErrorResponse(req.ID, ErrInvalidParams, "prompt name is required", nil)
    }

    if params.Arguments == nil {
        params.Arguments = make(map[string]string)
    }

    fmt.Fprintf(os.Stderr, "Getting prompt: %s with %d arguments\n", params.Name, len(params.Arguments))
    result, err := s.GetPrompt(params.Name, params.Arguments)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting prompt: %v\n", err)
        if strings.Contains(err.Error(), "unknown prompt") {
            return newErrorResponse(req.ID, ErrNotFound, "prompt not found", err)
        }
        return newErrorResponse(req.ID, ErrInternal, "internal error", err)
    }

    return &RPCResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  result,
    }
}

// handleListTools processes the list_tools RPC method.
// It returns a list of all available tools.
//
// The response contains:
//   - JSONRPC: Version string (always "2.0")
//   - ID: Request ID from the original request
//   - Result: Array of available tools
func (s *Server) handleListTools(req *RPCRequest) *RPCResponse {
    fmt.Fprintf(os.Stderr, "Handling list_tools request\n")
    tools := s.ListTools()
    return &RPCResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  tools,
    }
}

// handleCallTool processes the call_tool RPC method.
// It executes a specific tool with provided arguments.
//
// Parameters:
//   - name: String identifying the tool to execute
//   - arguments: Optional map of key-value pairs for tool execution
//
// Returns a response with the tool execution result or an error if:
//   - Name parameter is missing or invalid
//   - Tool is not found
//   - Invalid arguments are provided
//   - Internal error occurs during execution
func (s *Server) handleCallTool(req *RPCRequest) *RPCResponse {
    if req.Params == nil {
        return newErrorResponse(req.ID, ErrInvalidParams, "params required", nil)
    }

    var params struct {
        Name      string                 `json:"name"`      // Name of the tool to execute
        Arguments map[string]interface{} `json:"arguments"` // Tool arguments
    }
    if err := json.Unmarshal(req.Params, &params); err != nil {
        fmt.Fprintf(os.Stderr, "Error unmarshaling call_tool params: %v\n", err)
        return newErrorResponse(req.ID, ErrInvalidParams, "invalid tool parameters", err)
    }

    if params.Name == "" {
        return newErrorResponse(req.ID, ErrInvalidParams, "tool name is required", nil)
    }

    if params.Arguments == nil {
        params.Arguments = make(map[string]interface{})
    }

    fmt.Fprintf(os.Stderr, "Calling tool: %s with %d arguments\n", params.Name, len(params.Arguments))
    result, err := s.CallTool(params.Name, params.Arguments)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error calling tool: %v\n", err)
        if strings.Contains(err.Error(), "unknown tool") {
            return newErrorResponse(req.ID, ErrNotFound, "tool not found", err)
        }
        return newErrorResponse(req.ID, ErrInvalidParams, "invalid tool arguments", err)
    }

    return &RPCResponse{
        JSONRPC: "2.0",
        ID:      req.ID,
        Result:  result,
    }
}

// handleRequest is the main entry point for processing RPC requests.
// It routes requests to appropriate handlers based on the method name.
//
// Supported methods:
//   - list_resources: List available resources
//   - read_resource: Read a specific resource
//   - list_prompts: List available prompts
//   - get_prompt: Get and process a specific prompt
//   - list_tools: List available tools
//   - call_tool: Execute a specific tool
//
// Returns an error response if:
//   - Method is missing or invalid
//   - Required parameters are missing
//   - Method is not found
func (s *Server) handleRequest(req *RPCRequest) *RPCResponse {
    if req.Method == "" {
        return newErrorResponse(req.ID, ErrInvalidReq, "method is required", nil)
    }

    fmt.Fprintf(os.Stderr, "Handling request for method: %s\n", req.Method)

    switch req.Method {
    case "list_resources":
        return s.handleListResources(req)
    case "read_resource":
        if req.Params == nil {
            return newErrorResponse(req.ID, ErrInvalidParams, "params required", nil)
        }
        return s.handleReadResource(req)
    case "list_prompts":
        return s.handleListPrompts(req)
    case "get_prompt":
        if req.Params == nil {
            return newErrorResponse(req.ID, ErrInvalidParams, "params required", nil)
        }
        return s.handleGetPrompt(req)
    case "list_tools":
        return s.handleListTools(req)
    case "call_tool":
        if req.Params == nil {
            return newErrorResponse(req.ID, ErrInvalidParams, "params required", nil)
        }
        return s.handleCallTool(req)
    default:
        return newErrorResponse(req.ID, ErrMethodNotFound, "method not found", fmt.Errorf("unknown method: %s", req.Method))
    }
}

// newErrorResponse creates a new JSON-RPC 2.0 error response.
//
// Parameters:
//   - id: Request ID from the original request
//   - code: Error code as defined in JSON-RPC 2.0 spec
//   - message: Human-readable error message
//   - err: Optional underlying error
//
// Returns a properly formatted RPCResponse with error details.
// If err is provided, its message is included in the error data field.
func newErrorResponse(id interface{}, code int, message string, err error) *RPCResponse {
    data := message
    if err != nil {
        data = err.Error()
    }
    fmt.Fprintf(os.Stderr, "Creating error response: [%d] %s - %v\n", code, message, err)
    return &RPCResponse{
        JSONRPC: "2.0",
        ID:      id,
        Error: &RPCError{
            Code:    code,
            Message: message,
            Data:    data,
        },
    }
}