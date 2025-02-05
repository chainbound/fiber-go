// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.29.3
// source: api.proto

package api

import (
	context "context"
	eth "github.com/chainbound/fiber-go/protobuf/eth"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClient interface {
	// Opens a new transaction stream with the given filter.
	// Opens a new transaction stream with the given filter.
	SubscribeNewTxsV2(ctx context.Context, in *TxFilter, opts ...grpc.CallOption) (API_SubscribeNewTxsV2Client, error)
	// Opens a new blob transaction stream with the given filter.
	SubscribeNewBlobTxs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeNewBlobTxsClient, error)
	// Sends a signed transaction to the network.
	// Sends a signed, RLP encoded transaction to the network
	SendTransactionV2(ctx context.Context, opts ...grpc.CallOption) (API_SendTransactionV2Client, error)
	// Sends a sequence of signed transactions to the network.
	SendTransactionSequenceV2(ctx context.Context, opts ...grpc.CallOption) (API_SendTransactionSequenceV2Client, error)
	// Sends a sequence of signed, RLP encoded transactions to the network.
	SendRawTransactionSequence(ctx context.Context, opts ...grpc.CallOption) (API_SendRawTransactionSequenceClient, error)
	// Opens a stream of new execution payloads.
	SubscribeExecutionPayloadsV2(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeExecutionPayloadsV2Client, error)
	// Opens a stream of new execution payload headers.
	SubscribeExecutionHeaders(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeExecutionHeadersClient, error)
	// Opens a stream of new beacon blocks. The beacon blocks are "compacted", meaning that the
	// execution payload is not included.
	SubscribeBeaconBlocks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeBeaconBlocksClient, error)
	// Opens a stream of new beacon blocks.
	SubscribeBeaconBlocksV2(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeBeaconBlocksV2Client, error)
	// Opens a bi-directional stream for new block submissions. The client stream is used to send
	// SSZ-encoded beacon blocks, and the server stream is used to send back the state_root, slot and
	// a local timestamp as a confirmation that the block was seen and handled.
	SubmitBlockStream(ctx context.Context, opts ...grpc.CallOption) (API_SubmitBlockStreamClient, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) SubscribeNewTxsV2(ctx context.Context, in *TxFilter, opts ...grpc.CallOption) (API_SubscribeNewTxsV2Client, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[0], "/api.API/SubscribeNewTxsV2", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeNewTxsV2Client{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeNewTxsV2Client interface {
	Recv() (*TransactionWithSenderMsg, error)
	grpc.ClientStream
}

type aPISubscribeNewTxsV2Client struct {
	grpc.ClientStream
}

func (x *aPISubscribeNewTxsV2Client) Recv() (*TransactionWithSenderMsg, error) {
	m := new(TransactionWithSenderMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeNewBlobTxs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeNewBlobTxsClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[1], "/api.API/SubscribeNewBlobTxs", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeNewBlobTxsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeNewBlobTxsClient interface {
	Recv() (*TransactionWithSenderMsg, error)
	grpc.ClientStream
}

type aPISubscribeNewBlobTxsClient struct {
	grpc.ClientStream
}

func (x *aPISubscribeNewBlobTxsClient) Recv() (*TransactionWithSenderMsg, error) {
	m := new(TransactionWithSenderMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SendTransactionV2(ctx context.Context, opts ...grpc.CallOption) (API_SendTransactionV2Client, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[2], "/api.API/SendTransactionV2", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISendTransactionV2Client{stream}
	return x, nil
}

type API_SendTransactionV2Client interface {
	Send(*TransactionMsg) error
	Recv() (*TransactionResponse, error)
	grpc.ClientStream
}

type aPISendTransactionV2Client struct {
	grpc.ClientStream
}

func (x *aPISendTransactionV2Client) Send(m *TransactionMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPISendTransactionV2Client) Recv() (*TransactionResponse, error) {
	m := new(TransactionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SendTransactionSequenceV2(ctx context.Context, opts ...grpc.CallOption) (API_SendTransactionSequenceV2Client, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[3], "/api.API/SendTransactionSequenceV2", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISendTransactionSequenceV2Client{stream}
	return x, nil
}

type API_SendTransactionSequenceV2Client interface {
	Send(*TxSequenceMsgV2) error
	Recv() (*TxSequenceResponse, error)
	grpc.ClientStream
}

type aPISendTransactionSequenceV2Client struct {
	grpc.ClientStream
}

func (x *aPISendTransactionSequenceV2Client) Send(m *TxSequenceMsgV2) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPISendTransactionSequenceV2Client) Recv() (*TxSequenceResponse, error) {
	m := new(TxSequenceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SendRawTransactionSequence(ctx context.Context, opts ...grpc.CallOption) (API_SendRawTransactionSequenceClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[4], "/api.API/SendRawTransactionSequence", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISendRawTransactionSequenceClient{stream}
	return x, nil
}

type API_SendRawTransactionSequenceClient interface {
	Send(*RawTxSequenceMsg) error
	Recv() (*TxSequenceResponse, error)
	grpc.ClientStream
}

type aPISendRawTransactionSequenceClient struct {
	grpc.ClientStream
}

func (x *aPISendRawTransactionSequenceClient) Send(m *RawTxSequenceMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPISendRawTransactionSequenceClient) Recv() (*TxSequenceResponse, error) {
	m := new(TxSequenceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeExecutionPayloadsV2(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeExecutionPayloadsV2Client, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[5], "/api.API/SubscribeExecutionPayloadsV2", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeExecutionPayloadsV2Client{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeExecutionPayloadsV2Client interface {
	Recv() (*ExecutionPayloadMsg, error)
	grpc.ClientStream
}

type aPISubscribeExecutionPayloadsV2Client struct {
	grpc.ClientStream
}

func (x *aPISubscribeExecutionPayloadsV2Client) Recv() (*ExecutionPayloadMsg, error) {
	m := new(ExecutionPayloadMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeExecutionHeaders(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeExecutionHeadersClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[6], "/api.API/SubscribeExecutionHeaders", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeExecutionHeadersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeExecutionHeadersClient interface {
	Recv() (*eth.ExecutionPayloadHeader, error)
	grpc.ClientStream
}

type aPISubscribeExecutionHeadersClient struct {
	grpc.ClientStream
}

func (x *aPISubscribeExecutionHeadersClient) Recv() (*eth.ExecutionPayloadHeader, error) {
	m := new(eth.ExecutionPayloadHeader)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeBeaconBlocks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeBeaconBlocksClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[7], "/api.API/SubscribeBeaconBlocks", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeBeaconBlocksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeBeaconBlocksClient interface {
	Recv() (*eth.CompactBeaconBlock, error)
	grpc.ClientStream
}

type aPISubscribeBeaconBlocksClient struct {
	grpc.ClientStream
}

func (x *aPISubscribeBeaconBlocksClient) Recv() (*eth.CompactBeaconBlock, error) {
	m := new(eth.CompactBeaconBlock)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeBeaconBlocksV2(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (API_SubscribeBeaconBlocksV2Client, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[8], "/api.API/SubscribeBeaconBlocksV2", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeBeaconBlocksV2Client{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeBeaconBlocksV2Client interface {
	Recv() (*BeaconBlockMsg, error)
	grpc.ClientStream
}

type aPISubscribeBeaconBlocksV2Client struct {
	grpc.ClientStream
}

func (x *aPISubscribeBeaconBlocksV2Client) Recv() (*BeaconBlockMsg, error) {
	m := new(BeaconBlockMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubmitBlockStream(ctx context.Context, opts ...grpc.CallOption) (API_SubmitBlockStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[9], "/api.API/SubmitBlockStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubmitBlockStreamClient{stream}
	return x, nil
}

type API_SubmitBlockStreamClient interface {
	Send(*BlockSubmissionMsg) error
	Recv() (*BlockSubmissionResponse, error)
	grpc.ClientStream
}

type aPISubmitBlockStreamClient struct {
	grpc.ClientStream
}

func (x *aPISubmitBlockStreamClient) Send(m *BlockSubmissionMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPISubmitBlockStreamClient) Recv() (*BlockSubmissionResponse, error) {
	m := new(BlockSubmissionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// APIServer is the server API for API service.
// All implementations must embed UnimplementedAPIServer
// for forward compatibility
type APIServer interface {
	// Opens a new transaction stream with the given filter.
	// Opens a new transaction stream with the given filter.
	SubscribeNewTxsV2(*TxFilter, API_SubscribeNewTxsV2Server) error
	// Opens a new blob transaction stream with the given filter.
	SubscribeNewBlobTxs(*emptypb.Empty, API_SubscribeNewBlobTxsServer) error
	// Sends a signed transaction to the network.
	// Sends a signed, RLP encoded transaction to the network
	SendTransactionV2(API_SendTransactionV2Server) error
	// Sends a sequence of signed transactions to the network.
	SendTransactionSequenceV2(API_SendTransactionSequenceV2Server) error
	// Sends a sequence of signed, RLP encoded transactions to the network.
	SendRawTransactionSequence(API_SendRawTransactionSequenceServer) error
	// Opens a stream of new execution payloads.
	SubscribeExecutionPayloadsV2(*emptypb.Empty, API_SubscribeExecutionPayloadsV2Server) error
	// Opens a stream of new execution payload headers.
	SubscribeExecutionHeaders(*emptypb.Empty, API_SubscribeExecutionHeadersServer) error
	// Opens a stream of new beacon blocks. The beacon blocks are "compacted", meaning that the
	// execution payload is not included.
	SubscribeBeaconBlocks(*emptypb.Empty, API_SubscribeBeaconBlocksServer) error
	// Opens a stream of new beacon blocks.
	SubscribeBeaconBlocksV2(*emptypb.Empty, API_SubscribeBeaconBlocksV2Server) error
	// Opens a bi-directional stream for new block submissions. The client stream is used to send
	// SSZ-encoded beacon blocks, and the server stream is used to send back the state_root, slot and
	// a local timestamp as a confirmation that the block was seen and handled.
	SubmitBlockStream(API_SubmitBlockStreamServer) error
	mustEmbedUnimplementedAPIServer()
}

// UnimplementedAPIServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (UnimplementedAPIServer) SubscribeNewTxsV2(*TxFilter, API_SubscribeNewTxsV2Server) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeNewTxsV2 not implemented")
}
func (UnimplementedAPIServer) SubscribeNewBlobTxs(*emptypb.Empty, API_SubscribeNewBlobTxsServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeNewBlobTxs not implemented")
}
func (UnimplementedAPIServer) SendTransactionV2(API_SendTransactionV2Server) error {
	return status.Errorf(codes.Unimplemented, "method SendTransactionV2 not implemented")
}
func (UnimplementedAPIServer) SendTransactionSequenceV2(API_SendTransactionSequenceV2Server) error {
	return status.Errorf(codes.Unimplemented, "method SendTransactionSequenceV2 not implemented")
}
func (UnimplementedAPIServer) SendRawTransactionSequence(API_SendRawTransactionSequenceServer) error {
	return status.Errorf(codes.Unimplemented, "method SendRawTransactionSequence not implemented")
}
func (UnimplementedAPIServer) SubscribeExecutionPayloadsV2(*emptypb.Empty, API_SubscribeExecutionPayloadsV2Server) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeExecutionPayloadsV2 not implemented")
}
func (UnimplementedAPIServer) SubscribeExecutionHeaders(*emptypb.Empty, API_SubscribeExecutionHeadersServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeExecutionHeaders not implemented")
}
func (UnimplementedAPIServer) SubscribeBeaconBlocks(*emptypb.Empty, API_SubscribeBeaconBlocksServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeBeaconBlocks not implemented")
}
func (UnimplementedAPIServer) SubscribeBeaconBlocksV2(*emptypb.Empty, API_SubscribeBeaconBlocksV2Server) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeBeaconBlocksV2 not implemented")
}
func (UnimplementedAPIServer) SubmitBlockStream(API_SubmitBlockStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SubmitBlockStream not implemented")
}
func (UnimplementedAPIServer) mustEmbedUnimplementedAPIServer() {}

// UnsafeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServer will
// result in compilation errors.
type UnsafeAPIServer interface {
	mustEmbedUnimplementedAPIServer()
}

func RegisterAPIServer(s grpc.ServiceRegistrar, srv APIServer) {
	s.RegisterService(&API_ServiceDesc, srv)
}

func _API_SubscribeNewTxsV2_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TxFilter)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeNewTxsV2(m, &aPISubscribeNewTxsV2Server{stream})
}

type API_SubscribeNewTxsV2Server interface {
	Send(*TransactionWithSenderMsg) error
	grpc.ServerStream
}

type aPISubscribeNewTxsV2Server struct {
	grpc.ServerStream
}

func (x *aPISubscribeNewTxsV2Server) Send(m *TransactionWithSenderMsg) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubscribeNewBlobTxs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeNewBlobTxs(m, &aPISubscribeNewBlobTxsServer{stream})
}

type API_SubscribeNewBlobTxsServer interface {
	Send(*TransactionWithSenderMsg) error
	grpc.ServerStream
}

type aPISubscribeNewBlobTxsServer struct {
	grpc.ServerStream
}

func (x *aPISubscribeNewBlobTxsServer) Send(m *TransactionWithSenderMsg) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SendTransactionV2_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).SendTransactionV2(&aPISendTransactionV2Server{stream})
}

type API_SendTransactionV2Server interface {
	Send(*TransactionResponse) error
	Recv() (*TransactionMsg, error)
	grpc.ServerStream
}

type aPISendTransactionV2Server struct {
	grpc.ServerStream
}

func (x *aPISendTransactionV2Server) Send(m *TransactionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPISendTransactionV2Server) Recv() (*TransactionMsg, error) {
	m := new(TransactionMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _API_SendTransactionSequenceV2_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).SendTransactionSequenceV2(&aPISendTransactionSequenceV2Server{stream})
}

type API_SendTransactionSequenceV2Server interface {
	Send(*TxSequenceResponse) error
	Recv() (*TxSequenceMsgV2, error)
	grpc.ServerStream
}

type aPISendTransactionSequenceV2Server struct {
	grpc.ServerStream
}

func (x *aPISendTransactionSequenceV2Server) Send(m *TxSequenceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPISendTransactionSequenceV2Server) Recv() (*TxSequenceMsgV2, error) {
	m := new(TxSequenceMsgV2)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _API_SendRawTransactionSequence_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).SendRawTransactionSequence(&aPISendRawTransactionSequenceServer{stream})
}

type API_SendRawTransactionSequenceServer interface {
	Send(*TxSequenceResponse) error
	Recv() (*RawTxSequenceMsg, error)
	grpc.ServerStream
}

type aPISendRawTransactionSequenceServer struct {
	grpc.ServerStream
}

func (x *aPISendRawTransactionSequenceServer) Send(m *TxSequenceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPISendRawTransactionSequenceServer) Recv() (*RawTxSequenceMsg, error) {
	m := new(RawTxSequenceMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _API_SubscribeExecutionPayloadsV2_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeExecutionPayloadsV2(m, &aPISubscribeExecutionPayloadsV2Server{stream})
}

type API_SubscribeExecutionPayloadsV2Server interface {
	Send(*ExecutionPayloadMsg) error
	grpc.ServerStream
}

type aPISubscribeExecutionPayloadsV2Server struct {
	grpc.ServerStream
}

func (x *aPISubscribeExecutionPayloadsV2Server) Send(m *ExecutionPayloadMsg) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubscribeExecutionHeaders_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeExecutionHeaders(m, &aPISubscribeExecutionHeadersServer{stream})
}

type API_SubscribeExecutionHeadersServer interface {
	Send(*eth.ExecutionPayloadHeader) error
	grpc.ServerStream
}

type aPISubscribeExecutionHeadersServer struct {
	grpc.ServerStream
}

func (x *aPISubscribeExecutionHeadersServer) Send(m *eth.ExecutionPayloadHeader) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubscribeBeaconBlocks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeBeaconBlocks(m, &aPISubscribeBeaconBlocksServer{stream})
}

type API_SubscribeBeaconBlocksServer interface {
	Send(*eth.CompactBeaconBlock) error
	grpc.ServerStream
}

type aPISubscribeBeaconBlocksServer struct {
	grpc.ServerStream
}

func (x *aPISubscribeBeaconBlocksServer) Send(m *eth.CompactBeaconBlock) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubscribeBeaconBlocksV2_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeBeaconBlocksV2(m, &aPISubscribeBeaconBlocksV2Server{stream})
}

type API_SubscribeBeaconBlocksV2Server interface {
	Send(*BeaconBlockMsg) error
	grpc.ServerStream
}

type aPISubscribeBeaconBlocksV2Server struct {
	grpc.ServerStream
}

func (x *aPISubscribeBeaconBlocksV2Server) Send(m *BeaconBlockMsg) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubmitBlockStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).SubmitBlockStream(&aPISubmitBlockStreamServer{stream})
}

type API_SubmitBlockStreamServer interface {
	Send(*BlockSubmissionResponse) error
	Recv() (*BlockSubmissionMsg, error)
	grpc.ServerStream
}

type aPISubmitBlockStreamServer struct {
	grpc.ServerStream
}

func (x *aPISubmitBlockStreamServer) Send(m *BlockSubmissionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPISubmitBlockStreamServer) Recv() (*BlockSubmissionMsg, error) {
	m := new(BlockSubmissionMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// API_ServiceDesc is the grpc.ServiceDesc for API service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var API_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.API",
	HandlerType: (*APIServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeNewTxsV2",
			Handler:       _API_SubscribeNewTxsV2_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeNewBlobTxs",
			Handler:       _API_SubscribeNewBlobTxs_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SendTransactionV2",
			Handler:       _API_SendTransactionV2_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SendTransactionSequenceV2",
			Handler:       _API_SendTransactionSequenceV2_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SendRawTransactionSequence",
			Handler:       _API_SendRawTransactionSequence_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SubscribeExecutionPayloadsV2",
			Handler:       _API_SubscribeExecutionPayloadsV2_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeExecutionHeaders",
			Handler:       _API_SubscribeExecutionHeaders_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeBeaconBlocks",
			Handler:       _API_SubscribeBeaconBlocks_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeBeaconBlocksV2",
			Handler:       _API_SubscribeBeaconBlocksV2_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubmitBlockStream",
			Handler:       _API_SubmitBlockStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api.proto",
}
