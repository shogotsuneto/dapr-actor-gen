# Dapr Actor Code Generator

A standalone code generator for creating Go interfaces, types, and factory functions from OpenAPI 3.0 specifications for Dapr actors.

## Overview

This tool enables schema-first development for Dapr actors by generating Go code from OpenAPI specifications. It creates:

- **Actor interfaces** with proper Dapr actor method signatures
- **Type definitions** from OpenAPI schemas
- **Factory functions** for actor registration
- **Complete actor packages** ready for implementation

## Quick Start

### 1. Build the Generator

```bash
# Clone the repository
git clone https://github.com/shogotsuneto/dapr-actor-gen.git
cd dapr-actor-gen

# Build the generator binary
make build
```

This will build `dapr-actor-gen` and place it in `bin/`.

### 2. Generate Code from OpenAPI Schema

```bash
# Use the binary directly to generate from the example schema
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./output

# Or generate to a different output directory
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./generated
```

### 3. Use Generated Code

The generator creates actor-specific packages:

```
output/
├── counteractor/
│   ├── api.go          # Generated interfaces and constants
│   ├── factory.go      # Factory functions for registration
│   └── types.go        # Generated type definitions
└── bankaccountactor/
    ├── api.go
    ├── factory.go
    └── types.go
```

Implement your actor by embedding the generated interface:

```go
package main

import (
    "context"
    "github.com/dapr/go-sdk/actor"
    "./output/counteractor"
)

// Implementation struct embeds the generated API interface
type CounterActor struct {
    counteractor.CounterActorAPI  // Embeds actor.ServerContext
}

// Implement the methods defined in your OpenAPI schema
func (c *CounterActor) Increment(ctx context.Context) (*counteractor.CounterState, error) {
    // Your implementation here
    return &counteractor.CounterState{Count: 1}, nil
}

// Register with Dapr using generated factory
func main() {
    s := daprd.NewService(":8080")
    s.RegisterActorImplFactoryContext(counteractor.NewActorFactory())
    s.Start()
}
```

## Available Make Targets

The project uses Make for common development tasks:

```bash
# Show all available targets
make help

# Build the binary
make build

# Run tests
make test

# Clean build artifacts
make clean

# Tidy go modules
make tidy
```

## Using the Binary Directly

After building with `make build`, the binary will be available at `./bin/dapr-actor-gen`:

```bash
# Generate code from any OpenAPI schema
./bin/dapr-actor-gen path/to/schema.yaml ./output

# Example with the provided sample schema
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./generated
```

## OpenAPI Schema Requirements

Your OpenAPI specification should follow these conventions for Dapr actors:

```yaml
openapi: 3.0.0
info:
  title: My Actors API
  version: 1.0.0

paths:
  # Actor methods are defined as paths
  /counter/{id}/increment:
    post:
      operationId: increment
      x-actor-type: counter      # Required: specifies the actor type
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Counter incremented
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CounterState'

components:
  schemas:
    CounterState:
      type: object
      properties:
        count:
          type: integer
```

Key requirements:
- Use `x-actor-type` extension to specify which actor type the method belongs to
- Actor ID should be a path parameter named `id`
- Method names come from `operationId`
- Request/response schemas become Go types

## Examples

The `examples/` directory contains:

- **multi-actors/openapi.yaml** - Example OpenAPI spec with multiple actor types
- Generated code examples and documentation

## Command Line Usage

```bash
dapr-actor-gen <openapi-file> <output-directory>
```

### Arguments

- `openapi-file`: Path to your OpenAPI 3.0 specification file (YAML or JSON)
- `output-directory`: Directory where generated code will be placed

### Generated File Structure

For each actor type found in your OpenAPI spec, the generator creates:

- `{actortype}/api.go` - Main interface that embeds `actor.ServerContext`
- `{actortype}/types.go` - Type definitions from OpenAPI schemas
- `{actortype}/factory.go` - Factory function for Dapr registration

## Features

- ✅ **OpenAPI 3.0 Support** - Full support for OpenAPI specifications
- ✅ **Multiple Actor Types** - Generate multiple actors from one spec
- ✅ **Type Safety** - Generated types match your OpenAPI schemas exactly
- ✅ **Dapr Integration** - Ready-to-use with Dapr Go SDK
- ✅ **Factory Functions** - Automatic registration helpers
- 🔄 **Future**: Protocol Buffers, JSON Schema, GraphQL support

## Building from Source

```bash
git clone https://github.com/shogotsuneto/dapr-actor-gen.git
cd dapr-actor-gen
go mod tidy
# Build using make (recommended)
make build
# Or build directly
go build -o bin/dapr-actor-gen ./cmd
```

The built binary will be available at `bin/dapr-actor-gen`.

## Testing

```bash
# Run tests using make (recommended)
make test
# Or run directly
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

[Add your license information here]
