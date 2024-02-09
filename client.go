// package client is responsible for easily interacting with the protobuf API in Go.
// It contains wrappers for all the rpc methods that accept standard go-ethereum
// objects.
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/chainbound/fiber-go/filter"
	"github.com/chainbound/fiber-go/protobuf/api"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
}

// NewConfig creates a new config with sensible default values.
func NewConfig() *ClientConfig {
	return &ClientConfig{
		enableCompression: false,
		writeBufferSize:   1024 * 8,
		readBufferSize:    1024 * 8,
		connWindowSize:    1024 * 512,
		windowSize:        1024 * 256,
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

	conn, err := grpc.DialContext(ctx, c.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithWriteBufferSize(c.config.writeBufferSize),
		grpc.WithReadBufferSize(c.config.readBufferSize),
		grpc.WithInitialConnWindowSize(c.config.connWindowSize),
		grpc.WithInitialWindowSize(c.config.windowSize),
	)
	if err != nil {
		return err
	}

	c.conn = conn

	// Create the stub (client) with the channel
	c.client = api.NewAPIClient(conn)

	ctx = metadata.AppendToOutgoingContext(context.Background(), "x-api-key", c.key)
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

	return nil

}

// Close closes all the streams and then the underlying connection. IMPORTANT: you should call this
// to ensure correct API accounting.
func (c *Client) Close() error {
	c.txStream.CloseSend()
	c.txSeqStream.CloseSend()
	c.submitBlockStream.CloseSend()

	return c.conn.Close()
}

// SendTransaction sends the (signed) transaction to Fibernet and returns the hash and a timestamp (us).
// It blocks until the transaction was sent.
func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) (string, int64, error) {
	rlpTransaction, err := tx.MarshalBinary()
	if err != nil {
		return tx.Hash().Hex(), 0, err
	}

	errc := make(chan error)
	go func() {
		if err := c.txStream.Send(&api.TransactionMsg{RlpTransaction: rlpTransaction}); err != nil {
			errc <- err
		}
	}()

	for {
		select {
		case err := <-errc:
			return "", 0, err
		default:
		}

		res, err := c.txStream.Recv()
		if err != nil {
			return "", 0, err
		} else {
			return res.Hash, res.Timestamp, nil
		}
	}
}

// SendRawTransaction sends the RLP-encoded transaction to Fibernet and returns the hash and a timestamp (us).
func (c *Client) SendRawTransaction(ctx context.Context, rawTx []byte) (string, int64, error) {
	errc := make(chan error)
	go func() {
		if err := c.txStream.Send(&api.TransactionMsg{RlpTransaction: rawTx}); err != nil {
			errc <- err
		}
	}()

	for {
		select {
		case err := <-errc:
			return "", 0, err
		default:
		}

		res, err := c.txStream.Recv()
		if err != nil {
			return "", 0, err
		} else {
			return res.Hash, res.Timestamp, nil
		}
	}
}

func (c *Client) SendTransactionSequence(ctx context.Context, transactions ...*types.Transaction) ([]string, int64, error) {
	errc := make(chan error)

	rlpSequence := make([][]byte, len(transactions))

	for i, tx := range transactions {
		rlpTransaction, err := tx.MarshalBinary()
		if err != nil {
			fmt.Println("error marshalling transaction:", err)
			return nil, 0, err
		}

		rlpSequence[i] = rlpTransaction
	}

	go func() {
		if err := c.txSeqStream.Send(&api.TxSequenceMsgV2{Sequence: rlpSequence}); err != nil {
			errc <- err
		}
	}()

	select {
	case err := <-errc:
		return nil, 0, err
	default:
	}

	res, err := c.txSeqStream.Recv()
	if err != nil {
		return nil, 0, err
	}

	hashes := make([]string, len(res.SequenceResponse))
	ts := res.SequenceResponse[0].Timestamp

	for i, response := range res.SequenceResponse {
		hashes[i] = response.Hash
	}

	return hashes, ts, nil
}

// SendRawTransactionSequence sends a sequence of RLP-encoded transactions to Fibernet and returns the hashes and a timestamp (us).
func (c *Client) SendRawTransactionSequence(ctx context.Context, rawTransactions ...[]byte) ([]string, int64, error) {
	errc := make(chan error)

	go func() {
		if err := c.txSeqStream.Send(&api.TxSequenceMsgV2{Sequence: rawTransactions}); err != nil {
			errc <- err
		}
	}()

	select {
	case err := <-errc:
		return nil, 0, err
	default:
	}

	res, err := c.txSeqStream.Recv()
	if err != nil {
		return nil, 0, err
	}

	hashes := make([]string, len(res.SequenceResponse))

	ts := res.SequenceResponse[0].Timestamp

	for i, response := range res.SequenceResponse {
		hashes[i] = response.Hash
	}

	return hashes, ts, nil
}

