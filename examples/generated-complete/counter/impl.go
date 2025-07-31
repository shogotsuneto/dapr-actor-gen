// Package counter provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification.
// WARNING: Manual edits to this file will be overwritten by subsequent code generation.
// Implement your business logic and then avoid re-running code generation on this file.
package counter

import (
	"context"
	"fmt"
	"github.com/dapr/go-sdk/actor"
)

const stateKeyValue = "value"

// Counter is a working implementation of CounterAPI.
// This implementation demonstrates state-based actor patterns with Dapr state management.
type Counter struct {
	actor.ServerImplBaseCtx
}

// Type returns the actor type for Dapr registration
func (a *Counter) Type() string {
	return ActorTypeCounter
}

// getCurrentValue retrieves the current counter value from actor state
func (a *Counter) getCurrentValue(ctx context.Context) (int32, error) {
	var value int32
	err := a.GetStateManager().Get(ctx, stateKeyValue, &value)
	if err != nil {
		return 0, fmt.Errorf("failed to get counter state: %w", err)
	}
	
	return value, nil
}

// saveValue saves the counter value to actor state
func (a *Counter) saveValue(ctx context.Context, value int32) error {
	if err := a.GetStateManager().Set(ctx, stateKeyValue, value); err != nil {
		return fmt.Errorf("failed to save counter state: %w", err)
	}
	
	return a.GetStateManager().Save(ctx)
}

// Decrement decrements counter by 1
func (a *Counter) Decrement(ctx context.Context) (*CounterState, error) {
	currentValue, err := a.getCurrentValue(ctx)
	if err != nil {
		return nil, err
	}
	
	newValue := currentValue - 1
	if err := a.saveValue(ctx, newValue); err != nil {
		return nil, err
	}
	
	return &CounterState{Value: newValue}, nil
}

// Get gets current counter value
func (a *Counter) Get(ctx context.Context) (*CounterState, error) {
	value, err := a.getCurrentValue(ctx)
	if err != nil {
		return nil, err
	}
	
	return &CounterState{Value: value}, nil
}

// Increment increments counter by 1
func (a *Counter) Increment(ctx context.Context) (*CounterState, error) {
	currentValue, err := a.getCurrentValue(ctx)
	if err != nil {
		return nil, err
	}
	
	newValue := currentValue + 1
	if err := a.saveValue(ctx, newValue); err != nil {
		return nil, err
	}
	
	return &CounterState{Value: newValue}, nil
}

// Set sets counter to specific value
func (a *Counter) Set(ctx context.Context, request SetValueRequest) (*CounterState, error) {
	if err := a.saveValue(ctx, request.Value); err != nil {
		return nil, err
	}
	
	return &CounterState{Value: request.Value}, nil
}
