// package client is responsible for easily interacting with the protobuf API in Go.
// It contains wrappers for all the rpc methods that accept standard go-ethereum
// objects.
package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/chainbound/fiber-go/filter"
	"github.com/chainbound/fiber-go/protobuf/api"
	"github.com/chainbound/fiber-go/protobuf/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Transaction struct {
	ChainID     uint32
	To          *common.Address
	From        common.Address
	Gas         uint64
	GasPrice    *big.Int
	Hash        common.Hash
	Input       []byte
	Value       *big.Int
	Nonce       uint64
	Type        uint32
	MaxFee      *big.Int
	PriorityFee *big.Int
	V           uint64
	R           []byte
	S           []byte
	AccessList  types.AccessList
}

type Header struct {
	Number        uint64
	Hash          common.Hash
	ParentHash    common.Hash
	PrevRandao    common.Hash
	StateRoot     common.Hash
	ReceiptRoot   common.Hash
	FeeRecipient  common.Address
	ExtraData     []byte
	GasLimit      uint64
	GasUsed       uint64
	Timestamp     uint64
	LogsBloom     types.Bloom
	BaseFeePerGas *big.Int
}

type Block struct {
	Number        uint64
	Hash          common.Hash
	ParentHash    common.Hash
	PrevRandao    common.Hash
	StateRoot     common.Hash
	ReceiptRoot   common.Hash
	FeeRecipient  common.Address
	ExtraData     []byte
	GasLimit      uint64
	GasUsed       uint64
	Timestamp     uint64
	LogsBloom     types.Bloom
	BaseFeePerGas *big.Int
	Transactions  []*Transaction
}

func (tx *Transaction) ToNative() *types.Transaction {
	switch tx.Type {
	case 0:
		return types.NewTx(&types.LegacyTx{
			Nonce:    tx.Nonce,
			GasPrice: tx.GasPrice,
			Gas:      tx.Gas,
			To:       tx.To,
			Value:    tx.Value,
			Data:     tx.Input,
			V:        big.NewInt(int64(tx.V)),
			R:        new(big.Int).SetBytes(tx.R),
			S:        new(big.Int).SetBytes(tx.S),
		})
	case 1:
		return types.NewTx(&types.AccessListTx{
			ChainID:  big.NewInt(int64(tx.ChainID)),
			Nonce:    tx.Nonce,
			GasPrice: tx.GasPrice,
			Gas:      tx.Gas,
			To:       tx.To,
			Value:    tx.Value,
			Data:     tx.Input,
			V:        big.NewInt(int64(tx.V)),
			R:        new(big.Int).SetBytes(tx.R),
			S:        new(big.Int).SetBytes(tx.S),
		})
	case 2:
		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(int64(tx.ChainID)),
			Nonce:     tx.Nonce,
			GasFeeCap: tx.MaxFee,
			GasTipCap: tx.PriorityFee,
			Gas:       tx.Gas,
			To:        tx.To,
			Value:     tx.Value,
			Data:      tx.Input,
			V:         big.NewInt(int64(tx.V)),
			R:         new(big.Int).SetBytes(tx.R),
			S:         new(big.Int).SetBytes(tx.S),
		})
	}

	return nil
}

type Client struct {
	target string
	conn   *grpc.ClientConn
	client api.APIClient
	key    string

	// streams
	txStream       api.API_SendTransactionClient
	rawTxStream    api.API_SendRawTransactionClient
	txSeqStream    api.API_SendTransactionSequenceClient
	rawTxSeqStream api.API_SendRawTransactionSequenceClient
}

func NewClient(target, apiKey string) *Client {
	return &Client{
		target: target,
		key:    apiKey,
	}
}

