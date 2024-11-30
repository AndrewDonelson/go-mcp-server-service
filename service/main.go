// Package main implements the service wrapper for the notes server.
// It provides service management capabilities across different platforms
// using the kardianos/service library.
//
// The service can be installed, started, stopped, and uninstalled using
// standard service management commands.
//
// Usage:
//   - Install: notes-service install
//   - Start: notes-service start
//   - Stop: notes-service stop
//   - Uninstall: notes-service uninstall
//   - Run directly: notes-service
//
// The service maintains its own logging through the platform's service
// management system rather than writing directly to stdout/stderr.
package main

import (
    "context"
    "fmt"
    "notes-server/internal/server"
    "os"

    "github.com/kardianos/service"
)

var logger service.Logger

// program structures the note server for service management.
// It wraps the server instance and manages its lifecycle.
type program struct {
    srv    *server.Server
    ctx    context.Context
    cancel context.CancelFunc
}

func (p *program) Start(s service.Service) error {
    logger.Info("Starting notes service...")
    go p.run()
    return nil
}

func (p *program) run() {
    logger.Info("Notes service is now running")
    if err := p.srv.Run(p.ctx); err != nil {
        logger.Error(err)
    }
}

func (p *program) Stop(s service.Service) error {
    logger.Info("Stopping notes service...")
    p.cancel()
    return nil
}

// handleServiceCommand processes a service control command and provides user feedback
// through the service logger rather than directly to stdout/stderr.
func handleServiceCommand(s service.Service, command string) error {
    switch command {
    case "install":
        logger.Info("Installing service...")
        err := s.Install()
        if err != nil {
            return fmt.Errorf("failed to install service: %v", err)
        }
        logger.Info("Service installed successfully")

    case "uninstall":
        logger.Info("Uninstalling service...")
        err := s.Uninstall()
        if err != nil {
            return fmt.Errorf("failed to uninstall service: %v", err)
        }
        logger.Info("Service uninstalled successfully")

    case "start":
        logger.Info("Starting service...")
        err := s.Start()
        if err != nil {
            return fmt.Errorf("failed to start service: %v", err)
        }
        logger.Info("Service started successfully")

    case "stop":
        logger.Info("Stopping service...")
        err := s.Stop()
        if err != nil {
            return fmt.Errorf("failed to stop service: %v", err)
        }
        logger.Info("Service stopped successfully")

    case "restart":
        logger.Info("Restarting service...")
        err := s.Restart()
        if err != nil {
            return fmt.Errorf("failed to restart service: %v", err)
        }
        logger.Info("Service restarted successfully")

    case "status":
        status, err := s.Status()
        if err != nil {
            return fmt.Errorf("failed to get service status: %v", err)
        }
        switch status {
        case service.StatusRunning:
            logger.Info("Service is running")
        case service.StatusStopped:
            logger.Info("Service is stopped")
        default:
            logger.Infof("Service status: %v", status)
        }

    default:
        return fmt.Errorf("invalid command: %s", command)
    }
    return nil
}

func main() {
    svcConfig := &service.Config{
        Name:        "MCPServerNotes",
        DisplayName: "MCP Service - Notes",
        Description: "A service for running the notes MCP server",
        
        // Important: This option ensures service output is properly handled
        Option: map[string]interface{}{
            "LogOutput": true,
        },
    }

    ctx, cancel := context.WithCancel(context.Background())
    prg := &program{
        srv:    server.NewServer("notes-server"),
        ctx:    ctx,
        cancel: cancel,
    }

    s, err := service.New(prg, svcConfig)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create service: %v\n", err)
        os.Exit(1)
    }

    // Get the service logger
    logger, err = s.Logger(nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
        os.Exit(1)
    }

    // Handle command line arguments for service control
    if len(os.Args) > 1 {
        command := os.Args[1]
        if err := handleServiceCommand(s, command); err != nil {
            logger.Error(err)
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            fmt.Fprintf(os.Stderr, "\nAvailable commands:\n")
            fmt.Fprintf(os.Stderr, "  install  - Install the service\n")
            fmt.Fprintf(os.Stderr, "  uninstall - Remove the service\n")
            fmt.Fprintf(os.Stderr, "  start    - Start the service\n")
            fmt.Fprintf(os.Stderr, "  stop     - Stop the service\n")
            fmt.Fprintf(os.Stderr, "  restart  - Restart the service\n")
            fmt.Fprintf(os.Stderr, "  status   - Check service status\n")
            os.Exit(1)
        }
        os.Exit(0)
    }

    // Run the service
    logger.Info("Starting NotesServer service...")
    err = s.Run()
    if err != nil {
        logger.Error(err)
        fmt.Fprintf(os.Stderr, "Service failed: %v\n", err)
        os.Exit(1)
    }
}