# Examples

This directory contains examples showing how to use dapr-actor-gen with comprehensive enum support.

## OpenAPI Specifications

- `multi-actors/openapi.yaml` - **Enhanced OpenAPI spec** demonstrating comprehensive enum usage across Counter and BankAccount actors

## 🎯 Enum Features Demonstrated

The example showcases the full power of enum generation from OpenAPI specifications:

### **Type-Safe Enum Constants**
```go
// Generated enum types with constants
type AccountStatus string
const (
    AccountStatusActive AccountStatus = "Active"
    AccountStatusSuspended AccountStatus = "Suspended"
    AccountStatusClosed AccountStatus = "Closed"
    AccountStatusPending AccountStatus = "Pending"
    AccountStatusFrozen AccountStatus = "Frozen"
)
```

### **Runtime Validation**
```go
// Validate enum values at runtime
if !status.IsValid() {
    return fmt.Errorf("invalid status: %s", status)
}
```

### **Value Discovery**
```go
// Get all valid enum values for validation or UI
allStatuses := AllAccountStatusValues()
// Returns: [Active, Suspended, Closed, Pending, Frozen]
```

### **Business Logic Integration**
The implementation demonstrates how enums enable:
- **Status transitions** with validation (e.g., Pending → Active, but not Closed → Active)
- **Operation restrictions** based on account status
- **Priority-based processing** for counter operations
- **Transaction type validation** for banking operations
- **Currency handling** with type safety

## Generated Examples

The following directories contain generated code demonstrating different generation modes:

- `generated-complete/` - **🌟 FEATURED EXAMPLE** - Complete enum demonstration with implemented business logic
  - Full enum generation with constants, validation, and helper methods
  - Implemented counter and bank account actors showcasing enum usage
  - Enhanced main.go with enum documentation endpoint
  - Real business logic demonstrating enum-based validation and rules

## 📋 Enum Types Available

| Enum Type | Used In | Values | Purpose |
|-----------|---------|--------|---------|
| `CounterMode` | Counter Actor | Manual, Automatic, Scheduled, Triggered | Operation mode control |
| `Priority` | Counter Operations | Low, Medium, High, Critical | Priority-based processing |
| `AccountStatus` | Bank Account | Active, Suspended, Closed, Pending, Frozen | Account state management |
| `Currency` | Bank Account | USD, EUR, GBP, JPY, CAD, AUD | Multi-currency support |
| `TransactionType` | Transactions | Deposit, Withdrawal, Transfer, Fee, Interest, Refund, Adjustment | Transaction categorization |
| `EventType` | Event Sourcing | AccountCreated, StatusChanged, MoneyDeposited, MoneyWithdrawn, AccountClosed | Event categorization |

## Usage

Generate code from the enhanced enum-rich schema:

```bash
# Generate complete example with enum implementation
./bin/dapr-actor-gen --generate-impl --generate-example examples/multi-actors/openapi.yaml ./output

# Test the generated code
cd ./output
go mod tidy
go build ./...

# Run the demo service (shows enum documentation at /enum-demo)
go run main.go
```

## 🌐 Demo Endpoints

The generated example includes a live demo endpoint:

- **GET /enum-demo** - Interactive enum documentation with examples
- **Dapr Actor Endpoints** - All standard Dapr actor methods with enum parameters

Example requests:
```bash
# Configure counter with enums
curl -X POST http://localhost:3500/v1.0/actors/Counter/counter1/method/Configure \
  -d '{"mode": "Automatic", "priority": "High", "isEnabled": true}'

# Create bank account with currency and status
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account1/method/CreateAccount \
  -d '{"ownerName": "John Doe", "currency": "USD", "initialDeposit": 100, "initialStatus": "Pending"}'

# Make deposit with transaction type
curl -X POST http://localhost:3500/v1.0/actors/BankAccount/account1/method/Deposit \
  -d '{"amount": 250, "transactionType": "Deposit", "description": "Salary deposit"}'
```

## Key Benefits

1. **🔒 Type Safety**: Enum constants provide compile-time checking for valid values
2. **✅ Runtime Validation**: `IsValid()` methods for validating enum values from external sources
3. **🔍 Value Discovery**: `AllEnumValues()` functions for building UIs and validation
4. **📝 Go Conventions**: Generated constants follow standard Go naming patterns (`EnumTypeValue`)
5. **🏗️ Business Logic**: Enums integrate seamlessly with business rules and validation
6. **📚 Documentation**: Preserves OpenAPI descriptions in generated Go code
7. **🔄 State Machines**: Perfect for modeling state transitions and workflows
8. **🌍 Internationalization**: String-based enums work well with localization