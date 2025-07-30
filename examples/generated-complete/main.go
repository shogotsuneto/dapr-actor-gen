package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	daprd "github.com/dapr/go-sdk/service/http"

	"example-dapr-actors/bankaccount"
	"example-dapr-actors/counter"
)

func main() {
	// Create a Dapr service on port 8080
	s := daprd.NewService(":8080")

	// Register all generated actors

	// Register BankAccount actor
	s.RegisterActorImplFactoryContext(bankaccount.NewActorFactory())
	// Register Counter actor
	s.RegisterActorImplFactoryContext(counter.NewActorFactory())

	// Setup graceful shutdown
	go func() {
		// Wait for interrupt signal to gracefully shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		
		// The server shutdown is handled by Dapr service
		os.Exit(0)
	}()

	// Start the service
	log.Println("Starting Dapr actor service on :8080")
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start service: %v", err)
	}
}