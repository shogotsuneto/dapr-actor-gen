# Generated Complete Example

This directory contains a **complete, working implementation** of Dapr actors generated from the OpenAPI specification in `../multi-actors/openapi.yaml`.

## Purpose

This example demonstrates:
1. **Generated interfaces and types** (from `dapr-actor-gen --generate-example --generate-impl`)
2. **Fully implemented business logic** showing real-world actor patterns
3. **Runnable application** that users can test and learn from

## What's Included

### Generated Files (DO NOT EDIT)
- `counter/api.go` - Counter actor interface
- `counter/types.go` - Counter data types
- `counter/factory.go` - Counter actor factory
- `bankaccount/api.go` - BankAccount actor interface  
- `bankaccount/types.go` - BankAccount data types
- `bankaccount/factory.go` - BankAccount actor factory
- `main.go` - Example application entry point
- `go.mod` - Go module definition

### Implemented Files (USER IMPLEMENTED)
- `counter/impl.go` - **Working Counter actor implementation**
- `bankaccount/impl.go` - **Working BankAccount actor implementation**

## Key Features Demonstrated

### Counter Actor (State-Based)
- Simple state management with Dapr StateManager
- CRUD operations (Get, Set, Increment, Decrement)
- Default value handling

### BankAccount Actor (Event-Sourced)
- Event sourcing pattern with complete audit trail
- Business logic validation (sufficient funds, positive amounts)
- State reconstruction from events
- Transaction history

## Running the Example

See [IMPLEMENTATION.md](./IMPLEMENTATION.md) for detailed usage instructions and API examples.

## For Developers

This example shows how to:
1. Take generated stub implementations
2. Add real business logic using Dapr APIs
3. Implement different actor patterns (state-based vs event-sourced)
4. Handle validation, errors, and edge cases
5. Create production-ready actor implementations

Compare the implemented files with the stub versions to understand the progression from generated code to working applications.