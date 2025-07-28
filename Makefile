# Makefile for dapr-actor-gen
.PHONY: help build install test clean generate setup-env tidy

# Configuration
GO_VERSION := 1.19
BINARY_NAME := dapr-actor-gen
BIN_DIR := bin
CMD_DIR := cmd
PKG_DIR := pkg
OUTPUT_DIR := generated

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

install: build ## Install the generator (build binary and set up environment)
	@echo "$(GREEN)[INFO]$(NC) Installing API Generation Tools..."
	@echo "$(GREEN)[INFO]$(NC) Go version: $$(go version)"
	@echo "$(GREEN)[INFO]$(NC) ✓ $(BINARY_NAME) installed to $(BIN_DIR)/"
	@echo ""
	@echo "$(GREEN)[INFO]$(NC) To use the tools, either:"
	@echo "$(GREEN)[INFO]$(NC)   1. Add $$(pwd)/$(BIN_DIR) to your PATH"
	@echo "$(GREEN)[INFO]$(NC)   2. Run 'make setup-env' to export PATH"
	@echo "$(GREEN)[INFO]$(NC)   3. Use 'make generate' for code generation"

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
	@rm -rf $(OUTPUT_DIR)
	@rm -rf test-output
	@echo "$(GREEN)[INFO]$(NC) ✓ Clean completed"

setup-env: ## Export PATH to include the bin directory
	@echo "$(GREEN)[INFO]$(NC) API generation tools PATH setup:"
	@echo "export PATH=\"$$(pwd)/$(BIN_DIR):\$$PATH\""
	@echo ""
	@echo "$(GREEN)[INFO]$(NC) Available tools:"
	@ls -1 $(BIN_DIR) 2>/dev/null | sed 's/^/  - /' || echo "  (No tools found - run 'make install' first)"

generate: ## Generate code from OpenAPI schema (usage: make generate SCHEMA=path/to/schema.yaml OUTPUT=output/dir)
	@if [ -z "$(SCHEMA)" ]; then \
		echo "$(RED)[ERROR]$(NC) SCHEMA parameter is required"; \
		echo "$(YELLOW)[INFO]$(NC) Usage: make generate SCHEMA=examples/multi-actors/openapi.yaml [OUTPUT=./generated]"; \
		exit 1; \
	fi
	@if [ ! -f "$(SCHEMA)" ]; then \
		echo "$(RED)[ERROR]$(NC) Schema file not found: $(SCHEMA)"; \
		exit 1; \
	fi
	@if [ ! -f "$(BIN_DIR)/$(BINARY_NAME)" ]; then \
		echo "$(YELLOW)[WARN]$(NC) Binary not found, building first..."; \
		$(MAKE) build; \
	fi
	@OUTPUT_PATH=$${OUTPUT:-$(OUTPUT_DIR)}; \
	echo "$(GREEN)[INFO]$(NC) === API Code Generation ==="; \
	echo "$(GREEN)[INFO]$(NC) Schema File: $(SCHEMA)"; \
	echo "$(GREEN)[INFO]$(NC) Output Dir:  $$OUTPUT_PATH"; \
	echo ""; \
	mkdir -p "$$OUTPUT_PATH"; \
	./$(BIN_DIR)/$(BINARY_NAME) "$(SCHEMA)" "$$OUTPUT_PATH"; \
	echo ""; \
	echo "$(GREEN)[INFO]$(NC) Generated files:"; \
	find "$$OUTPUT_PATH" -type f -name "*.go" 2>/dev/null | sort | sed 's/^/  /' || echo "  (No files found)"; \
	echo ""; \
	echo "$(GREEN)[INFO]$(NC) ✓ Code generation completed successfully!"

# Example targets for common use cases
generate-example: ## Generate code from the example schema
	@$(MAKE) generate SCHEMA=examples/multi-actors/openapi.yaml OUTPUT=./generated

# Development targets
dev-build: tidy build ## Development build (tidy + build)

dev-test: tidy test ## Development test (tidy + test)

dev-all: clean dev-build dev-test ## Full development cycle (clean + build + test)

# Check if Go is installed
check-go:
	@if ! command -v go &> /dev/null; then \
		echo "$(RED)[ERROR]$(NC) Go is not installed. Please install Go $(GO_VERSION)+ and try again."; \
		exit 1; \
	fi

# Install with dependency check
install-safe: check-go install ## Install with Go dependency check