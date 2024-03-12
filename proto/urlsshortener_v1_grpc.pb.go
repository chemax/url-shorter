// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: proto/urlsshortener_v1.proto

package api

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
	URLShortenerV1_Ping_FullMethodName            = "/url_shortener.URLShortenerV1/Ping"
	URLShortenerV1_GetOriginalURL_FullMethodName  = "/url_shortener.URLShortenerV1/GetOriginalURL"
	URLShortenerV1_GetURLsByUserID_FullMethodName = "/url_shortener.URLShortenerV1/GetURLsByUserID"
	URLShortenerV1_CreateURL_FullMethodName       = "/url_shortener.URLShortenerV1/CreateURL"
	URLShortenerV1_CreateURLs_FullMethodName      = "/url_shortener.URLShortenerV1/CreateURLs"
	URLShortenerV1_DeleteURLs_FullMethodName      = "/url_shortener.URLShortenerV1/DeleteURLs"
	URLShortenerV1_Stat_FullMethodName            = "/url_shortener.URLShortenerV1/Stat"
)

// URLShortenerV1Client is the client API for URLShortenerV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerV1Client interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	GetOriginalURL(ctx context.Context, in *UnshortURLRequest, opts ...grpc.CallOption) (*UnshortURLResponse, error)
	GetURLsByUserID(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error)
	CreateURL(ctx context.Context, in *ShortURLRequest, opts ...grpc.CallOption) (*ShortURLResponse, error)
	CreateURLs(ctx context.Context, in *ShortURLsBatchRequest, opts ...grpc.CallOption) (*ShortURLsBatchResponse, error)
	DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*DeleteURLsResponse, error)
	Stat(ctx context.Context, in *StatRequest, opts ...grpc.CallOption) (*StatResponse, error)
}

type uRLShortenerV1Client struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerV1Client(cc grpc.ClientConnInterface) URLShortenerV1Client {
	return &uRLShortenerV1Client{cc}
}

