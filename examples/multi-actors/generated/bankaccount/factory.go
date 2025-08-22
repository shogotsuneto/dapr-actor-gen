// Package bankaccount provides primitives for OpenAPI-based schema validation.
//
// Factory functions generated from OpenAPI specification as examples.
// These can be customized for dependency injection or other initialization needs.
package bankaccount

import (
	"fmt"
	"github.com/dapr/go-sdk/actor"
)

// NewActorFactory creates a factory function for BankAccount.
// This is a generated example that can be customized for dependency injection.
// Returns a factory function compatible with Dapr's RegisterActorImplFactoryContext.
// Usage: s.RegisterActorImplFactoryContext(bankaccount.NewActorFactory())
func NewActorFactory() func() actor.ServerContext {
	return func() actor.ServerContext {
		// Create a new BankAccount instance
		impl := &BankAccount{}
		
		// Compile-time check ensures the implementation satisfies the schema
		var _ BankAccountAPI = impl
		
		// Verify the actor type matches the schema
		if impl.Type() != ActorTypeBankAccount {
			panic(fmt.Sprintf("actor implementation Type() returns '%s', expected '%s'", impl.Type(), ActorTypeBankAccount))
		}
		
		return impl
	}
}