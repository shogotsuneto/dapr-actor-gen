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


// CreateAccount Create new bank account
// TODO: Implement the actual business logic for this method
func (a *BankAccount) CreateAccount(ctx context.Context, request CreateAccountRequest) (*BankAccountState, error) {
	return nil, errors.New("CreateAccount method is not implemented")
}

// Deposit Deposit money to account
// TODO: Implement the actual business logic for this method
func (a *BankAccount) Deposit(ctx context.Context, request DepositRequest) (*BankAccountState, error) {
	return nil, errors.New("Deposit method is not implemented")
}

// GetBalance Get current account balance
// TODO: Implement the actual business logic for this method
func (a *BankAccount) GetBalance(ctx context.Context) (*BankAccountState, error) {
	return nil, errors.New("GetBalance method is not implemented")
}

// GetHistory Get transaction history
// TODO: Implement the actual business logic for this method
func (a *BankAccount) GetHistory(ctx context.Context) (*TransactionHistory, error) {
	return nil, errors.New("GetHistory method is not implemented")
}

// Withdraw Withdraw money from account
// TODO: Implement the actual business logic for this method
func (a *BankAccount) Withdraw(ctx context.Context, request WithdrawRequest) (*BankAccountState, error) {
	return nil, errors.New("Withdraw method is not implemented")
}
