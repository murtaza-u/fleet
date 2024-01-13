// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: fleet.proto

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
	Fleet_Listen_FullMethodName = "/pb.Fleet/Listen"
)

// FleetClient is the client API for Fleet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FleetClient interface {
	Listen(ctx context.Context, opts ...grpc.CallOption) (Fleet_ListenClient, error)
}

type fleetClient struct {
	cc grpc.ClientConnInterface
}

func NewFleetClient(cc grpc.ClientConnInterface) FleetClient {
	return &fleetClient{cc}
}

func (c *fleetClient) Listen(ctx context.Context, opts ...grpc.CallOption) (Fleet_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &Fleet_ServiceDesc.Streams[0], Fleet_Listen_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &fleetListenClient{stream}
	return x, nil
}

type Fleet_ListenClient interface {
	Send(*Reply) error
	Recv() (*Request, error)
	grpc.ClientStream
}

type fleetListenClient struct {
	grpc.ClientStream
}

func (x *fleetListenClient) Send(m *Reply) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fleetListenClient) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FleetServer is the server API for Fleet service.
// All implementations must embed UnimplementedFleetServer
// for forward compatibility
type FleetServer interface {
	Listen(Fleet_ListenServer) error
	mustEmbedUnimplementedFleetServer()
}

// UnimplementedFleetServer must be embedded to have forward compatible implementations.
type UnimplementedFleetServer struct {
}

func (UnimplementedFleetServer) Listen(Fleet_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedFleetServer) mustEmbedUnimplementedFleetServer() {}

// UnsafeFleetServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FleetServer will
// result in compilation errors.
type UnsafeFleetServer interface {
	mustEmbedUnimplementedFleetServer()
}

func RegisterFleetServer(s grpc.ServiceRegistrar, srv FleetServer) {
	s.RegisterService(&Fleet_ServiceDesc, srv)
}

func _Fleet_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FleetServer).Listen(&fleetListenServer{stream})
}

type Fleet_ListenServer interface {
	Send(*Request) error
	Recv() (*Reply, error)
	grpc.ServerStream
}

type fleetListenServer struct {
	grpc.ServerStream
}

func (x *fleetListenServer) Send(m *Request) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fleetListenServer) Recv() (*Reply, error) {
	m := new(Reply)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Fleet_ServiceDesc is the grpc.ServiceDesc for Fleet service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fleet_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Fleet",
	HandlerType: (*FleetServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _Fleet_Listen_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "fleet.proto",
}
