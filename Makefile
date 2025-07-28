# Makefile for dapr-actor-gen
.PHONY: help build test clean tidy build-all build-linux build-darwin build-windows

# Configuration
GO_VERSION := 1.19
BINARY_NAME := dapr-actor-gen
BIN_DIR := bin
CMD_DIR := cmd
PKG_DIR := pkg
OUTPUT_DIR := generated
DIST_DIR := dist

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build the dapr-actor-gen binary
	@echo "$(GREEN)[INFO]$(NC) Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)[INFO]$(NC) ✓ $(BINARY_NAME) built successfully"

build-all: build-linux build-darwin build-windows ## Build binaries for all platforms

build-linux: ## Build Linux binaries (amd64 and arm64)
	@echo "$(GREEN)[INFO]$(NC) Building Linux binaries..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)
	@echo "$(GREEN)[INFO]$(NC) ✓ Linux binaries built successfully"

build-darwin: ## Build macOS binaries (amd64 and arm64)
	@echo "$(GREEN)[INFO]$(NC) Building macOS binaries..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	@echo "$(GREEN)[INFO]$(NC) ✓ macOS binaries built successfully"

build-windows: ## Build Windows binary (amd64)
	@echo "$(GREEN)[INFO]$(NC) Building Windows binary..."
	@mkdir -p $(DIST_DIR)
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "$(GREEN)[INFO]$(NC) ✓ Windows binary built successfully"

test: ## Run all tests
	@echo "$(GREEN)[INFO]$(NC) Running tests..."
	@go test -v ./...
	@echo "$(GREEN)[INFO]$(NC) ✓ All tests passed"

tidy: ## Tidy go modules
	@echo "$(GREEN)[INFO]$(NC) Tidying go modules..."
	@go mod tidy
	@echo "$(GREEN)[INFO]$(NC) ✓ Go modules tidied"

clean: ## Clean build artifacts and generated files
	@echo "$(GREEN)[INFO]$(NC) Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(DIST_DIR)
	@rm -rf $(OUTPUT_DIR)
	@rm -rf test-output
	@echo "$(GREEN)[INFO]$(NC) ✓ Clean completed"

# Development targets
dev-build: tidy build ## Development build (tidy + build)

dev-test: tidy test ## Development test (tidy + test)

# Check if Go is installed
check-go:
	@if ! command -v go &> /dev/null; then \
		echo "$(RED)[ERROR]$(NC) Go is not installed. Please install Go $(GO_VERSION)+ and try again."; \
		exit 1; \
	fi

