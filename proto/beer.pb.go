// Code generated by protoc-gen-go.
// source: beer.proto
// DO NOT EDIT!

/*
Package beer is a generated protocol buffer package.

It is generated from these files:
	beer.proto

It has these top-level messages:
	Empty
	Token
	Beer
	Cellar
	BeerCellar
*/
package beer

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

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Token struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Secret string `protobuf:"bytes,2,opt,name=secret" json:"secret,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Beer struct {
	Id        int64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	DrinkDate int64  `protobuf:"varint,2,opt,name=drink_date,json=drinkDate" json:"drink_date,omitempty"`
	Size      string `protobuf:"bytes,3,opt,name=size" json:"size,omitempty"`
	Name      string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
}

func (m *Beer) Reset()                    { *m = Beer{} }
func (m *Beer) String() string            { return proto.CompactTextString(m) }
func (*Beer) ProtoMessage()               {}
func (*Beer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Cellar struct {
	Name  string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Beers []*Beer `protobuf:"bytes,2,rep,name=beers" json:"beers,omitempty"`
}

func (m *Cellar) Reset()                    { *m = Cellar{} }
func (m *Cellar) String() string            { return proto.CompactTextString(m) }
func (*Cellar) ProtoMessage()               {}
func (*Cellar) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Cellar) GetBeers() []*Beer {
	if m != nil {
		return m.Beers
	}
	return nil
}

type BeerCellar struct {
	Name          string    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	SyncTime      int64     `protobuf:"varint,2,opt,name=syncTime" json:"syncTime,omitempty"`
	UntappdKey    string    `protobuf:"bytes,3,opt,name=untappdKey" json:"untappdKey,omitempty"`
	UntappdSecret string    `protobuf:"bytes,4,opt,name=untappdSecret" json:"untappdSecret,omitempty"`
	Cellars       []*Cellar `protobuf:"bytes,5,rep,name=cellars" json:"cellars,omitempty"`
}

func (m *BeerCellar) Reset()                    { *m = BeerCellar{} }
func (m *BeerCellar) String() string            { return proto.CompactTextString(m) }
func (*BeerCellar) ProtoMessage()               {}
func (*BeerCellar) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *BeerCellar) GetCellars() []*Cellar {
	if m != nil {
		return m.Cellars
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "beer.Empty")
	proto.RegisterType((*Token)(nil), "beer.Token")
	proto.RegisterType((*Beer)(nil), "beer.Beer")
	proto.RegisterType((*Cellar)(nil), "beer.Cellar")
	proto.RegisterType((*BeerCellar)(nil), "beer.BeerCellar")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for BeerCellarService service

type BeerCellarServiceClient interface {
	AddBeer(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Cellar, error)
	GetBeer(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Beer, error)
	GetCellar(ctx context.Context, in *Cellar, opts ...grpc.CallOption) (*Cellar, error)
	RemoveBeer(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Beer, error)
	GetName(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Beer, error)
}

type beerCellarServiceClient struct {
	cc *grpc.ClientConn
}

func NewBeerCellarServiceClient(cc *grpc.ClientConn) BeerCellarServiceClient {
	return &beerCellarServiceClient{cc}
}

func (c *beerCellarServiceClient) AddBeer(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Cellar, error) {
	out := new(Cellar)
	err := grpc.Invoke(ctx, "/beer.BeerCellarService/AddBeer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beerCellarServiceClient) GetBeer(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Beer, error) {
	out := new(Beer)
	err := grpc.Invoke(ctx, "/beer.BeerCellarService/GetBeer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beerCellarServiceClient) GetCellar(ctx context.Context, in *Cellar, opts ...grpc.CallOption) (*Cellar, error) {
	out := new(Cellar)
	err := grpc.Invoke(ctx, "/beer.BeerCellarService/GetCellar", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beerCellarServiceClient) RemoveBeer(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Beer, error) {
	out := new(Beer)
	err := grpc.Invoke(ctx, "/beer.BeerCellarService/RemoveBeer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beerCellarServiceClient) GetName(ctx context.Context, in *Beer, opts ...grpc.CallOption) (*Beer, error) {
	out := new(Beer)
	err := grpc.Invoke(ctx, "/beer.BeerCellarService/GetName", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for BeerCellarService service

type BeerCellarServiceServer interface {
	AddBeer(context.Context, *Beer) (*Cellar, error)
	GetBeer(context.Context, *Beer) (*Beer, error)
	GetCellar(context.Context, *Cellar) (*Cellar, error)
	RemoveBeer(context.Context, *Beer) (*Beer, error)
	GetName(context.Context, *Beer) (*Beer, error)
}

func RegisterBeerCellarServiceServer(s *grpc.Server, srv BeerCellarServiceServer) {
	s.RegisterService(&_BeerCellarService_serviceDesc, srv)
}

func _BeerCellarService_AddBeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Beer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeerCellarServiceServer).AddBeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/beer.BeerCellarService/AddBeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeerCellarServiceServer).AddBeer(ctx, req.(*Beer))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeerCellarService_GetBeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Beer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeerCellarServiceServer).GetBeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/beer.BeerCellarService/GetBeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeerCellarServiceServer).GetBeer(ctx, req.(*Beer))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeerCellarService_GetCellar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cellar)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeerCellarServiceServer).GetCellar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/beer.BeerCellarService/GetCellar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeerCellarServiceServer).GetCellar(ctx, req.(*Cellar))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeerCellarService_RemoveBeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Beer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeerCellarServiceServer).RemoveBeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/beer.BeerCellarService/RemoveBeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeerCellarServiceServer).RemoveBeer(ctx, req.(*Beer))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeerCellarService_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Beer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeerCellarServiceServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/beer.BeerCellarService/GetName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeerCellarServiceServer).GetName(ctx, req.(*Beer))
	}
	return interceptor(ctx, in, info, handler)
}

var _BeerCellarService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "beer.BeerCellarService",
	HandlerType: (*BeerCellarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddBeer",
			Handler:    _BeerCellarService_AddBeer_Handler,
		},
		{
			MethodName: "GetBeer",
			Handler:    _BeerCellarService_GetBeer_Handler,
		},
		{
			MethodName: "GetCellar",
			Handler:    _BeerCellarService_GetCellar_Handler,
		},
		{
			MethodName: "RemoveBeer",
			Handler:    _BeerCellarService_RemoveBeer_Handler,
		},
		{
			MethodName: "GetName",
			Handler:    _BeerCellarService_GetName_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "beer.proto",
}

func init() { proto.RegisterFile("beer.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x92, 0xcd, 0x4a, 0xfb, 0x40,
	0x14, 0xc5, 0x9b, 0xaf, 0xe6, 0xdf, 0xf3, 0x57, 0xc1, 0x59, 0x48, 0x28, 0x28, 0x25, 0x6a, 0xa9,
	0x9b, 0x0a, 0x75, 0x2f, 0xf8, 0x45, 0x17, 0x82, 0x8b, 0xb4, 0x5b, 0x91, 0x34, 0x73, 0x17, 0xa1,
	0xcd, 0x07, 0x93, 0xb1, 0x50, 0x1f, 0xc9, 0x17, 0xf2, 0x75, 0x64, 0x66, 0xd2, 0x36, 0xa5, 0xa8,
	0xbb, 0x7b, 0xcf, 0xbd, 0x39, 0xf3, 0xbb, 0x87, 0x00, 0x33, 0x22, 0x31, 0x2c, 0x45, 0x21, 0x0b,
	0xe6, 0xaa, 0x3a, 0xf4, 0xe1, 0x3d, 0x65, 0xa5, 0x5c, 0x85, 0xd7, 0xf0, 0xa6, 0xc5, 0x9c, 0x72,
	0x76, 0x04, 0x3b, 0xe5, 0x81, 0xd5, 0xb3, 0x06, 0x9d, 0xc8, 0x4e, 0x39, 0x3b, 0x41, 0xbb, 0xa2,
	0x44, 0x90, 0x0c, 0x6c, 0xad, 0xd5, 0x5d, 0xf8, 0x0a, 0xf7, 0x9e, 0x48, 0x34, 0xf6, 0x1d, 0xbd,
	0x7f, 0x0a, 0x70, 0x91, 0xe6, 0xf3, 0x37, 0x1e, 0x4b, 0xd2, 0xdf, 0x38, 0x51, 0x47, 0x2b, 0x8f,
	0xb1, 0x24, 0xc6, 0xe0, 0x56, 0xe9, 0x07, 0x05, 0x8e, 0x36, 0xd3, 0xb5, 0xd2, 0xf2, 0x38, 0xa3,
	0xc0, 0x35, 0x9a, 0xaa, 0xc3, 0x5b, 0xb4, 0x1f, 0x68, 0xb1, 0x88, 0xc5, 0x66, 0x6a, 0x6d, 0xa7,
	0xac, 0x07, 0x4f, 0xe1, 0x57, 0x81, 0xdd, 0x73, 0x06, 0xff, 0x47, 0x18, 0xea, 0xc3, 0x14, 0x4f,
	0x64, 0x06, 0xe1, 0xa7, 0x05, 0xa8, 0xfe, 0x17, 0x93, 0x2e, 0xfe, 0x55, 0xab, 0x3c, 0x99, 0xa6,
	0xd9, 0x9a, 0x73, 0xd3, 0xb3, 0x33, 0xe0, 0x3d, 0x97, 0x71, 0x59, 0xf2, 0x67, 0x5a, 0xd5, 0xb0,
	0x0d, 0x85, 0x5d, 0xe0, 0xb0, 0xee, 0x26, 0x26, 0x1c, 0xc3, 0xbe, 0x2b, 0xb2, 0x3e, 0xfc, 0x44,
	0xbf, 0x5f, 0x05, 0x9e, 0x06, 0x3d, 0x30, 0xa0, 0x06, 0x2a, 0x5a, 0x0f, 0x47, 0x5f, 0x16, 0x8e,
	0xb7, 0xb0, 0x13, 0x12, 0xcb, 0x34, 0x21, 0x76, 0x09, 0xff, 0x8e, 0x73, 0x1d, 0x72, 0xe3, 0xc0,
	0xee, 0x8e, 0x47, 0xd8, 0x62, 0xe7, 0xf0, 0xc7, 0x24, 0xf7, 0xd6, 0x1a, 0x75, 0xd8, 0x62, 0x57,
	0xe8, 0x8c, 0x49, 0xd6, 0x61, 0xec, 0x38, 0xec, 0xf9, 0xf5, 0x81, 0x88, 0xb2, 0x62, 0x49, 0x7f,
	0x58, 0x9a, 0x77, 0x5f, 0x54, 0x92, 0x3f, 0x2e, 0xcd, 0xda, 0xfa, 0x67, 0xbb, 0xf9, 0x0e, 0x00,
	0x00, 0xff, 0xff, 0x1f, 0xef, 0xfa, 0xcb, 0x7a, 0x02, 0x00, 0x00,
}
