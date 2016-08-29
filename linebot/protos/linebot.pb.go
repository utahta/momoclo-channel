// Code generated by protoc-gen-go.
// source: linebot/protos/linebot.proto
// DO NOT EDIT!

/*
Package linebot is a generated protocol buffer package.

It is generated from these files:
	linebot/protos/linebot.proto

It has these top-level messages:
	NotifyMessageRequest
	NotifyMessageResponse
*/
package linebot

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NotifyMessageRequest struct {
	To        []string `protobuf:"bytes,1,rep,name=to" json:"to,omitempty"`
	Text      string   `protobuf:"bytes,2,opt,name=text" json:"text,omitempty"`
	ImageUrls []string `protobuf:"bytes,3,rep,name=imageUrls" json:"imageUrls,omitempty"`
	VideoUrls []string `protobuf:"bytes,4,rep,name=videoUrls" json:"videoUrls,omitempty"`
}

func (m *NotifyMessageRequest) Reset()                    { *m = NotifyMessageRequest{} }
func (m *NotifyMessageRequest) String() string            { return proto.CompactTextString(m) }
func (*NotifyMessageRequest) ProtoMessage()               {}
func (*NotifyMessageRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type NotifyMessageResponse struct {
}

func (m *NotifyMessageResponse) Reset()                    { *m = NotifyMessageResponse{} }
func (m *NotifyMessageResponse) String() string            { return proto.CompactTextString(m) }
func (*NotifyMessageResponse) ProtoMessage()               {}
func (*NotifyMessageResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*NotifyMessageRequest)(nil), "linebot.NotifyMessageRequest")
	proto.RegisterType((*NotifyMessageResponse)(nil), "linebot.NotifyMessageResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for LineBot service

type LineBotClient interface {
	NotifyMessage(ctx context.Context, in *NotifyMessageRequest, opts ...grpc.CallOption) (*NotifyMessageResponse, error)
}

type lineBotClient struct {
	cc *grpc.ClientConn
}

func NewLineBotClient(cc *grpc.ClientConn) LineBotClient {
	return &lineBotClient{cc}
}

func (c *lineBotClient) NotifyMessage(ctx context.Context, in *NotifyMessageRequest, opts ...grpc.CallOption) (*NotifyMessageResponse, error) {
	out := new(NotifyMessageResponse)
	err := grpc.Invoke(ctx, "/linebot.LineBot/NotifyMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LineBot service

type LineBotServer interface {
	NotifyMessage(context.Context, *NotifyMessageRequest) (*NotifyMessageResponse, error)
}

func RegisterLineBotServer(s *grpc.Server, srv LineBotServer) {
	s.RegisterService(&_LineBot_serviceDesc, srv)
}

func _LineBot_NotifyMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LineBotServer).NotifyMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/linebot.LineBot/NotifyMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LineBotServer).NotifyMessage(ctx, req.(*NotifyMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LineBot_serviceDesc = grpc.ServiceDesc{
	ServiceName: "linebot.LineBot",
	HandlerType: (*LineBotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifyMessage",
			Handler:    _LineBot_NotifyMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("linebot/protos/linebot.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x92, 0xc9, 0xc9, 0xcc, 0x4b,
	0x4d, 0xca, 0x2f, 0xd1, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x2f, 0xd6, 0x87, 0x72, 0xf5, 0xc0, 0x5c,
	0x21, 0x76, 0x28, 0x57, 0xa9, 0x8c, 0x4b, 0xc4, 0x2f, 0xbf, 0x24, 0x33, 0xad, 0xd2, 0x37, 0xb5,
	0xb8, 0x38, 0x31, 0x3d, 0x35, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x88, 0x8f, 0x8b, 0xa9,
	0x24, 0x5f, 0x82, 0x51, 0x81, 0x59, 0x83, 0x33, 0x08, 0xc8, 0x12, 0x12, 0xe2, 0x62, 0x29, 0x49,
	0xad, 0x28, 0x91, 0x60, 0x52, 0x60, 0x04, 0x8a, 0x80, 0xd9, 0x42, 0x32, 0x5c, 0x9c, 0x99, 0xb9,
	0x40, 0x3d, 0xa1, 0x45, 0x39, 0xc5, 0x12, 0xcc, 0x60, 0xa5, 0x08, 0x01, 0x90, 0x6c, 0x59, 0x66,
	0x4a, 0x6a, 0x3e, 0x58, 0x96, 0x05, 0x22, 0x0b, 0x17, 0x50, 0x12, 0xe7, 0x12, 0x45, 0xb3, 0xb7,
	0xb8, 0x20, 0x3f, 0xaf, 0x38, 0xd5, 0x28, 0x9a, 0x8b, 0xdd, 0x07, 0xe8, 0x36, 0xa7, 0xfc, 0x12,
	0xa1, 0x00, 0x2e, 0x5e, 0x14, 0x35, 0x42, 0xb2, 0x7a, 0x30, 0x5f, 0x60, 0x73, 0xb3, 0x94, 0x1c,
	0x2e, 0x69, 0x88, 0xd1, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0xdf, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x04, 0xed, 0x88, 0xb3, 0x1d, 0x01, 0x00, 0x00,
}