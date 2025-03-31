// package client is responsible for easily interacting with the protobuf API in Go.
// It contains wrappers for all the rpc methods that accept standard go-ethereum
// objects.
package client

import (
	"context"
	"time"

	"github.com/chainbound/fiber-go/protobuf/api"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// Client version string that is appended to each stream request.
const Version string = "fiber-go/1.9.3"

type Client struct {
	target string
	conn   *grpc.ClientConn
	config *ClientConfig
	client api.APIClient
	key    string
	logger *zap.SugaredLogger

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
	hcInterval        time.Duration
	logLevel          string
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
		hcInterval:        10 * time.Second,
		logLevel:          "fatal",
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

// SetIdleTimeout sets the duration after which idle connections will be restarted.
// Set to 0 to disable idle timeout (default).
// Note: If set below 10s, gRPC will use a minimum value of 10s instead.
func (c *ClientConfig) SetIdleTimeout(timeout time.Duration) *ClientConfig {
	c.idleTimeout = timeout
	return c
}

// SetHealthCheckInterval sets the interval for health checks.
func (c *ClientConfig) SetHealthCheckInterval(interval time.Duration) *ClientConfig {
	c.hcInterval = interval
	return c
}

func (c *ClientConfig) SetLogLevel(level string) *ClientConfig {
	c.logLevel = level
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
		logger: GetLogger(config.logLevel),
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
			Timeout: 10 * time.Second,
			// Allow keepalive pings even when there are no active streams
			PermitWithoutStream: true,
		}
		opts = append(opts, grpc.WithKeepaliveParams(kaParams))
	}

	conn, err := grpc.NewClient(c.target, opts...)
	if err != nil {
		return err
	}

	c.conn = conn

	// Create the stub (client) with the channel
	c.client = api.NewAPIClient(conn)

	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key, "x-client-version", Version)
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

	c.logger.Infof("Connected to %s", c.target)

	if c.config.hcInterval > 0 {
		c.logger.Debugw("Starting health check", "interval", c.config.hcInterval)
		// Start the health check
		go c.startHealthCheck()
	}

	return nil
}

func (c *Client) startHealthCheck() {
	ticker := time.NewTicker(c.config.hcInterval)
	defer ticker.Stop()

	for range ticker.C {
		if c.conn.GetState() != connectivity.Ready {
			c.logger.Warn("Connection is not ready, reconnecting...")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			// Reconnect
			c.txStream.CloseSend()
			c.txSeqStream.CloseSend()
			c.submitBlockStream.CloseSend()

			c.conn.Close()

			// Reconnect, this will start a new health check goroutine so we
			// return from the current one.
			c.Connect(ctx)

			cancel()

			return
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
