// Code generated by protoc-gen-go. DO NOT EDIT.
// source: file.proto

package model // import "github.com/saichler/syncit/model"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type File struct {
	NameA                string   `protobuf:"bytes,1,opt,name=nameA,proto3" json:"nameA,omitempty"`
	SizeA                int64    `protobuf:"varint,2,opt,name=sizeA,proto3" json:"sizeA,omitempty"`
	DateA                int64    `protobuf:"varint,3,opt,name=dateA,proto3" json:"dateA,omitempty"`
	HashA                string   `protobuf:"bytes,4,opt,name=hashA,proto3" json:"hashA,omitempty"`
	NameZ                string   `protobuf:"bytes,5,opt,name=nameZ,proto3" json:"nameZ,omitempty"`
	SizeZ                int64    `protobuf:"varint,6,opt,name=sizeZ,proto3" json:"sizeZ,omitempty"`
	DateZ                int64    `protobuf:"varint,7,opt,name=dateZ,proto3" json:"dateZ,omitempty"`
	HashZ                string   `protobuf:"bytes,8,opt,name=hashZ,proto3" json:"hashZ,omitempty"`
	Files                []*File  `protobuf:"bytes,9,rep,name=files,proto3" json:"files,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *File) Reset()         { *m = File{} }
func (m *File) String() string { return proto.CompactTextString(m) }
func (*File) ProtoMessage()    {}
func (*File) Descriptor() ([]byte, []int) {
	return fileDescriptor_file_b1c937bf0d874d4a, []int{0}
}
func (m *File) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_File.Unmarshal(m, b)
}
func (m *File) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_File.Marshal(b, m, deterministic)
}
func (dst *File) XXX_Merge(src proto.Message) {
	xxx_messageInfo_File.Merge(dst, src)
}
func (m *File) XXX_Size() int {
	return xxx_messageInfo_File.Size(m)
}
func (m *File) XXX_DiscardUnknown() {
	xxx_messageInfo_File.DiscardUnknown(m)
}

var xxx_messageInfo_File proto.InternalMessageInfo

func (m *File) GetNameA() string {
	if m != nil {
		return m.NameA
	}
	return ""
}

func (m *File) GetSizeA() int64 {
	if m != nil {
		return m.SizeA
	}
	return 0
}

func (m *File) GetDateA() int64 {
	if m != nil {
		return m.DateA
	}
	return 0
}

func (m *File) GetHashA() string {
	if m != nil {
		return m.HashA
	}
	return ""
}

func (m *File) GetNameZ() string {
	if m != nil {
		return m.NameZ
	}
	return ""
}

func (m *File) GetSizeZ() int64 {
	if m != nil {
		return m.SizeZ
	}
	return 0
}

func (m *File) GetDateZ() int64 {
	if m != nil {
		return m.DateZ
	}
	return 0
}

func (m *File) GetHashZ() string {
	if m != nil {
		return m.HashZ
	}
	return ""
}

func (m *File) GetFiles() []*File {
	if m != nil {
		return m.Files
	}
	return nil
}

type Command struct {
	Cli                  string   `protobuf:"bytes,1,opt,name=cli,proto3" json:"cli,omitempty"`
	Args                 []string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
	Response             []byte   `protobuf:"bytes,3,opt,name=response,proto3" json:"response,omitempty"`
	Id                   string   `protobuf:"bytes,4,opt,name=id,proto3" json:"id,omitempty"`
	ResponseId           int32    `protobuf:"varint,5,opt,name=responseId,proto3" json:"responseId,omitempty"`
	ResponseCount        int32    `protobuf:"varint,6,opt,name=responseCount,proto3" json:"responseCount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_file_b1c937bf0d874d4a, []int{1}
}
func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (dst *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(dst, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetCli() string {
	if m != nil {
		return m.Cli
	}
	return ""
}

func (m *Command) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

func (m *Command) GetResponse() []byte {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *Command) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Command) GetResponseId() int32 {
	if m != nil {
		return m.ResponseId
	}
	return 0
}

func (m *Command) GetResponseCount() int32 {
	if m != nil {
		return m.ResponseCount
	}
	return 0
}

type UserPass struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Token                string   `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserPass) Reset()         { *m = UserPass{} }
func (m *UserPass) String() string { return proto.CompactTextString(m) }
func (*UserPass) ProtoMessage()    {}
func (*UserPass) Descriptor() ([]byte, []int) {
	return fileDescriptor_file_b1c937bf0d874d4a, []int{2}
}
func (m *UserPass) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserPass.Unmarshal(m, b)
}
func (m *UserPass) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserPass.Marshal(b, m, deterministic)
}
func (dst *UserPass) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserPass.Merge(dst, src)
}
func (m *UserPass) XXX_Size() int {
	return xxx_messageInfo_UserPass.Size(m)
}
func (m *UserPass) XXX_DiscardUnknown() {
	xxx_messageInfo_UserPass.DiscardUnknown(m)
}

var xxx_messageInfo_UserPass proto.InternalMessageInfo

func (m *UserPass) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserPass) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *UserPass) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*File)(nil), "file.File")
	proto.RegisterType((*Command)(nil), "file.Command")
	proto.RegisterType((*UserPass)(nil), "file.UserPass")
}

func init() { proto.RegisterFile("file.proto", fileDescriptor_file_b1c937bf0d874d4a) }

var fileDescriptor_file_b1c937bf0d874d4a = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xcd, 0x4a, 0xf4, 0x30,
	0x14, 0x86, 0xe9, 0xdf, 0xcc, 0xf4, 0x7c, 0x9f, 0x22, 0xc1, 0x45, 0x70, 0x21, 0xa5, 0xb8, 0xe8,
	0x6a, 0x06, 0xf4, 0x0a, 0xea, 0x80, 0xe0, 0x4e, 0x02, 0x82, 0x64, 0x97, 0x69, 0xe2, 0x34, 0xd8,
	0x36, 0x43, 0x4e, 0x07, 0xd1, 0xbb, 0xf1, 0xae, 0xbc, 0x1c, 0x49, 0x3a, 0x6d, 0x75, 0x77, 0x9e,
	0xf7, 0xc0, 0xd3, 0xbc, 0x69, 0x00, 0x5e, 0x75, 0xa3, 0xd6, 0x07, 0x6b, 0x7a, 0x43, 0x62, 0x37,
	0xe7, 0xdf, 0x01, 0xc4, 0x0f, 0xba, 0x51, 0xe4, 0x12, 0x92, 0x4e, 0xb4, 0xaa, 0xa4, 0x41, 0x16,
	0x14, 0x29, 0x1b, 0xc0, 0xa5, 0xa8, 0x3f, 0x55, 0x49, 0xc3, 0x2c, 0x28, 0x22, 0x36, 0x80, 0x4b,
	0xa5, 0xe8, 0x55, 0x49, 0xa3, 0x21, 0xf5, 0xe0, 0xd2, 0x5a, 0x60, 0x5d, 0xd2, 0x78, 0x30, 0x78,
	0x18, 0xbd, 0x9c, 0x26, 0xb3, 0x97, 0x8f, 0x5e, 0x4e, 0x17, 0xb3, 0x97, 0x8f, 0x5e, 0x4e, 0x97,
	0xb3, 0x97, 0x8f, 0x5e, 0x4e, 0x57, 0xb3, 0x97, 0x93, 0x0c, 0x12, 0x57, 0x00, 0x69, 0x9a, 0x45,
	0xc5, 0xbf, 0x5b, 0x58, 0xfb, 0x6a, 0xae, 0x0a, 0x1b, 0x16, 0xf9, 0x57, 0x00, 0xcb, 0xad, 0x69,
	0x5b, 0xd1, 0x49, 0x72, 0x01, 0x51, 0xd5, 0xe8, 0x53, 0x37, 0x37, 0x12, 0x02, 0xb1, 0xb0, 0x7b,
	0xa4, 0x61, 0x16, 0x15, 0x29, 0xf3, 0x33, 0xb9, 0x82, 0x95, 0x55, 0x78, 0x30, 0x1d, 0x2a, 0x5f,
	0xed, 0x3f, 0x9b, 0x98, 0x9c, 0x43, 0xa8, 0xe5, 0xa9, 0x5a, 0xa8, 0x25, 0xb9, 0x06, 0x18, 0x77,
	0x8f, 0xd2, 0x97, 0x4b, 0xd8, 0xaf, 0x84, 0xdc, 0xc0, 0xd9, 0x48, 0x5b, 0x73, 0xec, 0x7a, 0xdf,
	0x34, 0x61, 0x7f, 0xc3, 0xfc, 0x05, 0x56, 0xcf, 0xa8, 0xec, 0x93, 0x40, 0xff, 0xf5, 0x23, 0x2a,
	0xeb, 0x2e, 0xe8, 0x74, 0xd0, 0x89, 0xdd, 0xee, 0x20, 0x10, 0xdf, 0x8d, 0x95, 0xfe, 0x57, 0xa4,
	0x6c, 0x62, 0x77, 0x3f, 0xbd, 0x79, 0x53, 0x9d, 0x3f, 0x72, 0xca, 0x06, 0xb8, 0xcf, 0x79, 0xb6,
	0xd7, 0x7d, 0x7d, 0xdc, 0xad, 0x2b, 0xd3, 0x6e, 0x50, 0xe8, 0xaa, 0x6e, 0x94, 0xdd, 0xe0, 0x47,
	0x57, 0xe9, 0x7e, 0xd3, 0x1a, 0xa9, 0x9a, 0xdd, 0xc2, 0xbf, 0x84, 0xbb, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x26, 0x48, 0x84, 0xa3, 0x17, 0x02, 0x00, 0x00,
}
