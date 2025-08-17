// Package test provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package test


// CounterState Current state of the counter actor (state-based)
type CounterState struct {
	// Data payload for counter state
	Data CounterStateData `json:"data,omitempty"`
	// Error information
	Error Error `json:"error,omitempty"`
	// Whether the operation was successful
	Success bool `json:"success"`
}

// CounterStateData Data payload for counter state
type CounterStateData struct {
	// When the value was set
	Timestamp string `json:"timestamp,omitempty"`
	// The counter value
	Value int `json:"value,omitempty"`
}

// Error Error information
type Error struct {
	// Error code
	Code string `json:"code,omitempty"`
	// Error message
	Message string `json:"message,omitempty"`
}




