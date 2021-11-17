// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rpc_server

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

// OrgInfoServiceClient is the client API for OrgInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrgInfoServiceClient interface {
	Fetch(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type orgInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrgInfoServiceClient(cc grpc.ClientConnInterface) OrgInfoServiceClient {
	return &orgInfoServiceClient{cc}
}

func (c *orgInfoServiceClient) Fetch(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc_server.OrgInfoService/Fetch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrgInfoServiceServer is the server API for OrgInfoService service.
// All implementations must embed UnimplementedOrgInfoServiceServer
// for forward compatibility
type OrgInfoServiceServer interface {
	Fetch(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedOrgInfoServiceServer()
}

// UnimplementedOrgInfoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrgInfoServiceServer struct {
}

func (UnimplementedOrgInfoServiceServer) Fetch(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fetch not implemented")
}
func (UnimplementedOrgInfoServiceServer) mustEmbedUnimplementedOrgInfoServiceServer() {}

// UnsafeOrgInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrgInfoServiceServer will
// result in compilation errors.
type UnsafeOrgInfoServiceServer interface {
	mustEmbedUnimplementedOrgInfoServiceServer()
}

func RegisterOrgInfoServiceServer(s grpc.ServiceRegistrar, srv OrgInfoServiceServer) {
	s.RegisterService(&OrgInfoService_ServiceDesc, srv)
}

func _OrgInfoService_Fetch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrgInfoServiceServer).Fetch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc_server.OrgInfoService/Fetch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrgInfoServiceServer).Fetch(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// OrgInfoService_ServiceDesc is the grpc.ServiceDesc for OrgInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrgInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc_server.OrgInfoService",
	HandlerType: (*OrgInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Fetch",
			Handler:    _OrgInfoService_Fetch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/rpc_server/server.proto",
}