# Examples

This directory contains examples demonstrating how to use dapr-actor-gen to generate Dapr actors from OpenAPI specifications.

## Available Examples

### multi-actors/
Complete example showing multiple actor types (Counter and BankAccount) in a single OpenAPI specification. Includes:
- `openapi.yaml` - OpenAPI specification defining two actor types
- `generated/` - Complete working implementation with business logic
- `README.md` - Detailed documentation and usage instructions

This example demonstrates different actor patterns (state-based vs event-sourced) and shows the progression from generated stubs to production-ready implementations.

## Quick Start

1. **Build the tool:**
   ```bash
   make build
   ```

2. **Generate code:**
   ```bash
   # Basic generation (interfaces only)
   ./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./output
   
   # With implementation stubs
   ./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml ./output
   ```

3. **Compile generated code:**
   ```bash
   cd ./output
   go mod init your-module
   go mod tidy
   go build ./...
   ```


