// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: pub/v1/pub.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PubService_PushMessage_FullMethodName = "/pub.v1.PubService/PushMessage"
)

// PubServiceClient is the client API for PubService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PubServiceClient interface {
	PushMessage(ctx context.Context, in *PushMessageRequest, opts ...grpc.CallOption) (*PushMessageResponse, error)
}

type pubServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPubServiceClient(cc grpc.ClientConnInterface) PubServiceClient {
	return &pubServiceClient{cc}
}

func (c *pubServiceClient) PushMessage(ctx context.Context, in *PushMessageRequest, opts ...grpc.CallOption) (*PushMessageResponse, error) {
	out := new(PushMessageResponse)
	err := c.cc.Invoke(ctx, PubService_PushMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PubServiceServer is the server API for PubService service.
// All implementations must embed UnimplementedPubServiceServer
// for forward compatibility
type PubServiceServer interface {
	PushMessage(context.Context, *PushMessageRequest) (*PushMessageResponse, error)
	mustEmbedUnimplementedPubServiceServer()
}

// UnimplementedPubServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPubServiceServer struct {
}

func (UnimplementedPubServiceServer) PushMessage(context.Context, *PushMessageRequest) (*PushMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMessage not implemented")
}
func (UnimplementedPubServiceServer) mustEmbedUnimplementedPubServiceServer() {}

// UnsafePubServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PubServiceServer will
// result in compilation errors.
type UnsafePubServiceServer interface {
	mustEmbedUnimplementedPubServiceServer()
}

func RegisterPubServiceServer(s grpc.ServiceRegistrar, srv PubServiceServer) {
	s.RegisterService(&PubService_ServiceDesc, srv)
}

func _PubService_PushMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PubServiceServer).PushMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PubService_PushMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PubServiceServer).PushMessage(ctx, req.(*PushMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PubService_ServiceDesc is the grpc.ServiceDesc for PubService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PubService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pub.v1.PubService",
	HandlerType: (*PubServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushMessage",
			Handler:    _PubService_PushMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pub/v1/pub.proto",
}
