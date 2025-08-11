// Package counter provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification.
// WARNING: Manual edits to this file will be overwritten by subsequent code generation.
// Implement your business logic and then avoid re-running code generation on this file.
package counter

import (
	"context"
	"errors"
	"github.com/dapr/go-sdk/actor"
)

// Counter is a partial implementation of CounterAPI.
// This is a stub implementation with methods that return not-implemented errors.
// You should implement the actual business logic for each method.
type Counter struct {
	actor.ServerImplBaseCtx
}

// Type returns the actor type for Dapr registration
func (a *Counter) Type() string {
	return ActorTypeCounter
}


// Configure Configure counter mode
// TODO: Implement the actual business logic for this method
func (a *Counter) Configure(ctx context.Context, request ConfigureCounterRequest) (*CounterState, error) {
	return nil, errors.New("Configure method is not implemented")
}

// Decrement Decrement the counter
// TODO: Implement the actual business logic for this method
func (a *Counter) Decrement(ctx context.Context, request DecrementRequest) (*CounterState, error) {
	return nil, errors.New("Decrement method is not implemented")
}

// Increment Increment the counter
// TODO: Implement the actual business logic for this method
func (a *Counter) Increment(ctx context.Context, request IncrementRequest) (*CounterState, error) {
	return nil, errors.New("Increment method is not implemented")
}

// SetValue Set counter to specific value
// TODO: Implement the actual business logic for this method
func (a *Counter) SetValue(ctx context.Context, request SetValueRequest) (*CounterState, error) {
	return nil, errors.New("SetValue method is not implemented")
}
