// Package bankaccount provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package bankaccount


// BankAccountState Current state of the bank account actor with enum-based status
type BankAccountState struct {
	// Unique identifier for the bank account
	AccountId string `json:"accountId"`
	// Current account balance
	Balance float64 `json:"balance"`
	// Whether the account is currently enabled for transactions
	IsEnabled bool `json:"isEnabled,omitempty"`
	// When the account was last modified
	LastModified string `json:"lastModified,omitempty"`
	// Additional account metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Current status of the bank account
	Status AccountStatus `json:"status"`
}

// CreateAccountRequest Request to create a new bank account
type CreateAccountRequest struct {
	// Unique identifier for the new account
	AccountId string `json:"accountId"`
	// Initial balance for the account
	InitialBalance float64 `json:"initialBalance"`
	// Optional metadata for the account
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Current status of the bank account
	Status AccountStatus `json:"status,omitempty"`
}

// UpdateStatusRequest Request to update account status
type UpdateStatusRequest struct {
	// Optional metadata for the status change
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Reason for the status change
	Reason string `json:"reason,omitempty"`
	// Current status of the bank account
	Status AccountStatus `json:"status"`
}





// AccountStatus Current status of the bank account
type AccountStatus string

// AccountStatus constants
const (
	AccountStatusActive AccountStatus = "Active"
	AccountStatusSuspended AccountStatus = "Suspended"
	AccountStatusClosed AccountStatus = "Closed"
	AccountStatusPending AccountStatus = "Pending"
	AccountStatusFrozen AccountStatus = "Frozen"
)

// IsValid returns true if the AccountStatus value is valid
func (e AccountStatus) IsValid() bool {
	switch e {
	case AccountStatusActive:
		return true
	case AccountStatusSuspended:
		return true
	case AccountStatusClosed:
		return true
	case AccountStatusPending:
		return true
	case AccountStatusFrozen:
		return true
	default:
		return false
	}
}

// String returns the string representation of AccountStatus
func (e AccountStatus) String() string {
	return string(e)
}

// AllAccountStatusValues returns all valid AccountStatus values
func AllAccountStatusValues() []AccountStatus {
	return []AccountStatus{
		AccountStatusActive,
		AccountStatusSuspended,
		AccountStatusClosed,
		AccountStatusPending,
		AccountStatusFrozen,
	}
}
