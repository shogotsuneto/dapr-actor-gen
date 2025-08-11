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
	// Create new bank account with initial status
	CreateAccount(ctx context.Context, request CreateAccountRequest) (*BankAccountState, error)
	// Get current account information
	GetAccount(ctx context.Context) (*BankAccountState, error)
	// Update account status
	UpdateStatus(ctx context.Context, request UpdateStatusRequest) (*BankAccountState, error)
}