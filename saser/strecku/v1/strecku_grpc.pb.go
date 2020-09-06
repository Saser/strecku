// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package streckuv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StreckUClient is the client API for StreckU service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreckUClient interface {
	// GetUser gets a single user.
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error)
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error)
}

type streckUClient struct {
	cc grpc.ClientConnInterface
}

func NewStreckUClient(cc grpc.ClientConnInterface) StreckUClient {
	return &streckUClient{cc}
}

func (c *streckUClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/saser.strecku.v1.StreckU/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streckUClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/saser.strecku.v1.StreckU/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StreckUServer is the server API for StreckU service.
// All implementations must embed UnimplementedStreckUServer
// for forward compatibility
type StreckUServer interface {
	// GetUser gets a single user.
	GetUser(context.Context, *GetUserRequest) (*User, error)
	// CreateUser creates a new user.
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	mustEmbedUnimplementedStreckUServer()
}

// UnimplementedStreckUServer must be embedded to have forward compatible implementations.
type UnimplementedStreckUServer struct {
}

func (*UnimplementedStreckUServer) GetUser(context.Context, *GetUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (*UnimplementedStreckUServer) CreateUser(context.Context, *CreateUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (*UnimplementedStreckUServer) mustEmbedUnimplementedStreckUServer() {}

func RegisterStreckUServer(s *grpc.Server, srv StreckUServer) {
	s.RegisterService(&_StreckU_serviceDesc, srv)
}

func _StreckU_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreckUServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/saser.strecku.v1.StreckU/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreckUServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreckU_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreckUServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/saser.strecku.v1.StreckU/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreckUServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StreckU_serviceDesc = grpc.ServiceDesc{
	ServiceName: "saser.strecku.v1.StreckU",
	HandlerType: (*StreckUServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _StreckU_GetUser_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _StreckU_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "saser/strecku/v1/strecku.proto",
}
