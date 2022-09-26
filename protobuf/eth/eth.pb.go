// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: eth.proto

package eth

import (
	types "github.com/chainbound/fiber-go/protobuf/types"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BlockNumber struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to BlockNumber:
	//	*BlockNumber_Latest
	//	*BlockNumber_Pending
	//	*BlockNumber_Number
	BlockNumber isBlockNumber_BlockNumber `protobuf_oneof:"block_number"`
}

func (x *BlockNumber) Reset() {
	*x = BlockNumber{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockNumber) ProtoMessage() {}

func (x *BlockNumber) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockNumber.ProtoReflect.Descriptor instead.
func (*BlockNumber) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{0}
}

func (m *BlockNumber) GetBlockNumber() isBlockNumber_BlockNumber {
	if m != nil {
		return m.BlockNumber
	}
	return nil
}

func (x *BlockNumber) GetLatest() *emptypb.Empty {
	if x, ok := x.GetBlockNumber().(*BlockNumber_Latest); ok {
		return x.Latest
	}
	return nil
}

func (x *BlockNumber) GetPending() *emptypb.Empty {
	if x, ok := x.GetBlockNumber().(*BlockNumber_Pending); ok {
		return x.Pending
	}
	return nil
}

func (x *BlockNumber) GetNumber() uint64 {
	if x, ok := x.GetBlockNumber().(*BlockNumber_Number); ok {
		return x.Number
	}
	return 0
}

type isBlockNumber_BlockNumber interface {
	isBlockNumber_BlockNumber()
}

type BlockNumber_Latest struct {
	Latest *emptypb.Empty `protobuf:"bytes,1,opt,name=latest,proto3,oneof"`
}

type BlockNumber_Pending struct {
	Pending *emptypb.Empty `protobuf:"bytes,2,opt,name=pending,proto3,oneof"`
}

type BlockNumber_Number struct {
	Number uint64 `protobuf:"varint,3,opt,name=number,proto3,oneof"`
}

func (*BlockNumber_Latest) isBlockNumber_BlockNumber() {}

func (*BlockNumber_Pending) isBlockNumber_BlockNumber() {}

func (*BlockNumber_Number) isBlockNumber_BlockNumber() {}

type BlockId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Id:
	//	*BlockId_Hash
	//	*BlockId_Number
	Id isBlockId_Id `protobuf_oneof:"id"`
}

func (x *BlockId) Reset() {
	*x = BlockId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockId) ProtoMessage() {}

func (x *BlockId) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockId.ProtoReflect.Descriptor instead.
func (*BlockId) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{1}
}

func (m *BlockId) GetId() isBlockId_Id {
	if m != nil {
		return m.Id
	}
	return nil
}

func (x *BlockId) GetHash() *types.H256 {
	if x, ok := x.GetId().(*BlockId_Hash); ok {
		return x.Hash
	}
	return nil
}

func (x *BlockId) GetNumber() *BlockNumber {
	if x, ok := x.GetId().(*BlockId_Number); ok {
		return x.Number
	}
	return nil
}

type isBlockId_Id interface {
	isBlockId_Id()
}

type BlockId_Hash struct {
	Hash *types.H256 `protobuf:"bytes,1,opt,name=hash,proto3,oneof"`
}

type BlockId_Number struct {
	Number *BlockNumber `protobuf:"bytes,2,opt,name=number,proto3,oneof"`
}

func (*BlockId_Hash) isBlockId_Id() {}

func (*BlockId_Number) isBlockId_Id() {}

type CanonicalTransactionData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHash   *types.H256 `protobuf:"bytes,1,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	BlockNumber uint64      `protobuf:"varint,2,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	Index       uint64      `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *CanonicalTransactionData) Reset() {
	*x = CanonicalTransactionData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CanonicalTransactionData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CanonicalTransactionData) ProtoMessage() {}

