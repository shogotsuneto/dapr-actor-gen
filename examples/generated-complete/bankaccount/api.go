// Package bankaccount provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package bankaccount

import (
	"context"
	"github.com/dapr/go-sdk/actor"
)

// ActorTypeBankAccount is the Dapr actor type identifier for BankAccount
const ActorTypeBankAccount = "BankAccount"

// BankAccountAPI defines the interface that must be implemented to satisfy the OpenAPI schema for BankAccount.
// This interface enforces compile-time schema compliance and includes actor.ServerContext for proper Dapr actor implementation.
type BankAccountAPI interface {
	actor.ServerContext
	// Create new bank account with currency and initial status
	CreateAccount(ctx context.Context, request CreateAccountRequest) (*BankAccountState, error)
	// Deposit money with transaction type
	Deposit(ctx context.Context, request DepositRequest) (*BankAccountState, error)
	// Get current account balance and status
	GetBalance(ctx context.Context) (*BankAccountState, error)
	// Get transactions filtered by type
	GetTransactions(ctx context.Context) (*TransactionHistory, error)
	// Update account status
	UpdateStatus(ctx context.Context, request UpdateStatusRequest) (*BankAccountState, error)
	// Withdraw money with transaction type
	Withdraw(ctx context.Context, request WithdrawRequest) (*BankAccountState, error)
}