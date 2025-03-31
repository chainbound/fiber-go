package client

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

// TestConcurrentConnections tests for race conditions when creating multiple connections
func TestConcurrentConnections(t *testing.T) {
	// Skip unless explicitly enabled
	if os.Getenv("RUN_RECONNECTION_TEST") != "true" {
		t.Skip("Skipping reconnection test. Set RUN_RECONNECTION_TEST=true to run")
	}

	apiKey := os.Getenv("FIBER_API_KEY")
	if apiKey == "" {
		t.Skip("FIBER_API_KEY environment variable not set")
	}

	target := "beta.fiberapi.io:8080"

	var wg sync.WaitGroup
	numClients := 5

	// Create multiple clients concurrently
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			client := NewClient(target, apiKey)
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := client.Connect(ctx); err != nil {
				t.Errorf("Client %d failed to connect: %v", clientID, err)
			}
		}(i)
	}

	wg.Wait()
}