func (x *CanonicalTransactionData) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CanonicalTransactionData.ProtoReflect.Descriptor instead.
func (*CanonicalTransactionData) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{2}
}

func (x *CanonicalTransactionData) GetBlockHash() *types.H256 {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *CanonicalTransactionData) GetBlockNumber() uint64 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *CanonicalTransactionData) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

type AccessListItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address *types.H160   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Slots   []*types.H256 `protobuf:"bytes,2,rep,name=slots,proto3" json:"slots,omitempty"`
}

func (x *AccessListItem) Reset() {
	*x = AccessListItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessListItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessListItem) ProtoMessage() {}

func (x *AccessListItem) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessListItem.ProtoReflect.Descriptor instead.
func (*AccessListItem) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{3}
}

func (x *AccessListItem) GetAddress() *types.H160 {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *AccessListItem) GetSlots() []*types.H256 {
	if x != nil {
		return x.Slots
	}
	return nil
}

// TODO: make eip1559 compatible + type
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	To          []byte `protobuf:"bytes,1,opt,name=to,proto3,oneof" json:"to,omitempty"`
	Gas         uint64 `protobuf:"varint,2,opt,name=gas,proto3" json:"gas,omitempty"`
	GasPrice    uint64 `protobuf:"varint,3,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`
	Hash        []byte `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Input       []byte `protobuf:"bytes,5,opt,name=input,proto3" json:"input,omitempty"`
	Nonce       uint64 `protobuf:"varint,6,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Value       []byte `protobuf:"bytes,7,opt,name=value,proto3" json:"value,omitempty"`
	From        []byte `protobuf:"bytes,8,opt,name=from,proto3" json:"from,omitempty"`
	Type        uint32 `protobuf:"varint,9,opt,name=type,proto3" json:"type,omitempty"`
	MaxFee      uint64 `protobuf:"varint,10,opt,name=max_fee,json=maxFee,proto3" json:"max_fee,omitempty"`                // = maxFeePerGas = GasFeeCap
	PriorityFee uint64 `protobuf:"varint,11,opt,name=priority_fee,json=priorityFee,proto3" json:"priority_fee,omitempty"` // = maxPriorityFeePerGas = GasTipCap
	V           uint64 `protobuf:"varint,12,opt,name=v,proto3" json:"v,omitempty"`
	R           []byte `protobuf:"bytes,13,opt,name=r,proto3" json:"r,omitempty"`
	S           []byte `protobuf:"bytes,14,opt,name=s,proto3" json:"s,omitempty"`
	ChainId     uint32 `protobuf:"varint,15,opt,name=chainId,proto3" json:"chainId,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{4}
}

func (x *Transaction) GetTo() []byte {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *Transaction) GetGas() uint64 {
	if x != nil {
		return x.Gas
	}
	return 0
}

func (x *Transaction) GetGasPrice() uint64 {
	if x != nil {
		return x.GasPrice
	}
	return 0
}

func (x *Transaction) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Transaction) GetInput() []byte {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *Transaction) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *Transaction) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *Transaction) GetFrom() []byte {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Transaction) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Transaction) GetMaxFee() uint64 {
	if x != nil {
		return x.MaxFee
	}
	return 0
}

func (x *Transaction) GetPriorityFee() uint64 {
	if x != nil {
		return x.PriorityFee
	}
	return 0
}

func (x *Transaction) GetV() uint64 {
	if x != nil {
		return x.V
	}
	return 0
}

func (x *Transaction) GetR() []byte {
	if x != nil {
		return x.R
	}
	return nil
}

func (x *Transaction) GetS() []byte {
	if x != nil {
		return x.S
	}
	return nil
}

func (x *Transaction) GetChainId() uint32 {
	if x != nil {
		return x.ChainId
	}
	return 0
}

type StoredTransaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CanonicalData *CanonicalTransactionData `protobuf:"bytes,1,opt,name=canonical_data,json=canonicalData,proto3,oneof" json:"canonical_data,omitempty"`
	Transaction   *Transaction              `protobuf:"bytes,2,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (x *StoredTransaction) Reset() {
	*x = StoredTransaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoredTransaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoredTransaction) ProtoMessage() {}

