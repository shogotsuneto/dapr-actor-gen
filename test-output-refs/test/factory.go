// Package test provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package test

import (
	"fmt"
	"github.com/dapr/go-sdk/actor"
)

// NewActorFactory creates a factory function for Test with a cleaner API.
// Returns a factory function compatible with Dapr's RegisterActorImplFactoryContext.
// Usage: s.RegisterActorImplFactoryContext(test.NewActorFactory())
func NewActorFactory() func() actor.ServerContext {
	return func() actor.ServerContext {
		// Create a new Test instance
		impl := &Test{}
		
		// Compile-time check ensures the implementation satisfies the schema
		var _ TestAPI = impl
		
		// Verify the actor type matches the schema
		if impl.Type() != ActorTypeTest {
			panic(fmt.Sprintf("actor implementation Type() returns '%s', expected '%s'", impl.Type(), ActorTypeTest))
		}
		
		return impl
	}
}