// Package counter provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package counter


// CounterState Current state of the counter actor (state-based)
type CounterState struct {
	// Type of the last operation performed on the counter
	LastOperation CounterOperation `json:"lastOperation,omitempty"`
	// Current status of the counter
	Status CounterStatus `json:"status"`
	// The current counter value
	Value int32 `json:"value"`
}

// SetValueRequest Request to set the counter to a specific value
type SetValueRequest struct {
	// The value to set the counter to
	Value int32 `json:"value"`
}





// CounterOperation Type of the last operation performed on the counter
type CounterOperation string

// CounterOperation constants
const (
	CounterOperationincrement CounterOperation = "increment"
	CounterOperationdecrement CounterOperation = "decrement"
	CounterOperationset CounterOperation = "set"
	CounterOperationget CounterOperation = "get"
	CounterOperationreset CounterOperation = "reset"
)

// CounterStatus Current status of the counter
type CounterStatus string

// CounterStatus constants
const (
	CounterStatusactive CounterStatus = "active"
	CounterStatuspaused CounterStatus = "paused"
	CounterStatuserror CounterStatus = "error"
	CounterStatusreset CounterStatus = "reset"
)
