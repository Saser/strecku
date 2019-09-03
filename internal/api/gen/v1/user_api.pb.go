// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/user_api.proto

package streckuv1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Request message for UserAPI.ListUsers.
type ListUsersRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListUsersRequest) Reset()         { *m = ListUsersRequest{} }
func (m *ListUsersRequest) String() string { return proto.CompactTextString(m) }
func (*ListUsersRequest) ProtoMessage()    {}
func (*ListUsersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6beeca62d8388557, []int{0}
}

func (m *ListUsersRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListUsersRequest.Unmarshal(m, b)
}
func (m *ListUsersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListUsersRequest.Marshal(b, m, deterministic)
}
func (m *ListUsersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListUsersRequest.Merge(m, src)
}
func (m *ListUsersRequest) XXX_Size() int {
	return xxx_messageInfo_ListUsersRequest.Size(m)
}
func (m *ListUsersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListUsersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListUsersRequest proto.InternalMessageInfo

// Response message for UserAPI.ListUsers.
type ListUsersResponse struct {
	// The list of users.
	Users                []*User  `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListUsersResponse) Reset()         { *m = ListUsersResponse{} }
func (m *ListUsersResponse) String() string { return proto.CompactTextString(m) }
func (*ListUsersResponse) ProtoMessage()    {}
func (*ListUsersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6beeca62d8388557, []int{1}
}

func (m *ListUsersResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListUsersResponse.Unmarshal(m, b)
}
func (m *ListUsersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListUsersResponse.Marshal(b, m, deterministic)
}
func (m *ListUsersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListUsersResponse.Merge(m, src)
}
func (m *ListUsersResponse) XXX_Size() int {
	return xxx_messageInfo_ListUsersResponse.Size(m)
}
func (m *ListUsersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListUsersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListUsersResponse proto.InternalMessageInfo

func (m *ListUsersResponse) GetUsers() []*User {
	if m != nil {
		return m.Users
	}
	return nil
}

// Request message for UserAPI.CreateUser.
type CreateUserRequest struct {
	// The user to create.
	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// The password that the user should be created with.
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateUserRequest) Reset()         { *m = CreateUserRequest{} }
func (m *CreateUserRequest) String() string { return proto.CompactTextString(m) }
func (*CreateUserRequest) ProtoMessage()    {}
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6beeca62d8388557, []int{2}
}

func (m *CreateUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateUserRequest.Unmarshal(m, b)
}
func (m *CreateUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateUserRequest.Marshal(b, m, deterministic)
}
func (m *CreateUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateUserRequest.Merge(m, src)
}
func (m *CreateUserRequest) XXX_Size() int {
	return xxx_messageInfo_CreateUserRequest.Size(m)
}
func (m *CreateUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateUserRequest proto.InternalMessageInfo

func (m *CreateUserRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *CreateUserRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

// Response message for UserAPI.CreateUser.
type CreateUserResponse struct {
	// The user that was created.
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateUserResponse) Reset()         { *m = CreateUserResponse{} }
func (m *CreateUserResponse) String() string { return proto.CompactTextString(m) }
func (*CreateUserResponse) ProtoMessage()    {}
func (*CreateUserResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6beeca62d8388557, []int{3}
}

func (m *CreateUserResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateUserResponse.Unmarshal(m, b)
}
func (m *CreateUserResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateUserResponse.Marshal(b, m, deterministic)
}
func (m *CreateUserResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateUserResponse.Merge(m, src)
}
func (m *CreateUserResponse) XXX_Size() int {
	return xxx_messageInfo_CreateUserResponse.Size(m)
}
func (m *CreateUserResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateUserResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateUserResponse proto.InternalMessageInfo

func (m *CreateUserResponse) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func init() {
	proto.RegisterType((*ListUsersRequest)(nil), "strecku.v1.ListUsersRequest")
	proto.RegisterType((*ListUsersResponse)(nil), "strecku.v1.ListUsersResponse")
	proto.RegisterType((*CreateUserRequest)(nil), "strecku.v1.CreateUserRequest")
	proto.RegisterType((*CreateUserResponse)(nil), "strecku.v1.CreateUserResponse")
}

func init() { proto.RegisterFile("v1/user_api.proto", fileDescriptor_6beeca62d8388557) }

var fileDescriptor_6beeca62d8388557 = []byte{
	// 275 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x33, 0xd4, 0x2f,
	0x2d, 0x4e, 0x2d, 0x8a, 0x4f, 0x2c, 0xc8, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2a,
	0x2e, 0x29, 0x4a, 0x4d, 0xce, 0x2e, 0xd5, 0x2b, 0x33, 0x94, 0xe2, 0x85, 0x4a, 0x43, 0xa4, 0x94,
	0x84, 0xb8, 0x04, 0x7c, 0x32, 0x8b, 0x4b, 0x42, 0x8b, 0x53, 0x8b, 0x8a, 0x83, 0x52, 0x0b, 0x4b,
	0x53, 0x8b, 0x4b, 0x94, 0xac, 0xb9, 0x04, 0x91, 0xc4, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85,
	0xd4, 0xb8, 0x58, 0x41, 0xda, 0x8a, 0x25, 0x18, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0x04, 0xf4, 0x10,
	0x66, 0xea, 0x81, 0x54, 0x06, 0x41, 0xa4, 0x95, 0x42, 0xb9, 0x04, 0x9d, 0x8b, 0x52, 0x13, 0x4b,
	0x52, 0xc1, 0x82, 0x10, 0x13, 0x85, 0x54, 0xb8, 0x58, 0x40, 0xb2, 0x12, 0x8c, 0x0a, 0x8c, 0x58,
	0xf5, 0x82, 0x65, 0x85, 0xa4, 0xb8, 0x38, 0x0a, 0x12, 0x8b, 0x8b, 0xcb, 0xf3, 0x8b, 0x52, 0x24,
	0x98, 0x14, 0x18, 0x35, 0x38, 0x83, 0xe0, 0x7c, 0x25, 0x2b, 0x2e, 0x21, 0x64, 0x63, 0xa1, 0x8e,
	0x22, 0xca, 0x5c, 0xa3, 0x05, 0x8c, 0x5c, 0xec, 0x20, 0xae, 0x63, 0x80, 0xa7, 0x90, 0x07, 0x17,
	0x27, 0xdc, 0x6f, 0x42, 0x32, 0xc8, 0x1a, 0xd0, 0x83, 0x41, 0x4a, 0x16, 0x87, 0x2c, 0xd4, 0x6e,
	0x6f, 0x2e, 0x2e, 0x84, 0x8b, 0x84, 0x50, 0x14, 0x63, 0x04, 0x80, 0x94, 0x1c, 0x2e, 0x69, 0x88,
	0x61, 0x4e, 0xde, 0x5c, 0x7c, 0xc9, 0xf9, 0xb9, 0x48, 0x8a, 0x9c, 0x78, 0xc0, 0x2e, 0x2e, 0xc8,
	0x0c, 0x00, 0x45, 0x53, 0x00, 0x63, 0x14, 0x27, 0x54, 0xae, 0xcc, 0x70, 0x11, 0x13, 0x73, 0x70,
	0x44, 0xc4, 0x2a, 0x26, 0xae, 0x60, 0xa8, 0xea, 0x30, 0xc3, 0x53, 0x70, 0x4e, 0x4c, 0x98, 0x61,
	0x12, 0x1b, 0x38, 0x6a, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x08, 0x3b, 0x5a, 0x09, 0x0a,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UserAPIClient is the client API for UserAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserAPIClient interface {
	// ListUsers returns a list of users. The order is unspecified but deterministic.
	ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error)
	// CreateUser creates a new user. The new user must have a unique email address.
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
}

type userAPIClient struct {
	cc *grpc.ClientConn
}

func NewUserAPIClient(cc *grpc.ClientConn) UserAPIClient {
	return &userAPIClient{cc}
}

func (c *userAPIClient) ListUsers(ctx context.Context, in *ListUsersRequest, opts ...grpc.CallOption) (*ListUsersResponse, error) {
	out := new(ListUsersResponse)
	err := c.cc.Invoke(ctx, "/strecku.v1.UserAPI/ListUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAPIClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/strecku.v1.UserAPI/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserAPIServer is the server API for UserAPI service.
type UserAPIServer interface {
	// ListUsers returns a list of users. The order is unspecified but deterministic.
	ListUsers(context.Context, *ListUsersRequest) (*ListUsersResponse, error)
	// CreateUser creates a new user. The new user must have a unique email address.
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
}

// UnimplementedUserAPIServer can be embedded to have forward compatible implementations.
type UnimplementedUserAPIServer struct {
}

func (*UnimplementedUserAPIServer) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (*UnimplementedUserAPIServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}

func RegisterUserAPIServer(s *grpc.Server, srv UserAPIServer) {
	s.RegisterService(&_UserAPI_serviceDesc, srv)
}

func _UserAPI_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAPIServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strecku.v1.UserAPI/ListUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAPIServer).ListUsers(ctx, req.(*ListUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAPI_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAPIServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strecku.v1.UserAPI/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAPIServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "strecku.v1.UserAPI",
	HandlerType: (*UserAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListUsers",
			Handler:    _UserAPI_ListUsers_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _UserAPI_CreateUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/user_api.proto",
}
