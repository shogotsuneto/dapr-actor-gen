// Package counter provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package counter

import (
	"context"
	"github.com/dapr/go-sdk/actor"
)

// ActorTypeCounter is the Dapr actor type identifier for Counter
const ActorTypeCounter = "Counter"

// CounterAPI defines the interface that must be implemented to satisfy the OpenAPI schema for Counter.
// This interface enforces compile-time schema compliance and includes actor.ServerContext for proper Dapr actor implementation.
type CounterAPI interface {
	actor.ServerContext
	// Configure counter mode
	Configure(ctx context.Context, request ConfigureCounterRequest) (*CounterState, error)
	// Decrement the counter
	Decrement(ctx context.Context, request DecrementRequest) (*CounterState, error)
	// Increment the counter
	Increment(ctx context.Context, request IncrementRequest) (*CounterState, error)
	// Set counter to specific value
	SetValue(ctx context.Context, request SetValueRequest) (*CounterState, error)
}