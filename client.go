// package client is responsible for easily interacting with the protobuf API in Go.
// It contains wrappers for all the rpc methods that accept standard go-ethereum
// objects.
package client

import (
	"context"
	"time"

	"github.com/chainbound/fiber-go/protobuf/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// Client version string that is appended to each stream request.
const Version string = "fiber-go/1.9.2"

type Client struct {
	target string
	conn   *grpc.ClientConn
	config *ClientConfig
	client api.APIClient
	key    string

	txStream          api.API_SendTransactionV2Client
	txSeqStream       api.API_SendTransactionSequenceV2Client
	submitBlockStream api.API_SubmitBlockStreamClient
}

type ClientConfig struct {
	enableCompression bool
	writeBufferSize   int
	readBufferSize    int
	connWindowSize    int32
	windowSize        int32
	idleTimeout       time.Duration
}

// NewConfig creates a new config with sensible default values.
func NewConfig() *ClientConfig {
	return &ClientConfig{
		enableCompression: false,
		writeBufferSize:   1024 * 8,
		readBufferSize:    1024 * 8,
		connWindowSize:    1024 * 512,
		windowSize:        1024 * 256,
		idleTimeout:       0,
	}
}

func (c *ClientConfig) EnableCompression() *ClientConfig {
	c.enableCompression = true
	return c
}

func (c *ClientConfig) SetWriteBufferSize(size int) *ClientConfig {
	c.writeBufferSize = size
	return c
}

func (c *ClientConfig) SetReadBufferSize(size int) *ClientConfig {
	c.readBufferSize = size
	return c
}

func (c *ClientConfig) SetConnWindowSize(size int32) *ClientConfig {
	c.connWindowSize = size
	return c
}

func (c *ClientConfig) SetWindowSize(size int32) *ClientConfig {
	c.windowSize = size
	return c
}

// SetIdleTimeout sets the client idle timeout. This is the duration after which
// idle connections will be restarted. Setting to 0 disables the idle timeout.
func (c *ClientConfig) SetIdleTimeout(timeout time.Duration) *ClientConfig {
	c.idleTimeout = timeout
	return c
}

func NewClient(target, apiKey string) *Client {
	return NewClientWithConfig(target, apiKey, NewConfig())
}

// NewClientWithConfig creates a new client with the given config.
func NewClientWithConfig(target, apiKey string, config *ClientConfig) *Client {
	return &Client{
		config: config,
		target: target,
		key:    apiKey,
	}
}

// Connects sets up the gRPC channel and creates the stub. It blocks until connected or the given context expires.
// Always use a context with timeout.
func (c *Client) Connect(ctx context.Context) error {
	serviceConfig := `{
		"methodConfig": [{
			"name": [{"service": "API"}],
			"retryPolicy": {
				"MaxAttempts": 20,
				"InitialBackoff": "2s",
				"MaxBackoff": "60s",
				"BackoffMultiplier": 1.2,
				"RetryableStatusCodes": ["UNAVAILABLE", "ABORTED"]
			}
		}]
	}`

	if c.config.enableCompression {
		registerGzipCompression()
	}

	// Setup options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithWriteBufferSize(c.config.writeBufferSize),
		grpc.WithReadBufferSize(c.config.readBufferSize),
		grpc.WithInitialConnWindowSize(c.config.connWindowSize),
		grpc.WithInitialWindowSize(c.config.windowSize),
	}

	// Add keepalive parameters if idle timeout is set
	if c.config.idleTimeout > 0 {
		kaParams := keepalive.ClientParameters{
			// If no activity on the connection after this period, send a keepalive ping
			Time: c.config.idleTimeout,
			// Wait time for a keepalive ping response before closing the connection
			Timeout: 20 * time.Second,
			// Allow keepalive pings even when there are no active streams
			PermitWithoutStream: true,
		}
		opts = append(opts, grpc.WithKeepaliveParams(kaParams))
	}

	conn, err := grpc.DialContext(ctx, c.target, opts...)
	if err != nil {
		return err
	}

	c.conn = conn

	// Create the stub (client) with the channel
	c.client = api.NewAPIClient(conn)

	ctx = metadata.AppendToOutgoingContext(context.Background(), "x-api-key", c.key, "x-client-version", Version)
	c.txStream, err = c.client.SendTransactionV2(ctx)
	if err != nil {
		return err
	}

	c.txSeqStream, err = c.client.SendTransactionSequenceV2(ctx)
	if err != nil {
		return err
	}

	c.submitBlockStream, err = c.client.SubmitBlockStream(ctx)
	if err != nil {
		return err
	}

	// Start connection monitor in the background
	// This will help detect silent disconnects and trigger reconnection attempts
	go c.monitorConnection(2 * time.Second)

	return nil
}

// monitorConnection periodically checks the gRPC connection state
// and attempts to reconnect if the connection appears to be dead
func (c *Client) monitorConnection(checkInterval time.Duration) {
	// Use default interval if not specified
	if checkInterval <= 0 {
		checkInterval = 10 * time.Second
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		// Check connection state
		state := c.conn.GetState()
		if state != connectivity.Ready && state != connectivity.Connecting && state != connectivity.Idle {
			// If connection is in a bad state, attempt to reset it
			c.conn.ResetConnectBackoff()

			// Create a context with timeout for reconnection
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			// Try to force state change
			if !c.conn.WaitForStateChange(ctx, state) {
				// If state doesn't change, connection might be dead
				// We should close it to trigger reconnection in the stream methods
				c.conn.Close()

				// Attempt to create a new connection
				newCtx, newCancel := context.WithTimeout(context.Background(), 15*time.Second)
				newConn, err := grpc.DialContext(
					newCtx,
					c.target,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
					grpc.WithBlock(),
				)

				if err == nil {
					// Successfully created new connection
					c.conn = newConn
					c.client = api.NewAPIClient(newConn)
				}

				newCancel()
			}

			cancel()
		}
	}
}

// Close closes all the streams and then the underlying connection. IMPORTANT: you should call this
// to ensure correct API accounting.
func (c *Client) Close() error {
	c.txStream.CloseSend()
	c.txSeqStream.CloseSend()
	c.submitBlockStream.CloseSend()

	return c.conn.Close()
}
