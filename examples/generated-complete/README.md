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
- **State Management**: Stores a simple `int32` value using `StateManagerContext`
- **Default Value**: Returns 0 when no state exists

### BankAccount Actor (`bankaccount/impl.go`)
- **Pattern**: Event-sourced actor
- **Operations**:
  - `CreateAccount(ownerName, initialDeposit)` - Creates new account
  - `Deposit(amount, description)` - Deposits money
  - `Withdraw(amount, description)` - Withdraws money (with balance validation)
  - `GetBalance()` - Returns computed current state
  - `GetHistory()` - Returns complete transaction history
- **State Management**: Stores events and computes current state from event history
- **Event Types**: `AccountCreated`, `MoneyDeposited`, `MoneyWithdrawn`

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
curl -X GET http://localhost:3500/v1.0/actors/Counter/my-counter/method/Get

# Increment counter
curl -X POST http://localhost:3500/v1.0/actors/Counter/my-counter/method/Increment

# Set specific value
curl -X POST http://localhost:3500/v1.0/actors/Counter/my-counter/method/Set \
  -H "Content-Type: application/json" \
  -d '{"value": 42}'

# Decrement counter
curl -X POST http://localhost:3500/v1.0/actors/Counter/my-counter/method/Decrement
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

### State Management
Both actors use `a.GetStateManager()` to interact with Dapr's state store:
- `Set(ctx, key, value)` - Store value
- `Get(ctx, key, &variable)` - Retrieve value into variable
- `Save(ctx)` - Persist changes

### Error Handling
- Input validation (positive amounts, required fields)
- Business rule validation (sufficient funds for withdrawals)
- State operation error handling

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
    currentValue, err := a.getCurrentValue(ctx)
    if err != nil {
        return nil, err
    }
    
    newValue := currentValue + 1
    if err := a.saveValue(ctx, newValue); err != nil {
        return nil, err
    }
    
    return &CounterState{Value: newValue}, nil
}
```

## For Developers

This example shows how to:
1. Take generated stub implementations
2. Add real business logic using Dapr APIs
3. Implement different actor patterns (state-based vs event-sourced)
4. Handle validation, errors, and edge cases
5. Create production-ready actor implementations

Compare the implemented files with the stub versions to understand the progression from generated code to working applications.