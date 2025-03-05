package client

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/chainbound/fiber-go/filter"
	"github.com/chainbound/fiber-go/protobuf/api"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// All Available streams:
//
// - sendTransaction
// - sendTransactionSequence
// - sendRawTransaction
// - sendRawTransactionSequence
// - submitBlock
// - subscribeNewTxs
// - subscribeNewRawTxs
// - subscribeNewBlobTxs
// - subscribeNewExecutionPayloads
// - subscribeNewRawExecutionPayloads
// - subscribeNewBeaconBlocks
// - subscribeNewRawBeaconBlocks

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

// SendTransactionSequence sends a sequence of transactions to Fibernet and returns the hashes and a timestamp (us).
func (c *Client) SendTransactionSequence(ctx context.Context, transactions ...*types.Transaction) ([]string, int64, error) {
	errc := make(chan error)

	rlpSequence := make([][]byte, len(transactions))

	for i, tx := range transactions {
		rlpTransaction, err := tx.MarshalBinary()
		if err != nil {
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

// min returns the smaller of x or y.
func min(x, y time.Duration) time.Duration {
	if x < y {
		return x
	}
	return y
}

// createHeartbeat starts a goroutine that monitors the connection state
// and triggers reconnection if the connection is in a bad state.
// Returns a channel that should be closed when the monitoring should stop.
func (c *Client) createHeartbeat(checkInterval time.Duration) chan struct{} {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Check connection state
				state := c.conn.GetState()
				if state != connectivity.Ready && state != connectivity.Idle && state != connectivity.Connecting {
					// Connection is in a bad state, force close it to trigger reconnection
					fmt.Printf("Detected bad connection state: %s, forcing reconnection\n", state)
					c.conn.Close()
					return
				}
			}
		}
	}()
	return done
}

// Helper function to wait for a ready connection state
func (c *Client) waitForReadyConnection(timeout time.Duration) {
	state := c.conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle {
		// Wait for connection to become ready or timeout
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
		c.conn.WaitForStateChange(timeoutCtx, state)
		timeoutCancel()
	}
}

// backoffSubscriptionError handles exponential backoff logic for subscription retries.
// It returns false if all retries have been attempted.
func (c *Client) backoffSubscriptionError(err error, attempts int, backoff *time.Duration, maxBackoff time.Duration, subscriptionName string) bool {
	// After too many attempts, give up
	if attempts > 50 {
		return false
	}

	// Use exponential backoff with jitter
	sleepTime := *backoff + time.Duration(rand.Int63n(int64(*backoff/2)))
	fmt.Printf("%s subscription error, retrying in %v: %v\n", subscriptionName, sleepTime, err)
	time.Sleep(sleepTime)

	// Increase backoff for next attempt, up to the maximum
	if *backoff < maxBackoff {
		*backoff = time.Duration(float64(*backoff) * 1.5)
		if *backoff > maxBackoff {
			*backoff = maxBackoff
		}
	}
	return true
}

// SubscribeNewTxs subscribes to new transactions, and sends transactions on the given
// channel according to the filter. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewTxs(filter *filter.Filter, ch chan<- *TransactionWithSender) error {
	attempts := 0
	backoff := time.Second * 1
	maxBackoff := time.Second * 30
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		protoFilter := &api.TxFilter{}
		if filter != nil {
			protoFilter.Encoded = filter.Encode()
		}

		// Wait for connection to be ready before attempting to subscribe
		c.waitForReadyConnection(5 * time.Second)

		res, err := c.client.SubscribeNewTxsV2(ctx, protoFilter)
		if err != nil {
			cancel() // Cancel this context since we're going to retry

			if !c.backoffSubscriptionError(err, attempts, &backoff, maxBackoff, "Transaction") {
				return fmt.Errorf("failed to subscribe to transactions after 50 attempts: %w", err)
			}
			continue outer
		}

		defer cancel()

		for {
			msg, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Stream completed")
					break
				}

				fmt.Printf("Stream error, reconnecting: %v\n", err)
				break
			}

			// Decode tx
			tx := new(types.Transaction)
			if err := tx.UnmarshalBinary(msg.RlpTransaction); err != nil {
				fmt.Printf("Error decoding transaction: %v\n", err)
				continue
			}

			// Skip transactions if we're in a bad state (e.g. reused stream)
			if tx.Hash().String() == "" {
				continue
			}

			sender := common.BytesToAddress(msg.Sender)
			txWithSender := &TransactionWithSender{
				Transaction: tx,
				Sender:      &sender,
			}

			ch <- txWithSender
		}
	}
}

