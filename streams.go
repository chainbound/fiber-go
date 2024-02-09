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
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

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
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

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
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

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

// SubscribeNewRawExecutionPayloads subscribes to new SSZ-encoded raw execution payloads, and sends
// blocks on the given channel. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewRawExecutionPayloads(ch chan<- []byte) error {
	attempts := 0
outer:
	for {
		attempts++
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		res, err := c.client.SubscribeExecutionPayloadsV2(ctx, &emptypb.Empty{})
		if err != nil {
			if attempts > 50 {
				return fmt.Errorf("subscribing to raw execution payloads after 50 attempts: %w", err)
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

			ch <- proto.SszPayload
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
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

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

// SubscribeNewRawBeaconBlocks subscribes to new SSZ-encoded raw signed beacon blocks, and sends
// blocks on the given channel. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewRawBeaconBlocks(ch chan<- []byte) error {
	attempts := 0
outer:
	for {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)
		ctx = metadata.AppendToOutgoingContext(ctx, "x-client-version", Version)

		res, err := c.client.SubscribeBeaconBlocksV2(ctx, &emptypb.Empty{})
		if err != nil {
			if attempts > 50 {
				return fmt.Errorf("subscribing to raw beacon blocks after 50 attempts: %w", err)
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

			ch <- proto.SszBlock
		}
	}
}
