// Package bankaccount provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification.
// WARNING: Manual edits to this file will be overwritten by subsequent code generation.
// Implement your business logic and then avoid re-running code generation on this file.
package bankaccount

import (
	"context"
	"fmt"
	"time"
	"github.com/dapr/go-sdk/actor"
	"github.com/google/uuid"
)

const (
	stateKeyEvents = "events"
	stateKeyAccount = "account"
	
	// Event types
	EventTypeAccountCreated = "AccountCreated"
	EventTypeMoneyDeposited = "MoneyDeposited"
	EventTypeMoneyWithdrawn = "MoneyWithdrawn"
)

// AccountEvent represents a single account event for event sourcing
type AccountEvent struct {
	EventID   string                 `json:"eventId"`
	EventType string                 `json:"eventType"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// BankAccount is a working implementation of BankAccountAPI using event sourcing patterns.
// This implementation demonstrates event-sourced actor patterns with complete audit trails.
type BankAccount struct {
	actor.ServerImplBaseCtx
}

// Type returns the actor type for Dapr registration
func (a *BankAccount) Type() string {
	return ActorTypeBankAccount
}

// getEvents retrieves all events from actor state
func (a *BankAccount) getEvents(ctx context.Context) ([]AccountEvent, error) {
	var events []AccountEvent
	
	// Check if state manager is available
	stateManager := a.GetStateManager()
	if stateManager == nil {
		// If state manager is not available, return empty events
		return []AccountEvent{}, nil
	}
	
	err := stateManager.Get(ctx, stateKeyEvents, &events)
	if err != nil {
		return []AccountEvent{}, nil // Return empty if no events exist yet
	}
	
	return events, nil
}

// appendEvent adds a new event and updates the stored events
func (a *BankAccount) appendEvent(ctx context.Context, eventType string, eventData map[string]interface{}) error {
	events, err := a.getEvents(ctx)
	if err != nil {
		return err
	}
	
	event := AccountEvent{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      eventData,
	}
	
	events = append(events, event)
	
	stateManager := a.GetStateManager()
	if stateManager == nil {
		return fmt.Errorf("state manager not available")
	}
	
	if err := stateManager.Set(ctx, stateKeyEvents, events); err != nil {
		return fmt.Errorf("failed to save events state: %w", err)
	}
	
	return stateManager.Save(ctx)
}

// computeCurrentState computes the current account state from all events
func (a *BankAccount) computeCurrentState(ctx context.Context) (*BankAccountState, error) {
	events, err := a.getEvents(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(events) == 0 {
		return nil, fmt.Errorf("account not found - no events exist")
	}
	
	var state BankAccountState
	state.IsActive = true
	
	for _, event := range events {
		switch event.EventType {
		case EventTypeAccountCreated:
			state.AccountId = a.ID()
			if ownerName, ok := event.Data["ownerName"].(string); ok {
				state.OwnerName = ownerName
			}
			if initialDeposit, ok := event.Data["initialDeposit"].(float64); ok {
				state.Balance = initialDeposit
			}
			state.CreatedAt = event.Timestamp
			
		case EventTypeMoneyDeposited:
			if amount, ok := event.Data["amount"].(float64); ok {
				state.Balance += amount
			}
			
		case EventTypeMoneyWithdrawn:
			if amount, ok := event.Data["amount"].(float64); ok {
				state.Balance -= amount
			}
		}
	}
	
	return &state, nil
}

// CreateAccount creates a new bank account
func (a *BankAccount) CreateAccount(ctx context.Context, request CreateAccountRequest) (*BankAccountState, error) {
	// Check if account already exists
	events, err := a.getEvents(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(events) > 0 {
		return nil, fmt.Errorf("account already exists")
	}
	
	// Validate request
	if request.OwnerName == "" {
		return nil, fmt.Errorf("owner name is required")
	}
	if request.InitialDeposit < 0 {
		return nil, fmt.Errorf("initial deposit cannot be negative")
	}
	
	eventData := map[string]interface{}{
		"ownerName":      request.OwnerName,
		"initialDeposit": request.InitialDeposit,
	}
	
	if err := a.appendEvent(ctx, EventTypeAccountCreated, eventData); err != nil {
		return nil, err
	}
	
	return a.computeCurrentState(ctx)
}

// Deposit deposits money to account
func (a *BankAccount) Deposit(ctx context.Context, request DepositRequest) (*BankAccountState, error) {
	// Validate request
	if request.Amount <= 0 {
		return nil, fmt.Errorf("deposit amount must be positive")
	}
	
	eventData := map[string]interface{}{
		"amount":      request.Amount,
		"description": request.Description,
	}
	
	if err := a.appendEvent(ctx, EventTypeMoneyDeposited, eventData); err != nil {
		return nil, err
	}
	
	return a.computeCurrentState(ctx)
}

// GetBalance gets current account balance
func (a *BankAccount) GetBalance(ctx context.Context) (*BankAccountState, error) {
	return a.computeCurrentState(ctx)
}

// GetHistory gets transaction history
func (a *BankAccount) GetHistory(ctx context.Context) (*TransactionHistory, error) {
	events, err := a.getEvents(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert AccountEvent to interface{} for the response
	interfaceEvents := make([]interface{}, len(events))
	for i, event := range events {
		interfaceEvents[i] = event
	}
	
	return &TransactionHistory{
		AccountId: a.ID(),
		Events:    interfaceEvents,
	}, nil
}

// Withdraw withdraws money from account
func (a *BankAccount) Withdraw(ctx context.Context, request WithdrawRequest) (*BankAccountState, error) {
	// Validate request
	if request.Amount <= 0 {
		return nil, fmt.Errorf("withdraw amount must be positive")
	}
	
	// Check current balance
	currentState, err := a.computeCurrentState(ctx)
	if err != nil {
		return nil, err
	}
	
	if currentState.Balance < request.Amount {
		return nil, fmt.Errorf("insufficient funds - current balance: %.2f, requested: %.2f", 
			currentState.Balance, request.Amount)
	}
	
	eventData := map[string]interface{}{
		"amount":      request.Amount,
		"description": request.Description,
	}
	
	if err := a.appendEvent(ctx, EventTypeMoneyWithdrawn, eventData); err != nil {
		return nil, err
	}
	
	return a.computeCurrentState(ctx)
}
