// Code generated by protoc-gen-go. DO NOT EDIT.
// source: beer.proto

package beer

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Token struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Secret               string   `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	Rtoken               string   `protobuf:"bytes,3,opt,name=rtoken,proto3" json:"rtoken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{1}
}

func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Token) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *Token) GetRtoken() string {
	if m != nil {
		return m.Rtoken
	}
	return ""
}

type Config struct {
	Token                *Token   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Drunk                []*Beer  `protobuf:"bytes,2,rep,name=drunk,proto3" json:"drunk,omitempty"`
	Cellar               *Cellar  `protobuf:"bytes,3,opt,name=cellar,proto3" json:"cellar,omitempty"`
	LastSync             int64    `protobuf:"varint,4,opt,name=last_sync,json=lastSync,proto3" json:"last_sync,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{2}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

func (m *Config) GetDrunk() []*Beer {
	if m != nil {
		return m.Drunk
	}
	return nil
}

func (m *Config) GetCellar() *Cellar {
	if m != nil {
		return m.Cellar
	}
	return nil
}

func (m *Config) GetLastSync() int64 {
	if m != nil {
		return m.LastSync
	}
	return 0
}

type Beer struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DrinkDate            int64    `protobuf:"varint,2,opt,name=drink_date,json=drinkDate,proto3" json:"drink_date,omitempty"`
	Size                 string   `protobuf:"bytes,3,opt,name=size,proto3" json:"size,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	OnDeck               bool     `protobuf:"varint,5,opt,name=on_deck,json=onDeck,proto3" json:"on_deck,omitempty"`
	Abv                  float32  `protobuf:"fixed32,6,opt,name=abv,proto3" json:"abv,omitempty"`
	Index                int32    `protobuf:"varint,7,opt,name=index,proto3" json:"index,omitempty"`
	CheckinId            int32    `protobuf:"varint,8,opt,name=checkin_id,json=checkinId,proto3" json:"checkin_id,omitempty"`
	InCellar             int32    `protobuf:"varint,9,opt,name=in_cellar,json=inCellar,proto3" json:"in_cellar,omitempty"`
	Uid                  int64    `protobuf:"varint,10,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Beer) Reset()         { *m = Beer{} }
func (m *Beer) String() string { return proto.CompactTextString(m) }
func (*Beer) ProtoMessage()    {}
func (*Beer) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{3}
}

func (m *Beer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Beer.Unmarshal(m, b)
}
func (m *Beer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Beer.Marshal(b, m, deterministic)
}
func (m *Beer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Beer.Merge(m, src)
}
func (m *Beer) XXX_Size() int {
	return xxx_messageInfo_Beer.Size(m)
}
func (m *Beer) XXX_DiscardUnknown() {
	xxx_messageInfo_Beer.DiscardUnknown(m)
}

var xxx_messageInfo_Beer proto.InternalMessageInfo

func (m *Beer) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Beer) GetDrinkDate() int64 {
	if m != nil {
		return m.DrinkDate
	}
	return 0
}

func (m *Beer) GetSize() string {
	if m != nil {
		return m.Size
	}
	return ""
}

func (m *Beer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Beer) GetOnDeck() bool {
	if m != nil {
		return m.OnDeck
	}
	return false
}

func (m *Beer) GetAbv() float32 {
	if m != nil {
		return m.Abv
	}
	return 0
}

func (m *Beer) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Beer) GetCheckinId() int32 {
	if m != nil {
		return m.CheckinId
	}
	return 0
}

func (m *Beer) GetInCellar() int32 {
	if m != nil {
		return m.InCellar
	}
	return 0
}

func (m *Beer) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

type CellarSlot struct {
	Accepts              string   `protobuf:"bytes,1,opt,name=accepts,proto3" json:"accepts,omitempty"`
	NumSlots             int32    `protobuf:"varint,2,opt,name=num_slots,json=numSlots,proto3" json:"num_slots,omitempty"`
	Beers                []*Beer  `protobuf:"bytes,3,rep,name=beers,proto3" json:"beers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CellarSlot) Reset()         { *m = CellarSlot{} }
func (m *CellarSlot) String() string { return proto.CompactTextString(m) }
func (*CellarSlot) ProtoMessage()    {}
func (*CellarSlot) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{4}
}

