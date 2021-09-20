// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package queue_pb

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

// QueueClient is the client API for Queue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueueClient interface {
	SaveMessages(ctx context.Context, in *Messages, opts ...grpc.CallOption) (*ErrorReply, error)
	DeleteMessages(ctx context.Context, in *MessageUids, opts ...grpc.CallOption) (*ErrorReply, error)
	GetMessage(ctx context.Context, in *RequestSingle, opts ...grpc.CallOption) (*Message, error)
	GetMessages(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Messages, error)
	GetTopicList(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*TopicListReply, error)
	GetMessageCount(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*TopicCount, error)
}

type queueClient struct {
	cc grpc.ClientConnInterface
}

func NewQueueClient(cc grpc.ClientConnInterface) QueueClient {
	return &queueClient{cc}
}

func (c *queueClient) SaveMessages(ctx context.Context, in *Messages, opts ...grpc.CallOption) (*ErrorReply, error) {
	out := new(ErrorReply)
	err := c.cc.Invoke(ctx, "/queue_pb.Queue/SaveMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) DeleteMessages(ctx context.Context, in *MessageUids, opts ...grpc.CallOption) (*ErrorReply, error) {
	out := new(ErrorReply)
	err := c.cc.Invoke(ctx, "/queue_pb.Queue/DeleteMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) GetMessage(ctx context.Context, in *RequestSingle, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/queue_pb.Queue/GetMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) GetMessages(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Messages, error) {
	out := new(Messages)
	err := c.cc.Invoke(ctx, "/queue_pb.Queue/GetMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) GetTopicList(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*TopicListReply, error) {
	out := new(TopicListReply)
	err := c.cc.Invoke(ctx, "/queue_pb.Queue/GetTopicList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queueClient) GetMessageCount(ctx context.Context, in *CountRequest, opts ...grpc.CallOption) (*TopicCount, error) {
	out := new(TopicCount)
	err := c.cc.Invoke(ctx, "/queue_pb.Queue/GetMessageCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueueServer is the server API for Queue service.
// All implementations must embed UnimplementedQueueServer
// for forward compatibility
type QueueServer interface {
	SaveMessages(context.Context, *Messages) (*ErrorReply, error)
	DeleteMessages(context.Context, *MessageUids) (*ErrorReply, error)
	GetMessage(context.Context, *RequestSingle) (*Message, error)
	GetMessages(context.Context, *Request) (*Messages, error)
	GetTopicList(context.Context, *EmptyRequest) (*TopicListReply, error)
	GetMessageCount(context.Context, *CountRequest) (*TopicCount, error)
	mustEmbedUnimplementedQueueServer()
}

// UnimplementedQueueServer must be embedded to have forward compatible implementations.
type UnimplementedQueueServer struct {
}

func (UnimplementedQueueServer) SaveMessages(context.Context, *Messages) (*ErrorReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveMessages not implemented")
}
func (UnimplementedQueueServer) DeleteMessages(context.Context, *MessageUids) (*ErrorReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessages not implemented")
}
func (UnimplementedQueueServer) GetMessage(context.Context, *RequestSingle) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessage not implemented")
}
func (UnimplementedQueueServer) GetMessages(context.Context, *Request) (*Messages, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}
func (UnimplementedQueueServer) GetTopicList(context.Context, *EmptyRequest) (*TopicListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopicList not implemented")
}
func (UnimplementedQueueServer) GetMessageCount(context.Context, *CountRequest) (*TopicCount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessageCount not implemented")
}
func (UnimplementedQueueServer) mustEmbedUnimplementedQueueServer() {}

// UnsafeQueueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueueServer will
// result in compilation errors.
type UnsafeQueueServer interface {
	mustEmbedUnimplementedQueueServer()
}

func RegisterQueueServer(s grpc.ServiceRegistrar, srv QueueServer) {
	s.RegisterService(&Queue_ServiceDesc, srv)
}

func _Queue_SaveMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Messages)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).SaveMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue_pb.Queue/SaveMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).SaveMessages(ctx, req.(*Messages))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_DeleteMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageUids)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).DeleteMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue_pb.Queue/DeleteMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).DeleteMessages(ctx, req.(*MessageUids))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestSingle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue_pb.Queue/GetMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).GetMessage(ctx, req.(*RequestSingle))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_GetMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).GetMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue_pb.Queue/GetMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).GetMessages(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_GetTopicList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).GetTopicList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue_pb.Queue/GetTopicList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).GetTopicList(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queue_GetMessageCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueServer).GetMessageCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/queue_pb.Queue/GetMessageCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueServer).GetMessageCount(ctx, req.(*CountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Queue_ServiceDesc is the grpc.ServiceDesc for Queue service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Queue_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "queue_pb.Queue",
	HandlerType: (*QueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveMessages",
			Handler:    _Queue_SaveMessages_Handler,
		},
		{
			MethodName: "DeleteMessages",
			Handler:    _Queue_DeleteMessages_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _Queue_GetMessage_Handler,
		},
		{
			MethodName: "GetMessages",
			Handler:    _Queue_GetMessages_Handler,
		},
		{
			MethodName: "GetTopicList",
			Handler:    _Queue_GetTopicList_Handler,
		},
		{
			MethodName: "GetMessageCount",
			Handler:    _Queue_GetMessageCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "queue.proto",
}