// SubscribeNewRawTxs subscribes to new transactions, and sends the raw transaction data on the given
// channel according to the filter. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewRawTxs(filter *filter.Filter, ch chan<- *RawTransactionWithSender) error {
	attempts := 0
	backoff := time.Second * 1
	maxBackoff := time.Second * 30
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		protoFilter := &api.TxFilter{}
		if filter != nil {
			protoFilter.Encoded = filter.Encode()
		}

		// Wait for connection to be ready before attempting to subscribe
		c.waitForReadyConnection(5 * time.Second)

		res, err := c.client.SubscribeNewTxsV2(ctx, protoFilter)
		if err != nil {
			cancel() // Cancel this context since we're going to retry

			if !c.backoffSubscriptionError(err, attempts, &backoff, maxBackoff, "Raw transaction") {
				return fmt.Errorf("failed to subscribe to raw transactions after 50 attempts: %w", err)
			}
			continue outer
		}

		defer cancel()

		for {
			msg, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Stream completed")
					break
				}

				fmt.Printf("Stream error, reconnecting: %v\n", err)
				break
			}

			// Skip empty or invalid transactions
			if len(msg.RlpTransaction) == 0 {
				continue
			}

			sender := common.BytesToAddress(msg.Sender)
			ch <- &RawTransactionWithSender{
				Rlp:    msg.RlpTransaction,
				Sender: &sender,
			}
		}
	}
}

// SubscribeNewBlobTxs subscribes to new blob transactions, and sends transactions on the given channel.
// This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewBlobTxs(ch chan<- *TransactionWithSender) error {
	attempts := 0
	backoff := time.Second * 1
	maxBackoff := time.Second * 30
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		// Wait for connection to be ready before attempting to subscribe
		c.waitForReadyConnection(5 * time.Second)

		res, err := c.client.SubscribeNewBlobTxs(ctx, &emptypb.Empty{})
		if err != nil {
			cancel() // Cancel this context since we're going to retry

			if !c.backoffSubscriptionError(err, attempts, &backoff, maxBackoff, "Blob transaction") {
				return fmt.Errorf("failed to subscribe to blob transactions after 50 attempts: %w", err)
			}
			continue outer
		}

		defer cancel()

		// Create a heartbeat goroutine to detect silent disconnects
		heartbeatDone := c.createHeartbeat(5 * time.Second)
		defer close(heartbeatDone)

		// Receive transactions in this loop
		for {
			msg, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Stream completed")
					break
				}

				fmt.Printf("Stream error, reconnecting: %v\n", err)
				break
			}

			tx := new(types.Transaction)
			if err := tx.UnmarshalBinary(msg.RlpTransaction); err != nil {
				// Bad transaction data, but connection is still good
				continue
			}

			sender := common.BytesToAddress(msg.Sender)
			txWithSender := TransactionWithSender{
				Sender:      &sender,
				Transaction: tx,
			}

			ch <- &txWithSender
		}
	}
}

// SubscribeNewExecutionPayloads subscribes to new execution payloads, and sends them on the given channel.
func (c *Client) SubscribeNewExecutionPayloads(ch chan<- *Block) error {
	attempts := 0
	backoff := time.Second * 1
	maxBackoff := time.Second * 30
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		// Wait for connection to be ready before attempting to subscribe
		c.waitForReadyConnection(5 * time.Second)

		res, err := c.client.SubscribeExecutionPayloadsV2(ctx, &emptypb.Empty{})
		if err != nil {
			cancel() // Cancel this context since we're going to retry

			if !c.backoffSubscriptionError(err, attempts, &backoff, maxBackoff, "Execution payload") {
				return fmt.Errorf("subscribing to execution payloads after 50 attempts: %w", err)
			}
			continue outer
		}

		defer cancel()

		// Create a heartbeat goroutine to detect silent disconnects
		heartbeatDone := c.createHeartbeat(5 * time.Second)
		defer close(heartbeatDone)

		// Receive blocks in this loop
		for {
			msg, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Stream completed")
					break
				}

				fmt.Printf("Block stream error, reconnecting: %v\n", err)
				break
			}

			var block *Block
			var decodeErr error

			switch msg.DataVersion {
			case DataVersionBellatrix:
				block, decodeErr = DecodeBellatrixExecutionPayload(msg)
			case DataVersionCapella:
				block, decodeErr = DecodeCapellaExecutionPayload(msg)
			case DataVersionDeneb:
				block, decodeErr = DecodeDenebExecutionPayload(msg)
			}

			if decodeErr != nil {
				fmt.Printf("Failed to decode execution payload: %v\n", decodeErr)
				continue
			}

			// Send the block to the channel
			ch <- block
		}

		// If we break out of the loop, reconnect
		continue outer
	}
}

