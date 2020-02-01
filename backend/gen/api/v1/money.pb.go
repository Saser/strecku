// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/money.proto

package streckuv1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// A `Money` represents an amount of money that is negative, zero, or
// positive. It is assumed to be given in the currency of `SEK`.
type Money struct {
	// The `units` field represents the whole units of `SEK` of the amount. For
	// example, if the amount is `5.50 SEK`, then `units` is 5. As another
	// example, if the amount is `-4.75 SEK`, then `units` is -4.
	Units int64 `protobuf:"varint,1,opt,name=units,proto3" json:"units,omitempty"`
	// The `cents` field represents the number of cents ("ören") of the
	// amount. The value must be between -99 and +99, inclusive. The value can be
	// negative, zero, or positive:
	// * if `units` is negative, then `cents` must be zero or negative.
	// * if `units` is zero, then `cents` can be negative, zero, or positive.
	// * if `units` is positive, then `cents` must be zero or positive.
	Cents                int32    `protobuf:"varint,2,opt,name=cents,proto3" json:"cents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Money) Reset()         { *m = Money{} }
func (m *Money) String() string { return proto.CompactTextString(m) }
func (*Money) ProtoMessage()    {}
func (*Money) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d61d2520a142386, []int{0}
}

func (m *Money) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Money.Unmarshal(m, b)
}
func (m *Money) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Money.Marshal(b, m, deterministic)
}
func (m *Money) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Money.Merge(m, src)
}
func (m *Money) XXX_Size() int {
	return xxx_messageInfo_Money.Size(m)
}
func (m *Money) XXX_DiscardUnknown() {
	xxx_messageInfo_Money.DiscardUnknown(m)
}

var xxx_messageInfo_Money proto.InternalMessageInfo

func (m *Money) GetUnits() int64 {
	if m != nil {
		return m.Units
	}
	return 0
}

func (m *Money) GetCents() int32 {
	if m != nil {
		return m.Cents
	}
	return 0
}

func init() {
	proto.RegisterType((*Money)(nil), "strecku.v1.Money")
}

func init() { proto.RegisterFile("v1/money.proto", fileDescriptor_6d61d2520a142386) }

var fileDescriptor_6d61d2520a142386 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x33, 0xd4, 0xcf,
	0xcd, 0xcf, 0x4b, 0xad, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2a, 0x2e, 0x29, 0x4a,
	0x4d, 0xce, 0x2e, 0xd5, 0x2b, 0x33, 0x54, 0x32, 0xe6, 0x62, 0xf5, 0x05, 0x49, 0x09, 0x89, 0x70,
	0xb1, 0x96, 0xe6, 0x65, 0x96, 0x14, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x30, 0x07, 0x41, 0x38, 0x20,
	0xd1, 0xe4, 0xd4, 0xbc, 0x92, 0x62, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x08, 0xc7, 0xc9,
	0x93, 0x8b, 0x2f, 0x39, 0x3f, 0x57, 0x0f, 0x61, 0x8c, 0x13, 0x17, 0xd8, 0x90, 0x00, 0x90, 0xf1,
	0x01, 0x8c, 0x51, 0x9c, 0x50, 0x99, 0x32, 0xc3, 0x45, 0x4c, 0xcc, 0xc1, 0x11, 0x11, 0xab, 0x98,
	0xb8, 0x82, 0xa1, 0x6a, 0xc3, 0x0c, 0x4f, 0xc1, 0x39, 0x31, 0x61, 0x86, 0x49, 0x6c, 0x60, 0x27,
	0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x82, 0xbb, 0x06, 0xa4, 0x00, 0x00, 0x00,
}
