# go-mcp-server-service
A JSON-RPC 2.0 compliant server implementing the Model Context Protocol (MCP) for note management (as an example)

This is an example. You can modify this and use as a Boilerplate for your own project. It supports cross-platform development and includes a command-line interface and service component for both development and release builds.

## Features

- JSON-RPC 2.0 compliant API
- Cross-platform support (Windows, Linux, macOS)
- Thread-safe note management
- Development and release build configurations
- Service and command-line interface components

## Components

### Command Line Interface (`cmd`)

The command-line interface provides direct access to the notes server functionality.

### Service (`service`)

The service component enables system-level integration and background operation.

### Resources

The server implements a note storage system with:

- Custom `note://` URI scheme for accessing individual notes
- Resource metadata including name, description, and MIME type
- Thread-safe concurrent access

### Prompts

Available prompts:

- `summarize-notes`: Creates summaries of all stored notes
  - Optional `style` argument ("brief"/"detailed")
  - Combines all current notes with style preference
  - Thread-safe note access

### Tools

Available tools:

- `add-note`: Adds a new note to the server
  - Required arguments: `name` (string), `content` (string)
  - Thread-safe state updates
  - Returns confirmation message

## Building

### Prerequisites

- Go 1.21 or later
- GNU Make or compatible build tool
- Git (for version information)

### Build Commands

Development builds (includes debug symbols and race detection):

```bash
# Build all components for all platforms
make dev

# Build for specific platform
make dev-windows
make dev-linux
make dev-darwin

# Build specific components
make build-cmd
make build-service
```

Release builds (optimized and stripped):

```bash
# Build all components for all platforms
make release-all

# Build for specific platform
make release-windows
make release-linux
make release-darwin
```

Run locally:

```bash
# Run command-line interface
make run-cmd

# Run service
make run-service
```

View all available targets:

```bash
make help
```

### Build Output

Binaries are created in the `bin` directory:

- Development builds: `bin/dev/<platform>/`
- Release builds: `bin/release/<platform>/`

## Configuration

### Claude Desktop Integration

Configure the notes server in Claude Desktop's configuration file:

#### Location

- MacOS: `~/Library/Application\ Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%/Claude/claude_desktop_config.json`

#### Development Configuration

```json
{
  "mcpServers": {
    "notes-server": {
      "command": "./bin/dev/<platform>/notes-server",
      "args": []
    }
  }
}
```

#### Release Configuration (Service Example)

```json
{
  "mcpServers": {
    "notes-server": {
      "command": "./bin/release/<platform>/notes-service",
      "args": []
    }
  }
}
```

## Development

### Project Structure

```
.
├── cmd/                    # Command-line interface
├── service/               # Service implementation
├── internal/
│   └── server/           # Core server implementation
│       ├── operations.go # Server operations
│       ├── server.go    # Main server logic
│       └── types.go     # Type definitions
├── Makefile              # Build configuration
└── README.md
```

### Debugging

Since the server runs over stdio, we recommend using the [MCP Inspector](https://github.com/modelcontextprotocol/inspector) for debugging:

```bash
npx @modelcontextprotocol/inspector ./bin/dev/<platform>/notes-server
```

The Inspector will provide a URL for the debugging interface.

## Error Codes

The server implements standard JSON-RPC 2.0 error codes plus custom codes:

| Code   | Description           | Standard |
| ------ | --------------------- | -------- |
| -32700 | Parse error           | Yes      |
| -32600 | Invalid request       | Yes      |
| -32601 | Method not found      | Yes      |
| -32602 | Invalid params        | Yes      |
| -32603 | Internal error        | Yes      |
| -32001 | Resource not found    | No       |
| -32002 | Unsupported operation | No       |

## License

MIT License

Copyright (c) 2024 Andrew Lee Donelson