func (m *CellarSlot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CellarSlot.Unmarshal(m, b)
}
func (m *CellarSlot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CellarSlot.Marshal(b, m, deterministic)
}
func (m *CellarSlot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CellarSlot.Merge(m, src)
}
func (m *CellarSlot) XXX_Size() int {
	return xxx_messageInfo_CellarSlot.Size(m)
}
func (m *CellarSlot) XXX_DiscardUnknown() {
	xxx_messageInfo_CellarSlot.DiscardUnknown(m)
}

var xxx_messageInfo_CellarSlot proto.InternalMessageInfo

func (m *CellarSlot) GetAccepts() string {
	if m != nil {
		return m.Accepts
	}
	return ""
}

func (m *CellarSlot) GetNumSlots() int32 {
	if m != nil {
		return m.NumSlots
	}
	return 0
}

func (m *CellarSlot) GetBeers() []*Beer {
	if m != nil {
		return m.Beers
	}
	return nil
}

type Cellar struct {
	Slots                []*CellarSlot `protobuf:"bytes,1,rep,name=slots,proto3" json:"slots,omitempty"`
	OnDeck               []*Beer       `protobuf:"bytes,2,rep,name=on_deck,json=onDeck,proto3" json:"on_deck,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Cellar) Reset()         { *m = Cellar{} }
func (m *Cellar) String() string { return proto.CompactTextString(m) }
func (*Cellar) ProtoMessage()    {}
func (*Cellar) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{5}
}

func (m *Cellar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Cellar.Unmarshal(m, b)
}
func (m *Cellar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Cellar.Marshal(b, m, deterministic)
}
func (m *Cellar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Cellar.Merge(m, src)
}
func (m *Cellar) XXX_Size() int {
	return xxx_messageInfo_Cellar.Size(m)
}
func (m *Cellar) XXX_DiscardUnknown() {
	xxx_messageInfo_Cellar.DiscardUnknown(m)
}

var xxx_messageInfo_Cellar proto.InternalMessageInfo

func (m *Cellar) GetSlots() []*CellarSlot {
	if m != nil {
		return m.Slots
	}
	return nil
}

func (m *Cellar) GetOnDeck() []*Beer {
	if m != nil {
		return m.OnDeck
	}
	return nil
}

type AddBeerRequest struct {
	Beer                 *Beer    `protobuf:"bytes,1,opt,name=beer,proto3" json:"beer,omitempty"`
	Quantity             int32    `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddBeerRequest) Reset()         { *m = AddBeerRequest{} }
func (m *AddBeerRequest) String() string { return proto.CompactTextString(m) }
func (*AddBeerRequest) ProtoMessage()    {}
func (*AddBeerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{6}
}

func (m *AddBeerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddBeerRequest.Unmarshal(m, b)
}
func (m *AddBeerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddBeerRequest.Marshal(b, m, deterministic)
}
func (m *AddBeerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddBeerRequest.Merge(m, src)
}
func (m *AddBeerRequest) XXX_Size() int {
	return xxx_messageInfo_AddBeerRequest.Size(m)
}
func (m *AddBeerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddBeerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddBeerRequest proto.InternalMessageInfo

func (m *AddBeerRequest) GetBeer() *Beer {
	if m != nil {
		return m.Beer
	}
	return nil
}

func (m *AddBeerRequest) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type AddBeerResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddBeerResponse) Reset()         { *m = AddBeerResponse{} }
func (m *AddBeerResponse) String() string { return proto.CompactTextString(m) }
func (*AddBeerResponse) ProtoMessage()    {}
func (*AddBeerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{7}
}

