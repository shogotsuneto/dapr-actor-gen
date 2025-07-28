# Examples

This directory contains examples showing how to use the Dapr Actor Code Generator.

## Example OpenAPI Schemas

### `multi-actors/openapi.yaml`

A complete example showing how to define multiple actor types in a single OpenAPI specification:

- **Counter Actor**: Simple counter with increment/decrement operations
- **BankAccount Actor**: Bank account with deposit/withdraw operations and balance tracking

This example demonstrates:
- Multiple actor types in one schema
- Different method signatures (void returns vs. object returns)
- Complex data types with nested objects
- Proper use of `x-actor-type` extension

## Usage

Generate code from the example schema:

```bash
# Using make targets
make generate SCHEMA=examples/multi-actors/openapi.yaml OUTPUT=./output

# Or using the convenient example target
make generate-example

# Or using the binary directly  
./bin/dapr-actor-gen examples/multi-actors/openapi.yaml ./output
```

This will create actor packages in the `./output` directory:

```
output/
├── counter/
│   ├── api.go          # Counter actor interface
│   ├── factory.go      # Registration factory
│   └── types.go        # CounterState type
└── bankaccount/
    ├── api.go          # BankAccount actor interface
    ├── factory.go      # Registration factory
    └── types.go        # BankAccountState, DepositRequest types
```

## Implementation Example

Once you've generated the code, implement your actors:

```go
package main

import (
    "context"
    "github.com/dapr/go-sdk/actor"
    "./output/counter"
)

// CounterActor implements the generated CounterAPI interface
type CounterActor struct {
    counter.CounterAPI  // Embeds actor.ServerContext
}

// Implement the interface methods
func (c *CounterActor) Increment(ctx context.Context) (*counter.CounterState, error) {
    // Get current state
    state, err := c.GetStateManager().Get(ctx, "count")
    if err != nil {
        return nil, err
    }
    
    var count int
    if state != nil {
        count = state.(int)
    }
    
    // Increment and save
    count++
    err = c.GetStateManager().Set(ctx, "count", count)
    if err != nil {
        return nil, err
    }
    
    return &counter.CounterState{Count: count}, nil
}

func (c *CounterActor) Decrement(ctx context.Context) (*counter.CounterState, error) {
    // Similar implementation...
}

func (c *CounterActor) GetCount(ctx context.Context) (*counter.CounterState, error) {
    // Similar implementation...
}

// Register with Dapr
func main() {
    s := daprd.NewService(":8080")
    s.RegisterActorImplFactoryContext(counter.NewActorFactory())
    s.Start()
}
```

## Key Benefits

1. **Type Safety**: Generated types match your OpenAPI schema exactly
2. **Compile-time Validation**: Interface ensures you implement all required methods
3. **Automatic Factories**: Ready-to-use registration functions
4. **Documentation**: Generated code includes comments from your OpenAPI spec
5. **Consistency**: All actors follow the same patterns and conventions