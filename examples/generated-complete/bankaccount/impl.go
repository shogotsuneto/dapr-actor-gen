// Package bankaccount provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification.
// WARNING: Manual edits to this file will be overwritten by subsequent code generation.
// Implement your business logic and then avoid re-running code generation on this file.
package bankaccount

import (
	"context"
	"errors"
	"github.com/dapr/go-sdk/actor"
)

// BankAccount is a partial implementation of BankAccountAPI.
// This is a stub implementation with methods that return not-implemented errors.
// You should implement the actual business logic for each method.
type BankAccount struct {
	actor.ServerImplBaseCtx
}

// Type returns the actor type for Dapr registration
func (a *BankAccount) Type() string {
	return ActorTypeBankAccount
}


// CreateAccount Create new bank account with initial status
// TODO: Implement the actual business logic for this method
func (a *BankAccount) CreateAccount(ctx context.Context, request CreateAccountRequest) (*BankAccountState, error) {
	return nil, errors.New("CreateAccount method is not implemented")
}

// GetAccount Get current account information
// TODO: Implement the actual business logic for this method
func (a *BankAccount) GetAccount(ctx context.Context) (*BankAccountState, error) {
	return nil, errors.New("GetAccount method is not implemented")
}

// UpdateStatus Update account status
// TODO: Implement the actual business logic for this method
func (a *BankAccount) UpdateStatus(ctx context.Context, request UpdateStatusRequest) (*BankAccountState, error) {
	return nil, errors.New("UpdateStatus method is not implemented")
}
