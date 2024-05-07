// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: loms/loms.proto

package loms

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LomsClient is the client API for Loms service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LomsClient interface {
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	OrderInfo(ctx context.Context, in *OrderInfoRequest, opts ...grpc.CallOption) (*OrderInfoResponse, error)
	OrderPay(ctx context.Context, in *OrderPayRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	OrderCancel(ctx context.Context, in *OrderCancelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	StocksInfo(ctx context.Context, in *StocksInfoRequest, opts ...grpc.CallOption) (*StocksInfoResponse, error)
}

type lomsClient struct {
	cc grpc.ClientConnInterface
}

func NewLomsClient(cc grpc.ClientConnInterface) LomsClient {
	return &lomsClient{cc}
}

func (c *lomsClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) OrderInfo(ctx context.Context, in *OrderInfoRequest, opts ...grpc.CallOption) (*OrderInfoResponse, error) {
	out := new(OrderInfoResponse)
	err := c.cc.Invoke(ctx, "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/OrderInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) OrderPay(ctx context.Context, in *OrderPayRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/OrderPay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) OrderCancel(ctx context.Context, in *OrderCancelRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/OrderCancel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lomsClient) StocksInfo(ctx context.Context, in *StocksInfoRequest, opts ...grpc.CallOption) (*StocksInfoResponse, error) {
	out := new(StocksInfoResponse)
	err := c.cc.Invoke(ctx, "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/StocksInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LomsServer is the server API for Loms service.
// All implementations must embed UnimplementedLomsServer
// for forward compatibility
type LomsServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	OrderInfo(context.Context, *OrderInfoRequest) (*OrderInfoResponse, error)
	OrderPay(context.Context, *OrderPayRequest) (*emptypb.Empty, error)
	OrderCancel(context.Context, *OrderCancelRequest) (*emptypb.Empty, error)
	StocksInfo(context.Context, *StocksInfoRequest) (*StocksInfoResponse, error)
	mustEmbedUnimplementedLomsServer()
}

// UnimplementedLomsServer must be embedded to have forward compatible implementations.
type UnimplementedLomsServer struct {
}

func (UnimplementedLomsServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}
func (UnimplementedLomsServer) OrderInfo(context.Context, *OrderInfoRequest) (*OrderInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderInfo not implemented")
}
func (UnimplementedLomsServer) OrderPay(context.Context, *OrderPayRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderPay not implemented")
}
func (UnimplementedLomsServer) OrderCancel(context.Context, *OrderCancelRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderCancel not implemented")
}
func (UnimplementedLomsServer) StocksInfo(context.Context, *StocksInfoRequest) (*StocksInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StocksInfo not implemented")
}
func (UnimplementedLomsServer) mustEmbedUnimplementedLomsServer() {}

// UnsafeLomsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LomsServer will
// result in compilation errors.
type UnsafeLomsServer interface {
	mustEmbedUnimplementedLomsServer()
}

func RegisterLomsServer(s grpc.ServiceRegistrar, srv LomsServer) {
	s.RegisterService(&Loms_ServiceDesc, srv)
}

func _Loms_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/CreateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_OrderInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).OrderInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/OrderInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).OrderInfo(ctx, req.(*OrderInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_OrderPay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderPayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).OrderPay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/OrderPay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).OrderPay(ctx, req.(*OrderPayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_OrderCancel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).OrderCancel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/OrderCancel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).OrderCancel(ctx, req.(*OrderCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Loms_StocksInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StocksInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LomsServer).StocksInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitlab.ozon.dev.max_lorriess.student_project.loms.Loms/StocksInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LomsServer).StocksInfo(ctx, req.(*StocksInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Loms_ServiceDesc is the grpc.ServiceDesc for Loms service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Loms_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gitlab.ozon.dev.max_lorriess.student_project.loms.Loms",
	HandlerType: (*LomsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOrder",
			Handler:    _Loms_CreateOrder_Handler,
		},
		{
			MethodName: "OrderInfo",
			Handler:    _Loms_OrderInfo_Handler,
		},
		{
			MethodName: "OrderPay",
			Handler:    _Loms_OrderPay_Handler,
		},
		{
			MethodName: "OrderCancel",
			Handler:    _Loms_OrderCancel_Handler,
		},
		{
			MethodName: "StocksInfo",
			Handler:    _Loms_StocksInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "loms/loms.proto",
}