func (m *AddBeerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddBeerResponse.Unmarshal(m, b)
}
func (m *AddBeerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddBeerResponse.Marshal(b, m, deterministic)
}
func (m *AddBeerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddBeerResponse.Merge(m, src)
}
func (m *AddBeerResponse) XXX_Size() int {
	return xxx_messageInfo_AddBeerResponse.Size(m)
}
func (m *AddBeerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddBeerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddBeerResponse proto.InternalMessageInfo

type ListBeerRequest struct {
	OnDeck               bool     `protobuf:"varint,1,opt,name=on_deck,json=onDeck,proto3" json:"on_deck,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListBeerRequest) Reset()         { *m = ListBeerRequest{} }
func (m *ListBeerRequest) String() string { return proto.CompactTextString(m) }
func (*ListBeerRequest) ProtoMessage()    {}
func (*ListBeerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{8}
}

func (m *ListBeerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListBeerRequest.Unmarshal(m, b)
}
func (m *ListBeerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListBeerRequest.Marshal(b, m, deterministic)
}
func (m *ListBeerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListBeerRequest.Merge(m, src)
}
func (m *ListBeerRequest) XXX_Size() int {
	return xxx_messageInfo_ListBeerRequest.Size(m)
}
func (m *ListBeerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListBeerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListBeerRequest proto.InternalMessageInfo

func (m *ListBeerRequest) GetOnDeck() bool {
	if m != nil {
		return m.OnDeck
	}
	return false
}

type ListBeerResponse struct {
	Beers                []*Beer  `protobuf:"bytes,1,rep,name=beers,proto3" json:"beers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListBeerResponse) Reset()         { *m = ListBeerResponse{} }
func (m *ListBeerResponse) String() string { return proto.CompactTextString(m) }
func (*ListBeerResponse) ProtoMessage()    {}
func (*ListBeerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{9}
}

func (m *ListBeerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListBeerResponse.Unmarshal(m, b)
}
func (m *ListBeerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListBeerResponse.Marshal(b, m, deterministic)
}
func (m *ListBeerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListBeerResponse.Merge(m, src)
}
func (m *ListBeerResponse) XXX_Size() int {
	return xxx_messageInfo_ListBeerResponse.Size(m)
}
func (m *ListBeerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListBeerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListBeerResponse proto.InternalMessageInfo

func (m *ListBeerResponse) GetBeers() []*Beer {
	if m != nil {
		return m.Beers
	}
	return nil
}

type DeleteBeerRequest struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteBeerRequest) Reset()         { *m = DeleteBeerRequest{} }
func (m *DeleteBeerRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteBeerRequest) ProtoMessage()    {}
func (*DeleteBeerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{10}
}

func (m *DeleteBeerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteBeerRequest.Unmarshal(m, b)
}
func (m *DeleteBeerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteBeerRequest.Marshal(b, m, deterministic)
}
func (m *DeleteBeerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteBeerRequest.Merge(m, src)
}
func (m *DeleteBeerRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteBeerRequest.Size(m)
}
func (m *DeleteBeerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteBeerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteBeerRequest proto.InternalMessageInfo

func (m *DeleteBeerRequest) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

type DeleteBeerResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteBeerResponse) Reset()         { *m = DeleteBeerResponse{} }
func (m *DeleteBeerResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteBeerResponse) ProtoMessage()    {}
func (*DeleteBeerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f9c2065d50c4def, []int{11}
}

func (m *DeleteBeerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteBeerResponse.Unmarshal(m, b)
}
func (m *DeleteBeerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteBeerResponse.Marshal(b, m, deterministic)
}
func (m *DeleteBeerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteBeerResponse.Merge(m, src)
}
func (m *DeleteBeerResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteBeerResponse.Size(m)
}
func (m *DeleteBeerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteBeerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteBeerResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Empty)(nil), "beer.Empty")
	proto.RegisterType((*Token)(nil), "beer.Token")
	proto.RegisterType((*Config)(nil), "beer.Config")
	proto.RegisterType((*Beer)(nil), "beer.Beer")
	proto.RegisterType((*CellarSlot)(nil), "beer.CellarSlot")
	proto.RegisterType((*Cellar)(nil), "beer.Cellar")
	proto.RegisterType((*AddBeerRequest)(nil), "beer.AddBeerRequest")
	proto.RegisterType((*AddBeerResponse)(nil), "beer.AddBeerResponse")
	proto.RegisterType((*ListBeerRequest)(nil), "beer.ListBeerRequest")
	proto.RegisterType((*ListBeerResponse)(nil), "beer.ListBeerResponse")
	proto.RegisterType((*DeleteBeerRequest)(nil), "beer.DeleteBeerRequest")
	proto.RegisterType((*DeleteBeerResponse)(nil), "beer.DeleteBeerResponse")
}