// Connects sets up the gRPC channel and creates the stub. It blocks until connected or the given context expires.
// Always use a context with timeout.
func (c *Client) Connect(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, c.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithReadBufferSize(0),
		grpc.WithWriteBufferSize(0),
	)
	if err != nil {
		return err
	}

	c.conn = conn

	// Create the stub (client) with the channel
	c.client = api.NewAPIClient(conn)

	ctx = metadata.AppendToOutgoingContext(context.Background(), "x-api-key", c.key)
	c.txStream, err = c.client.SendTransaction(ctx)
	if err != nil {
		return err
	}

	c.rawTxStream, err = c.client.SendRawTransaction(ctx)
	if err != nil {
		return err
	}

	c.txSeqStream, err = c.client.SendTransactionSequence(ctx)
	if err != nil {
		return err
	}

	c.rawTxSeqStream, err = c.client.SendRawTransactionSequence(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Close closes all the streams and then the underlying connection. IMPORTANT: you should call this
// to ensure correct API accounting.
func (c *Client) Close() error {
	c.txStream.CloseSend()
	c.rawTxStream.CloseSend()
	c.txSeqStream.CloseSend()
	c.rawTxSeqStream.CloseSend()

	return c.conn.Close()
}

// SendTransaction sends the (signed) transaction to Fibernet and returns the hash and a timestamp (us).
// It blocks until the transaction was sent.
func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) (string, int64, error) {
	proto, err := TxToProto(tx)
	if err != nil {
		return "", 0, fmt.Errorf("converting to protobuf: %w", err)
	}

	errc := make(chan error)
	go func() {
		if err := c.txStream.Send(proto); err != nil {
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

func (c *Client) SendRawTransaction(ctx context.Context, rawTx []byte) (string, int64, error) {
	errc := make(chan error)
	go func() {
		if err := c.rawTxStream.Send(&api.RawTxMsg{RawTx: rawTx}); err != nil {
			errc <- err
		}
	}()

	for {
		select {
		case err := <-errc:
			return "", 0, err
		default:
		}

		res, err := c.rawTxStream.Recv()
		if err != nil {
			return "", 0, err
		} else {
			return res.Hash, res.Timestamp, nil
		}
	}
}

func (c *Client) SendTransactionSequence(ctx context.Context, transactions ...*types.Transaction) ([]string, int64, error) {
	errc := make(chan error)

	protoSeq := make([]*eth.Transaction, len(transactions))

	for i, tx := range transactions {
		proto, err := TxToProto(tx)
		if err != nil {
			return nil, 0, err
		}

		protoSeq[i] = proto
	}

	go func() {
		if err := c.txSeqStream.Send(&api.TxSequenceMsg{Sequence: protoSeq}); err != nil {
			errc <- err
		}
	}()

	for {
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
}

func (c *Client) SendRawTransactionSequence(ctx context.Context, rawTransactions ...[]byte) ([]string, int64, error) {
	errc := make(chan error)

	go func() {
		if err := c.rawTxSeqStream.Send(&api.RawTxSequenceMsg{RawTxs: rawTransactions}); err != nil {
			errc <- err
		}
	}()

	for {
		select {
		case err := <-errc:
			return nil, 0, err
		default:
		}

		res, err := c.rawTxSeqStream.Recv()
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
}

// SubscribeNewTxs subscribes to new transactions, and sends transactions on the given
// channel according to the filter. This function blocks and should be called in a goroutine.
// If there's an error receiving the new message it will close the channel and return the error.
func (c *Client) SubscribeNewTxs(filter *filter.Filter, ch chan<- *Transaction) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)

	protoFilter := &api.TxFilter{}
	if filter != nil {
		protoFilter.Encoded = filter.Encode()
	}

	res, err := c.client.SubscribeNewTxs(ctx, protoFilter)
	if err != nil {
		return fmt.Errorf("subscribing to transactions: %w", err)
	}

	for {
		proto, err := res.Recv()
		if err != nil {
			close(ch)
			return err
		}

		ch <- ProtoToTx(proto)
	}
}

func (c *Client) SubscribeNewHeaders(ch chan<- *Header) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)

	res, err := c.client.SubscribeNewBlocks(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("subscribing to blocks: %w", err)
	}

	for {
		proto, err := res.Recv()
		if err != nil {
			close(ch)
			return err
		}

		ch <- ProtoToHeader(proto)
	}
}

func (c *Client) SubscribeNewBlocks(ch chan<- *Block) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.key)

	res, err := c.client.SubscribeNewBlocks(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("subscribing to blocks: %w", err)
	}

	for {
		proto, err := res.Recv()
		if err != nil {
			close(ch)
			return err
		}

		ch <- ProtoToBlock(proto)
	}
}

