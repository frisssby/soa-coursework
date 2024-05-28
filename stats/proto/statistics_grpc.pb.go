// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: statistics.proto

package statspb

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

const (
	StatisticsService_GetTaskStats_FullMethodName = "/statspb.StatisticsService/GetTaskStats"
	StatisticsService_GetTasksTop_FullMethodName  = "/statspb.StatisticsService/GetTasksTop"
	StatisticsService_GetUsersTop_FullMethodName  = "/statspb.StatisticsService/GetUsersTop"
)

// StatisticsServiceClient is the client API for StatisticsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatisticsServiceClient interface {
	GetTaskStats(ctx context.Context, in *GetTaskStatsRequest, opts ...grpc.CallOption) (*GetTaskStatsResponse, error)
	GetTasksTop(ctx context.Context, in *GetTasksTopRequest, opts ...grpc.CallOption) (*GetTasksTopResponse, error)
	GetUsersTop(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUsersTopResponse, error)
}

type statisticsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStatisticsServiceClient(cc grpc.ClientConnInterface) StatisticsServiceClient {
	return &statisticsServiceClient{cc}
}

func (c *statisticsServiceClient) GetTaskStats(ctx context.Context, in *GetTaskStatsRequest, opts ...grpc.CallOption) (*GetTaskStatsResponse, error) {
	out := new(GetTaskStatsResponse)
	err := c.cc.Invoke(ctx, StatisticsService_GetTaskStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticsServiceClient) GetTasksTop(ctx context.Context, in *GetTasksTopRequest, opts ...grpc.CallOption) (*GetTasksTopResponse, error) {
	out := new(GetTasksTopResponse)
	err := c.cc.Invoke(ctx, StatisticsService_GetTasksTop_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticsServiceClient) GetUsersTop(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUsersTopResponse, error) {
	out := new(GetUsersTopResponse)
	err := c.cc.Invoke(ctx, StatisticsService_GetUsersTop_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatisticsServiceServer is the server API for StatisticsService service.
// All implementations must embed UnimplementedStatisticsServiceServer
// for forward compatibility
type StatisticsServiceServer interface {
	GetTaskStats(context.Context, *GetTaskStatsRequest) (*GetTaskStatsResponse, error)
	GetTasksTop(context.Context, *GetTasksTopRequest) (*GetTasksTopResponse, error)
	GetUsersTop(context.Context, *emptypb.Empty) (*GetUsersTopResponse, error)
	mustEmbedUnimplementedStatisticsServiceServer()
}

// UnimplementedStatisticsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStatisticsServiceServer struct {
}

func (UnimplementedStatisticsServiceServer) GetTaskStats(context.Context, *GetTaskStatsRequest) (*GetTaskStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaskStats not implemented")
}
func (UnimplementedStatisticsServiceServer) GetTasksTop(context.Context, *GetTasksTopRequest) (*GetTasksTopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTasksTop not implemented")
}
func (UnimplementedStatisticsServiceServer) GetUsersTop(context.Context, *emptypb.Empty) (*GetUsersTopResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersTop not implemented")
}
func (UnimplementedStatisticsServiceServer) mustEmbedUnimplementedStatisticsServiceServer() {}

// UnsafeStatisticsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatisticsServiceServer will
// result in compilation errors.
type UnsafeStatisticsServiceServer interface {
	mustEmbedUnimplementedStatisticsServiceServer()
}

func RegisterStatisticsServiceServer(s grpc.ServiceRegistrar, srv StatisticsServiceServer) {
	s.RegisterService(&StatisticsService_ServiceDesc, srv)
}

func _StatisticsService_GetTaskStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaskStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticsServiceServer).GetTaskStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatisticsService_GetTaskStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticsServiceServer).GetTaskStats(ctx, req.(*GetTaskStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatisticsService_GetTasksTop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTasksTopRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticsServiceServer).GetTasksTop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatisticsService_GetTasksTop_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticsServiceServer).GetTasksTop(ctx, req.(*GetTasksTopRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatisticsService_GetUsersTop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticsServiceServer).GetUsersTop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatisticsService_GetUsersTop_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticsServiceServer).GetUsersTop(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// StatisticsService_ServiceDesc is the grpc.ServiceDesc for StatisticsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatisticsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "statspb.StatisticsService",
	HandlerType: (*StatisticsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTaskStats",
			Handler:    _StatisticsService_GetTaskStats_Handler,
		},
		{
			MethodName: "GetTasksTop",
			Handler:    _StatisticsService_GetTasksTop_Handler,
		},
		{
			MethodName: "GetUsersTop",
			Handler:    _StatisticsService_GetUsersTop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statistics.proto",
}
