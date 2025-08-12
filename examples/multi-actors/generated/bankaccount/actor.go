// Package bankaccount provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification.
// WARNING: Manual edits to this file will be overwritten by subsequent code generation.
// Implement your business logic and then avoid re-running code generation on this file.
package bankaccount

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dapr/go-sdk/actor"
	"github.com/google/uuid"
)

// BankAccount is a working implementation of BankAccountAPI using event sourcing patterns.
// This implementation demonstrates event-sourced actor patterns with in-memory storage.
type BankAccount struct {
	actor.ServerImplBaseCtx
	mu     sync.RWMutex   // Protects events from concurrent access
	events []AccountEvent // In-memory event storage
}

// Type returns the actor type for Dapr registration
func (a *BankAccount) Type() string {
	return ActorTypeBankAccount
}

// getEvents retrieves all events from in-memory storage
func (a *BankAccount) getEvents(ctx context.Context) []AccountEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	
	// Return a copy to prevent external modification
	eventsCopy := make([]AccountEvent, len(a.events))
	copy(eventsCopy, a.events)
	return eventsCopy
}

// appendEvent adds a new event to in-memory storage
func (a *BankAccount) appendEvent(ctx context.Context, eventType AccountEventEventType, eventData map[string]interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	event := AccountEvent{
		EventId:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      eventData,
	}
	
	a.events = append(a.events, event)
}

// computeCurrentState computes the current account state from all events
func (a *BankAccount) computeCurrentState(ctx context.Context) (*BankAccountState, error) {
	events := a.getEvents(ctx)
	
	if len(events) == 0 {
		return nil, fmt.Errorf("account not found - no events exist")
	}
	
	var state BankAccountState
	state.IsActive = true
	
	for _, event := range events {
		switch event.EventType {
		case AccountEventEventTypeAccountCreated:
			state.AccountId = a.ID()
			if ownerName, ok := event.Data["ownerName"].(string); ok {
				state.OwnerName = ownerName
			}
			if initialDeposit, ok := event.Data["initialDeposit"].(float64); ok {
				state.Balance = initialDeposit
			}
			state.CreatedAt = event.Timestamp
			
		case AccountEventEventTypeMoneyDeposited:
			if amount, ok := event.Data["amount"].(float64); ok {
				state.Balance += amount
			}
			
		case AccountEventEventTypeMoneyWithdrawn:
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
	events := a.getEvents(ctx)
	
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
	
	a.appendEvent(ctx, AccountEventEventTypeAccountCreated, eventData)
	
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
	
	a.appendEvent(ctx, AccountEventEventTypeMoneyDeposited, eventData)
	
	return a.computeCurrentState(ctx)
}

// GetBalance gets current account balance
func (a *BankAccount) GetBalance(ctx context.Context) (*BankAccountState, error) {
	return a.computeCurrentState(ctx)
}

// GetHistory gets transaction history
func (a *BankAccount) GetHistory(ctx context.Context) (*TransactionHistory, error) {
	events := a.getEvents(ctx)
	
	if len(events) == 0 {
		return nil, fmt.Errorf("account not found - no events exist")
	}
	
	return &TransactionHistory{
		AccountId: a.ID(),
		Events:    events,
	}, nil
}

// Withdraw withdraws money from account
func (a *BankAccount) Withdraw(ctx context.Context, request WithdrawRequest) (*BankAccountState, error) {
	// Get current state to check balance
	currentState, err := a.computeCurrentState(ctx)
	if err != nil {
		return nil, err
	}
	
	// Validate request
	if request.Amount <= 0 {
		return nil, fmt.Errorf("withdrawal amount must be positive")
	}
	if request.Amount > currentState.Balance {
		return nil, fmt.Errorf("insufficient funds: current balance %.2f, requested %.2f", 
			currentState.Balance, request.Amount)
	}
	
	eventData := map[string]interface{}{
		"amount":      request.Amount,
		"description": request.Description,
	}
	
	a.appendEvent(ctx, AccountEventEventTypeMoneyWithdrawn, eventData)
	
	return a.computeCurrentState(ctx)
}
