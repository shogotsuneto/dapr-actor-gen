# Actor YAML Format Example

This directory demonstrates the new Actor YAML format that provides a more intuitive way to define Dapr actors compared to OpenAPI specifications.

## What's Different About Actor YAML?

### Before (OpenAPI Format)
```yaml
paths:
  /Counter/{actorId}/method/Increment:
    post:
      summary: Increment counter by 1
      parameters:
        - name: actorId
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
```

### After (Actor YAML Format)
```yaml
actors:
  Counter:
    description: Simple state-based counter actor
    methods:
      Increment:
        description: Increment counter by 1
        returns: CounterState
```

## Key Benefits

1. **Actor-Centric**: Define actors directly instead of as REST API paths
2. **More Concise**: Less boilerplate than OpenAPI specifications
3. **Intuitive**: Natural way to think about actor methods and types
4. **Same Output**: Generates identical Go code as OpenAPI format
5. **Backward Compatible**: OpenAPI format still works alongside Actor YAML

## File Structure

- `actors.yaml` - Example Actor YAML definition (equivalent to `openapi.yaml`)
- `openapi.yaml` - Original OpenAPI definition for comparison

## Usage Examples

### Generate from Actor YAML
```bash
# Generate interfaces + implementations from Actor YAML
./bin/dapr-actor-gen --generate-impl examples/multi-actors/actors.yaml ./output

# Generate complete example application
./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/actors.yaml ./output
```

### Still Works with OpenAPI
```bash
# Original OpenAPI format still supported
./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml ./output
```

## Actor YAML Schema

### Basic Structure
```yaml
actors:
  ActorTypeName:
    description: Optional description of the actor
    methods:
      MethodName:
        description: Optional method description
        request: RequestTypeName    # Optional - for methods that take input
        returns: ResponseTypeName   # Optional - for methods that return data

types:
  TypeName:
    type: object|string|integer|number|boolean|array
    # ... standard OpenAPI schema properties
```

### Method Patterns

#### Query Method (no input, returns data)
```yaml
GetBalance:
  description: Get current account balance
  returns: BankAccountState
```

#### Command Method (takes input, returns data)
```yaml
Deposit:
  description: Deposit money to account
  request: DepositRequest
  returns: BankAccountState
```

#### Action Method (takes input, no return)
```yaml
ProcessPayment:
  description: Process a payment
  request: PaymentRequest
```

#### Simple Action (no input, no return)
```yaml
Reset:
  description: Reset the counter
```

### Type Definitions

Actor YAML uses the same type definition syntax as OpenAPI schemas:

```yaml
types:
  # Object types
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

  # Enum types
  CounterStatus:
    type: string
    enum: [active, paused, error, reset]

  # Array types
  EventList:
    type: array
    items:
      $ref: "#/types/AccountEvent"

  # Simple type aliases
  UserId:
    type: string
    pattern: '^[a-zA-Z0-9_-]+$'
```

## Format Detection

The tool automatically detects the format:

- **Actor YAML**: Files containing `actors:` top-level key
- **OpenAPI**: Files containing `openapi:` or `swagger:` or REST-style `paths:`

You can mix both formats in the same project - the tool will handle each file appropriately.

## Migration from OpenAPI

To migrate existing OpenAPI specifications to Actor YAML:

1. **Extract Actor Types**: Look for path patterns like `/{actorType}/{actorId}/method/{methodName}`
2. **Group by Actor**: Collect all methods for each actor type
3. **Simplify Method Definitions**: Convert HTTP operations to method definitions
4. **Keep Type Definitions**: Copy `components.schemas` to `types` (minimal changes needed)

## Comparison Output

Both formats generate identical Go code:

```bash
# Both commands produce the same output structure
./bin/dapr-actor-gen --generate-impl examples/multi-actors/openapi.yaml ./output-openapi
./bin/dapr-actor-gen --generate-impl examples/multi-actors/actors.yaml ./output-actors

# Verify they're equivalent
diff -r ./output-openapi ./output-actors
```

## When to Use Each Format

### Use Actor YAML When:
- Defining new actors from scratch
- Working primarily with actor patterns
- Want cleaner, more readable definitions
- Team prefers actor-centric thinking

### Use OpenAPI When:
- Migrating existing REST APIs to actors
- Need full OpenAPI ecosystem tooling
- Want to expose actors as REST endpoints
- Working with teams familiar with OpenAPI

## Example Actors in This Directory

The `actors.yaml` file demonstrates:

1. **Counter Actor**: Simple state-based actor with basic operations
2. **BankAccount Actor**: Event-sourced actor with complex business logic
3. **Rich Type System**: Objects, enums, arrays, and type references
4. **Various Method Patterns**: Queries, commands, actions with different signatures

Both generate the same complete, compilable Go actor packages.