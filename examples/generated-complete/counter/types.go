// Package counter provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package counter


// ConfigureCounterRequest Request to configure counter mode
type ConfigureCounterRequest struct {
	// Enable or disable the counter
	IsEnabled bool `json:"isEnabled,omitempty"`
	// Optional metadata for the counter
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Operating mode for the counter actor
	Mode CounterMode `json:"mode"`
}

// CounterState Current state of the counter actor with enum-based configuration
type CounterState struct {
	// Whether the counter is currently enabled
	IsEnabled bool `json:"isEnabled"`
	// When the counter was last modified
	LastModified string `json:"lastModified,omitempty"`
	// Additional counter metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Operating mode for the counter actor
	Mode CounterMode `json:"mode"`
	// The current counter value
	Value int32 `json:"value"`
}

// DecrementRequest Request to decrement counter
type DecrementRequest struct {
	// Reason for the decrement operation
	Reason string `json:"reason,omitempty"`
	// How much to decrement by
	Step int32 `json:"step,omitempty"`
}

// IncrementRequest Request to increment counter
type IncrementRequest struct {
	// Reason for the increment operation
	Reason string `json:"reason,omitempty"`
	// How much to increment by
	Step int32 `json:"step,omitempty"`
}

// SetValueRequest Request to set counter to specific value
type SetValueRequest struct {
	// Reason for setting the value
	Reason string `json:"reason,omitempty"`
	// The value to set the counter to
	Value int32 `json:"value"`
}





// CounterMode Operating mode for the counter actor
type CounterMode string

// CounterMode constants
const (
	CounterModeManual CounterMode = "Manual"
	CounterModeAutomatic CounterMode = "Automatic"
	CounterModeScheduled CounterMode = "Scheduled"
	CounterModeTriggered CounterMode = "Triggered"
)

// IsValid returns true if the CounterMode value is valid
func (e CounterMode) IsValid() bool {
	switch e {
	case CounterModeManual:
		return true
	case CounterModeAutomatic:
		return true
	case CounterModeScheduled:
		return true
	case CounterModeTriggered:
		return true
	default:
		return false
	}
}

// String returns the string representation of CounterMode
func (e CounterMode) String() string {
	return string(e)
}

// AllCounterModeValues returns all valid CounterMode values
func AllCounterModeValues() []CounterMode {
	return []CounterMode{
		CounterModeManual,
		CounterModeAutomatic,
		CounterModeScheduled,
		CounterModeTriggered,
	}
}
