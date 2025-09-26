# Dapr Actor Code Generator

A standalone code generator for creating Go interfaces, types, and factory functions from Actor Schema definitions for Dapr actors.

## Overview

This tool enables schema-first development for Dapr actors by generating Go code from **Actor Schema** definitions (default) or OpenAPI specifications. It creates:

- **Actor interfaces** with proper Dapr actor method signatures
- **Type definitions** from schema specifications
- **Factory functions** for actor registration
- **Complete actor packages** ready for implementation

The tool uses an intuitive **Actor Schema format** by default, which is specifically designed for defining actors rather than REST APIs.

## Quick Start

### Option 1: Using Docker (Recommended)

```bash
# Pull and use the latest pre-built image
docker pull ghcr.io/shogotsuneto/dapr-actor-gen:latest

# Generate code from Actor Schema (default format)
docker run --rm \
  -v $(pwd)/examples:/examples \
  -v $(pwd)/output:/output \
  ghcr.io/shogotsuneto/dapr-actor-gen:latest \
  /examples/multi-actors/actors.yaml /output
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/shogotsuneto/dapr-actor-gen.git
cd dapr-actor-gen

# Build the generator binary
make build

# Use the binary to generate from Actor Schema (default format)
./bin/dapr-actor-gen examples/multi-actors/actors.yaml ./generated
```

### 3. Use Generated Code

The generator creates actor-specific packages:

```
generated/
├── counteractor/
│   ├── api.go          # Generated interfaces and constants
│   └── types.go        # Generated type definitions
└── bankaccountactor/
    ├── api.go
    └── types.go

# With --generate-example flag:
generated/
├── counteractor/
│   ├── api.go          # Generated interfaces and constants
│   ├── factory.go      # Factory functions for registration (example only)
│   └── types.go        # Generated type definitions
├── bankaccountactor/
│   ├── api.go
│   ├── factory.go      # Factory functions for registration (example only)
│   └── types.go
├── main.go             # Example application
└── go.mod              # Go module file
```

Implement your actor by embedding the generated interface:

```go
package main

import (
    "context"
    "github.com/dapr/go-sdk/actor"
    "./generated/counteractor"
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

## Actor Factory Registration

### Using Generated Factories

When using `--generate-example`, factory functions are generated for convenience. You can use them as-is or customize them for your needs:

```go
// Generated factory (available with --generate-example)
s.RegisterActorImplFactoryContext(counteractor.NewActorFactory())
```

The generated factory functions can also be customized by modifying the `factory.go` files in your generated code.

### Using Anonymous Factories

For cases where you want to customize factories (e.g., dependency injection), register actors using anonymous factory functions:

```go
// Anonymous factory with dependency injection
func main() {
    s := daprd.NewService(":8080")
    
    // Create your dependencies
    database := setupDatabase()
    logger := setupLogger()
    
    // Register actor with custom factory
    s.RegisterActorImplFactoryContext(func() actor.ServerContext {
        return &CounterActor{
            Database: database,
            Logger:   logger,
        }
    })
    
    s.Start()
}
```

This approach allows you to:
- Inject dependencies into your actors
- Customize actor initialization
- Avoid regenerating factory code when updating your OpenAPI schema

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

## Actor Schema Format (Default)

The Actor Schema format is the default and most intuitive way to define Dapr actors. It uses a clean, actor-centric YAML structure:

```yaml
actors:
  Counter:
    description: Simple state-based counter actor
    methods:
      Increment:
        description: Increment counter by 1
        returns: CounterState
      GetValue:
        description: Get current counter value
        returns: CounterState
      SetValue:
        description: Set counter to specific value
        request: SetValueRequest
        returns: CounterState

types:
  CounterState:
    type: object
    properties:
      value:
        type: integer
        format: int32
      status:
        type: string
        enum: [active, paused, error]
    required: [value, status]

  SetValueRequest:
    type: object
    properties:
      value:
        type: integer
        format: int32
    required: [value]
```

### Using the Binary

After building with `make build`, generate code using the Actor Schema format (default):

```bash
# Generate interfaces only
./bin/dapr-actor-gen examples/multi-actors/actors.yaml ./generated

# Generate with implementation stubs
./bin/dapr-actor-gen --generate-impl examples/multi-actors/actors.yaml ./generated

