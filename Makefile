# DevGen CLI - Simple Makefile for building

.PHONY: build clean install test run help

# Variables
BINARY_NAME := devgen
BUILD_DIR := build
VERSION := v1.0.0

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "DevGen CLI - Build Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the CLI binary
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) src/*.go
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

install: build ## Install binary to /usr/local/bin
	@echo "📥 Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ $(BINARY_NAME) installed successfully"

install-user: build ## Install binary to ~/.local/bin (no sudo required)
	@echo "📥 Installing $(BINARY_NAME) to user directory..."
	@mkdir -p ~/.local/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) ~/.local/bin/
	@echo "✅ $(BINARY_NAME) installed to ~/.local/bin"
	@echo "💡 Make sure ~/.local/bin is in your PATH"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test ./...

run: build ## Build and run the CLI
	@echo "🏃 Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy

version: ## Show version
	@echo "DevGen CLI $(VERSION)"

list: build ## List MCP servers
	@./$(BUILD_DIR)/$(BINARY_NAME) server list

dashboard: build ## Run interactive dashboard
	@./$(BUILD_DIR)/$(BINARY_NAME) dashboard

status: build ## Show server status
	@./$(BUILD_DIR)/$(BINARY_NAME) server status

init: build ## Initialize configuration
	@./$(BUILD_DIR)/$(BINARY_NAME) config init
