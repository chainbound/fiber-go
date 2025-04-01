package client

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// TestReconnection verifies that the client automatically reconnects after a disconnect
func TestReconnection(t *testing.T) {
	// Skip unless explicitly enabled
	if os.Getenv("RUN_RECONNECTION_TEST") != "true" {
		t.Skip("Skipping reconnection test. Set RUN_RECONNECTION_TEST=true to run")
	}

	apiKey := os.Getenv("FIBER_API_KEY")
	if apiKey == "" {
		t.Skip("FIBER_API_KEY environment variable not set")
	}

	target := "beta.fiberapi.io:8080"

	// Simplified to only test reconnection behavior
	testReconnection(t, target, apiKey)
}

// testReconnection is a helper function that tests reconnection behavior
func testReconnection(t *testing.T, target, apiKey string) {
	t.Logf("========== CONNECTION SETUP ==========")

	// Create configuration
	config := NewConfig().SetIdleTimeout(10 * time.Second).SetHealthCheckInterval(10 * time.Second).SetLogLevel("debug")

	// Create client
	fiber := NewClientWithConfig(target, apiKey, config)
	defer func() {
		if err := fiber.Close(); err != nil {
			t.Logf("Error closing fiber client: %v", err)
		}
	}()

	// Connect to the API
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := fiber.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Create a channel for transactions
	txs := make(chan *TransactionWithSender)
	subDone := make(chan error, 1)

	// Subscribe to new transactions
	go func() {
		err := fiber.SubscribeNewTxs(nil, txs)
		subDone <- err
		close(subDone)
	}()

	// Wait for initial data to confirm subscription is working
	t.Log("Waiting for initial transactions before testing reconnection...")
	initialRxCount := 0
	initialWaitTimeout := time.After(15 * time.Second)

initialWait:
	for {
		select {
		case <-initialWaitTimeout:
			t.Fatal("Did not receive any initial transactions within timeout")
		case tx, ok := <-txs:
			if !ok {
				t.Fatal("Transaction channel closed during initial wait")
			}
			initialRxCount++
			t.Logf("Received initial tx %d: %s", initialRxCount, tx.Transaction.Hash())
			if initialRxCount >= 1 {
				t.Log("Received sufficient initial transactions")
				break initialWait
			}
		case err := <-subDone:
			t.Fatalf("Subscription ended during initial wait: %v", err)
		}
	}

	// Record state before we simulate a disconnection
	t.Logf("Current connection state before disconnect: %v", fiber.conn.GetState())

	// Simulate a network interruption by forcibly closing the connection
	t.Log("Simulating network interruption...")
	disconnectTime := time.Now()
	if err := fiber.conn.Close(); err != nil {
		t.Logf("Error forcibly closing connection: %v", err)
	}

	// Now wait for transactions to appear after reconnection
	t.Log("Waiting for reconnection...")
	txAfterDisconnect := atomic.Int32{}
	reconnected := false
	var reconnectTime time.Time
	waiting := true

	reconnectTimeout := time.After(30 * time.Second)

	for waiting {
		select {
		case <-reconnectTimeout:
			t.Fatal("Did not reconnect within the expected timeout period")
		case tx, ok := <-txs:
			if !ok {
				t.Log("Transaction channel was closed")
				waiting = false
				break
			}

			// Only count transactions received after a reasonable delay
			timeSinceDisconnect := time.Since(disconnectTime)
			if timeSinceDisconnect > 3*time.Second {
				if !reconnected {
					reconnected = true
					reconnectTime = time.Now()
					t.Logf("RECONNECTION DETECTED after %v - Transaction received after disconnect",
						reconnectTime.Sub(disconnectTime))
					t.Logf("Current connection state after reconnect: %v", fiber.conn.GetState())
				}

				// Log transaction hash
				t.Logf("Post-reconnect tx: %s", tx.Transaction.Hash())
				txAfterDisconnect.Add(1)

				// If we've verified reconnection with multiple transactions, we can exit
				if txAfterDisconnect.Load() >= 3 {
					t.Logf("Successfully confirmed reconnection with %d transactions",
						txAfterDisconnect.Load())
					waiting = false
				}
			}
		case <-time.After(500 * time.Millisecond):
			// Poll interval
		case err := <-subDone:
			if err != nil {
				t.Fatalf("Subscription ended with error: %v", err)
			}
			waiting = false
		}
	}

	// Verify results
	if !reconnected {
		t.Errorf("Expected client to reconnect, but it did not")
	} else {
		t.Logf("SUCCESS: Client reconnected after %v", reconnectTime.Sub(disconnectTime))
	}
}
