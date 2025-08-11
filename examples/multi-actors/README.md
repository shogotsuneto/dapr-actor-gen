# Multi-Actors Example

This directory contains an example demonstrating how to generate Dapr actors from an OpenAPI specification that defines multiple actor types.

## What This Example Demonstrates

The `openapi.yaml` specification defines two different actor types with distinct patterns:

- **Counter Actor** - Simple state-based actor with increment/decrement operations
- **BankAccount Actor** - Event-sourced actor with transaction history

This demonstrates:
1. **Multiple actor types** in a single OpenAPI specification
2. **Different actor patterns** (state-based vs event-sourced)
3. **Complete code generation** including interfaces, types, factories, and implementations
4. **Runnable example application** with custom middleware

## Files in This Example

- `openapi.yaml` - OpenAPI 3.0 specification defining Counter and BankAccount actors
- `generated/` - Complete generated code including working implementations

## Regenerating the Code

To regenerate the code from the OpenAPI specification:

### 1. Build the Generator

```bash
# From the repository root
go mod download
make build
```

### 2. Generate the Code

```bash
# Basic generation (interfaces only - will not compile)
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./output

# Generate with implementation stubs (compilable but not-implemented errors)
./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml ./output

# Generate with example application
./bin/dapr-actor-gen --generate-example examples/multi-actors/openapi.yaml ./output

# Generate everything (what's in the generated/ directory)
./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/openapi.yaml ./output
```

### 3. Prepare the Generated Code

```bash
# Navigate to the output directory
cd ./output

# Initialize Go module and download dependencies
go mod init your-module-name
go mod tidy

# Verify the code compiles
go build ./...
```

## Running the Example

The `generated/` directory contains a complete, working implementation that you can run immediately:

### 1. Start the Application

```bash
# Navigate to the generated example
cd examples/multi-actors/generated

# Install dependencies (if not already done)
go mod tidy

# Start Dapr sidecar in one terminal
dapr run --app-id example-actors --app-port 8080 --dapr-http-port 3500

# Start the application in another terminal
go run .
```

### 2. Test the Actors

**Counter Actor:**
```bash
# Get current value (returns 0 initially)
curl -X GET http://localhost:3500/v1.0/actors/Counter/my-counter/method/Get

# Increment counter
curl -X POST http://localhost:3500/v1.0/actors/Counter/my-counter/method/Increment

# Set specific value
curl -X POST http://localhost:3500/v1.0/actors/Counter/my-counter/method/Set \
  -H "Content-Type: application/json" \
  -d '{"value": 42}'
```

**BankAccount Actor:**
```bash
# Create account
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account-123/method/CreateAccount \
  -H "Content-Type: application/json" \
  -d '{"ownerName": "John Doe", "initialDeposit": 100.0}'

# Get balance
curl -X GET http://localhost:3500/v1.0/actors/BankAccount/account-123/method/GetBalance

# Deposit money
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account-123/method/Deposit \
  -H "Content-Type: application/json" \
  -d '{"amount": 50.0, "description": "Salary deposit"}'
```

## Key Features Demonstrated

- **Schema-first development** - Actors defined in OpenAPI, implemented in Go
- **Type safety** - Generated types match OpenAPI schema exactly
- **Multiple patterns** - State-based and event-sourced actor implementations
- **Custom middleware** - Chi router with logging and context enrichment
- **Real business logic** - Working implementations, not just stubs
- **Thread safety** - Proper concurrency handling with mutexes

## Comparison with Stub Generation

- **Basic generation**: Creates interfaces and types only (will not compile)
- **With --generate-impl**: Adds stub implementations that return "not implemented" errors
- **This example**: Shows fully working implementations with real business logic

The `generated/` directory demonstrates how to evolve from generated stubs to production-ready actor implementations.