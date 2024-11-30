// File: cmd/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"notes-server/internal/server"
	"os"

	"github.com/kardianos/service"
)

var logger service.Logger

// Program structures.
type program struct {
	srv    *server.Server
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *program) Start(s service.Service) error {
	logger.Info("Starting service...")
	go p.run()
	return nil
}

func (p *program) run() {
	logger.Info("Service is now running")
	if err := p.srv.Run(p.ctx); err != nil {
		logger.Error(err)
	}
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stopping service...")
	p.cancel()
	return nil
}

// handleServiceCommand processes a service control command and provides user feedback
func handleServiceCommand(s service.Service, command string) error {
	switch command {
	case "install":
		fmt.Println("Installing service...")
		err := s.Install()
		if err != nil {
			return fmt.Errorf("failed to install service: %v", err)
		}
		fmt.Println("Service installed successfully")

	case "uninstall":
		fmt.Println("Uninstalling service...")
		err := s.Uninstall()
		if err != nil {
			return fmt.Errorf("failed to uninstall service: %v", err)
		}
		fmt.Println("Service uninstalled successfully")

	case "start":
		fmt.Println("Starting service...")
		err := s.Start()
		if err != nil {
			return fmt.Errorf("failed to start service: %v", err)
		}
		fmt.Println("Service started successfully")

	case "stop":
		fmt.Println("Stopping service...")
		err := s.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop service: %v", err)
		}
		fmt.Println("Service stopped successfully")

	case "restart":
		fmt.Println("Restarting service...")
		err := s.Restart()
		if err != nil {
			return fmt.Errorf("failed to restart service: %v", err)
		}
		fmt.Println("Service restarted successfully")

	case "status":
		status, err := s.Status()
		if err != nil {
			return fmt.Errorf("failed to get service status: %v", err)
		}
		switch status {
		case service.StatusRunning:
			fmt.Println("Service is running")
		case service.StatusStopped:
			fmt.Println("Service is stopped")
		default:
			fmt.Printf("Service status: %v\n", status)
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
		Description: "A service for running the notes MCP",
	}

	ctx, cancel := context.WithCancel(context.Background())
	prg := &program{
		srv:    server.NewServer("notes-server"),
		ctx:    ctx,
		cancel: cancel,
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalf("Failed to create service: %v\n", err)
	}

	// Get the logger
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatalf("Failed to create logger: %v\n", err)
	}

	// Handle command line arguments for service control
	if len(os.Args) > 1 {
		command := os.Args[1]
		if err := handleServiceCommand(s, command); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("\nAvailable commands:")
			fmt.Println("  install  - Install the service")
			fmt.Println("  uninstall - Remove the service")
			fmt.Println("  start    - Start the service")
			fmt.Println("  stop     - Stop the service")
			fmt.Println("  restart  - Restart the service")
			fmt.Println("  status   - Check service status")
			os.Exit(1)
		}
		return
	}

	// Run the service (this block is reached when no command-line arguments are provided)
	fmt.Println("Starting NotesServer service...")
	err = s.Run()
	if err != nil {
		logger.Error(err)
		fmt.Printf("Service failed: %v\n", err)
		os.Exit(1)
	}
}