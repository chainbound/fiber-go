package client

import (
	"fmt"
	"math/big"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/chainbound/fiber-go/protobuf/api"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Helper type that wraps a transaction and its sender.
//
// The sender address is included separately even though it could be calculated
// via ECDSA recovery from the raw RLP transaction bytes. The reason for this is
// that the recovery process is CPU-bound and takes time (order of magnitude of ~200 μs)
type TransactionWithSender struct {
	Sender      *common.Address
	Transaction *types.Transaction
}

// Helper type that wraps a raw, RLP-encoded transaction  and its sender.
//
// The sender address is included separately even though it could be calculated
// via ECDSA recovery from the raw RLP transaction bytes. The reason for this is
// that the recovery process is CPU-bound and takes time (order of magnitude of ~200 μs)
type RawTransactionWithSender struct {
	Sender *common.Address
	Rlp    []byte
}

// Helper type that wraps go-ethereum core primitives for an Ethereum block,
// such as Header, Transactions and Withdrawals.
type Block struct {
	Hash         common.Hash
	Header       types.Header
	Transactions []types.Transaction
	Withdrawals  []types.Withdrawal
}

// Helper type that wraps a signed beacon block from any of the supported hard-forks.
//
// This type will either contain a Bellatrix, Capella, or Deneb signed beacon block.
//
// DataVersion is used to indicate which type of payload is contained in the struct:
// 3: Bellatrix, 4: Capella, 5: Deneb
type SignedBeaconBlock struct {
	DataVersion uint32
	Bellatrix   *bellatrix.SignedBeaconBlock
	Capella     *capella.SignedBeaconBlock
	Deneb       *deneb.SignedBeaconBlock
}

const DataVersionBellatrix uint32 = 3
const DataVersionCapella uint32 = 4
const DataVersionDeneb uint32 = 5

func DecodeBellatrixExecutionPayload(input *api.ExecutionPayloadMsg) (*Block, error) {
	payload := new(bellatrix.ExecutionPayload)

	if err := payload.UnmarshalSSZ(input.SszPayload); err != nil {
		fmt.Println("error unmarshalling execution payload:", err)
		return nil, err
	}

	transactions := make([]types.Transaction, len(payload.Transactions))
	for _, rawTx := range payload.Transactions {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(rawTx); err != nil {
			continue
		}
		transactions = append(transactions, *tx)
	}

	basefee := new(big.Int).SetBytes(payload.BaseFeePerGas[:])

	header := types.Header{
		ParentHash:  common.Hash(payload.ParentHash),
		UncleHash:   common.Hash([32]byte{}),
		Coinbase:    common.Address(payload.FeeRecipient),
		Root:        payload.StateRoot,
		ReceiptHash: common.Hash(payload.ReceiptsRoot),
		Bloom:       payload.LogsBloom,
		Difficulty:  nil,
		Number:      new(big.Int).SetUint64(payload.BlockNumber),
		GasLimit:    payload.GasLimit,
		GasUsed:     payload.GasUsed,
		Time:        payload.Timestamp,
		Extra:       payload.ExtraData,
		MixDigest:   payload.PrevRandao,
		BaseFee:     basefee,
		Nonce:       [8]byte{},
		TxHash:      common.Hash{}, // TODO: this is not present in block
	}

	block := &Block{
		Hash:         common.Hash(payload.BlockHash),
		Header:       header,
		Transactions: transactions,
		// No withdrawals pre Capella
		Withdrawals: []types.Withdrawal{},
	}

	return block, nil
}

func DecodeCapellaExecutionPayload(input *api.ExecutionPayloadMsg) (*Block, error) {
	payload := new(capella.ExecutionPayload)

	if err := payload.UnmarshalSSZ(input.SszPayload); err != nil {
		fmt.Println("error unmarshalling execution payload:", err)
		return nil, err
	}

	transactions := make([]types.Transaction, len(payload.Transactions))
	for _, rawTx := range payload.Transactions {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(rawTx); err != nil {
			continue
		}
		transactions = append(transactions, *tx)
	}

	basefee := new(big.Int).SetBytes(payload.BaseFeePerGas[:])

	header := types.Header{
		ParentHash:  common.Hash(payload.ParentHash),
		UncleHash:   common.Hash([32]byte{}),
		Coinbase:    common.Address(payload.FeeRecipient),
		Root:        payload.StateRoot,
		ReceiptHash: common.Hash(payload.ReceiptsRoot),
		Bloom:       payload.LogsBloom,
		Difficulty:  nil,
		Number:      new(big.Int).SetUint64(payload.BlockNumber),
		GasLimit:    payload.GasLimit,
		GasUsed:     payload.GasUsed,
		Time:        payload.Timestamp,
		Extra:       payload.ExtraData,
		MixDigest:   payload.PrevRandao,
		BaseFee:     basefee,
		Nonce:       [8]byte{},
		TxHash:      common.Hash{}, // TODO: this is not present in block
	}

	withdrawals := make([]types.Withdrawal, len(payload.Withdrawals))
	for _, withdrawal := range payload.Withdrawals {
		wd := types.Withdrawal{
			Index:     uint64(withdrawal.Index),
			Validator: uint64(withdrawal.ValidatorIndex),
			Address:   common.Address(withdrawal.Address),
			Amount:    uint64(withdrawal.Amount),
		}

		withdrawals = append(withdrawals, wd)
	}

	block := &Block{
		Hash:         common.Hash(payload.BlockHash),
		Header:       header,
		Transactions: transactions,
		Withdrawals:  withdrawals,
	}

	return block, nil
}

func DecodeDenebExecutionPayload(input *api.ExecutionPayloadMsg) (*Block, error) {
	payload := new(deneb.ExecutionPayload)

	if err := payload.UnmarshalSSZ(input.SszPayload); err != nil {
		fmt.Println("error unmarshalling execution payload:", err)
		return nil, err
	}

	transactions := make([]types.Transaction, len(payload.Transactions))
	for _, rawTx := range payload.Transactions {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(rawTx); err != nil {
			continue
		}
		transactions = append(transactions, *tx)
	}

	header := types.Header{
		ParentHash:       common.Hash(payload.ParentHash),
		UncleHash:        common.Hash([32]byte{}),
		Coinbase:         common.Address(payload.FeeRecipient),
		Root:             common.Hash(payload.StateRoot),
		ReceiptHash:      common.Hash(payload.ReceiptsRoot),
		Bloom:            payload.LogsBloom,
		Difficulty:       nil,
		Number:           new(big.Int).SetUint64(payload.BlockNumber),
		GasLimit:         payload.GasLimit,
		GasUsed:          payload.GasUsed,
		Time:             payload.Timestamp,
		Extra:            payload.ExtraData,
		MixDigest:        payload.PrevRandao,
		BaseFee:          payload.BaseFeePerGas.ToBig(),
		BlobGasUsed:      &payload.BlobGasUsed,
		ExcessBlobGas:    &payload.ExcessBlobGas,
		Nonce:            [8]byte{},
		TxHash:           common.Hash{},  // TODO: this is not present in block
		ParentBeaconRoot: &common.Hash{}, // TODO: this is not present in block
	}

	withdrawals := make([]types.Withdrawal, len(payload.Withdrawals))
	for _, withdrawal := range payload.Withdrawals {
		wd := types.Withdrawal{
			Index:     uint64(withdrawal.Index),
			Validator: uint64(withdrawal.ValidatorIndex),
			Address:   common.Address(withdrawal.Address),
			Amount:    uint64(withdrawal.Amount),
		}

		withdrawals = append(withdrawals, wd)
	}

	block := &Block{
		Hash:         common.Hash(payload.BlockHash),
		Header:       header,
		Transactions: transactions,
		Withdrawals:  withdrawals,
	}

	return block, nil
}
