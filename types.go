package client

import (
	"math/big"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
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
	Header       *types.Header
	Transactions []*types.Transaction
	Withdrawals  []*types.Withdrawal
}

// Helper type that wraps a signed beacon block from any of the supported hard-forks.
//
// This type will either contain a Bellatrix, Capella, Deneb, or Electra signed beacon block.
//
// DataVersion is used to indicate which type of payload is contained in the struct:
// 3: Bellatrix, 4: Capella, 5: Deneb, 6: Electra
type SignedBeaconBlock struct {
	DataVersion uint32
	Bellatrix   *bellatrix.SignedBeaconBlock
	Capella     *capella.SignedBeaconBlock
	Deneb       *deneb.SignedBeaconBlock
	Electra     *electra.SignedBeaconBlock
}

const DataVersionBellatrix uint32 = 3
const DataVersionCapella uint32 = 4
const DataVersionDeneb uint32 = 5
const DataVersionElectra uint32 = 6

func (bb *SignedBeaconBlock) StateRoot() common.Hash {
	switch bb.DataVersion {
	case DataVersionBellatrix:
		return common.Hash(bb.Bellatrix.Message.StateRoot)
	case DataVersionCapella:
		return common.Hash(bb.Capella.Message.StateRoot)
	case DataVersionDeneb:
		return common.Hash(bb.Deneb.Message.StateRoot)
	case DataVersionElectra:
		return common.Hash(bb.Electra.Message.StateRoot)
	default:
		return common.Hash{}
	}
}

func (bb *SignedBeaconBlock) Slot() phase0.Slot {
	switch bb.DataVersion {
	case DataVersionBellatrix:
		return bb.Bellatrix.Message.Slot
	case DataVersionCapella:
		return bb.Capella.Message.Slot
	case DataVersionDeneb:
		return bb.Deneb.Message.Slot
	case DataVersionElectra:
		return bb.Electra.Message.Slot
	default:
		return 0
	}
}

func (bb *SignedBeaconBlock) BlockHash() []byte {
	switch bb.DataVersion {
	case DataVersionBellatrix:
		return bb.Bellatrix.Message.Body.ETH1Data.BlockHash
	case DataVersionCapella:
		return bb.Capella.Message.Body.ETH1Data.BlockHash
	case DataVersionDeneb:
		return bb.Deneb.Message.Body.ETH1Data.BlockHash
	case DataVersionElectra:
		return bb.Electra.Message.Body.ETH1Data.BlockHash
	default:
		return []byte{}
	}
}

func DecodeBellatrixExecutionPayload(input *api.ExecutionPayloadMsg) (*Block, error) {
	payload := new(bellatrix.ExecutionPayload)

	if err := payload.UnmarshalSSZ(input.SszPayload); err != nil {
		return nil, err
	}

	transactions := make([]*types.Transaction, len(payload.Transactions))
	for i, rawTx := range payload.Transactions {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(rawTx); err != nil {
			continue
		}

		transactions[i] = tx
	}

	basefee := new(big.Int).SetBytes(reverseBytes(payload.BaseFeePerGas[:]))
	diff, _ := new(big.Int).SetString("58750003716598352816469", 10)

	header := &types.Header{
		ParentHash:  common.Hash(payload.ParentHash),
		Coinbase:    common.Address(payload.FeeRecipient),
		Root:        payload.StateRoot,
		ReceiptHash: common.Hash(payload.ReceiptsRoot),
		Bloom:       payload.LogsBloom,
		Difficulty:  diff,
		Number:      new(big.Int).SetUint64(payload.BlockNumber),
		GasLimit:    payload.GasLimit,
		GasUsed:     payload.GasUsed,
		Time:        payload.Timestamp,
		Extra:       payload.ExtraData,
		MixDigest:   payload.PrevRandao,
		BaseFee:     basefee,
		UncleHash:   common.Hash([32]byte{}), // Uncle hashes are always empty after merge
		Nonce:       [8]byte{},               // Nonce is always 0 after merge
		TxHash:      common.Hash{},           // TODO: this is not present in block
	}

	block := &Block{
		Hash:         common.Hash(payload.BlockHash),
		Header:       header,
		Transactions: transactions,
		// No withdrawals pre Capella
		Withdrawals: []*types.Withdrawal{},
	}

	return block, nil
}

