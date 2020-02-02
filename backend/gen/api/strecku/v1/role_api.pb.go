// Code generated by protoc-gen-go. DO NOT EDIT.
// source: strecku/v1/role_api.proto

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

type ListRolesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRolesRequest) Reset()         { *m = ListRolesRequest{} }
func (m *ListRolesRequest) String() string { return proto.CompactTextString(m) }
func (*ListRolesRequest) ProtoMessage()    {}
func (*ListRolesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_440647d807553baf, []int{0}
}

func (m *ListRolesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRolesRequest.Unmarshal(m, b)
}
func (m *ListRolesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRolesRequest.Marshal(b, m, deterministic)
}
func (m *ListRolesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRolesRequest.Merge(m, src)
}
func (m *ListRolesRequest) XXX_Size() int {
	return xxx_messageInfo_ListRolesRequest.Size(m)
}
func (m *ListRolesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRolesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRolesRequest proto.InternalMessageInfo

type ListRolesResponse struct {
	// The `roles` field contains all `Role` resources.
	Roles                []*Role  `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRolesResponse) Reset()         { *m = ListRolesResponse{} }
func (m *ListRolesResponse) String() string { return proto.CompactTextString(m) }
func (*ListRolesResponse) ProtoMessage()    {}
func (*ListRolesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_440647d807553baf, []int{1}
}

func (m *ListRolesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRolesResponse.Unmarshal(m, b)
}
func (m *ListRolesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRolesResponse.Marshal(b, m, deterministic)
}
func (m *ListRolesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRolesResponse.Merge(m, src)
}
func (m *ListRolesResponse) XXX_Size() int {
	return xxx_messageInfo_ListRolesResponse.Size(m)
}
func (m *ListRolesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRolesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListRolesResponse proto.InternalMessageInfo

func (m *ListRolesResponse) GetRoles() []*Role {
	if m != nil {
		return m.Roles
	}
	return nil
}

type GetRoleRequest struct {
	// The `name` field contains the resource name of a role.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRoleRequest) Reset()         { *m = GetRoleRequest{} }
func (m *GetRoleRequest) String() string { return proto.CompactTextString(m) }
func (*GetRoleRequest) ProtoMessage()    {}
func (*GetRoleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_440647d807553baf, []int{2}
}

func (m *GetRoleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRoleRequest.Unmarshal(m, b)
}
func (m *GetRoleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRoleRequest.Marshal(b, m, deterministic)
}
func (m *GetRoleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRoleRequest.Merge(m, src)
}
func (m *GetRoleRequest) XXX_Size() int {
	return xxx_messageInfo_GetRoleRequest.Size(m)
}
func (m *GetRoleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRoleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRoleRequest proto.InternalMessageInfo

func (m *GetRoleRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetRoleResponse struct {
	// The `role` field contains the `Role` resource.
	Role                 *Role    `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRoleResponse) Reset()         { *m = GetRoleResponse{} }
func (m *GetRoleResponse) String() string { return proto.CompactTextString(m) }
func (*GetRoleResponse) ProtoMessage()    {}
func (*GetRoleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_440647d807553baf, []int{3}
}

func (m *GetRoleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRoleResponse.Unmarshal(m, b)
}
func (m *GetRoleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRoleResponse.Marshal(b, m, deterministic)
}
func (m *GetRoleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRoleResponse.Merge(m, src)
}
func (m *GetRoleResponse) XXX_Size() int {
	return xxx_messageInfo_GetRoleResponse.Size(m)
}
func (m *GetRoleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRoleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetRoleResponse proto.InternalMessageInfo

func (m *GetRoleResponse) GetRole() *Role {
	if m != nil {
		return m.Role
	}
	return nil
}

type CreateRoleRequest struct {
	// The `Role` resource to be created. The `user_name` field of the `Role` must
	// contain a valid `User` resource name, and the `store_name` field must
	// contain a valid `Store` resource name. Anything else is an error.
	Role                 *Role    `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRoleRequest) Reset()         { *m = CreateRoleRequest{} }
func (m *CreateRoleRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRoleRequest) ProtoMessage()    {}
func (*CreateRoleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_440647d807553baf, []int{4}
}

func (m *CreateRoleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRoleRequest.Unmarshal(m, b)
}
func (m *CreateRoleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRoleRequest.Marshal(b, m, deterministic)
}
func (m *CreateRoleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRoleRequest.Merge(m, src)
}
func (m *CreateRoleRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRoleRequest.Size(m)
}
func (m *CreateRoleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRoleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRoleRequest proto.InternalMessageInfo

func (m *CreateRoleRequest) GetRole() *Role {
	if m != nil {
		return m.Role
	}
	return nil
}

type CreateRoleResponse struct {
	// The `Role` resource that was created. It is equal to the `Role` that was
	// provided in the request, except for the `name` field which has now been set
	// to the resource name of the newly created `Role` resource.
	Role                 *Role    `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRoleResponse) Reset()         { *m = CreateRoleResponse{} }
func (m *CreateRoleResponse) String() string { return proto.CompactTextString(m) }
func (*CreateRoleResponse) ProtoMessage()    {}
func (*CreateRoleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_440647d807553baf, []int{5}
}

func (m *CreateRoleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRoleResponse.Unmarshal(m, b)
}
func (m *CreateRoleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRoleResponse.Marshal(b, m, deterministic)
}
func (m *CreateRoleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRoleResponse.Merge(m, src)
}
func (m *CreateRoleResponse) XXX_Size() int {
	return xxx_messageInfo_CreateRoleResponse.Size(m)
}
func (m *CreateRoleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRoleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRoleResponse proto.InternalMessageInfo

func (m *CreateRoleResponse) GetRole() *Role {
	if m != nil {
		return m.Role
	}
	return nil
}

func init() {
	proto.RegisterType((*ListRolesRequest)(nil), "strecku.v1.ListRolesRequest")
	proto.RegisterType((*ListRolesResponse)(nil), "strecku.v1.ListRolesResponse")
	proto.RegisterType((*GetRoleRequest)(nil), "strecku.v1.GetRoleRequest")
	proto.RegisterType((*GetRoleResponse)(nil), "strecku.v1.GetRoleResponse")
	proto.RegisterType((*CreateRoleRequest)(nil), "strecku.v1.CreateRoleRequest")
	proto.RegisterType((*CreateRoleResponse)(nil), "strecku.v1.CreateRoleResponse")
}

func init() { proto.RegisterFile("strecku/v1/role_api.proto", fileDescriptor_440647d807553baf) }

var fileDescriptor_440647d807553baf = []byte{
	// 306 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0x2e, 0x29, 0x4a,
	0x4d, 0xce, 0x2e, 0xd5, 0x2f, 0x33, 0xd4, 0x2f, 0xca, 0xcf, 0x49, 0x8d, 0x4f, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x82, 0x4a, 0xe9, 0x95, 0x19, 0x4a, 0x89, 0xa2, 0x29,
	0x83, 0x28, 0x51, 0x12, 0xe2, 0x12, 0xf0, 0xc9, 0x2c, 0x2e, 0x09, 0xca, 0xcf, 0x49, 0x2d, 0x0e,
	0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x51, 0xb2, 0xe6, 0x12, 0x44, 0x12, 0x2b, 0x2e, 0xc8, 0xcf,
	0x2b, 0x4e, 0x15, 0x52, 0xe3, 0x62, 0x05, 0x69, 0x2b, 0x96, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x36,
	0x12, 0xd0, 0x43, 0x98, 0xad, 0x07, 0x52, 0x19, 0x04, 0x91, 0x56, 0x52, 0xe1, 0xe2, 0x73, 0x4f,
	0x05, 0xeb, 0x85, 0x1a, 0x27, 0x24, 0xc4, 0xc5, 0x92, 0x97, 0x98, 0x9b, 0x2a, 0xc1, 0xa8, 0xc0,
	0xa8, 0xc1, 0x19, 0x04, 0x66, 0x2b, 0x99, 0x73, 0xf1, 0xc3, 0x55, 0x41, 0x2d, 0x50, 0xe1, 0x62,
	0x01, 0x99, 0x00, 0x56, 0x86, 0xcd, 0x7c, 0xb0, 0xac, 0x92, 0x25, 0x97, 0xa0, 0x73, 0x51, 0x6a,
	0x62, 0x49, 0x2a, 0xb2, 0x0d, 0xc4, 0x69, 0xb5, 0xe2, 0x12, 0x42, 0xd6, 0x4a, 0x8a, 0xb5, 0x46,
	0x4f, 0x18, 0xb9, 0xd8, 0x41, 0x5c, 0xc7, 0x00, 0x4f, 0x21, 0x0f, 0x2e, 0x4e, 0x78, 0xf0, 0x08,
	0xc9, 0x20, 0x6b, 0x40, 0x0f, 0x49, 0x29, 0x59, 0x1c, 0xb2, 0x50, 0xbb, 0x9d, 0xb8, 0xd8, 0xa1,
	0xa1, 0x20, 0x24, 0x85, 0xac, 0x12, 0x35, 0x00, 0xa5, 0xa4, 0xb1, 0xca, 0x41, 0xcd, 0xf0, 0xe6,
	0xe2, 0x42, 0xf8, 0x4a, 0x08, 0xc5, 0x42, 0x8c, 0x80, 0x92, 0x92, 0xc3, 0x25, 0x0d, 0x31, 0xcc,
	0xc9, 0x9b, 0x8b, 0x2f, 0x39, 0x3f, 0x17, 0x49, 0x91, 0x13, 0x0f, 0xd8, 0xd7, 0x05, 0x99, 0x01,
	0xa0, 0xd4, 0x12, 0xc0, 0x18, 0xc5, 0x09, 0x95, 0x2b, 0x33, 0x5c, 0xc4, 0xc4, 0x1c, 0x1c, 0x11,
	0xb1, 0x8a, 0x89, 0x2b, 0x18, 0xaa, 0x3a, 0xcc, 0xf0, 0x14, 0x9c, 0x13, 0x13, 0x66, 0x98, 0xc4,
	0x06, 0x4e, 0x61, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdd, 0x2f, 0xe4, 0xbd, 0xa1, 0x02,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RoleAPIClient is the client API for RoleAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RoleAPIClient interface {
	// The `ListRoles` RPC lists `Role` resources.
	ListRoles(ctx context.Context, in *ListRolesRequest, opts ...grpc.CallOption) (*ListRolesResponse, error)
	// The `GetRole` RPC gets a single role by its resource name.
	GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*GetRoleResponse, error)
	// The `CreateRole` RPC allows for the creation of new `Role` resources.
	CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error)
}

type roleAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewRoleAPIClient(cc grpc.ClientConnInterface) RoleAPIClient {
	return &roleAPIClient{cc}
}

func (c *roleAPIClient) ListRoles(ctx context.Context, in *ListRolesRequest, opts ...grpc.CallOption) (*ListRolesResponse, error) {
	out := new(ListRolesResponse)
	err := c.cc.Invoke(ctx, "/strecku.v1.RoleAPI/ListRoles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAPIClient) GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*GetRoleResponse, error) {
	out := new(GetRoleResponse)
	err := c.cc.Invoke(ctx, "/strecku.v1.RoleAPI/GetRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleAPIClient) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error) {
	out := new(CreateRoleResponse)
	err := c.cc.Invoke(ctx, "/strecku.v1.RoleAPI/CreateRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoleAPIServer is the server API for RoleAPI service.
type RoleAPIServer interface {
	// The `ListRoles` RPC lists `Role` resources.
	ListRoles(context.Context, *ListRolesRequest) (*ListRolesResponse, error)
	// The `GetRole` RPC gets a single role by its resource name.
	GetRole(context.Context, *GetRoleRequest) (*GetRoleResponse, error)
	// The `CreateRole` RPC allows for the creation of new `Role` resources.
	CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleResponse, error)
}

// UnimplementedRoleAPIServer can be embedded to have forward compatible implementations.
type UnimplementedRoleAPIServer struct {
}

func (*UnimplementedRoleAPIServer) ListRoles(ctx context.Context, req *ListRolesRequest) (*ListRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoles not implemented")
}
func (*UnimplementedRoleAPIServer) GetRole(ctx context.Context, req *GetRoleRequest) (*GetRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRole not implemented")
}
func (*UnimplementedRoleAPIServer) CreateRole(ctx context.Context, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRole not implemented")
}

func RegisterRoleAPIServer(s *grpc.Server, srv RoleAPIServer) {
	s.RegisterService(&_RoleAPI_serviceDesc, srv)
}

func _RoleAPI_ListRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAPIServer).ListRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strecku.v1.RoleAPI/ListRoles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAPIServer).ListRoles(ctx, req.(*ListRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAPI_GetRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAPIServer).GetRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strecku.v1.RoleAPI/GetRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAPIServer).GetRole(ctx, req.(*GetRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleAPI_CreateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleAPIServer).CreateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strecku.v1.RoleAPI/CreateRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleAPIServer).CreateRole(ctx, req.(*CreateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RoleAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "strecku.v1.RoleAPI",
	HandlerType: (*RoleAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListRoles",
			Handler:    _RoleAPI_ListRoles_Handler,
		},
		{
			MethodName: "GetRole",
			Handler:    _RoleAPI_GetRole_Handler,
		},
		{
			MethodName: "CreateRole",
			Handler:    _RoleAPI_CreateRole_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "strecku/v1/role_api.proto",
}