func TxToProto(tx *types.Transaction) (*eth.Transaction, error) {
	signer := types.NewLondonSigner(common.Big1)
	sender, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}

	var to []byte
	if tx.To() != nil {
		to = tx.To().Bytes()
	}

	var acl []*eth.AccessTuple
	if tx.Type() != 0 {
		if len(tx.AccessList()) > 0 {
			acl = make([]*eth.AccessTuple, len(tx.AccessList()))
			for i, tuple := range tx.AccessList() {
				acl[i] = &eth.AccessTuple{
					Address: tuple.Address.Bytes(),
				}

				storageKeys := make([][]byte, len(tuple.StorageKeys))
				for j, key := range tuple.StorageKeys {
					storageKeys[j] = key.Bytes()
				}
			}
		}
	}

	v, r, s := tx.RawSignatureValues()
	return &eth.Transaction{
		ChainId:     uint32(tx.ChainId().Uint64()),
		To:          to,
		Gas:         tx.Gas(),
		GasPrice:    tx.GasPrice().Uint64(),
		MaxFee:      tx.GasFeeCap().Uint64(),
		PriorityFee: tx.GasTipCap().Uint64(),
		Hash:        tx.Hash().Bytes(),
		Input:       tx.Data(),
		Nonce:       tx.Nonce(),
		Value:       tx.Value().Bytes(),
		From:        sender.Bytes(),
		Type:        uint32(tx.Type()),
		V:           v.Uint64(),
		R:           r.Bytes(),
		S:           s.Bytes(),
		AccessList:  acl,
	}, nil
}

// ProtoToTx converts a protobuf transaction to a go-ethereum transaction.
// It does not include the AccessList.
func ProtoToTx(proto *eth.Transaction) *Transaction {
	to := new(common.Address)

	if len(proto.To) > 0 {
		to = (*common.Address)(proto.To)
	}

	var acl []types.AccessTuple
	if len(proto.AccessList) > 0 {
		acl = make([]types.AccessTuple, len(proto.AccessList))
		for i, tuple := range proto.AccessList {
			storageKeys := make([]common.Hash, len(tuple.StorageKeys))

			for j, key := range tuple.StorageKeys {
				storageKeys[j] = common.BytesToHash(key)
			}

			acl[i] = types.AccessTuple{
				Address:     common.BytesToAddress(tuple.Address),
				StorageKeys: storageKeys,
			}
		}
	}

	if proto.Type > 0 {
		if proto.V > 1 {
			proto.V = proto.V - 37
		}
	} else {
		if proto.V > 30 {
			proto.V = proto.V - 10
		}
	}

	return &Transaction{
		ChainID:     proto.ChainId,
		Type:        proto.Type,
		Nonce:       proto.Nonce,
		GasPrice:    big.NewInt(int64(proto.GasPrice)),
		MaxFee:      big.NewInt(int64(proto.MaxFee)),
		PriorityFee: big.NewInt(int64(proto.PriorityFee)),
		Gas:         proto.Gas,
		To:          to,
		From:        common.BytesToAddress(proto.From),
		Hash:        common.BytesToHash(proto.Hash),
		Value:       new(big.Int).SetBytes(proto.Value),
		Input:       proto.Input,
		V:           proto.V,
		R:           proto.R,
		S:           proto.S,
		AccessList:  acl,
	}
}

func ProtoToHeader(proto *eth.Block) *Header {
	return &Header{
		Number:        proto.Number,
		Hash:          common.BytesToHash(proto.Hash),
		ParentHash:    common.BytesToHash(proto.ParentHash),
		StateRoot:     common.BytesToHash(proto.StateRoot),
		ReceiptRoot:   common.BytesToHash(proto.ReceiptRoot),
		PrevRandao:    common.BytesToHash(proto.PrevRandao),
		LogsBloom:     types.BytesToBloom(proto.LogsBloom),
		GasLimit:      proto.GasLimit,
		GasUsed:       proto.GasUsed,
		Timestamp:     proto.Timestamp,
		ExtraData:     proto.ExtraData,
		FeeRecipient:  common.BytesToAddress(proto.FeeRecipient),
		BaseFeePerGas: new(big.Int).SetBytes(proto.BaseFeePerGas),
	}
}

func ProtoToBlock(proto *eth.Block) *Block {
	txs := make([]*Transaction, len(proto.Transactions))
	for i, proto := range proto.Transactions {
		txs[i] = ProtoToTx(proto)
	}

	return &Block{
		Number:        proto.Number,
		Hash:          common.BytesToHash(proto.Hash),
		StateRoot:     common.BytesToHash(proto.StateRoot),
		ParentHash:    common.BytesToHash(proto.ParentHash),
		ReceiptRoot:   common.BytesToHash(proto.ReceiptRoot),
		PrevRandao:    common.BytesToHash(proto.PrevRandao),
		LogsBloom:     types.BytesToBloom(proto.LogsBloom),
		GasLimit:      proto.GasLimit,
		GasUsed:       proto.GasUsed,
		Timestamp:     proto.Timestamp,
		ExtraData:     proto.ExtraData,
		FeeRecipient:  common.BytesToAddress(proto.FeeRecipient),
		BaseFeePerGas: new(big.Int).SetBytes(proto.BaseFeePerGas),
		Transactions:  txs,
	}
}