// SubmitBlock submits an SSZ encoded signed block to Fiber and returns the slot, state root and timestamp (us).
func (c *Client) SubmitBlock(ctx context.Context, sszBlock []byte) (uint64, []byte, uint64, error) {
	errc := make(chan error)

	go func() {
		if err := c.submitBlockStream.Send(&api.BlockSubmissionMsg{SszBlock: sszBlock}); err != nil {
			errc <- err
		}
	}()

	select {
	case err := <-errc:
		return 0, nil, 0, err
	default:
	}

	res, err := c.submitBlockStream.Recv()
	if err != nil {
		return 0, nil, 0, err
	}

	return res.Slot, res.StateRoot, res.Timestamp, nil
}

// SubscribeNewTxs subscribes to new transactions, and sends transactions on the given
// channel according to the filter. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewTxs(filter *filter.Filter, ch chan<- *TransactionWithSender) error {
	attempts := 0
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)

		protoFilter := &api.TxFilter{}
		if filter != nil {
			protoFilter.Encoded = filter.Encode()
		}

		res, err := c.client.SubscribeNewTxsV2(ctx, protoFilter)
		if err != nil {
			if attempts > 50 {
				return fmt.Errorf("subscribing to transactions after 50 attempts: %w", err)
			}
			time.Sleep(time.Second * 2)
			continue outer
		}

		for {
			proto, err := res.Recv()
			// For now, retry on every error.
			if err != nil {
				time.Sleep(time.Second * 2)
				continue outer
			}

			tx := new(types.Transaction)
			if err := tx.UnmarshalBinary(proto.RlpTransaction); err != nil {
				continue outer
			}

			sender := common.BytesToAddress(proto.Sender)

			txWithSender := &TransactionWithSender{
				Sender:      &sender,
				Transaction: tx,
			}

			ch <- txWithSender
		}
	}
}

// SubscribeNewRawTxs subscribes to new RLP-encoded transaction bytes, and sends transactions on the given
// channel according to the filter. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewRawTxs(filter *filter.Filter, ch chan<- *RawTransactionWithSender) error {
	attempts := 0
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)

		protoFilter := &api.TxFilter{}
		if filter != nil {
			protoFilter.Encoded = filter.Encode()
		}

		res, err := c.client.SubscribeNewTxsV2(ctx, protoFilter)
		if err != nil {
			if attempts > 50 {
				return fmt.Errorf("subscribing to raw transactions after 50 attempts: %w", err)
			}
			time.Sleep(time.Second * 2)
			continue outer
		}

		for {
			proto, err := res.Recv()
			// For now, retry on every error.
			if err != nil {
				time.Sleep(time.Second * 2)
				continue outer
			}

			sender := common.BytesToAddress(proto.Sender)

			rawTxWithSender := &RawTransactionWithSender{
				Sender: &sender,
				Data:   proto.RlpTransaction,
			}

			ch <- rawTxWithSender
		}
	}
}

// SubscribeNewBlocks subscribes to new execution payloads, and sends blocks on the given
// channel. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewExecutionPayloads(ch chan<- *capella.ExecutionPayload) error {
	attempts := 0
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)

		res, err := c.client.SubscribeExecutionPayloadsV2(ctx, &emptypb.Empty{})
		if err != nil {
			if attempts > 50 {
				return fmt.Errorf("subscribing to execution payloads after 50 attempts: %w", err)
			}
			time.Sleep(time.Second * 2)
			continue outer
		}

		for {
			proto, err := res.Recv()
			if err != nil {
				time.Sleep(time.Second * 2)
				continue outer
			}

			block := new(capella.ExecutionPayload)
			if err := block.UnmarshalSSZ(proto.SszPayload); err != nil {
				fmt.Println("error unmarshalling execution payload:", err)
				continue outer
			}

			ch <- block
		}
	}
}

// SubscribeNewBeaconBlocks subscribes to new beacon blocks, and sends blocks on the given
// channel. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewBeaconBlocks(ch chan<- *capella.SignedBeaconBlock) error {
	attempts := 0
outer:
	for {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)

		res, err := c.client.SubscribeBeaconBlocksV2(ctx, &emptypb.Empty{})
		if err != nil {
			if attempts > 50 {
				return fmt.Errorf("subscribing to beacon blocks after 50 attempts: %w", err)
			}
			time.Sleep(time.Second * 2)
			continue outer
		}

		for {
			proto, err := res.Recv()
			if err != nil {
				time.Sleep(time.Second * 2)
				continue outer
			}

			block := new(capella.SignedBeaconBlock)
			if err := block.UnmarshalSSZ(proto.SszBlock); err != nil {
				fmt.Println("error unmarshalling beacon block:", err)
				continue outer
			}

			ch <- block
		}
	}
}
