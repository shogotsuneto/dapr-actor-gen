# Enum Generation Example

The dapr-actor-gen tool now supports generating enums from OpenAPI specifications. Enums must be defined as separate schema components (not inline field enums).

## OpenAPI Schema Definition

```yaml
components:
  schemas:
    UserStatus:
      type: string
      enum:
        - active
        - inactive
        - suspended
        - pending
      description: User account status

    UserProfile:
      type: object
      properties:
        id:
          type: string
        status:
          $ref: '#/components/schemas/UserStatus'
        name:
          type: string
      required:
        - id
        - status
        - name
```

## Generated Go Code

The generator produces:

### Type Definition
```go
// UserStatus User account status
type UserStatus string
```

### Constants
```go
// UserStatus constants
const (
    UserStatusActive UserStatus = "active"
    UserStatusInactive UserStatus = "inactive"
    UserStatusSuspended UserStatus = "suspended"
    UserStatusPending UserStatus = "pending"
)
```

### Validation Method
```go
// IsValid returns true if the UserStatus value is valid
func (e UserStatus) IsValid() bool {
    switch e {
    case UserStatusActive:
        return true
    case UserStatusInactive:
        return true
    case UserStatusSuspended:
        return true
    case UserStatusPending:
        return true
    default:
        return false
    }
}
```

### String Method
```go
// String returns the string representation of UserStatus
func (e UserStatus) String() string {
    return string(e)
}
```

### All Values Function
```go
// AllUserStatusValues returns all valid UserStatus values
func AllUserStatusValues() []UserStatus {
    return []UserStatus{
        UserStatusActive,
        UserStatusInactive,
        UserStatusSuspended,
        UserStatusPending,
    }
}
```

## Usage Example

```go
// Using enum constants
status := user.UserStatusActive

// Validation
if status.IsValid() {
    fmt.Printf("Status %s is valid\n", status)
}

// Getting all values
allStatuses := user.AllUserStatusValues()
for _, s := range allStatuses {
    fmt.Printf("Valid status: %s\n", s)
}

// Type safety in structs
profile := user.UserProfile{
    Id:     "user-123",
    Status: user.UserStatusPending,
    Name:   "John Doe",
}
```

## Benefits

1. **Type Safety**: Compile-time checking for valid enum values
2. **Validation**: Runtime validation with `IsValid()` method
3. **Discovery**: `AllEnumValues()` function for getting all valid values
4. **Consistency**: Generated constants follow Go naming conventions
5. **Documentation**: Clear type definitions with descriptions from OpenAPI