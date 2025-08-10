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

## Implementation Details

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

This example now demonstrates how to use **Chi router explicitly** with **custom middleware** in a Dapr actor application.

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

## Running the Example

1. **Start Dapr sidecar**:
   ```bash
   dapr run --app-id example-actors --app-port 8080 --dapr-http-port 3500
   ```

2. **In another terminal, start the application**:
   ```bash
   cd examples/generated-complete
   go run .
   ```

3. **Test the actors** using HTTP calls:

### Counter Actor Examples

```bash
# Create/get counter (returns 0 initially)
# Note: First-time calls to new actors automatically initialize with default values
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

**Expected log output with custom headers:**
```
=== HTTP Headers for POST /actors/Counter/my-counter/method/Increment ===
Header: Content-Type: application/json
Header: User-Agent: curl/7.64.1
Header: X-Custom-Header: my-custom-value
Header: X-User-Id: user123
Header: X-Client-Version: 1.0.0
Header: Accept: */*
=== End Headers ===
Context enriched with RequestID: 20241201-143022.123
[Counter] Operation: Increment, RequestID: 20241201-143022.123
[Counter] Operation: Increment, User: example-user, Role: actor-service, Timestamp: 2024-12-01T14:30:22Z
[Counter] Incremented from 0 to 1
```

### BankAccount Actor Examples

```bash
# Create account
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account-123/method/CreateAccount \
  -H "Content-Type: application/json" \
  -d '{"ownerName": "John Doe", "initialDeposit": 100.0}'

# Deposit money
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account-123/method/Deposit \
  -H "Content-Type: application/json" \
  -d '{"amount": 50.0, "description": "Salary deposit"}'

# Get balance
curl -X GET http://localhost:3500/v1.0/actors/BankAccount/account-123/method/GetBalance

# Withdraw money
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account-123/method/Withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount": 25.0, "description": "ATM withdrawal"}'

# Get transaction history
curl -X GET http://localhost:3500/v1.0/actors/BankAccount/account-123/method/GetHistory
```

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

## Comparison with Generated Stubs

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

## For Developers

This example shows how to:
1. Take generated stub implementations
2. Add real business logic using in-memory storage
3. Implement different actor patterns (state-based vs event-sourced)
4. Handle validation, errors, and edge cases
5. Ensure thread-safe concurrent access
6. Create self-contained actor implementations

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

Compare the implemented files with the stub versions to understand the progression from generated code to working applications.