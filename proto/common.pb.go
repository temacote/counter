// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/common.proto

package counter

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

type EmptyMessage struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyMessage) Reset()         { *m = EmptyMessage{} }
func (m *EmptyMessage) String() string { return proto.CompactTextString(m) }
func (*EmptyMessage) ProtoMessage()    {}
func (*EmptyMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_1747d3070a2311a0, []int{0}
}

func (m *EmptyMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyMessage.Unmarshal(m, b)
}
func (m *EmptyMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyMessage.Marshal(b, m, deterministic)
}
func (m *EmptyMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyMessage.Merge(m, src)
}
func (m *EmptyMessage) XXX_Size() int {
	return xxx_messageInfo_EmptyMessage.Size(m)
}
func (m *EmptyMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyMessage.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyMessage proto.InternalMessageInfo

func init() {
	proto.RegisterType((*EmptyMessage)(nil), "counter.EmptyMessage")
}

func init() { proto.RegisterFile("proto/common.proto", fileDescriptor_1747d3070a2311a0) }

var fileDescriptor_1747d3070a2311a0 = []byte{
	// 71 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xd3, 0x03, 0x73, 0x84, 0xd8, 0x93, 0xf3, 0x4b,
	0xf3, 0x4a, 0x52, 0x8b, 0x94, 0xf8, 0xb8, 0x78, 0x5c, 0x73, 0x0b, 0x4a, 0x2a, 0x7d, 0x53, 0x8b,
	0x8b, 0x13, 0xd3, 0x53, 0x93, 0xd8, 0xc0, 0xf2, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xdf,
	0x66, 0xd2, 0x11, 0x35, 0x00, 0x00, 0x00,
}