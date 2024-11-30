// Package main is the entry point for the notes-server application.
// It provides a command-line interface to start the server and handles
// service management operations across different platforms.
//
// The server can be run directly or as a system service on Windows, Linux,
// and macOS. When run as a service, it supports standard service operations
// like install, uninstall, start, stop, and restart.
//
// Usage as a direct application:
//
//	$ notes-server
//
// Environment Variables:
//   - LOG_LEVEL: Set logging level (debug, info, warn, error). Default: info
//
// Exit Codes:
//   - 0: Successful execution
//   - 1: Fatal error during execution
package main

import (
    "context"
    "fmt"
    "os"
    "notes-server/internal/server"
)

// main is the entry point of the notes-server application.
// It initializes and runs the server with a background context.
// If the server encounters an error during execution, it will
// log the error and exit with status code 1.
//
// The server will continue running until it receives a termination
// signal (SIGTERM, SIGINT) or encounters a fatal error.
func main() {
    // Write all startup logging to stderr
    fmt.Fprintf(os.Stderr, "Starting notes-server...\n")

    // Create a new server instance with the default name
    srv := server.NewServer("notes-server")

    // Run the server with a background context
    // This will block until the server is shutdown or encounters an error
    if err := srv.Run(context.Background()); err != nil {
        // Log any fatal errors to stderr and exit with status code 1
        fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
        os.Exit(1)
    }
}