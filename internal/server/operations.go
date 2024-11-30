// Package server provides functionality for managing and serving notes through
// a server implementation. It includes operations for resource management,
// prompt handling, and tool execution.
package server

import (
    "encoding/json"
    "fmt"
    "net/url"
    "os"
)

// ListResources returns a slice of all available resources in the server.
// Each resource represents a note with its URI, name, description, and MIME type.
// The resources are returned in an unspecified order.
//
// The URI format follows the scheme: note://internal/{name}
// where {name} is the unique identifier of the note.
//
// The function acquires a read lock on the notes map to ensure thread safety.
func (s *Server) ListResources() []Resource {
    s.notesMap.RLock()
    defer s.notesMap.RUnlock()

    fmt.Fprintf(os.Stderr, "Listing %d resources\n", len(s.notes))
    resources := make([]Resource, 0, len(s.notes))
    for name := range s.notes {
        resources = append(resources, Resource{
            URI:         fmt.Sprintf("note://internal/%s", name),
            Name:        fmt.Sprintf("Note: %s", name),
            Description: fmt.Sprintf("A simple note named %s", name),
            MimeType:    "text/plain",
        })
    }
    return resources
}

// ReadResource retrieves the content of a resource identified by the given URI.
// The URI must follow the format: note://{path} where path is the note identifier.
//
// Parameters:
//   - uri: The URI of the resource to read
//
// Returns:
//   - string: The content of the resource
//   - error: An error if the URI is invalid, the scheme is unsupported,
//     or the resource is not found
//
// Examples:
//
//	content, err := server.ReadResource("note://internal/example-note")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *Server) ReadResource(uri string) (string, error) {
    parsedURI, err := url.Parse(uri)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to parse URI %s: %v\n", uri, err)
        return "", fmt.Errorf("invalid URI: %w", err)
    }

    if parsedURI.Scheme != "note" {
        fmt.Fprintf(os.Stderr, "Unsupported URI scheme: %s\n", parsedURI.Scheme)
        return "", fmt.Errorf("unsupported URI scheme: %s", parsedURI.Scheme)
    }

    name := parsedURI.Path
    if name != "" {
        name = name[1:]
    }

    fmt.Fprintf(os.Stderr, "Reading resource: %s\n", name)

    s.notesMap.RLock()
    content, ok := s.notes[name]
    s.notesMap.RUnlock()

    if !ok {
        fmt.Fprintf(os.Stderr, "Note not found: %s\n", name)
        return "", fmt.Errorf("note not found: %s", name)
    }

    return content, nil
}

// ListPrompts returns a slice of all available prompts in the server.
// Currently, it only supports the "summarize-notes" prompt, which creates
// a summary of all notes with optional style configuration.
func (s *Server) ListPrompts() []Prompt {
    fmt.Fprintf(os.Stderr, "Listing available prompts\n")
    return []Prompt{{
        Name:        "summarize-notes",
        Description: "Creates a summary of all notes",
        Arguments: []PromptArgument{{
            Name:        "style",
            Description: "Style of the summary (brief/detailed)",
            Required:    false,
        }},
    }}
}

// GetPrompt retrieves the prompt configuration and generates the appropriate
// messages for the specified prompt name and arguments.
//
// Parameters:
//   - name: The name of the prompt to retrieve
//   - arguments: A map of argument names to their values
//
// Returns:
//   - GetPromptResult: The result containing the prompt description and messages
//   - error: An error if the prompt name is unknown
//
// Currently supported prompts:
//   - "summarize-notes": Generates a summary of all notes
//     Arguments:
//   - "style": Optional. Values: "brief" (default) or "detailed"
func (s *Server) GetPrompt(name string, arguments map[string]string) (GetPromptResult, error) {
    fmt.Fprintf(os.Stderr, "Getting prompt %s with arguments: %v\n", name, arguments)
    
    if name != "summarize-notes" {
        return GetPromptResult{}, fmt.Errorf("unknown prompt: %s", name)
    }

    style := arguments["style"]
    if style == "" {
        style = "brief"
    }

    detailPrompt := ""
    if style == "detailed" {
        detailPrompt = " Give extensive details."
    }

    s.notesMap.RLock()
    var notesList string
    for name, content := range s.notes {
        notesList += fmt.Sprintf("- %s: %s\n", name, content)
    }
    s.notesMap.RUnlock()

    fmt.Fprintf(os.Stderr, "Generated prompt with style: %s\n", style)

    return GetPromptResult{
        Description: "Summarize the current notes",
        Messages: []PromptMessage{{
            Role: "user",
            Content: TextContent{
                Type: "text",
                Text: fmt.Sprintf("Here are the current notes to summarize:%s\n\n%s", detailPrompt, notesList),
            },
        }},
    }, nil
}

// ListTools returns a slice of all available tools in the server.
// Currently, it only supports the "add-note" tool, which allows adding
// new notes to the server.
func (s *Server) ListTools() []Tool {
    fmt.Fprintf(os.Stderr, "Listing available tools\n")
    return []Tool{{
        Name:        "add-note",
        Description: "Add a new note",
        InputSchema: json.RawMessage(`{
            "type": "object",
            "properties": {
                "name": {"type": "string"},
                "content": {"type": "string"}
            },
            "required": ["name", "content"]
        }`),
    }}
}

// CallTool executes the specified tool with the given arguments.
//
// Parameters:
//   - name: The name of the tool to call
//   - arguments: A map of argument names to their interface{} values
//
// Returns:
//   - []TextContent: A slice of text content responses from the tool execution
//   - error: An error if the tool name is unknown or if required arguments are missing
//
// Currently supported tools:
//   - "add-note": Adds a new note to the server
//     Required arguments:
//   - "name": string - The name of the note
//   - "content": string - The content of the note
//
// Thread safety:
// The function uses appropriate locking mechanisms when modifying the notes map.
func (s *Server) CallTool(name string, arguments map[string]interface{}) ([]TextContent, error) {
    fmt.Fprintf(os.Stderr, "Calling tool %s with arguments: %v\n", name, arguments)
    
    if name != "add-note" {
        return nil, fmt.Errorf("unknown tool: %s", name)
    }

    noteName, ok := arguments["name"].(string)
    if !ok || noteName == "" {
        fmt.Fprintf(os.Stderr, "Missing or invalid name argument\n")
        return nil, fmt.Errorf("missing or invalid name")
    }

    content, ok := arguments["content"].(string)
    if !ok || content == "" {
        fmt.Fprintf(os.Stderr, "Missing or invalid content argument\n")
        return nil, fmt.Errorf("missing or invalid content")
    }

    s.notesMap.Lock()
    s.notes[noteName] = content
    s.notesMap.Unlock()

    fmt.Fprintf(os.Stderr, "Added note '%s'\n", noteName)

    return []TextContent{{
        Type: "text",
        Text: fmt.Sprintf("Added note '%s' with content: %s", noteName, content),
    }}, nil
}