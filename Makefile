# Build configuration
BINARY_NAME=notes-server
SERVICE_NAME=notes-service
BUILD_DIR=bin
VERSION ?= 0.1.0

ifeq ($(OS),Windows_NT)
    # On Windows, delegate to build.bat
    SHELL=cmd.exe
    .SHELLFLAGS=/c

clean:
	@build.bat clean

dev:
	@build.bat dev

release-all:
	@build.bat release

release-windows:
	@build.bat release-windows

help:
	@build.bat help

else
    # Unix system, use direct commands
    RM_CMD = rm -rf
    MKDIR_CMD = mkdir -p

clean:
	$(RM_CMD) $(BUILD_DIR)

dev:
	$(MKDIR_CMD) $(BUILD_DIR)/dev/linux
	$(MKDIR_CMD) $(BUILD_DIR)/dev/darwin
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/dev/linux/$(BINARY_NAME) ./cmd
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/dev/linux/$(SERVICE_NAME) ./service
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/dev/darwin/$(BINARY_NAME) ./cmd
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/dev/darwin/$(SERVICE_NAME) ./service

release-all: release-linux release-darwin

release-linux:
	$(MKDIR_CMD) $(BUILD_DIR)/release/linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/release/linux/$(BINARY_NAME) ./cmd
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/release/linux/$(SERVICE_NAME) ./service

release-darwin:
	$(MKDIR_CMD) $(BUILD_DIR)/release/darwin
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/release/darwin/$(BINARY_NAME) ./cmd
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/release/darwin/$(SERVICE_NAME) ./service

help:
	@echo "Available commands:"
	@echo "  clean         - Remove build artifacts"
	@echo "  dev          - Build development versions"
	@echo "  release-all  - Build all release versions"
	@echo "  help         - Show this help"

endif

.PHONY: clean dev release-all release-windows help