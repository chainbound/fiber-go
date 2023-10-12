package client

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/chainbound/fiber-go/protobuf/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// ==================== TRANSACTION ====================

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

// TxToProto converts a go-ethereum transaction to a protobuf transaction.
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
	}

	value := new(big.Int)
	if len(proto.Value) > 0 {
		if err := rlp.Decode(bytes.NewReader(proto.Value), value); err != nil {
			fmt.Println(err)
		}
	} else {
		value = big.NewInt(0)
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
		Value:       value,
		Input:       proto.Input,
		V:           proto.V,
		R:           proto.R,
		S:           proto.S,
		AccessList:  acl,
	}
}

// ==================== EXECUTION PAYLOAD ====================

type ExecutionPayloadHeader struct {
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

type ExecutionPayload struct {
	Header       *ExecutionPayloadHeader
	Transactions []*Transaction
}

func ProtoToHeader(proto *eth.ExecutionPayloadHeader) *ExecutionPayloadHeader {
	return &ExecutionPayloadHeader{
		Number:        proto.BlockNumber,
		Hash:          common.BytesToHash(proto.BlockHash),
		ParentHash:    common.BytesToHash(proto.ParentHash),
		StateRoot:     common.BytesToHash(proto.StateRoot),
		ReceiptRoot:   common.BytesToHash(proto.ReceiptsRoot),
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

func ProtoToBlock(proto *eth.ExecutionPayload) *ExecutionPayload {
	header := proto.Header
	txs := make([]*Transaction, len(proto.Transactions))
	for i, proto := range proto.Transactions {
		txs[i] = ProtoToTx(proto)
	}

	return &ExecutionPayload{
		Header:       ProtoToHeader(header),
		Transactions: txs,
	}
}

// ==================== BEACON BLOCK ====================

type BeaconBlock struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    common.Hash
	StateRoot     common.Hash
	Body          *BeaconBlockBody
}

type BeaconBlockBody struct {
	RandaoReveal              common.Hash
	Eth1Data                  *Eth1Data
	Graffiti                  common.Hash
	ProposerSlashingsList     []ProposerSlashing
	AttesterSlashingsList     []AttesterSlashing
	AttestationsList          []Attestation
	DepositsList              []Deposit
	VoluntaryExitsList        []VoluntaryExit
	SyncAggregate             *SyncAggregate
	BlsToExecutionChangesList []ExecutionChange
}

type Eth1Data struct {
	DepositRoot  common.Hash
	DepositCount uint64
	BlockHash    common.Hash
}

type ProposerSlashing struct {
	Header1 *SignedBeaconBlockHeader
	Header2 *SignedBeaconBlockHeader
}

type SignedBeaconBlockHeader struct {
	Message   *BeaconBlockHeader
	Signature common.Hash
}

type BeaconBlockHeader struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    common.Hash
	StateRoot     common.Hash
	BodyRoot      common.Hash
}

type AttesterSlashing struct {
	Attestation1 *IndexedAttestation
	Attestation2 *IndexedAttestation
}

type IndexedAttestation struct {
	AttestingIndicesList []uint64
	Data                 *AttestationData
	Signature            common.Hash
}

type AttestationData struct {
	Slot            uint64
	Index           uint64
	BeaconBlockRoot common.Hash
	Source          *Checkpoint
	Target          *Checkpoint
}

type Checkpoint struct {
	Epoch uint64
	Root  common.Hash
}

type Attestation struct {
	AggregationBits common.Hash
	Data            *AttestationData
	Signature       common.Hash
}

type Deposit struct {
	ProofList []common.Hash
	Data      *DepositData
}

type DepositData struct {
	Pubkey                common.Hash
	WithdrawalCredentials common.Hash
	Amount                uint64
	Signature             common.Hash
}

type VoluntaryExit struct {
	Message   *VoluntaryExitMessage
	Signature common.Hash
}

type VoluntaryExitMessage struct {
	Epoch          uint64
	ValidatorIndex uint64
}

type SyncAggregate struct {
	SyncCommitteeBits      common.Hash
	SyncCommitteeSignature common.Hash
}

type ExecutionChange struct {
	Message   *ExecutionChangeMessage
	Signature common.Hash
}

type ExecutionChangeMessage struct {
	ValidatorIndex     uint64
	FromBlsPubkey      common.Hash
	ToExecutionAddress common.Hash
}

func ProtoToBeaconBlock(block *eth.CompactBeaconBlock) *BeaconBlock {
	body := block.GetBody()

	beacon := BeaconBlock{
		Slot:          block.GetSlot(),
		ProposerIndex: block.GetProposerIndex(),
		ParentRoot:    common.BytesToHash(block.GetParentRoot()),
		StateRoot:     common.BytesToHash(block.GetStateRoot()),
		Body: &BeaconBlockBody{
			RandaoReveal: common.BytesToHash(body.GetRandaoReveal()),
			Eth1Data: &Eth1Data{
				DepositRoot:  common.BytesToHash(body.GetEth1Data().GetDepositRoot()),
				DepositCount: body.GetEth1Data().GetDepositCount(),
				BlockHash:    common.BytesToHash(body.GetEth1Data().GetBlockHash()),
			},
			Graffiti: common.BytesToHash(body.GetGraffiti()),
		},
	}

	for _, slashing := range body.GetProposerSlashings() {
		beacon.Body.ProposerSlashingsList = append(beacon.Body.ProposerSlashingsList, fromProtoProposerSlashing(slashing))
	}
	for _, slashing := range body.GetAttesterSlashings() {
		beacon.Body.AttesterSlashingsList = append(beacon.Body.AttesterSlashingsList, fromProtoAttesterSlashing(slashing))
	}
	for _, attestation := range body.GetAttestations() {
		beacon.Body.AttestationsList = append(beacon.Body.AttestationsList, fromProtoAttestation(attestation))
	}
	for _, deposit := range body.GetDeposits() {
		beacon.Body.DepositsList = append(beacon.Body.DepositsList, fromProtoDeposit(deposit))
	}
	for _, exit := range body.GetVoluntaryExits() {
		beacon.Body.VoluntaryExitsList = append(beacon.Body.VoluntaryExitsList, fromProtoVoluntaryExit(exit))
	}

	syncAggregate := body.GetSyncAggregate()
	if syncAggregate != nil {
		beacon.Body.SyncAggregate = &SyncAggregate{
			SyncCommitteeBits:      common.BytesToHash(syncAggregate.GetSyncCommitteeBits()),
			SyncCommitteeSignature: common.BytesToHash(syncAggregate.GetSyncCommitteeSignature()),
		}
	}

	for _, change := range body.GetBlsToExecutionChanges() {
		beacon.Body.BlsToExecutionChangesList = append(beacon.Body.BlsToExecutionChangesList, ExecutionChange{
			Message: &ExecutionChangeMessage{
				ValidatorIndex:     change.GetMessage().GetValidatorIndex(),
				FromBlsPubkey:      common.BytesToHash(change.GetMessage().GetFromBlsPubkey()),
				ToExecutionAddress: common.BytesToHash(change.GetMessage().GetToExecutionAddress()),
			},
			Signature: common.BytesToHash(change.GetSignature()),
		})
	}

	return &beacon
}

func fromProtoProposerSlashing(slashing *eth.ProposerSlashing) ProposerSlashing {
	return ProposerSlashing{
		Header1: fromProtoSignedBeaconBlockHeader(slashing.GetHeader_1()),
		Header2: fromProtoSignedBeaconBlockHeader(slashing.GetHeader_2()),
	}
}

func fromProtoSignedBeaconBlockHeader(header *eth.SignedBeaconBlockHeader) *SignedBeaconBlockHeader {
	return &SignedBeaconBlockHeader{
		Message:   fromProtoBeaconBlockHeader(header.GetMessage()),
		Signature: common.BytesToHash(header.GetSignature()),
	}
}

func fromProtoBeaconBlockHeader(header *eth.BeaconBlockHeader) *BeaconBlockHeader {
	return &BeaconBlockHeader{
		Slot:          header.GetSlot(),
		ProposerIndex: header.GetProposerIndex(),
		ParentRoot:    common.BytesToHash(header.GetParentRoot()),
		StateRoot:     common.BytesToHash(header.GetStateRoot()),
		BodyRoot:      common.BytesToHash(header.GetBodyRoot()),
	}
}

func fromProtoAttesterSlashing(slashing *eth.AttesterSlashing) AttesterSlashing {
	return AttesterSlashing{
		Attestation1: fromProtoIndexedAttestation(slashing.GetAttestation_1()),
		Attestation2: fromProtoIndexedAttestation(slashing.GetAttestation_2()),
	}
}

func fromProtoIndexedAttestation(attestation *eth.IndexedAttestation) *IndexedAttestation {
	return &IndexedAttestation{
		AttestingIndicesList: attestation.GetAttestingIndices(),
		Data:                 fromProtoAttestationData(attestation.GetData()),
		Signature:            common.BytesToHash(attestation.GetSignature()),
	}
}

func fromProtoAttestationData(data *eth.AttestationData) *AttestationData {
	return &AttestationData{
		Slot:            data.GetSlot(),
		Index:           data.GetIndex(),
		BeaconBlockRoot: common.BytesToHash(data.GetBeaconBlockRoot()),
		Source:          fromProtoCheckpoint(data.GetSource()),
		Target:          fromProtoCheckpoint(data.GetTarget()),
	}
}

func fromProtoCheckpoint(checkpoint *eth.Checkpoint) *Checkpoint {
	return &Checkpoint{
		Epoch: checkpoint.GetEpoch(),
		Root:  common.BytesToHash(checkpoint.GetRoot()),
	}
}

func fromProtoAttestation(attestation *eth.Attestation) Attestation {
	return Attestation{
		AggregationBits: common.BytesToHash(attestation.GetAggregationBits()),
		Data:            fromProtoAttestationData(attestation.GetData()),
		Signature:       common.BytesToHash(attestation.GetSignature()),
	}
}

func fromProtoDeposit(deposit *eth.Deposit) Deposit {
	var proofs []common.Hash
	for _, proof := range deposit.GetProof() {
		proofs = append(proofs, common.BytesToHash(proof))
	}

	return Deposit{
		ProofList: proofs,
		Data:      fromProtoDepositData(deposit.GetData()),
	}
}

func fromProtoDepositData(data *eth.DepositData) *DepositData {
	return &DepositData{
		Pubkey:                common.BytesToHash(data.GetPubkey()),
		WithdrawalCredentials: common.BytesToHash(data.GetWithdrawalCredentials()),
		Amount:                data.GetAmount(),
		Signature:             common.BytesToHash(data.GetSignature()),
	}
}

func fromProtoVoluntaryExit(exit *eth.SignedVoluntaryExit) VoluntaryExit {
	return VoluntaryExit{
		Message: &VoluntaryExitMessage{
			Epoch:          exit.GetMessage().GetEpoch(),
			ValidatorIndex: exit.GetMessage().GetValidatorIndex(),
		},
		Signature: common.BytesToHash(exit.GetSignature()),
	}
}
