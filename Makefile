# Makefile for the notes application build system
# Supports development and release builds across multiple platforms

# Build configuration
BINARY_NAME=notes-server
SERVICE_NAME=notes-service
BUILD_DIR=bin

# Version information (should be updated for releases)
VERSION ?= 0.1.0
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
GO=go
# Debug builds: Include symbols and race detection
DEBUG_FLAGS=-race -gcflags="all=-N -l" -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"
# Release builds: Optimized, stripped, with version info
RELEASE_FLAGS=-trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# OS-specific commands and paths
MKDIR_P = mkdir -p
RM_RF = rm -rf
PATHSEP = /

# Default target builds everything in release mode
.PHONY: all
all: release-all

# Clean build artifacts
.PHONY: clean
clean:
	$(RM_RF) $(BUILD_DIR)

# Development builds
.PHONY: dev
dev: dev-windows dev-linux dev-darwin

.PHONY: dev-windows
dev-windows:
	$(MKDIR_P) $(BUILD_DIR)/dev/windows
	GOOS=windows GOARCH=amd64 $(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/dev/windows/$(BINARY_NAME).exe ./cmd
	GOOS=windows GOARCH=amd64 $(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/dev/windows/$(SERVICE_NAME).exe ./service

.PHONY: dev-linux
dev-linux:
	$(MKDIR_P) $(BUILD_DIR)/dev/linux
	GOOS=linux GOARCH=amd64 $(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/dev/linux/$(BINARY_NAME) ./cmd
	GOOS=linux GOARCH=amd64 $(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/dev/linux/$(SERVICE_NAME) ./service

.PHONY: dev-darwin
dev-darwin:
	$(MKDIR_P) $(BUILD_DIR)/dev/darwin
	GOOS=darwin GOARCH=amd64 $(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/dev/darwin/$(BINARY_NAME) ./cmd
	GOOS=darwin GOARCH=amd64 $(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/dev/darwin/$(SERVICE_NAME) ./service

# Release builds
.PHONY: release-all
release-all: release-windows release-linux release-darwin

.PHONY: release-windows
release-windows:
	$(MKDIR_P) $(BUILD_DIR)/release/windows
	GOOS=windows GOARCH=amd64 $(GO) build $(RELEASE_FLAGS) -o $(BUILD_DIR)/release/windows/$(BINARY_NAME).exe ./cmd
	GOOS=windows GOARCH=amd64 $(GO) build $(RELEASE_FLAGS) -o $(BUILD_DIR)/release/windows/$(SERVICE_NAME).exe ./service

.PHONY: release-linux
release-linux:
	$(MKDIR_P) $(BUILD_DIR)/release/linux
	GOOS=linux GOARCH=amd64 $(GO) build $(RELEASE_FLAGS) -o $(BUILD_DIR)/release/linux/$(BINARY_NAME) ./cmd
	GOOS=linux GOARCH=amd64 $(GO) build $(RELEASE_FLAGS) -o $(BUILD_DIR)/release/linux/$(SERVICE_NAME) ./service

.PHONY: release-darwin
release-darwin:
	$(MKDIR_P) $(BUILD_DIR)/release/darwin
	GOOS=darwin GOARCH=amd64 $(GO) build $(RELEASE_FLAGS) -o $(BUILD_DIR)/release/darwin/$(BINARY_NAME) ./cmd
	GOOS=darwin GOARCH=amd64 $(GO) build $(RELEASE_FLAGS) -o $(BUILD_DIR)/release/darwin/$(SERVICE_NAME) ./service

# Development environment targets
.PHONY: run-cmd run-service
run-cmd:
	$(GO) run $(DEBUG_FLAGS) ./cmd

run-service:
	$(GO) run $(DEBUG_FLAGS) ./service

# Build specific components (development mode)
.PHONY: build-cmd build-service
build-cmd:
	$(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

build-service:
	$(GO) build $(DEBUG_FLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME) ./service

# Help target
.PHONY: help
help:
	@echo Available targets:
	@echo   all            - Build release versions for all platforms
	@echo   dev            - Build development versions for all platforms
	@echo   clean          - Remove build artifacts
	@echo   
	@echo Development builds:
	@echo   dev-windows    - Build development version for Windows
	@echo   dev-linux      - Build development version for Linux
	@echo   dev-darwin     - Build development version for macOS
	@echo   run-cmd        - Run the command-line app in development mode
	@echo   run-service    - Run the service in development mode
	@echo   build-cmd      - Build only the command-line app (development)
	@echo   build-service  - Build only the service (development)
	@echo
	@echo Release builds:
	@echo   release-all    - Build release versions for all platforms
	@echo   release-windows- Build release version for Windows
	@echo   release-linux  - Build release version for Linux
	@echo   release-darwin - Build release version for macOS
	@echo
	@echo Variables:
	@echo   VERSION        - Set version number (current: $(VERSION))