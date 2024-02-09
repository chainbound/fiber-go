package client

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Helper type that wraps a transaction and its sender.
type TransactionWithSender struct {
	Sender      *common.Address
	Transaction *types.Transaction
}

// Helper type that wraps a raw, RLP-encoded transaction  and its sender.
type RawTransactionWithSender struct {
	Sender *common.Address
	Data   []byte
}
