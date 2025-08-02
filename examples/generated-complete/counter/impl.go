// Package counter provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification.
// WARNING: Manual edits to this file will be overwritten by subsequent code generation.
// Implement your business logic and then avoid re-running code generation on this file.
package counter

import (
	"context"
	"sync"
	"github.com/dapr/go-sdk/actor"
)

// Counter is a working implementation of CounterAPI.
// This implementation demonstrates state-based actor patterns with in-memory storage.
type Counter struct {
	actor.ServerImplBaseCtx
	mu    sync.RWMutex // Protects value from concurrent access
	value int32        // In-memory counter value
}

// Type returns the actor type for Dapr registration
func (a *Counter) Type() string {
	return ActorTypeCounter
}

// getCurrentValue retrieves the current counter value from in-memory storage
func (a *Counter) getCurrentValue(ctx context.Context) int32 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.value
}

// setValue saves the counter value to in-memory storage
func (a *Counter) setValue(ctx context.Context, value int32) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.value = value
}

// Decrement decrements counter by 1
func (a *Counter) Decrement(ctx context.Context) (*CounterState, error) {
	currentValue := a.getCurrentValue(ctx)
	newValue := currentValue - 1
	a.setValue(ctx, newValue)
	
	return &CounterState{Value: newValue}, nil
}

// Get gets current counter value
func (a *Counter) Get(ctx context.Context) (*CounterState, error) {
	value := a.getCurrentValue(ctx)
	return &CounterState{Value: value}, nil
}

// Increment increments counter by 1
func (a *Counter) Increment(ctx context.Context) (*CounterState, error) {
	currentValue := a.getCurrentValue(ctx)
	newValue := currentValue + 1
	a.setValue(ctx, newValue)
	
	return &CounterState{Value: newValue}, nil
}

// Set sets counter to specific value
func (a *Counter) Set(ctx context.Context, request SetValueRequest) (*CounterState, error) {
	a.setValue(ctx, request.Value)
	return &CounterState{Value: request.Value}, nil
}
