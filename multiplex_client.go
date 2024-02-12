package client

import (
	"context"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/chainbound/fiber-go/filter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/sync/errgroup"
)

type MultiplexClient struct {
	targets      []string
	clients      []*Client
	txCache      *lru.Cache[common.Hash, struct{}]
	headerCache  *lru.Cache[common.Hash, struct{}]
	payloadCache *lru.Cache[common.Hash, struct{}]
	beaconCache  *lru.Cache[common.Hash, struct{}]
}

// Returns a new multiplex client with the given targets.
func NewMultiplexClient(targets []string, apiKey string) *MultiplexClient {
	clients := make([]*Client, len(targets))
	for i, target := range targets {
		clients[i] = NewClient(target, apiKey)
	}

	txCache, err := lru.New[common.Hash, struct{}](8192)
	if err != nil {
		panic(err)
	}

	headerCache, err := lru.New[common.Hash, struct{}](128)
	if err != nil {
		panic(err)
	}

	payloadCache, err := lru.New[common.Hash, struct{}](128)
	if err != nil {
		panic(err)
	}

	beaconCache, err := lru.New[common.Hash, struct{}](128)
	if err != nil {
		panic(err)
	}

	return &MultiplexClient{
		targets:      targets,
		clients:      clients,
		txCache:      txCache,
		headerCache:  headerCache,
		payloadCache: payloadCache,
		beaconCache:  beaconCache,
	}
}

// Connects all the underlying clients concurrently. It blocks until the connection succeeds / fails, or the given context expires.
func (mc *MultiplexClient) Connect(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	// Connect all clients concurrently
	for _, client := range mc.clients {
		// Closure stuff
		client := client
		errGroup.Go(func() error {
			if err := client.Connect(ctx); err != nil {
				return fmt.Errorf("connecting to %s: %w", client.target, err)
			}

			return nil
		})
	}

	return errGroup.Wait()
}

type txRes struct {
	hash string
	ts   int64
	err  error
}

// SendTransaction sends the (signed) transaction on all clients, and returns the first response it receives.
func (mc *MultiplexClient) SendTransaction(ctx context.Context, tx *types.Transaction) (string, int64, error) {
	resCh := make(chan txRes, len(mc.clients))

	for _, client := range mc.clients {
		// Closure stuff
		client := client
		go func() {
			hash, ts, err := client.SendTransaction(ctx, tx)
			resCh <- txRes{hash, ts, err}
		}()
	}

	res := <-resCh
	return res.hash, res.ts, res.err
}

// SendRawTransaction sends the RLP encoded, signed transaction on all clients, and returns the first response it receives.
func (mc *MultiplexClient) SendRawTransaction(ctx context.Context, rawTx []byte) (string, int64, error) {
	resCh := make(chan txRes, len(mc.clients))

	for _, client := range mc.clients {
		// Closure stuff
		client := client
		go func() {
			hash, ts, err := client.SendRawTransaction(ctx, rawTx)
			resCh <- txRes{hash, ts, err}
		}()
	}

	res := <-resCh
	return res.hash, res.ts, res.err
}

type seqRes struct {
	hashes []string
	ts     int64
	err    error
}

func (mc *MultiplexClient) SendTransactionSequence(ctx context.Context, transactions ...*types.Transaction) ([]string, int64, error) {
	resCh := make(chan seqRes, len(mc.clients))

	for _, client := range mc.clients {
		// Closure stuff
		client := client
		go func() {
			hashes, ts, err := client.SendTransactionSequence(ctx, transactions...)
			resCh <- seqRes{hashes, ts, err}
		}()
	}

	res := <-resCh
	return res.hashes, res.ts, res.err
}

func (mc *MultiplexClient) SendRawTransactionSequence(ctx context.Context, transactions ...[]byte) ([]string, int64, error) {
	resCh := make(chan seqRes, len(mc.clients))

	for _, client := range mc.clients {
		// Closure stuff
		client := client
		go func() {
			hashes, ts, err := client.SendRawTransactionSequence(ctx, transactions...)
			resCh <- seqRes{hashes, ts, err}
		}()
	}

	res := <-resCh
	return res.hashes, res.ts, res.err
}

// SubscribeNewTxs subscribes to new transactions, and sends transactions on the given
// channel according to the filter. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
// It multiplexes the subscription across all clients.
func (mc *MultiplexClient) SubscribeNewTxs(filter *filter.Filter, ch chan<- *TransactionWithSender) error {
	mch := make(chan *TransactionWithSender)
	errc := make(chan error, len(mc.clients))

	for _, client := range mc.clients {
		client := client
		go func() {
			if err := client.SubscribeNewTxs(filter, mch); err != nil {
				errc <- err
			}
		}()
	}

	for {
		select {
		case tx := <-mch:
			if _, ok := mc.txCache.Get(tx.Transaction.Hash()); ok {
				continue
			}

			mc.txCache.Add(tx.Transaction.Hash(), struct{}{})
			ch <- tx
		case err := <-errc:
			return err
		}
	}
}

// SubscribeNewBeaconHeaders subscribes to new beacon headers, and sends headers on the given channel. This function blocks
// and should be called in a goroutine.
// It multiplexes the subscription across all clients.
func (mc *MultiplexClient) SubscribeNewExecutionPayloads(ch chan<- *Block) error {
	mch := make(chan *Block)
	errc := make(chan error, len(mc.clients))

	for _, client := range mc.clients {
		client := client
		go func() {
			if err := client.SubscribeNewExecutionPayloads(mch); err != nil {
				errc <- err
			}
		}()
	}

	for {
		select {
		case payload := <-mch:
			if _, ok := mc.payloadCache.Get(payload.Hash); ok {
				continue
			}

			mc.payloadCache.Add(payload.Hash, struct{}{})
			ch <- payload
		case err := <-errc:
			return err
		}
	}
}

// SubscribeNewBeaconHeaders subscribes to new beacon headers, and sends headers on the given channel. This function blocks
// and should be called in a goroutine.
// It multiplexes the subscription across all clients.
func (mc *MultiplexClient) SubscribeNewBeaconBlocks(ch chan<- *capella.SignedBeaconBlock) error {
	mch := make(chan *capella.SignedBeaconBlock)
	errc := make(chan error, len(mc.clients))

	for _, client := range mc.clients {
		client := client
		go func() {
			if err := client.SubscribeNewBeaconBlocks(mch); err != nil {
				errc <- err
			}
		}()
	}

	for {
		select {
		case block := <-mch:
			if _, ok := mc.beaconCache.Get(block.Message.Body.ExecutionPayload.StateRoot); ok {
				continue
			}

			mc.beaconCache.Add(block.Message.Body.ExecutionPayload.StateRoot, struct{}{})
			ch <- block
		case err := <-errc:
			return err
		}
	}
}

func (mc *MultiplexClient) Close() error {
	errGroup, _ := errgroup.WithContext(context.Background())

	// Connect all clients concurrently
	for _, client := range mc.clients {
		// Closure stuff
		client := client
		errGroup.Go(func() error {
			if err := client.Close(); err != nil {
				return fmt.Errorf("closing %s: %w", client.target, err)
			}

			return nil
		})
	}

	return errGroup.Wait()

}
