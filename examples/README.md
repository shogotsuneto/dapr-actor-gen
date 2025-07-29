# Examples

This directory contains examples showing how to use dapr-actor-gen.

## OpenAPI Specifications

- `multi-actors/openapi.yaml` - Example OpenAPI spec defining Counter and BankAccount actors

## Generated Examples

The following directories contain generated code demonstrating different generation modes:

- `generated-interfaces-only/` - Basic interface generation (default behavior)
- `generated-with-impl/` - Interfaces + partial implementation stubs (`--generate-impl`)
- `generated-with-example/` - Interfaces + example application (`--generate-example`)
- `generated-complete/` - All features combined (`--generate-impl --generate-example`)

These generated examples serve as:
- Reference implementations for users
- Regression tests for validating code generation changes
- Documentation of the tool's capabilities

All generated examples compile successfully and demonstrate proper Dapr actor integration.

## Usage

Generate code from the example schema:

```bash
# Basic interface generation (default behavior)
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./output

# Generate with partial implementation stubs
./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml ./output

# Generate with example application
./bin/dapr-actor-gen --generate-example examples/multi-actors/openapi.yaml ./output

# Generate everything
./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/openapi.yaml ./output
```

## Key Benefits

1. **Type Safety**: Generated types match your OpenAPI schema exactly
2. **Compile-time Validation**: Interface ensures you implement all required methods
3. **Automatic Factories**: Ready-to-use registration functions
4. **Documentation**: Generated code includes comments from your OpenAPI spec
5. **Consistency**: All actors follow the same patterns and conventions