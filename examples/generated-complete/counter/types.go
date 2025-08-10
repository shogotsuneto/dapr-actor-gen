// Package counter provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package counter


// ConfigureCounterRequest Request to configure counter mode and priority
type ConfigureCounterRequest struct {
	// Enable or disable the counter
	IsEnabled bool `json:"isEnabled,omitempty"`
	// Optional metadata for the counter
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Operating mode for the counter actor
	Mode CounterMode `json:"mode"`
	// Priority level for counter operations
	Priority Priority `json:"priority"`
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
	// Priority level for counter operations
	Priority Priority `json:"priority"`
	// The current counter value
	Value int32 `json:"value"`
}

// DecrementRequest Request to decrement counter with priority
type DecrementRequest struct {
	// Priority level for counter operations
	Priority Priority `json:"priority"`
	// Reason for the decrement operation
	Reason string `json:"reason,omitempty"`
	// How much to decrement by
	Step int32 `json:"step,omitempty"`
}

// IncrementRequest Request to increment counter with priority
type IncrementRequest struct {
	// Priority level for counter operations
	Priority Priority `json:"priority"`
	// Reason for the increment operation
	Reason string `json:"reason,omitempty"`
	// How much to increment by
	Step int32 `json:"step,omitempty"`
}

// SetValueRequest Request to set counter to specific value with priority
type SetValueRequest struct {
	// Priority level for counter operations
	Priority Priority `json:"priority"`
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

// Priority Priority level for counter operations
type Priority string

// Priority constants
const (
	PriorityLow Priority = "Low"
	PriorityMedium Priority = "Medium"
	PriorityHigh Priority = "High"
	PriorityCritical Priority = "Critical"
)

// IsValid returns true if the Priority value is valid
func (e Priority) IsValid() bool {
	switch e {
	case PriorityLow:
		return true
	case PriorityMedium:
		return true
	case PriorityHigh:
		return true
	case PriorityCritical:
		return true
	default:
		return false
	}
}

// String returns the string representation of Priority
func (e Priority) String() string {
	return string(e)
}

// AllPriorityValues returns all valid Priority values
func AllPriorityValues() []Priority {
	return []Priority{
		PriorityLow,
		PriorityMedium,
		PriorityHigh,
		PriorityCritical,
	}
}