func (c *uRLShortenerV1Client) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, URLShortenerV1_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerV1Client) GetOriginalURL(ctx context.Context, in *UnshortURLRequest, opts ...grpc.CallOption) (*UnshortURLResponse, error) {
	out := new(UnshortURLResponse)
	err := c.cc.Invoke(ctx, URLShortenerV1_GetOriginalURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerV1Client) GetURLsByUserID(ctx context.Context, in *GetUserURLsRequest, opts ...grpc.CallOption) (*GetUserURLsResponse, error) {
	out := new(GetUserURLsResponse)
	err := c.cc.Invoke(ctx, URLShortenerV1_GetURLsByUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerV1Client) CreateURL(ctx context.Context, in *ShortURLRequest, opts ...grpc.CallOption) (*ShortURLResponse, error) {
	out := new(ShortURLResponse)
	err := c.cc.Invoke(ctx, URLShortenerV1_CreateURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerV1Client) CreateURLs(ctx context.Context, in *ShortURLsBatchRequest, opts ...grpc.CallOption) (*ShortURLsBatchResponse, error) {
	out := new(ShortURLsBatchResponse)
	err := c.cc.Invoke(ctx, URLShortenerV1_CreateURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerV1Client) DeleteURLs(ctx context.Context, in *DeleteURLsRequest, opts ...grpc.CallOption) (*DeleteURLsResponse, error) {
	out := new(DeleteURLsResponse)
	err := c.cc.Invoke(ctx, URLShortenerV1_DeleteURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerV1Client) Stat(ctx context.Context, in *StatRequest, opts ...grpc.CallOption) (*StatResponse, error) {
	out := new(StatResponse)
	err := c.cc.Invoke(ctx, URLShortenerV1_Stat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerV1Server is the server API for URLShortenerV1 service.
// All implementations must embed UnimplementedURLShortenerV1Server
// for forward compatibility
type URLShortenerV1Server interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	GetOriginalURL(context.Context, *UnshortURLRequest) (*UnshortURLResponse, error)
	GetURLsByUserID(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error)
	CreateURL(context.Context, *ShortURLRequest) (*ShortURLResponse, error)
	CreateURLs(context.Context, *ShortURLsBatchRequest) (*ShortURLsBatchResponse, error)
	DeleteURLs(context.Context, *DeleteURLsRequest) (*DeleteURLsResponse, error)
	Stat(context.Context, *StatRequest) (*StatResponse, error)
	mustEmbedUnimplementedURLShortenerV1Server()
}

// UnimplementedURLShortenerV1Server must be embedded to have forward compatible implementations.
type UnimplementedURLShortenerV1Server struct {
}

func (UnimplementedURLShortenerV1Server) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedURLShortenerV1Server) GetOriginalURL(context.Context, *UnshortURLRequest) (*UnshortURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedURLShortenerV1Server) GetURLsByUserID(context.Context, *GetUserURLsRequest) (*GetUserURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURLsByUserID not implemented")
}
func (UnimplementedURLShortenerV1Server) CreateURL(context.Context, *ShortURLRequest) (*ShortURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateURL not implemented")
}
func (UnimplementedURLShortenerV1Server) CreateURLs(context.Context, *ShortURLsBatchRequest) (*ShortURLsBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateURLs not implemented")
}
func (UnimplementedURLShortenerV1Server) DeleteURLs(context.Context, *DeleteURLsRequest) (*DeleteURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURLs not implemented")
}
func (UnimplementedURLShortenerV1Server) Stat(context.Context, *StatRequest) (*StatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stat not implemented")
}
func (UnimplementedURLShortenerV1Server) mustEmbedUnimplementedURLShortenerV1Server() {}

// UnsafeURLShortenerV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerV1Server will
// result in compilation errors.
type UnsafeURLShortenerV1Server interface {
	mustEmbedUnimplementedURLShortenerV1Server()
}

func RegisterURLShortenerV1Server(s grpc.ServiceRegistrar, srv URLShortenerV1Server) {
	s.RegisterService(&URLShortenerV1_ServiceDesc, srv)
}

func _URLShortenerV1_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerV1Server).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerV1_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerV1Server).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerV1_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnshortURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerV1Server).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerV1_GetOriginalURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerV1Server).GetOriginalURL(ctx, req.(*UnshortURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerV1_GetURLsByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerV1Server).GetURLsByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerV1_GetURLsByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerV1Server).GetURLsByUserID(ctx, req.(*GetUserURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerV1_CreateURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerV1Server).CreateURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerV1_CreateURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerV1Server).CreateURL(ctx, req.(*ShortURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerV1_CreateURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortURLsBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerV1Server).CreateURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerV1_CreateURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerV1Server).CreateURLs(ctx, req.(*ShortURLsBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerV1_DeleteURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerV1Server).DeleteURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerV1_DeleteURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerV1Server).DeleteURLs(ctx, req.(*DeleteURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortenerV1_Stat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerV1Server).Stat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortenerV1_Stat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerV1Server).Stat(ctx, req.(*StatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortenerV1_ServiceDesc is the grpc.ServiceDesc for URLShortenerV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortenerV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "url_shortener.URLShortenerV1",
	HandlerType: (*URLShortenerV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _URLShortenerV1_Ping_Handler,
		},
		{
			MethodName: "GetOriginalURL",
			Handler:    _URLShortenerV1_GetOriginalURL_Handler,
		},
		{
			MethodName: "GetURLsByUserID",
			Handler:    _URLShortenerV1_GetURLsByUserID_Handler,
		},
		{
			MethodName: "CreateURL",
			Handler:    _URLShortenerV1_CreateURL_Handler,
		},
		{
			MethodName: "CreateURLs",
			Handler:    _URLShortenerV1_CreateURLs_Handler,
		},
		{
			MethodName: "DeleteURLs",
			Handler:    _URLShortenerV1_DeleteURLs_Handler,
		},
		{
			MethodName: "Stat",
			Handler:    _URLShortenerV1_Stat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/urlsshortener_v1.proto",
}