// SubscribeNewBeaconBlocks subscribes to new beacon blocks, and sends them on the given channel.
func (c *Client) SubscribeNewBeaconBlocks(ch chan<- *SignedBeaconBlock) error {
	attempts := 0
	backoff := time.Second * 1
	maxBackoff := time.Second * 30
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		// Wait for connection to be ready before attempting to subscribe
		c.waitForReadyConnection(5 * time.Second)

		res, err := c.client.SubscribeBeaconBlocksV2(ctx, &emptypb.Empty{})
		if err != nil {
			cancel() // Cancel this context since we're going to retry

			if !c.backoffSubscriptionError(err, attempts, &backoff, maxBackoff, "Beacon block") {
				return fmt.Errorf("subscribing to beacon blocks after 50 attempts: %w", err)
			}
			continue outer
		}

		defer cancel()

		// Create a heartbeat goroutine to detect silent disconnects
		heartbeatDone := c.createHeartbeat(5 * time.Second)
		defer close(heartbeatDone)

		for {
			proto, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Stream completed")
					break
				}

				fmt.Printf("Beacon block stream error, reconnecting: %v\n", err)
				break
			}

			signedBeaconBlock := new(SignedBeaconBlock)
			signedBeaconBlock.DataVersion = proto.DataVersion

			switch proto.DataVersion {
			case DataVersionBellatrix:
				block := new(bellatrix.SignedBeaconBlock)
				if err := block.UnmarshalSSZ(proto.SszBlock); err != nil {
					continue outer
				}
				signedBeaconBlock.Bellatrix = block
				ch <- signedBeaconBlock
			case DataVersionCapella:
				block := new(capella.SignedBeaconBlock)
				if err := block.UnmarshalSSZ(proto.SszBlock); err != nil {
					continue outer
				}
				signedBeaconBlock.Capella = block
				ch <- signedBeaconBlock
			case DataVersionDeneb:
				block := new(deneb.SignedBeaconBlock)
				if err := block.UnmarshalSSZ(proto.SszBlock); err != nil {
					continue outer
				}
				signedBeaconBlock.Deneb = block
				ch <- signedBeaconBlock
			}
		}
	}
}

// SubscribeNewRawBeaconBlocks subscribes to new SSZ-encoded raw signed beacon blocks, and sends
// blocks on the given channel. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewRawBeaconBlocks(ch chan<- []byte) error {
	attempts := 0
	backoff := time.Second * 1
	maxBackoff := time.Second * 30
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		// Wait for connection to be ready before attempting to subscribe
		c.waitForReadyConnection(5 * time.Second)

		res, err := c.client.SubscribeBeaconBlocksV2(ctx, &emptypb.Empty{})
		if err != nil {
			cancel() // Cancel this context since we're going to retry

			if !c.backoffSubscriptionError(err, attempts, &backoff, maxBackoff, "Raw beacon block") {
				return fmt.Errorf("subscribing to raw beacon blocks after 50 attempts: %w", err)
			}
			continue outer
		}

		defer cancel()

		// Create a heartbeat goroutine to detect silent disconnects
		heartbeatDone := c.createHeartbeat(5 * time.Second)
		defer close(heartbeatDone)

		for {
			proto, err := res.Recv()
			if err != nil {
				// Instead of simple sleep, use backoff logic for reconnection
				if !c.backoffSubscriptionError(err, attempts, &backoff, maxBackoff, "Raw beacon block stream") {
					return fmt.Errorf("error receiving beacon block after 50 attempts: %w", err)
				}
				continue outer
			}

			ch <- proto.SszBlock
		}
	}
}
