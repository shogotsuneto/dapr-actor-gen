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
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml examples/multi-actors/generated

# Generate with implementation stubs (compilable but not-implemented errors)
./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml examples/multi-actors/generated

# Generate with example application
./bin/dapr-actor-gen --generate-example examples/multi-actors/openapi.yaml examples/multi-actors/generated

# Generate everything (what's in the generated/ directory)
./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/openapi.yaml examples/multi-actors/generated
```



## Running the Example

The `generated/` directory contains a complete, working implementation that you can run immediately:

### 1. Start the Application

```bash
# Navigate to the generated example
cd examples/multi-actors/generated

# Install dependencies (if not already done)
go mod tidy

# Start Dapr sidecar in one terminal (run 'dapr init' if Dapr is not initialized, 'dapr uninstall --all' to clean up)
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

# Decrement counter
curl -X POST http://localhost:3500/v1.0/actors/Counter/my-counter/method/Decrement

# Example with custom headers (demonstrates header logging middleware)
curl -X POST http://localhost:3500/v1.0/actors/Counter/my-counter/method/Increment \
  -H "Content-Type: application/json" \
  -H "X-Custom-Header: my-custom-value" \
  -H "X-User-Id: user123" \
  -H "X-Client-Version: 1.0.0"
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

# Withdraw money
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account-123/method/Withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount": 25.0, "description": "ATM withdrawal"}'

# Get transaction history
curl -X GET http://localhost:3500/v1.0/actors/BankAccount/account-123/method/GetHistory
```

## Implementation Details

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

### Counter Actor (`counter/impl.go`)
- **Pattern**: State-based actor
- **Operations**: 
  - `Get()` - Retrieves current counter value
  - `Increment()` - Adds 1 to counter
  - `Decrement()` - Subtracts 1 from counter  
  - `Set(value)` - Sets counter to specific value
- **Storage**: In-memory storage with thread-safe access using sync.RWMutex
- **Default Value**: Returns 0 when newly created

### BankAccount Actor (`bankaccount/impl.go`)
- **Pattern**: Event-sourced actor
- **Operations**:
  - `CreateAccount(ownerName, initialDeposit)` - Creates new account
  - `Deposit(amount, description)` - Deposits money
  - `Withdraw(amount, description)` - Withdraws money (with balance validation)
  - `GetBalance()` - Returns computed current state
  - `GetHistory()` - Returns complete transaction history
- **Storage**: In-memory event storage with thread-safe access using sync.RWMutex
- **Event Types**: `AccountCreated`, `MoneyDeposited`, `MoneyWithdrawn`

## Middleware and Chi Router

This example demonstrates how to use **Chi router explicitly** with **custom middleware** in a Dapr actor application.

### Middleware Features

The `main.go` file shows how to:

1. **Use Chi router explicitly** with `daprd.NewServiceWithMux()` instead of the default router
2. **Add built-in Chi middleware**:
   - `middleware.Logger` - Logs all HTTP requests with response times
   - `middleware.Recoverer` - Recovers from panics and returns 500 status
   - `middleware.RequestID` - Adds X-Request-Id header to responses
   - `middleware.RealIP` - Sets the real IP address from X-Forwarded-For headers

3. **Add custom middleware**:
   - `headerLoggingMiddleware` - **Logs all HTTP headers** from incoming requests for debugging
   - `contextEnrichmentMiddleware` - **Adds custom values to request context** that actors can access

### Custom Context Values

The middleware adds these values to the request context:
- **RequestID**: Unique identifier for each request (timestamp-based)
- **UserInfo**: Map containing user, role, and timestamp information

### Actor Context Usage

The **Counter actor** (`counter/impl.go`) demonstrates how to:
- **Retrieve context values** set by middleware using `ctx.Value(key)`
- **Log context information** in actor methods (`Get`, `Increment`)
- **Access middleware-provided data** within actor business logic

## Key Implementation Patterns

### In-Memory Storage
Both actors use in-memory storage for simplicity and self-contained operation:
- **Counter**: Thread-safe `int32` value protected by `sync.RWMutex`
- **BankAccount**: Thread-safe event slice protected by `sync.RWMutex`
- **Concurrency**: Proper locking ensures data consistency across concurrent requests

### Error Handling
- Input validation (positive amounts, required fields)
- Business rule validation (sufficient funds for withdrawals)
- Thread-safe access to shared data

### Actor Lifecycle
- **Instance Creation**: Each actor ID gets its own in-memory storage
- **Data Persistence**: Data persists for the lifetime of the actor instance
- **State Reset**: Restarting the application resets all actor state

### Actor ID Access
Use `a.ID()` to get the current actor instance ID.

## Key Features Demonstrated

- **Schema-first development** - Actors defined in OpenAPI, implemented in Go
- **Type safety** - Generated types match OpenAPI schema exactly
- **Multiple patterns** - State-based and event-sourced actor implementations
- **Custom middleware** - Chi router with logging and context enrichment

## Comparison with Stub Generation

- **Basic generation**: Creates interfaces and types only (will not compile)
- **With --generate-impl**: Adds stub implementations that return "not implemented" errors
- **This example**: Shows fully working implementations with real business logic

The original generated `impl.go` files contained:
```go
func (a *Counter) Increment(ctx context.Context) (*CounterState, error) {
    return nil, errors.New("Increment method is not implemented")
}
```

The working implementation shows:
```go
func (a *Counter) Increment(ctx context.Context) (*CounterState, error) {
    currentValue := a.getCurrentValue(ctx)
    newValue := currentValue + 1
    a.setValue(ctx, newValue)
    
    return &CounterState{Value: newValue}, nil
}
```



## Troubleshooting

### Common Issues

- **EOF errors**: Usually indicate Dapr sidecar connectivity issues - make sure Dapr sidecar is running
- **Actor not found**: Make sure the actor type is properly registered in main.go
- **Concurrent access**: In-memory storage is thread-safe using mutexes

### Expected Responses

- **New Counter**: `{"value": 0}` (default state)
- **New BankAccount**: Error "account not found - no events exist" until `CreateAccount` is called

### Data Persistence Notes

- **Actor lifetime**: Data persists for the lifetime of each actor instance
- **Application restart**: All actor state is reset when the application restarts
- **Production use**: For persistent storage, consider using Dapr state management components

## For Developers

This example shows how to:
1. Take generated stub implementations
2. Add real business logic using in-memory storage
3. Create self-contained actor implementations