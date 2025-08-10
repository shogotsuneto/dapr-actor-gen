// Package user provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package user

import (
	"fmt"
	"github.com/dapr/go-sdk/actor"
)

// NewActorFactory creates a factory function for User with a cleaner API.
// Returns a factory function compatible with Dapr's RegisterActorImplFactoryContext.
// Usage: s.RegisterActorImplFactoryContext(user.NewActorFactory())
func NewActorFactory() func() actor.ServerContext {
	return func() actor.ServerContext {
		// Create a new User instance
		impl := &User{}
		
		// Compile-time check ensures the implementation satisfies the schema
		var _ UserAPI = impl
		
		// Verify the actor type matches the schema
		if impl.Type() != ActorTypeUser {
			panic(fmt.Sprintf("actor implementation Type() returns '%s', expected '%s'", impl.Type(), ActorTypeUser))
		}
		
		return impl
	}
}