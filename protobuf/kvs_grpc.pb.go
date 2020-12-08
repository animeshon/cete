// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protobuf

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// KVSClient is the client API for KVS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KVSClient interface {
	LivenessCheck(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*LivenessCheckResponse, error)
	ReadinessCheck(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ReadinessCheckResponse, error)
	Node(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*NodeResponse, error)
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Cluster(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ClusterResponse, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Snapshot(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	SetConditional(ctx context.Context, in *SetConditionalRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Watch(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (KVS_WatchClient, error)
	Metrics(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*MetricsResponse, error)
}

type kVSClient struct {
	cc grpc.ClientConnInterface
}

func NewKVSClient(cc grpc.ClientConnInterface) KVSClient {
	return &kVSClient{cc}
}

func (c *kVSClient) LivenessCheck(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*LivenessCheckResponse, error) {
	out := new(LivenessCheckResponse)
	err := c.cc.Invoke(ctx, "/kvs.KVS/LivenessCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) ReadinessCheck(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ReadinessCheckResponse, error) {
	out := new(ReadinessCheckResponse)
	err := c.cc.Invoke(ctx, "/kvs.KVS/ReadinessCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Node(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*NodeResponse, error) {
	out := new(NodeResponse)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Node", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Cluster(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ClusterResponse, error) {
	out := new(ClusterResponse)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Cluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Snapshot(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Snapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error) {
	out := new(ScanResponse)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Scan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) SetConditional(ctx context.Context, in *SetConditionalRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/kvs.KVS/SetConditional", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kVSClient) Watch(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (KVS_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &KVS_ServiceDesc.Streams[0], "/kvs.KVS/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &kVSWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KVS_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type kVSWatchClient struct {
	grpc.ClientStream
}

func (x *kVSWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *kVSClient) Metrics(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*MetricsResponse, error) {
	out := new(MetricsResponse)
	err := c.cc.Invoke(ctx, "/kvs.KVS/Metrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KVSServer is the server API for KVS service.
// All implementations must embed UnimplementedKVSServer
// for forward compatibility
type KVSServer interface {
	LivenessCheck(context.Context, *empty.Empty) (*LivenessCheckResponse, error)
	ReadinessCheck(context.Context, *empty.Empty) (*ReadinessCheckResponse, error)
	Node(context.Context, *empty.Empty) (*NodeResponse, error)
	Join(context.Context, *JoinRequest) (*empty.Empty, error)
	Cluster(context.Context, *empty.Empty) (*ClusterResponse, error)
	Leave(context.Context, *LeaveRequest) (*empty.Empty, error)
	Snapshot(context.Context, *empty.Empty) (*empty.Empty, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Scan(context.Context, *ScanRequest) (*ScanResponse, error)
	Set(context.Context, *SetRequest) (*empty.Empty, error)
	SetConditional(context.Context, *SetConditionalRequest) (*empty.Empty, error)
	Delete(context.Context, *DeleteRequest) (*empty.Empty, error)
	Watch(*empty.Empty, KVS_WatchServer) error
	Metrics(context.Context, *empty.Empty) (*MetricsResponse, error)
	mustEmbedUnimplementedKVSServer()
}

// UnimplementedKVSServer must be embedded to have forward compatible implementations.
type UnimplementedKVSServer struct {
}

func (UnimplementedKVSServer) LivenessCheck(context.Context, *empty.Empty) (*LivenessCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LivenessCheck not implemented")
}
func (UnimplementedKVSServer) ReadinessCheck(context.Context, *empty.Empty) (*ReadinessCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadinessCheck not implemented")
}
func (UnimplementedKVSServer) Node(context.Context, *empty.Empty) (*NodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Node not implemented")
}
func (UnimplementedKVSServer) Join(context.Context, *JoinRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedKVSServer) Cluster(context.Context, *empty.Empty) (*ClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cluster not implemented")
}
func (UnimplementedKVSServer) Leave(context.Context, *LeaveRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedKVSServer) Snapshot(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Snapshot not implemented")
}
func (UnimplementedKVSServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKVSServer) Scan(context.Context, *ScanRequest) (*ScanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Scan not implemented")
}
func (UnimplementedKVSServer) Set(context.Context, *SetRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedKVSServer) SetConditional(context.Context, *SetConditionalRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetConditional not implemented")
}
func (UnimplementedKVSServer) Delete(context.Context, *DeleteRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedKVSServer) Watch(*empty.Empty, KVS_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedKVSServer) Metrics(context.Context, *empty.Empty) (*MetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Metrics not implemented")
}
func (UnimplementedKVSServer) mustEmbedUnimplementedKVSServer() {}

// UnsafeKVSServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KVSServer will
// result in compilation errors.
type UnsafeKVSServer interface {
	mustEmbedUnimplementedKVSServer()
}

func RegisterKVSServer(s grpc.ServiceRegistrar, srv KVSServer) {
	s.RegisterService(&KVS_ServiceDesc, srv)
}

func _KVS_LivenessCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).LivenessCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/LivenessCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).LivenessCheck(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_ReadinessCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).ReadinessCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/ReadinessCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).ReadinessCheck(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Node_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Node(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Node",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Node(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Cluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Cluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Cluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Cluster(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Snapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Snapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Snapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Snapshot(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Scan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Scan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Scan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Scan(ctx, req.(*ScanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_SetConditional_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetConditionalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).SetConditional(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/SetConditional",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).SetConditional(ctx, req.(*SetConditionalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KVS_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KVSServer).Watch(m, &kVSWatchServer{stream})
}

type KVS_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type kVSWatchServer struct {
	grpc.ServerStream
}

func (x *kVSWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _KVS_Metrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KVSServer).Metrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvs.KVS/Metrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVSServer).Metrics(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// KVS_ServiceDesc is the grpc.ServiceDesc for KVS service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KVS_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kvs.KVS",
	HandlerType: (*KVSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LivenessCheck",
			Handler:    _KVS_LivenessCheck_Handler,
		},
		{
			MethodName: "ReadinessCheck",
			Handler:    _KVS_ReadinessCheck_Handler,
		},
		{
			MethodName: "Node",
			Handler:    _KVS_Node_Handler,
		},
		{
			MethodName: "Join",
			Handler:    _KVS_Join_Handler,
		},
		{
			MethodName: "Cluster",
			Handler:    _KVS_Cluster_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _KVS_Leave_Handler,
		},
		{
			MethodName: "Snapshot",
			Handler:    _KVS_Snapshot_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KVS_Get_Handler,
		},
		{
			MethodName: "Scan",
			Handler:    _KVS_Scan_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _KVS_Set_Handler,
		},
		{
			MethodName: "SetConditional",
			Handler:    _KVS_SetConditional_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KVS_Delete_Handler,
		},
		{
			MethodName: "Metrics",
			Handler:    _KVS_Metrics_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _KVS_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protobuf/kvs.proto",
}
