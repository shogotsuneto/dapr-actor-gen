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
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./generated
```

### 3. Use Generated Code

The generator creates actor-specific packages:

```
generated/
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
./bin/dapr-actor-gen path/to/schema.yaml ./generated

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

Pre-built Docker images are available from GitHub Container Registry for immediate use, or you can build locally.

### Using Pre-built Images

Pull the latest image:
```bash
docker pull ghcr.io/shogotsuneto/dapr-actor-gen:latest
```

Or pull a specific version:
```bash
docker pull ghcr.io/shogotsuneto/dapr-actor-gen:v0.0.2
```

### Usage Examples

#### Basic Usage
Generate code from an OpenAPI specification:
```bash
# Using latest image
docker run --rm \
  -v $(pwd)/examples:/examples \
  -v $(pwd)/output:/output \
  ghcr.io/shogotsuneto/dapr-actor-gen:latest \
  /examples/multi-actors/openapi.yaml /output

# Using specific version
docker run --rm \
  -v $(pwd)/examples:/examples \
  -v $(pwd)/output:/output \
  ghcr.io/shogotsuneto/dapr-actor-gen:v0.0.2 \
  /examples/multi-actors/openapi.yaml /output
```

#### Development Workflow
Generate code from your project's OpenAPI spec:
```bash
# Mount your project directory and specify full paths
docker run --rm \
  -v $(pwd):/workspace \
  ghcr.io/shogotsuneto/dapr-actor-gen:latest \
  /workspace/api/openapi.yaml /workspace/generated
```

#### CI/CD Integration
Use in continuous integration pipelines:
```bash
# GitHub Actions, Jenkins, etc.
docker run --rm \
  -v $GITHUB_WORKSPACE:/workspace \
  ghcr.io/shogotsuneto/dapr-actor-gen:latest \
  /workspace/schema/actors.yaml /workspace/src/generated
```

#### Multiple Output Directories
Generate different actor types to separate directories:
```bash
# Generate to specific subdirectories
docker run --rm \
  -v $(pwd)/schemas:/schemas \
  -v $(pwd)/generated:/generated \
  ghcr.io/shogotsuneto/dapr-actor-gen:latest \
  /schemas/payment-actors.yaml /generated/payment

docker run --rm \
  -v $(pwd)/schemas:/schemas \
  -v $(pwd)/generated:/generated \
  ghcr.io/shogotsuneto/dapr-actor-gen:latest \
  /schemas/user-actors.yaml /generated/user
```

### Building Locally

If you need to build the Docker image yourself:

```bash
# Build Docker image
docker build -t dapr-actor-gen .

# Run locally built image
docker run --rm -v $(pwd)/examples:/examples -v $(pwd)/output:/output \
  dapr-actor-gen /examples/multi-actors/openapi.yaml /output
```

### Notes

- **File Permissions**: Generated files are created with the container user's permissions. Ensure your output directory has appropriate write permissions.
- **Volume Paths**: Use absolute paths in volume mounts (`/workspace/file.yaml` instead of `./file.yaml`) when mounting your project directory.

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