func init() { proto.RegisterFile("beer.proto", fileDescriptor_3f9c2065d50c4def) }

var fileDescriptor_3f9c2065d50c4def = []byte{
	// 577 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0xcb, 0x6e, 0xd3, 0x4c,
	0x14, 0xfe, 0x27, 0x8e, 0x9d, 0xf8, 0xe4, 0x57, 0x2f, 0xa3, 0xd2, 0x8e, 0x82, 0x40, 0xc6, 0x5c,
	0x14, 0xb1, 0xe8, 0x22, 0xb0, 0x60, 0xc1, 0xa6, 0xb4, 0x08, 0x21, 0x75, 0x35, 0x85, 0x75, 0xe4,
	0x7a, 0x0e, 0x30, 0xb2, 0x33, 0x4e, 0x3d, 0xe3, 0x8a, 0xf0, 0x14, 0x3c, 0x19, 0x4f, 0xc2, 0x43,
	0xa0, 0xb9, 0xa4, 0x31, 0x8d, 0xd8, 0x9d, 0xf3, 0x9d, 0xcb, 0x37, 0xe7, 0xfb, 0x2c, 0x03, 0x5c,
	0x23, 0xb6, 0xa7, 0xab, 0xb6, 0x31, 0x0d, 0x1d, 0xda, 0x38, 0x1f, 0x41, 0xfc, 0x7e, 0xb9, 0x32,
	0xeb, 0xfc, 0x03, 0xc4, 0x9f, 0x9a, 0x0a, 0x15, 0xdd, 0x83, 0x81, 0x14, 0x8c, 0x64, 0x64, 0x96,
	0xf2, 0x81, 0x14, 0xf4, 0x18, 0x12, 0x8d, 0x65, 0x8b, 0x86, 0x0d, 0x1c, 0x16, 0x32, 0x8b, 0xb7,
	0xc6, 0x4e, 0xb0, 0xc8, 0xe3, 0x3e, 0xcb, 0x7f, 0x12, 0x48, 0xce, 0x1b, 0xf5, 0x45, 0x7e, 0xa5,
	0x4f, 0x20, 0xf6, 0x1d, 0x76, 0xdb, 0x64, 0x3e, 0x39, 0x75, 0xf4, 0x8e, 0x86, 0xfb, 0x0a, 0xcd,
	0x20, 0x16, 0x6d, 0xa7, 0x2a, 0x36, 0xc8, 0xa2, 0xd9, 0x64, 0x0e, 0xbe, 0xe5, 0x1d, 0x62, 0xcb,
	0x7d, 0x81, 0x3e, 0x83, 0xa4, 0xc4, 0xba, 0x2e, 0x5a, 0xc7, 0x33, 0x99, 0xff, 0xef, 0x5b, 0xce,
	0x1d, 0xc6, 0x43, 0x8d, 0x3e, 0x84, 0xb4, 0x2e, 0xb4, 0x59, 0xe8, 0xb5, 0x2a, 0xd9, 0x30, 0x23,
	0xb3, 0x88, 0x8f, 0x2d, 0x70, 0xb5, 0x56, 0x65, 0xfe, 0x9b, 0xc0, 0xd0, 0xae, 0xec, 0xdd, 0x16,
	0xb9, 0xdb, 0x1e, 0x01, 0x88, 0x56, 0xaa, 0x6a, 0x21, 0x0a, 0x83, 0xee, 0xbe, 0x88, 0xa7, 0x0e,
	0xb9, 0x28, 0x0c, 0x52, 0x0a, 0x43, 0x2d, 0x7f, 0x60, 0x38, 0xd0, 0xc5, 0x16, 0x53, 0xc5, 0x12,
	0x1d, 0x47, 0xca, 0x5d, 0x4c, 0x4f, 0x60, 0xd4, 0xa8, 0x85, 0xc0, 0xb2, 0x62, 0x71, 0x46, 0x66,
	0x63, 0x9e, 0x34, 0xea, 0x02, 0xcb, 0x8a, 0x1e, 0x40, 0x54, 0x5c, 0xdf, 0xb2, 0x24, 0x23, 0xb3,
	0x01, 0xb7, 0x21, 0x3d, 0x82, 0x58, 0x2a, 0x81, 0xdf, 0xd9, 0x28, 0x23, 0xb3, 0x98, 0xfb, 0xc4,
	0xbe, 0xa3, 0xfc, 0x86, 0x65, 0x25, 0xd5, 0x42, 0x0a, 0x36, 0x76, 0xa5, 0x34, 0x20, 0x1f, 0x85,
	0x3d, 0x4e, 0xaa, 0x45, 0x50, 0x21, 0x75, 0xd5, 0xb1, 0x54, 0x5e, 0x01, 0xcb, 0xd1, 0x49, 0xc1,
	0xc0, 0x3d, 0xde, 0x86, 0x39, 0x02, 0xf8, 0xda, 0x55, 0xdd, 0x18, 0xca, 0x60, 0x54, 0x94, 0x25,
	0xae, 0x8c, 0x0e, 0xa6, 0x6e, 0x52, 0xbb, 0x56, 0x75, 0xcb, 0x85, 0xae, 0x1b, 0xa3, 0xdd, 0xf1,
	0x31, 0x1f, 0xab, 0x6e, 0x69, 0xa7, 0xb4, 0x35, 0xc6, 0xea, 0xac, 0x59, 0xb4, 0x6b, 0x8c, 0x2b,
	0xe4, 0x9f, 0x21, 0x09, 0x4f, 0x78, 0x01, 0xb1, 0x5f, 0x42, 0x5c, 0xef, 0x41, 0xdf, 0x21, 0xbb,
	0x8d, 0xfb, 0x32, 0x7d, 0xba, 0xd5, 0x69, 0xd7, 0xee, 0xa0, 0x59, 0x7e, 0x09, 0x7b, 0x67, 0x42,
	0x38, 0x08, 0x6f, 0x3a, 0xd4, 0x86, 0x3e, 0x06, 0xf7, 0xad, 0x86, 0xaf, 0xa8, 0x3f, 0xe3, 0x70,
	0x3a, 0x85, 0xf1, 0x4d, 0x57, 0x28, 0x23, 0xcd, 0x7a, 0x73, 0xc6, 0x26, 0xcf, 0x0f, 0x61, 0xff,
	0x6e, 0x9b, 0x5e, 0x35, 0x4a, 0x63, 0xfe, 0x12, 0xf6, 0x2f, 0xa5, 0x36, 0x7d, 0x86, 0x9e, 0x81,
	0xa4, 0x6f, 0x60, 0xfe, 0x1a, 0x0e, 0xb6, 0xbd, 0x7e, 0x7e, 0xab, 0x0c, 0xf9, 0x97, 0x32, 0xcf,
	0xe1, 0xf0, 0x02, 0x6b, 0x34, 0xd8, 0xe7, 0x08, 0x3e, 0x91, 0xad, 0x4f, 0x47, 0x40, 0xfb, 0x6d,
	0x7e, 0xfd, 0xfc, 0x17, 0x81, 0x43, 0x0b, 0x04, 0xf9, 0xb0, 0xbd, 0x95, 0x25, 0xd2, 0x37, 0x30,
	0x0a, 0x77, 0xd0, 0x23, 0x4f, 0xf8, 0xb7, 0x48, 0xd3, 0x07, 0xf7, 0xd0, 0x70, 0xec, 0x7f, 0xf4,
	0x0c, 0x60, 0xcb, 0x42, 0x4f, 0x7c, 0xdb, 0xce, 0xf3, 0xa6, 0x6c, 0xb7, 0x70, 0xb7, 0xe2, 0x2d,
	0xa4, 0x1b, 0x15, 0x34, 0x0d, 0x44, 0xf7, 0x24, 0x9c, 0x1e, 0xdf, 0x87, 0x37, 0xd3, 0xd7, 0x89,
	0xfb, 0xdf, 0xbc, 0xfa, 0x13, 0x00, 0x00, 0xff, 0xff, 0xfe, 0xb4, 0x6f, 0x1d, 0x7d, 0x04, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BeerCellarServiceClient is the client API for BeerCellarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BeerCellarServiceClient interface {
	AddBeer(ctx context.Context, in *AddBeerRequest, opts ...grpc.CallOption) (*AddBeerResponse, error)
	DeleteBeer(ctx context.Context, in *DeleteBeerRequest, opts ...grpc.CallOption) (*DeleteBeerResponse, error)
	ListBeers(ctx context.Context, in *ListBeerRequest, opts ...grpc.CallOption) (*ListBeerResponse, error)
}

type beerCellarServiceClient struct {
	cc *grpc.ClientConn
}

func NewBeerCellarServiceClient(cc *grpc.ClientConn) BeerCellarServiceClient {
	return &beerCellarServiceClient{cc}
}

func (c *beerCellarServiceClient) AddBeer(ctx context.Context, in *AddBeerRequest, opts ...grpc.CallOption) (*AddBeerResponse, error) {
	out := new(AddBeerResponse)
	err := c.cc.Invoke(ctx, "/beer.BeerCellarService/AddBeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beerCellarServiceClient) DeleteBeer(ctx context.Context, in *DeleteBeerRequest, opts ...grpc.CallOption) (*DeleteBeerResponse, error) {
	out := new(DeleteBeerResponse)
	err := c.cc.Invoke(ctx, "/beer.BeerCellarService/DeleteBeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beerCellarServiceClient) ListBeers(ctx context.Context, in *ListBeerRequest, opts ...grpc.CallOption) (*ListBeerResponse, error) {
	out := new(ListBeerResponse)
	err := c.cc.Invoke(ctx, "/beer.BeerCellarService/ListBeers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BeerCellarServiceServer is the server API for BeerCellarService service.
type BeerCellarServiceServer interface {
	AddBeer(context.Context, *AddBeerRequest) (*AddBeerResponse, error)
	DeleteBeer(context.Context, *DeleteBeerRequest) (*DeleteBeerResponse, error)
	ListBeers(context.Context, *ListBeerRequest) (*ListBeerResponse, error)
}

func RegisterBeerCellarServiceServer(s *grpc.Server, srv BeerCellarServiceServer) {
	s.RegisterService(&_BeerCellarService_serviceDesc, srv)
}

func _BeerCellarService_AddBeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBeerRequest)
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
		return srv.(BeerCellarServiceServer).AddBeer(ctx, req.(*AddBeerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeerCellarService_DeleteBeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteBeerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeerCellarServiceServer).DeleteBeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/beer.BeerCellarService/DeleteBeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeerCellarServiceServer).DeleteBeer(ctx, req.(*DeleteBeerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeerCellarService_ListBeers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBeerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeerCellarServiceServer).ListBeers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/beer.BeerCellarService/ListBeers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeerCellarServiceServer).ListBeers(ctx, req.(*ListBeerRequest))
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
			MethodName: "DeleteBeer",
			Handler:    _BeerCellarService_DeleteBeer_Handler,
		},
		{
			MethodName: "ListBeers",
			Handler:    _BeerCellarService_ListBeers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "beer.proto",
}
