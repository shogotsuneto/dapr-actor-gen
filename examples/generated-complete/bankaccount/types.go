// Package bankaccount provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package bankaccount


// BankAccountState Current state of bank account with comprehensive enum usage
type BankAccountState struct {
	// Unique account identifier
	AccountId string `json:"accountId"`
	// Current account balance
	Balance float64 `json:"balance"`
	// Account creation timestamp
	CreatedAt string `json:"createdAt,omitempty"`
	// Supported currencies for bank accounts
	Currency Currency `json:"currency"`
	// Last transaction timestamp
	LastActivity string `json:"lastActivity,omitempty"`
	// Additional account metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Account owner name
	OwnerName string `json:"ownerName"`
	// Current status of the bank account
	Status AccountStatus `json:"status"`
	// Total number of transactions
	TransactionCount int `json:"transactionCount,omitempty"`
}

// CreateAccountRequest Request to create a new bank account with currency and status
type CreateAccountRequest struct {
	// Supported currencies for bank accounts
	Currency Currency `json:"currency"`
	// Initial deposit amount in the specified currency
	InitialDeposit float64 `json:"initialDeposit"`
	// Current status of the bank account
	InitialStatus AccountStatus `json:"initialStatus,omitempty"`
	// Optional account metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// Name of the account owner
	OwnerName string `json:"ownerName"`
}

// DepositRequest Request to deposit money with transaction type
type DepositRequest struct {
	// Amount to deposit in account currency
	Amount float64 `json:"amount"`
	// Description of the deposit
	Description string `json:"description"`
	// Additional transaction metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// External reference ID for the transaction
	ReferenceId string `json:"referenceId,omitempty"`
	// Type of financial transaction
	TransactionType TransactionType `json:"transactionType"`
}

// Transaction A single financial transaction with enum-based categorization
type Transaction struct {
	// Transaction amount (positive for credits, negative for debits)
	Amount float64 `json:"amount"`
	// Account balance after this transaction
	BalanceAfter float64 `json:"balanceAfter,omitempty"`
	// Transaction description
	Description string `json:"description"`
	// Additional transaction metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// External reference ID
	ReferenceId string `json:"referenceId,omitempty"`
	// When the transaction occurred
	Timestamp string `json:"timestamp"`
	// Unique transaction identifier
	TransactionId string `json:"transactionId"`
	// Type of financial transaction
	Type TransactionType `json:"type"`
}

// TransactionHistory Transaction history with enum-based filtering capabilities
type TransactionHistory struct {
	// Account identifier
	AccountId string `json:"accountId"`
	// Supported currencies for bank accounts
	Currency Currency `json:"currency,omitempty"`
	// Type of financial transaction
	FilteredBy TransactionType `json:"filteredBy,omitempty"`
	// Summary statistics by transaction type
	Summary interface{} `json:"summary,omitempty"`
	// Total number of transactions (may be more than returned)
	TotalCount int `json:"totalCount"`
	// List of transactions in chronological order
	Transactions []Transaction `json:"transactions"`
}

// UpdateStatusRequest Request to update account status
type UpdateStatusRequest struct {
	// When the status change takes effect (defaults to now)
	EffectiveDate string `json:"effectiveDate,omitempty"`
	// Reason for the status change
	Reason string `json:"reason"`
	// Current status of the bank account
	Status AccountStatus `json:"status"`
}

// WithdrawRequest Request to withdraw money with transaction type
type WithdrawRequest struct {
	// Amount to withdraw in account currency
	Amount float64 `json:"amount"`
	// Description of the withdrawal
	Description string `json:"description"`
	// Additional transaction metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// External reference ID for the transaction
	ReferenceId string `json:"referenceId,omitempty"`
	// Type of financial transaction
	TransactionType TransactionType `json:"transactionType"`
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

// Currency Supported currencies for bank accounts
type Currency string

// Currency constants
const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyJPY Currency = "JPY"
	CurrencyCAD Currency = "CAD"
	CurrencyAUD Currency = "AUD"
)

// IsValid returns true if the Currency value is valid
func (e Currency) IsValid() bool {
	switch e {
	case CurrencyUSD:
		return true
	case CurrencyEUR:
		return true
	case CurrencyGBP:
		return true
	case CurrencyJPY:
		return true
	case CurrencyCAD:
		return true
	case CurrencyAUD:
		return true
	default:
		return false
	}
}

// String returns the string representation of Currency
func (e Currency) String() string {
	return string(e)
}

// AllCurrencyValues returns all valid Currency values
func AllCurrencyValues() []Currency {
	return []Currency{
		CurrencyUSD,
		CurrencyEUR,
		CurrencyGBP,
		CurrencyJPY,
		CurrencyCAD,
		CurrencyAUD,
	}
}

// TransactionType Type of financial transaction
type TransactionType string

// TransactionType constants
const (
	TransactionTypeDeposit TransactionType = "Deposit"
	TransactionTypeWithdrawal TransactionType = "Withdrawal"
	TransactionTypeTransfer TransactionType = "Transfer"
	TransactionTypeFee TransactionType = "Fee"
	TransactionTypeInterest TransactionType = "Interest"
	TransactionTypeRefund TransactionType = "Refund"
	TransactionTypeAdjustment TransactionType = "Adjustment"
)

// IsValid returns true if the TransactionType value is valid
func (e TransactionType) IsValid() bool {
	switch e {
	case TransactionTypeDeposit:
		return true
	case TransactionTypeWithdrawal:
		return true
	case TransactionTypeTransfer:
		return true
	case TransactionTypeFee:
		return true
	case TransactionTypeInterest:
		return true
	case TransactionTypeRefund:
		return true
	case TransactionTypeAdjustment:
		return true
	default:
		return false
	}
}

// String returns the string representation of TransactionType
func (e TransactionType) String() string {
	return string(e)
}

// AllTransactionTypeValues returns all valid TransactionType values
func AllTransactionTypeValues() []TransactionType {
	return []TransactionType{
		TransactionTypeDeposit,
		TransactionTypeWithdrawal,
		TransactionTypeTransfer,
		TransactionTypeFee,
		TransactionTypeInterest,
		TransactionTypeRefund,
		TransactionTypeAdjustment,
	}
}
