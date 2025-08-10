// Package user provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package user

import (
	"context"
	"github.com/dapr/go-sdk/actor"
)

// ActorTypeUser is the Dapr actor type identifier for User
const ActorTypeUser = "User"

// UserAPI defines the interface that must be implemented to satisfy the OpenAPI schema for User.
// This interface enforces compile-time schema compliance and includes actor.ServerContext for proper Dapr actor implementation.
type UserAPI interface {
	actor.ServerContext
	// Generated method from OpenAPI operation
	GetUser(ctx context.Context) (*UserProfile, error)
	// Generated method from OpenAPI operation
	UpdateEmail(ctx context.Context, request UpdateEmailRequest) (*interface{}, error)
	// Generated method from OpenAPI operation
	UpdateStatus(ctx context.Context, request UpdateStatusRequest) (*interface{}, error)
}