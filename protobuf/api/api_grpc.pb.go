// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: api.proto

package api

import (
	context "context"
	eth "github.com/chainbound/fiber-go/protobuf/eth"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClient interface {
	SubscribeNewTxs(ctx context.Context, in *TxFilter, opts ...grpc.CallOption) (API_SubscribeNewTxsClient, error)
	SubscribeNewTxsV2(ctx context.Context, in *TxFilterV2, opts ...grpc.CallOption) (API_SubscribeNewTxsV2Client, error)
	SubscribeNewBlocks(ctx context.Context, in *BlockFilter, opts ...grpc.CallOption) (API_SubscribeNewBlocksClient, error)
	SendTransaction(ctx context.Context, in *eth.Transaction, opts ...grpc.CallOption) (*TransactionResponse, error)
	SendRawTransaction(ctx context.Context, in *RawTxMsg, opts ...grpc.CallOption) (*TransactionResponse, error)
	// Backrun is the RPC method for backrunning a transaction.
	Backrun(ctx context.Context, in *BackrunMsg, opts ...grpc.CallOption) (*TransactionResponse, error)
	RawBackrun(ctx context.Context, in *RawBackrunMsg, opts ...grpc.CallOption) (*TransactionResponse, error)
	SendTransactionStream(ctx context.Context, opts ...grpc.CallOption) (API_SendTransactionStreamClient, error)
	SendRawTransactionStream(ctx context.Context, opts ...grpc.CallOption) (API_SendRawTransactionStreamClient, error)
	// Backrun is the RPC method for backrunning a transaction.
	BackrunStream(ctx context.Context, opts ...grpc.CallOption) (API_BackrunStreamClient, error)
	RawBackrunStream(ctx context.Context, opts ...grpc.CallOption) (API_RawBackrunStreamClient, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) SubscribeNewTxs(ctx context.Context, in *TxFilter, opts ...grpc.CallOption) (API_SubscribeNewTxsClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[0], "/api.API/SubscribeNewTxs", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeNewTxsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeNewTxsClient interface {
	Recv() (*eth.Transaction, error)
	grpc.ClientStream
}

type aPISubscribeNewTxsClient struct {
	grpc.ClientStream
}

func (x *aPISubscribeNewTxsClient) Recv() (*eth.Transaction, error) {
	m := new(eth.Transaction)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeNewTxsV2(ctx context.Context, in *TxFilterV2, opts ...grpc.CallOption) (API_SubscribeNewTxsV2Client, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[1], "/api.API/SubscribeNewTxsV2", opts...)
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
	Recv() (*eth.Transaction, error)
	grpc.ClientStream
}

type aPISubscribeNewTxsV2Client struct {
	grpc.ClientStream
}

func (x *aPISubscribeNewTxsV2Client) Recv() (*eth.Transaction, error) {
	m := new(eth.Transaction)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeNewBlocks(ctx context.Context, in *BlockFilter, opts ...grpc.CallOption) (API_SubscribeNewBlocksClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[2], "/api.API/SubscribeNewBlocks", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeNewBlocksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeNewBlocksClient interface {
	Recv() (*eth.Block, error)
	grpc.ClientStream
}

type aPISubscribeNewBlocksClient struct {
	grpc.ClientStream
}

func (x *aPISubscribeNewBlocksClient) Recv() (*eth.Block, error) {
	m := new(eth.Block)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SendTransaction(ctx context.Context, in *eth.Transaction, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/api.API/SendTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) SendRawTransaction(ctx context.Context, in *RawTxMsg, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/api.API/SendRawTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) Backrun(ctx context.Context, in *BackrunMsg, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/api.API/Backrun", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) RawBackrun(ctx context.Context, in *RawBackrunMsg, opts ...grpc.CallOption) (*TransactionResponse, error) {
	out := new(TransactionResponse)
	err := c.cc.Invoke(ctx, "/api.API/RawBackrun", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) SendTransactionStream(ctx context.Context, opts ...grpc.CallOption) (API_SendTransactionStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[3], "/api.API/SendTransactionStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISendTransactionStreamClient{stream}
	return x, nil
}

type API_SendTransactionStreamClient interface {
	Send(*eth.Transaction) error
	Recv() (*TransactionResponse, error)
	grpc.ClientStream
}

type aPISendTransactionStreamClient struct {
	grpc.ClientStream
}

func (x *aPISendTransactionStreamClient) Send(m *eth.Transaction) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPISendTransactionStreamClient) Recv() (*TransactionResponse, error) {
	m := new(TransactionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SendRawTransactionStream(ctx context.Context, opts ...grpc.CallOption) (API_SendRawTransactionStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[4], "/api.API/SendRawTransactionStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISendRawTransactionStreamClient{stream}
	return x, nil
}

type API_SendRawTransactionStreamClient interface {
	Send(*RawTxMsg) error
	Recv() (*TransactionResponse, error)
	grpc.ClientStream
}

type aPISendRawTransactionStreamClient struct {
	grpc.ClientStream
}

func (x *aPISendRawTransactionStreamClient) Send(m *RawTxMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPISendRawTransactionStreamClient) Recv() (*TransactionResponse, error) {
	m := new(TransactionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) BackrunStream(ctx context.Context, opts ...grpc.CallOption) (API_BackrunStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[5], "/api.API/BackrunStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIBackrunStreamClient{stream}
	return x, nil
}

type API_BackrunStreamClient interface {
	Send(*BackrunMsg) error
	Recv() (*TransactionResponse, error)
	grpc.ClientStream
}

type aPIBackrunStreamClient struct {
	grpc.ClientStream
}

func (x *aPIBackrunStreamClient) Send(m *BackrunMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPIBackrunStreamClient) Recv() (*TransactionResponse, error) {
	m := new(TransactionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) RawBackrunStream(ctx context.Context, opts ...grpc.CallOption) (API_RawBackrunStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[6], "/api.API/RawBackrunStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIRawBackrunStreamClient{stream}
	return x, nil
}

type API_RawBackrunStreamClient interface {
	Send(*RawBackrunMsg) error
	Recv() (*TransactionResponse, error)
	grpc.ClientStream
}

type aPIRawBackrunStreamClient struct {
	grpc.ClientStream
}

func (x *aPIRawBackrunStreamClient) Send(m *RawBackrunMsg) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aPIRawBackrunStreamClient) Recv() (*TransactionResponse, error) {
	m := new(TransactionResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// APIServer is the server API for API service.
// All implementations must embed UnimplementedAPIServer
// for forward compatibility
type APIServer interface {
	SubscribeNewTxs(*TxFilter, API_SubscribeNewTxsServer) error
	SubscribeNewTxsV2(*TxFilterV2, API_SubscribeNewTxsV2Server) error
	SubscribeNewBlocks(*BlockFilter, API_SubscribeNewBlocksServer) error
	SendTransaction(context.Context, *eth.Transaction) (*TransactionResponse, error)
	SendRawTransaction(context.Context, *RawTxMsg) (*TransactionResponse, error)
	// Backrun is the RPC method for backrunning a transaction.
	Backrun(context.Context, *BackrunMsg) (*TransactionResponse, error)
	RawBackrun(context.Context, *RawBackrunMsg) (*TransactionResponse, error)
	SendTransactionStream(API_SendTransactionStreamServer) error
	SendRawTransactionStream(API_SendRawTransactionStreamServer) error
	// Backrun is the RPC method for backrunning a transaction.
	BackrunStream(API_BackrunStreamServer) error
	RawBackrunStream(API_RawBackrunStreamServer) error
	mustEmbedUnimplementedAPIServer()
}

// UnimplementedAPIServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (UnimplementedAPIServer) SubscribeNewTxs(*TxFilter, API_SubscribeNewTxsServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeNewTxs not implemented")
}
func (UnimplementedAPIServer) SubscribeNewTxsV2(*TxFilterV2, API_SubscribeNewTxsV2Server) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeNewTxsV2 not implemented")
}
func (UnimplementedAPIServer) SubscribeNewBlocks(*BlockFilter, API_SubscribeNewBlocksServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeNewBlocks not implemented")
}
func (UnimplementedAPIServer) SendTransaction(context.Context, *eth.Transaction) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTransaction not implemented")
}
func (UnimplementedAPIServer) SendRawTransaction(context.Context, *RawTxMsg) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRawTransaction not implemented")
}
func (UnimplementedAPIServer) Backrun(context.Context, *BackrunMsg) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Backrun not implemented")
}
func (UnimplementedAPIServer) RawBackrun(context.Context, *RawBackrunMsg) (*TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawBackrun not implemented")
}
func (UnimplementedAPIServer) SendTransactionStream(API_SendTransactionStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SendTransactionStream not implemented")
}
func (UnimplementedAPIServer) SendRawTransactionStream(API_SendRawTransactionStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SendRawTransactionStream not implemented")
}
func (UnimplementedAPIServer) BackrunStream(API_BackrunStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method BackrunStream not implemented")
}
func (UnimplementedAPIServer) RawBackrunStream(API_RawBackrunStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method RawBackrunStream not implemented")
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

func _API_SubscribeNewTxs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TxFilter)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeNewTxs(m, &aPISubscribeNewTxsServer{stream})
}

type API_SubscribeNewTxsServer interface {
	Send(*eth.Transaction) error
	grpc.ServerStream
}

type aPISubscribeNewTxsServer struct {
	grpc.ServerStream
}

func (x *aPISubscribeNewTxsServer) Send(m *eth.Transaction) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubscribeNewTxsV2_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TxFilterV2)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeNewTxsV2(m, &aPISubscribeNewTxsV2Server{stream})
}

type API_SubscribeNewTxsV2Server interface {
	Send(*eth.Transaction) error
	grpc.ServerStream
}

type aPISubscribeNewTxsV2Server struct {
	grpc.ServerStream
}

func (x *aPISubscribeNewTxsV2Server) Send(m *eth.Transaction) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubscribeNewBlocks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BlockFilter)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeNewBlocks(m, &aPISubscribeNewBlocksServer{stream})
}

type API_SubscribeNewBlocksServer interface {
	Send(*eth.Block) error
	grpc.ServerStream
}

type aPISubscribeNewBlocksServer struct {
	grpc.ServerStream
}

func (x *aPISubscribeNewBlocksServer) Send(m *eth.Block) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SendTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(eth.Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).SendTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.API/SendTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).SendTransaction(ctx, req.(*eth.Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_SendRawTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawTxMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).SendRawTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.API/SendRawTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).SendRawTransaction(ctx, req.(*RawTxMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_Backrun_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackrunMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Backrun(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.API/Backrun",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Backrun(ctx, req.(*BackrunMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_RawBackrun_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawBackrunMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).RawBackrun(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.API/RawBackrun",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).RawBackrun(ctx, req.(*RawBackrunMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_SendTransactionStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).SendTransactionStream(&aPISendTransactionStreamServer{stream})
}

type API_SendTransactionStreamServer interface {
	Send(*TransactionResponse) error
	Recv() (*eth.Transaction, error)
	grpc.ServerStream
}

type aPISendTransactionStreamServer struct {
	grpc.ServerStream
}

func (x *aPISendTransactionStreamServer) Send(m *TransactionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPISendTransactionStreamServer) Recv() (*eth.Transaction, error) {
	m := new(eth.Transaction)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _API_SendRawTransactionStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).SendRawTransactionStream(&aPISendRawTransactionStreamServer{stream})
}

type API_SendRawTransactionStreamServer interface {
	Send(*TransactionResponse) error
	Recv() (*RawTxMsg, error)
	grpc.ServerStream
}

type aPISendRawTransactionStreamServer struct {
	grpc.ServerStream
}

func (x *aPISendRawTransactionStreamServer) Send(m *TransactionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPISendRawTransactionStreamServer) Recv() (*RawTxMsg, error) {
	m := new(RawTxMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _API_BackrunStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).BackrunStream(&aPIBackrunStreamServer{stream})
}

type API_BackrunStreamServer interface {
	Send(*TransactionResponse) error
	Recv() (*BackrunMsg, error)
	grpc.ServerStream
}

type aPIBackrunStreamServer struct {
	grpc.ServerStream
}

func (x *aPIBackrunStreamServer) Send(m *TransactionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPIBackrunStreamServer) Recv() (*BackrunMsg, error) {
	m := new(BackrunMsg)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _API_RawBackrunStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(APIServer).RawBackrunStream(&aPIRawBackrunStreamServer{stream})
}

type API_RawBackrunStreamServer interface {
	Send(*TransactionResponse) error
	Recv() (*RawBackrunMsg, error)
	grpc.ServerStream
}

type aPIRawBackrunStreamServer struct {
	grpc.ServerStream
}

func (x *aPIRawBackrunStreamServer) Send(m *TransactionResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aPIRawBackrunStreamServer) Recv() (*RawBackrunMsg, error) {
	m := new(RawBackrunMsg)
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
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendTransaction",
			Handler:    _API_SendTransaction_Handler,
		},
		{
			MethodName: "SendRawTransaction",
			Handler:    _API_SendRawTransaction_Handler,
		},
		{
			MethodName: "Backrun",
			Handler:    _API_Backrun_Handler,
		},
		{
			MethodName: "RawBackrun",
			Handler:    _API_RawBackrun_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeNewTxs",
			Handler:       _API_SubscribeNewTxs_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeNewTxsV2",
			Handler:       _API_SubscribeNewTxsV2_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeNewBlocks",
			Handler:       _API_SubscribeNewBlocks_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SendTransactionStream",
			Handler:       _API_SendTransactionStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SendRawTransactionStream",
			Handler:       _API_SendRawTransactionStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "BackrunStream",
			Handler:       _API_BackrunStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "RawBackrunStream",
			Handler:       _API_RawBackrunStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api.proto",
}