# Generate complete example application
./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/actors.yaml ./generated
```

### Method Patterns

The Actor Schema format supports different method patterns:

```yaml
actors:
  BankAccount:
    methods:
      # Query method (no input, returns data)
      GetBalance:
        description: Get current account balance
        returns: BankAccountState
      
      # Command method (takes input, returns data)
      Deposit:
        description: Deposit money to account
        request: DepositRequest
        returns: BankAccountState
      
      # Action method (takes input, no return)
      ProcessPayment:
        description: Process a payment
        request: PaymentRequest
      
      # Simple action (no input, no return)
      Reset:
        description: Reset the account
```

## Examples

The `examples/` directory contains:

- **multi-actors/openapi.yaml** - Example OpenAPI spec with multiple actor types
- Generated code examples and documentation

## Command Line Usage

```bash
dapr-actor-gen [flags] <openapi-file> <output-directory>
```

### Arguments

- `openapi-file`: Path to your OpenAPI 3.0 specification file (YAML or JSON)
- `output-directory`: Directory where generated code will be placed

### Command Options

- `--generate-impl`: Generate partial implementation stubs with not-implemented errors  
- `--generate-example`: Generate example main.go, go.mod and other files for a complete app
- `-format`: Specify input format ('actor-schema' default, 'openapi' for OpenAPI 3.0)

### Usage Examples

```bash
# Generate interfaces only (Actor Schema format - default)
dapr-actor-gen actors.yaml ./output

# Generate interfaces + partial implementations  
dapr-actor-gen --generate-impl actors.yaml ./output

# Generate interfaces + example application
dapr-actor-gen --generate-example actors.yaml ./output

# Generate everything together
dapr-actor-gen --generate-impl --generate-example actors.yaml ./output
```

#### Partial Implementation Generation (`--generate-impl`)

Generates stub implementations alongside the existing API definitions. This creates `impl.go` files with method stubs that return not-implemented errors.

#### Example Application Generation (`--generate-example`)

Creates a complete, compilable Dapr application with `main.go` and `go.mod` that demonstrates how to register and use the generated actors.

### Generated File Structure

For each actor type found in your OpenAPI spec, the generator creates:

- `{actortype}/api.go` - Main interface that embeds `actor.ServerContext`
- `{actortype}/types.go` - Type definitions from OpenAPI schemas

When using `--generate-example`:
- `{actortype}/factory.go` - Factory function for Dapr registration

## Features

- ✅ **Actor Schema Format** - Native actor-centric YAML format (default and recommended)
- ✅ **OpenAPI 3.0 Support** - Alternative support for existing OpenAPI specifications
- ✅ **Multiple Actor Types** - Generate multiple actors from one specification  
- ✅ **Type Safety** - Generated types match your schemas exactly
- ✅ **Dapr Integration** - Ready-to-use with Dapr Go SDK
- ✅ **Factory Functions** - Automatic registration helpers
- ✅ **Explicit Format Selection** - Clean `-format` flag specification (defaults to actor-schema)
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

### Cross-Compilation

Build binaries for multiple platforms:

```bash
# Build for all platforms
make build-all

# Or build for specific platforms
make build-linux   # Linux amd64/arm64
make build-darwin  # macOS amd64/arm64
make build-windows # Windows amd64
```

Binaries will be available in the `dist/` directory.

## Testing

```bash
# Run tests using make (recommended)
make test
# Or run directly
go test ./...
```

## Docker

For local Docker builds:

```bash
# Build Docker image
docker build -t dapr-actor-gen .

# Run locally built image (using Actor Schema default format)
docker run --rm -v $(pwd)/examples:/examples -v $(pwd)/output:/output \
  dapr-actor-gen /examples/multi-actors/actors.yaml /output
```

## Releases

Releases are created through GitHub Actions with manual triggers. Each release includes:

- **Multi-platform binaries**: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)
- **Docker images**: Multi-architecture images published to GitHub Container Registry
- **Release notes**: Automated generation with changelog and installation instructions

### Release Process

Releases can only be created from the `main` branch by maintainers:

1. Go to the [Actions tab](../../actions/workflows/release.yml) in GitHub
2. Click "Run workflow"
3. Enter the version in `v*.*.*` format (e.g., `v1.0.0`)
4. The workflow will create a draft release with all artifacts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