func DecodeCapellaExecutionPayload(input *api.ExecutionPayloadMsg) (*Block, error) {
	payload := new(capella.ExecutionPayload)

	if err := payload.UnmarshalSSZ(input.SszPayload); err != nil {
		return nil, err
	}

	transactions := make([]*types.Transaction, len(payload.Transactions))
	for i, rawTx := range payload.Transactions {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(rawTx); err != nil {
			continue
		}

		transactions[i] = tx
	}

	basefee := new(big.Int).SetBytes(reverseBytes(payload.BaseFeePerGas[:]))
	diff, _ := new(big.Int).SetString("58750003716598352816469", 10)

	header := &types.Header{
		ParentHash:  common.Hash(payload.ParentHash),
		Coinbase:    common.Address(payload.FeeRecipient),
		Root:        payload.StateRoot,
		ReceiptHash: common.Hash(payload.ReceiptsRoot),
		Bloom:       payload.LogsBloom,
		Difficulty:  diff,
		Number:      new(big.Int).SetUint64(payload.BlockNumber),
		GasLimit:    payload.GasLimit,
		GasUsed:     payload.GasUsed,
		Time:        payload.Timestamp,
		Extra:       payload.ExtraData,
		MixDigest:   payload.PrevRandao,
		BaseFee:     basefee,
		UncleHash:   common.Hash([32]byte{}), // Uncle hashes are always empty after merge
		Nonce:       [8]byte{},               // Nonce is always 0 after merge
		TxHash:      common.Hash{},           // TODO: this is not present in block
	}

	withdrawals := make([]*types.Withdrawal, len(payload.Withdrawals))
	for i, withdrawal := range payload.Withdrawals {
		withdrawals[i] = &types.Withdrawal{
			Index:     uint64(withdrawal.Index),
			Validator: uint64(withdrawal.ValidatorIndex),
			Address:   common.Address(withdrawal.Address),
			Amount:    uint64(withdrawal.Amount),
		}
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
		return nil, err
	}

	transactions := make([]*types.Transaction, len(payload.Transactions))
	for i, rawTx := range payload.Transactions {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(rawTx); err != nil {
			continue
		}

		transactions[i] = tx
	}

	diff, _ := new(big.Int).SetString("58750003716598352816469", 10)

	header := &types.Header{
		ParentHash:       common.Hash(payload.ParentHash),
		Coinbase:         common.Address(payload.FeeRecipient),
		Root:             common.Hash(payload.StateRoot),
		ReceiptHash:      common.Hash(payload.ReceiptsRoot),
		Bloom:            payload.LogsBloom,
		Difficulty:       diff,
		Number:           new(big.Int).SetUint64(payload.BlockNumber),
		GasLimit:         payload.GasLimit,
		GasUsed:          payload.GasUsed,
		Time:             payload.Timestamp,
		Extra:            payload.ExtraData,
		MixDigest:        payload.PrevRandao,
		BaseFee:          payload.BaseFeePerGas.ToBig(),
		BlobGasUsed:      &payload.BlobGasUsed,
		ExcessBlobGas:    &payload.ExcessBlobGas,
		Nonce:            [8]byte{},               // Nonce is always 0 after merge
		UncleHash:        common.Hash([32]byte{}), // Uncle hashes are always empty after merge
		TxHash:           common.Hash{},           // TODO: this is not present in block
		ParentBeaconRoot: &common.Hash{},          // TODO: this is not present in block
	}

	withdrawals := make([]*types.Withdrawal, len(payload.Withdrawals))
	for i, withdrawal := range payload.Withdrawals {
		withdrawals[i] = &types.Withdrawal{
			Index:     uint64(withdrawal.Index),
			Validator: uint64(withdrawal.ValidatorIndex),
			Address:   common.Address(withdrawal.Address),
			Amount:    uint64(withdrawal.Amount),
		}
	}

	block := &Block{
		Hash:         common.Hash(payload.BlockHash),
		Header:       header,
		Transactions: transactions,
		Withdrawals:  withdrawals,
	}

	return block, nil
}

func DecodeElectraExecutionPayload(input *api.ExecutionPayloadMsg) (*Block, error) {
	// Electra uses the same execution payload as Deneb
	payload := new(deneb.ExecutionPayload)

	if err := payload.UnmarshalSSZ(input.SszPayload); err != nil {
		return nil, err
	}

	transactions := make([]*types.Transaction, len(payload.Transactions))
	for i, rawTx := range payload.Transactions {
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(rawTx); err != nil {
			continue
		}

		transactions[i] = tx
	}

	diff, _ := new(big.Int).SetString("58750003716598352816469", 10)

	header := &types.Header{
		ParentHash:       common.Hash(payload.ParentHash),
		Coinbase:         common.Address(payload.FeeRecipient),
		Root:             common.Hash(payload.StateRoot),
		ReceiptHash:      common.Hash(payload.ReceiptsRoot),
		Bloom:            payload.LogsBloom,
		Difficulty:       diff,
		Number:           new(big.Int).SetUint64(payload.BlockNumber),
		GasLimit:         payload.GasLimit,
		GasUsed:          payload.GasUsed,
		Time:             payload.Timestamp,
		Extra:            payload.ExtraData,
		MixDigest:        payload.PrevRandao,
		BaseFee:          payload.BaseFeePerGas.ToBig(),
		BlobGasUsed:      &payload.BlobGasUsed,
		ExcessBlobGas:    &payload.ExcessBlobGas,
		Nonce:            [8]byte{},               // Nonce is always 0 after merge
		UncleHash:        common.Hash([32]byte{}), // Uncle hashes are always empty after merge
		TxHash:           common.Hash{},           // TODO: this is not present in block
		ParentBeaconRoot: &common.Hash{},          // TODO: this is not present in block,
	}

	withdrawals := make([]*types.Withdrawal, len(payload.Withdrawals))
	for i, withdrawal := range payload.Withdrawals {
		withdrawals[i] = &types.Withdrawal{
			Index:     uint64(withdrawal.Index),
			Validator: uint64(withdrawal.ValidatorIndex),
			Address:   common.Address(withdrawal.Address),
			Amount:    uint64(withdrawal.Amount),
		}
	}

	block := &Block{
		Hash:         common.Hash(payload.BlockHash),
		Header:       header,
		Transactions: transactions,
		Withdrawals:  withdrawals,
	}

	return block, nil
}

// reverseBytes reverses a byte slice.
func reverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}
