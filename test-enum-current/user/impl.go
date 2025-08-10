// Package user provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification.
// WARNING: Manual edits to this file will be overwritten by subsequent code generation.
// Implement your business logic and then avoid re-running code generation on this file.
package user

import (
	"context"
	"errors"
	"github.com/dapr/go-sdk/actor"
)

// User is a partial implementation of UserAPI.
// This is a stub implementation with methods that return not-implemented errors.
// You should implement the actual business logic for each method.
type User struct {
	actor.ServerImplBaseCtx
}

// Type returns the actor type for Dapr registration
func (a *User) Type() string {
	return ActorTypeUser
}


// GetUser Generated method from OpenAPI operation
// TODO: Implement the actual business logic for this method
func (a *User) GetUser(ctx context.Context) (*UserProfile, error) {
	return nil, errors.New("GetUser method is not implemented")
}

// UpdateEmail Generated method from OpenAPI operation
// TODO: Implement the actual business logic for this method
func (a *User) UpdateEmail(ctx context.Context, request UpdateEmailRequest) (*interface{}, error) {
	return nil, errors.New("UpdateEmail method is not implemented")
}

// UpdateStatus Generated method from OpenAPI operation
// TODO: Implement the actual business logic for this method
func (a *User) UpdateStatus(ctx context.Context, request UpdateStatusRequest) (*interface{}, error) {
	return nil, errors.New("UpdateStatus method is not implemented")
}
