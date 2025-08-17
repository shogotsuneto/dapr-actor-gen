// Package test provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package test

import (
	"context"
	"github.com/dapr/go-sdk/actor"
)

// ActorTypeTest is the Dapr actor type identifier for Test
const ActorTypeTest = "Test"

// TestAPI defines the interface that must be implemented to satisfy the OpenAPI schema for Test.
// This interface enforces compile-time schema compliance and includes actor.ServerContext for proper Dapr actor implementation.
type TestAPI interface {
	actor.ServerContext
	// Get state with optional references
	GetState(ctx context.Context) (*CounterState, error)
}