func (x *StoredTransaction) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoredTransaction.ProtoReflect.Descriptor instead.
func (*StoredTransaction) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{5}
}

func (x *StoredTransaction) GetCanonicalData() *CanonicalTransactionData {
	if x != nil {
		return x.CanonicalData
	}
	return nil
}

func (x *StoredTransaction) GetTransaction() *Transaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Number          uint64  `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Hash            []byte  `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	ParentHash      []byte  `protobuf:"bytes,3,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	Nonce           uint64  `protobuf:"varint,4,opt,name=nonce,proto3" json:"nonce,omitempty"`
	UncleHash       []byte  `protobuf:"bytes,5,opt,name=uncle_hash,json=uncleHash,proto3" json:"uncle_hash,omitempty"`
	StateRoot       []byte  `protobuf:"bytes,6,opt,name=state_root,json=stateRoot,proto3" json:"state_root,omitempty"`
	ReceiptRoot     []byte  `protobuf:"bytes,7,opt,name=receipt_root,json=receiptRoot,proto3" json:"receipt_root,omitempty"`
	Coinbase        []byte  `protobuf:"bytes,8,opt,name=coinbase,proto3" json:"coinbase,omitempty"`
	Difficulty      uint64  `protobuf:"varint,9,opt,name=difficulty,proto3" json:"difficulty,omitempty"`
	TotalDifficulty *uint64 `protobuf:"varint,10,opt,name=total_difficulty,json=totalDifficulty,proto3,oneof" json:"total_difficulty,omitempty"`
	ExtraData       []byte  `protobuf:"bytes,11,opt,name=extra_data,json=extraData,proto3,oneof" json:"extra_data,omitempty"`
	Size            *uint64 `protobuf:"varint,12,opt,name=size,proto3,oneof" json:"size,omitempty"`
	GasLimit        uint64  `protobuf:"varint,13,opt,name=gas_limit,json=gasLimit,proto3" json:"gas_limit,omitempty"`
	GasUsed         uint64  `protobuf:"varint,14,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Timestamp       uint64  `protobuf:"varint,15,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{6}
}

func (x *Header) GetNumber() uint64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Header) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Header) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *Header) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *Header) GetUncleHash() []byte {
	if x != nil {
		return x.UncleHash
	}
	return nil
}

func (x *Header) GetStateRoot() []byte {
	if x != nil {
		return x.StateRoot
	}
	return nil
}

func (x *Header) GetReceiptRoot() []byte {
	if x != nil {
		return x.ReceiptRoot
	}
	return nil
}

func (x *Header) GetCoinbase() []byte {
	if x != nil {
		return x.Coinbase
	}
	return nil
}

func (x *Header) GetDifficulty() uint64 {
	if x != nil {
		return x.Difficulty
	}
	return 0
}

func (x *Header) GetTotalDifficulty() uint64 {
	if x != nil && x.TotalDifficulty != nil {
		return *x.TotalDifficulty
	}
	return 0
}

func (x *Header) GetExtraData() []byte {
	if x != nil {
		return x.ExtraData
	}
	return nil
}

func (x *Header) GetSize() uint64 {
	if x != nil && x.Size != nil {
		return *x.Size
	}
	return 0
}

func (x *Header) GetGasLimit() uint64 {
	if x != nil {
		return x.GasLimit
	}
	return 0
}

func (x *Header) GetGasUsed() uint64 {
	if x != nil {
		return x.GasUsed
	}
	return 0
}

func (x *Header) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type Body struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transactions []*Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	Ommers       [][]byte       `protobuf:"bytes,2,rep,name=ommers,proto3" json:"ommers,omitempty"`
}

