// Code generated by protoc-gen-go. DO NOT EDIT.
// source: testing/v1/reset_api.proto

package testingv1

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

type ResetRequest struct {
	Reason               string   `protobuf:"bytes,1,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResetRequest) Reset()         { *m = ResetRequest{} }
func (m *ResetRequest) String() string { return proto.CompactTextString(m) }
func (*ResetRequest) ProtoMessage()    {}
func (*ResetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac59e1f98999492, []int{0}
}

func (m *ResetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResetRequest.Unmarshal(m, b)
}
func (m *ResetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResetRequest.Marshal(b, m, deterministic)
}
func (m *ResetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResetRequest.Merge(m, src)
}
func (m *ResetRequest) XXX_Size() int {
	return xxx_messageInfo_ResetRequest.Size(m)
}
func (m *ResetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ResetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ResetRequest proto.InternalMessageInfo

func (m *ResetRequest) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type ResetResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResetResponse) Reset()         { *m = ResetResponse{} }
func (m *ResetResponse) String() string { return proto.CompactTextString(m) }
func (*ResetResponse) ProtoMessage()    {}
func (*ResetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac59e1f98999492, []int{1}
}

func (m *ResetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResetResponse.Unmarshal(m, b)
}
func (m *ResetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResetResponse.Marshal(b, m, deterministic)
}
func (m *ResetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResetResponse.Merge(m, src)
}
func (m *ResetResponse) XXX_Size() int {
	return xxx_messageInfo_ResetResponse.Size(m)
}
func (m *ResetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ResetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ResetResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ResetRequest)(nil), "strecku.testing.v1.ResetRequest")
	proto.RegisterType((*ResetResponse)(nil), "strecku.testing.v1.ResetResponse")
}

func init() { proto.RegisterFile("testing/v1/reset_api.proto", fileDescriptor_3ac59e1f98999492) }

var fileDescriptor_3ac59e1f98999492 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0x49, 0x2d, 0x2e,
	0xc9, 0xcc, 0x4b, 0xd7, 0x2f, 0x33, 0xd4, 0x2f, 0x4a, 0x2d, 0x4e, 0x2d, 0x89, 0x4f, 0x2c, 0xc8,
	0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x2a, 0x2e, 0x29, 0x4a, 0x4d, 0xce, 0x2e, 0xd5,
	0x83, 0xaa, 0xd1, 0x2b, 0x33, 0x54, 0x52, 0xe3, 0xe2, 0x09, 0x02, 0x29, 0x0b, 0x4a, 0x2d, 0x2c,
	0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe3, 0x62, 0x2b, 0x4a, 0x4d, 0x2c, 0xce, 0xcf, 0x93, 0x60, 0x54,
	0x60, 0xd4, 0xe0, 0x0c, 0x82, 0xf2, 0x94, 0xf8, 0xb9, 0x78, 0xa1, 0xea, 0x8a, 0x0b, 0xf2, 0xf3,
	0x8a, 0x53, 0x8d, 0x22, 0xb8, 0x38, 0xc0, 0x02, 0x8e, 0x01, 0x9e, 0x42, 0x3e, 0x5c, 0xac, 0x60,
	0xb6, 0x90, 0x82, 0x1e, 0xa6, 0x15, 0x7a, 0xc8, 0xe6, 0x4b, 0x29, 0xe2, 0x51, 0x01, 0x31, 0xd9,
	0x29, 0x85, 0x4b, 0x2c, 0x39, 0x3f, 0x17, 0x8b, 0x3a, 0x27, 0x88, 0x13, 0x1c, 0x0b, 0x32, 0x03,
	0x40, 0xfe, 0x09, 0x60, 0x8c, 0xe2, 0x84, 0x4a, 0x96, 0x19, 0x2e, 0x62, 0x62, 0x0e, 0x0e, 0x89,
	0x58, 0xc5, 0x24, 0x14, 0x0c, 0xd5, 0x16, 0x02, 0xd5, 0x16, 0x66, 0x78, 0x0a, 0x2e, 0x18, 0x03,
	0x15, 0x8c, 0x09, 0x33, 0x4c, 0x62, 0x03, 0x87, 0x89, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xf5,
	0xfd, 0x3d, 0xfb, 0x31, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ResetAPIClient is the client API for ResetAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ResetAPIClient interface {
	// The `Reset` RPC causes the state to be reset as per the description of this
	// service.
	Reset(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*ResetResponse, error)
}

type resetAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewResetAPIClient(cc grpc.ClientConnInterface) ResetAPIClient {
	return &resetAPIClient{cc}
}

func (c *resetAPIClient) Reset(ctx context.Context, in *ResetRequest, opts ...grpc.CallOption) (*ResetResponse, error) {
	out := new(ResetResponse)
	err := c.cc.Invoke(ctx, "/strecku.testing.v1.ResetAPI/Reset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResetAPIServer is the server API for ResetAPI service.
type ResetAPIServer interface {
	// The `Reset` RPC causes the state to be reset as per the description of this
	// service.
	Reset(context.Context, *ResetRequest) (*ResetResponse, error)
}

// UnimplementedResetAPIServer can be embedded to have forward compatible implementations.
type UnimplementedResetAPIServer struct {
}

func (*UnimplementedResetAPIServer) Reset(ctx context.Context, req *ResetRequest) (*ResetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reset not implemented")
}

func RegisterResetAPIServer(s *grpc.Server, srv ResetAPIServer) {
	s.RegisterService(&_ResetAPI_serviceDesc, srv)
}

func _ResetAPI_Reset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResetAPIServer).Reset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/strecku.testing.v1.ResetAPI/Reset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResetAPIServer).Reset(ctx, req.(*ResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ResetAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "strecku.testing.v1.ResetAPI",
	HandlerType: (*ResetAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Reset",
			Handler:    _ResetAPI_Reset_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "testing/v1/reset_api.proto",
}
