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


// Decrement Decrement counter by 1
// TODO: Implement the actual business logic for this method
func (a *Counter) Decrement(ctx context.Context) (*CounterState, error) {
	return nil, errors.New("Decrement method is not implemented")
}

// Get Get current counter value
// TODO: Implement the actual business logic for this method
func (a *Counter) Get(ctx context.Context) (*CounterState, error) {
	return nil, errors.New("Get method is not implemented")
}

// Increment Increment counter by 1
// TODO: Implement the actual business logic for this method
func (a *Counter) Increment(ctx context.Context) (*CounterState, error) {
	return nil, errors.New("Increment method is not implemented")
}

// Set Set counter to specific value
// TODO: Implement the actual business logic for this method
func (a *Counter) Set(ctx context.Context, request SetValueRequest) (*CounterState, error) {
	return nil, errors.New("Set method is not implemented")
}