func (x *Body) Reset() {
	*x = Body{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Body) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Body) ProtoMessage() {}

func (x *Body) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Body.ProtoReflect.Descriptor instead.
func (*Body) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{7}
}

func (x *Body) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *Body) GetOmmers() [][]byte {
	if x != nil {
		return x.Ommers
	}
	return nil
}

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header *Header `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Body   *Body   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eth_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_eth_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_eth_proto_rawDescGZIP(), []int{8}
}

func (x *Block) GetHeader() *Header {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Block) GetBody() *Body {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_eth_proto protoreflect.FileDescriptor

var file_eth_proto_rawDesc = []byte{
	0x0a, 0x09, 0x65, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x65, 0x74, 0x68,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x01, 0x0a, 0x0b, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x06, 0x6c, 0x61,
	0x74, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x48, 0x00, 0x52, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x07,
	0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x48, 0x00, 0x52, 0x07, 0x70, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67,
	0x12, 0x18, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x48, 0x00, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x42, 0x0e, 0x0a, 0x0c, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x5e, 0x0a, 0x07, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x48, 0x32, 0x35, 0x36,
	0x48, 0x00, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x2a, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x42,
	0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x48, 0x00, 0x52, 0x06, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x42, 0x04, 0x0a, 0x02, 0x69, 0x64, 0x22, 0x7f, 0x0a, 0x18, 0x43, 0x61,
	0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2a, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x48, 0x32, 0x35, 0x36, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x5a, 0x0a, 0x0e, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x25, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x48, 0x31, 0x36, 0x30, 0x52, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x21, 0x0a, 0x05, 0x73, 0x6c, 0x6f, 0x74, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x48, 0x32, 0x35, 0x36,
	0x52, 0x05, 0x73, 0x6c, 0x6f, 0x74, 0x73, 0x22, 0xd6, 0x02, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x02, 0x74, 0x6f, 0x88, 0x01, 0x01, 0x12, 0x10, 0x0a, 0x03,
	0x67, 0x61, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x67, 0x61, 0x73, 0x12, 0x1b,
	0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x08, 0x67, 0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x6d, 0x61, 0x78,
	0x5f, 0x66, 0x65, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6d, 0x61, 0x78, 0x46,
	0x65, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x5f, 0x66,
	0x65, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x46, 0x65, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x01, 0x76, 0x12, 0x0c, 0x0a, 0x01, 0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01,
	0x72, 0x12, 0x0c, 0x0a, 0x01, 0x73, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x74, 0x6f,
	0x22, 0xa5, 0x01, 0x0a, 0x11, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x49, 0x0a, 0x0e, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69,
	0x63, 0x61, 0x6c, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x65, 0x74, 0x68, 0x2e, 0x43, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52,
	0x0d, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x88, 0x01,
	0x01, 0x12, 0x32, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69,
	0x63, 0x61, 0x6c, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x22, 0xf8, 0x03, 0x0a, 0x06, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x6e, 0x63, 0x6c, 0x65, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x75, 0x6e, 0x63, 0x6c,
	0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72,
	0x6f, 0x6f, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x6f, 0x6f, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x5f,
	0x72, 0x6f, 0x6f, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x65,
	0x69, 0x70, 0x74, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x69, 0x6e, 0x62,
	0x61, 0x73, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6f, 0x69, 0x6e, 0x62,
	0x61, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74,
	0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75,
	0x6c, 0x74, 0x79, 0x12, 0x2e, 0x0a, 0x10, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x64, 0x69, 0x66,
	0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52,
	0x0f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x44, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79,
	0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x65, 0x78, 0x74, 0x72, 0x61, 0x5f, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x01, 0x52, 0x09, 0x65, 0x78, 0x74, 0x72, 0x61,
	0x44, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x04, 0x48, 0x02, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x67, 0x61, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x19, 0x0a,
	0x08, 0x67, 0x61, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x07, 0x67, 0x61, 0x73, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x5f, 0x64, 0x69, 0x66, 0x66, 0x69, 0x63, 0x75, 0x6c, 0x74, 0x79, 0x42, 0x0d, 0x0a, 0x0b, 0x5f,
	0x65, 0x78, 0x74, 0x72, 0x61, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x22, 0x54, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x34, 0x0a, 0x0c, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0c, 0x52, 0x06, 0x6f, 0x6d, 0x6d, 0x65, 0x72, 0x73, 0x22, 0x4b, 0x0a, 0x05, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x12, 0x23, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52,
	0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x42, 0x6f, 0x64, 0x79,
	0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_eth_proto_rawDescOnce sync.Once
	file_eth_proto_rawDescData = file_eth_proto_rawDesc
)

func file_eth_proto_rawDescGZIP() []byte {
	file_eth_proto_rawDescOnce.Do(func() {
		file_eth_proto_rawDescData = protoimpl.X.CompressGZIP(file_eth_proto_rawDescData)
	})
	return file_eth_proto_rawDescData
}

var file_eth_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_eth_proto_goTypes = []interface{}{
	(*BlockNumber)(nil),              // 0: eth.BlockNumber
	(*BlockId)(nil),                  // 1: eth.BlockId
	(*CanonicalTransactionData)(nil), // 2: eth.CanonicalTransactionData
	(*AccessListItem)(nil),           // 3: eth.AccessListItem
	(*Transaction)(nil),              // 4: eth.Transaction
	(*StoredTransaction)(nil),        // 5: eth.StoredTransaction
	(*Header)(nil),                   // 6: eth.Header
	(*Body)(nil),                     // 7: eth.Body
	(*Block)(nil),                    // 8: eth.Block
	(*emptypb.Empty)(nil),            // 9: google.protobuf.Empty
	(*types.H256)(nil),               // 10: types.H256
	(*types.H160)(nil),               // 11: types.H160
}
var file_eth_proto_depIdxs = []int32{
	9,  // 0: eth.BlockNumber.latest:type_name -> google.protobuf.Empty
	9,  // 1: eth.BlockNumber.pending:type_name -> google.protobuf.Empty
	10, // 2: eth.BlockId.hash:type_name -> types.H256
	0,  // 3: eth.BlockId.number:type_name -> eth.BlockNumber
	10, // 4: eth.CanonicalTransactionData.block_hash:type_name -> types.H256
	11, // 5: eth.AccessListItem.address:type_name -> types.H160
	10, // 6: eth.AccessListItem.slots:type_name -> types.H256
	2,  // 7: eth.StoredTransaction.canonical_data:type_name -> eth.CanonicalTransactionData
	4,  // 8: eth.StoredTransaction.transaction:type_name -> eth.Transaction
	4,  // 9: eth.Body.transactions:type_name -> eth.Transaction
	6,  // 10: eth.Block.header:type_name -> eth.Header
	7,  // 11: eth.Block.body:type_name -> eth.Body
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_eth_proto_init() }
func file_eth_proto_init() {
	if File_eth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_eth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockNumber); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CanonicalTransactionData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessListItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoredTransaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Body); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_eth_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_eth_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*BlockNumber_Latest)(nil),
		(*BlockNumber_Pending)(nil),
		(*BlockNumber_Number)(nil),
	}
	file_eth_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*BlockId_Hash)(nil),
		(*BlockId_Number)(nil),
	}
	file_eth_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_eth_proto_msgTypes[5].OneofWrappers = []interface{}{}
	file_eth_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_eth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_eth_proto_goTypes,
		DependencyIndexes: file_eth_proto_depIdxs,
		MessageInfos:      file_eth_proto_msgTypes,
	}.Build()
	File_eth_proto = out.File
	file_eth_proto_rawDesc = nil
	file_eth_proto_goTypes = nil
	file_eth_proto_depIdxs = nil
}
