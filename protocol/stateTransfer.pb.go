// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: stateTransfer.proto

package protocol;

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 定义协议种类
type CommitType int32

const (
	CommitType_TWO_PHASE_COMMIT   CommitType = 0
	CommitType_THREE_PHASE_COMMIT CommitType = 1
)

// Enum value maps for CommitType.
var (
	CommitType_name = map[int32]string{
		0: "TWO_PHASE_COMMIT",
		1: "THREE_PHASE_COMMIT",
	}
	CommitType_value = map[string]int32{
		"TWO_PHASE_COMMIT":   0,
		"THREE_PHASE_COMMIT": 1,
	}
)

func (x CommitType) Enum() *CommitType {
	p := new(CommitType)
	*p = x
	return p
}

func (x CommitType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommitType) Descriptor() protoreflect.EnumDescriptor {
	return file_stateTransfer_proto_enumTypes[0].Descriptor()
}

func (CommitType) Type() protoreflect.EnumType {
	return &file_stateTransfer_proto_enumTypes[0]
}

func (x CommitType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommitType.Descriptor instead.
func (CommitType) EnumDescriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{0}
}

// 定义Acknowledge character确认字符种类
type AckType int32

const (
	AckType_ACK  AckType = 0
	AckType_NACK AckType = 1
)

// Enum value maps for AckType.
var (
	AckType_name = map[int32]string{
		0: "ACK",
		1: "NACK",
	}
	AckType_value = map[string]int32{
		"ACK":  0,
		"NACK": 1,
	}
)

func (x AckType) Enum() *AckType {
	p := new(AckType)
	*p = x
	return p
}

func (x AckType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AckType) Descriptor() protoreflect.EnumDescriptor {
	return file_stateTransfer_proto_enumTypes[1].Descriptor()
}

func (AckType) Type() protoreflect.EnumType {
	return &file_stateTransfer_proto_enumTypes[1]
}

func (x AckType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AckType.Descriptor instead.
func (AckType) EnumDescriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{1}
}

// 定义follower种类
type FollowerType int32

const (
	FollowerType_SOURCE FollowerType = 0
	FollowerType_TARGET FollowerType = 1
)

// Enum value maps for FollowerType.
var (
	FollowerType_name = map[int32]string{
		0: "SOURCE",
		1: "TARGET",
	}
	FollowerType_value = map[string]int32{
		"SOURCE": 0,
		"TARGET": 1,
	}
)

func (x FollowerType) Enum() *FollowerType {
	p := new(FollowerType)
	*p = x
	return p
}

func (x FollowerType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FollowerType) Descriptor() protoreflect.EnumDescriptor {
	return file_stateTransfer_proto_enumTypes[2].Descriptor()
}

func (FollowerType) Type() protoreflect.EnumType {
	return &file_stateTransfer_proto_enumTypes[2]
}

func (x FollowerType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FollowerType.Descriptor instead.
func (FollowerType) EnumDescriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{2}
}

//*
//keylist_ori = [
//{address: xxxx; blockheight: 105; balance: 93;},
//{address: xxxx; blockheight: 105; balance: 32;},
//...
//]
type TransferRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source  uint32 `protobuf:"varint,1,opt,name=Source,proto3" json:"Source,omitempty"`
	Target  uint32 `protobuf:"varint,2,opt,name=Target,proto3" json:"Target,omitempty"`
	Keylist []byte `protobuf:"bytes,3,opt,name=Keylist,proto3" json:"Keylist,omitempty"`
	Index   uint64 `protobuf:"varint,4,opt,name=Index,proto3" json:"Index,omitempty"`
}

func (x *TransferRequest) Reset() {
	*x = TransferRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransferRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransferRequest) ProtoMessage() {}

func (x *TransferRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransferRequest.ProtoReflect.Descriptor instead.
func (*TransferRequest) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{0}
}

func (x *TransferRequest) GetSource() uint32 {
	if x != nil {
		return x.Source
	}
	return 0
}

func (x *TransferRequest) GetTarget() uint32 {
	if x != nil {
		return x.Target
	}
	return 0
}

func (x *TransferRequest) GetKeylist() []byte {
	if x != nil {
		return x.Keylist
	}
	return nil
}

func (x *TransferRequest) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

// 定义2PC, 3PC的request, response消息
type ProposeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source     uint32     `protobuf:"varint,1,opt,name=Source,proto3" json:"Source,omitempty"`
	Target     uint32     `protobuf:"varint,2,opt,name=Target,proto3" json:"Target,omitempty"`
	Keylist    []byte     `protobuf:"bytes,3,opt,name=Keylist,proto3" json:"Keylist,omitempty"`
	CommitType CommitType `protobuf:"varint,4,opt,name=CommitType,proto3,enum=protocol.CommitType" json:"CommitType,omitempty"`
	Index      uint64     `protobuf:"varint,5,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *ProposeRequest) Reset() {
	*x = ProposeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProposeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposeRequest) ProtoMessage() {}

func (x *ProposeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposeRequest.ProtoReflect.Descriptor instead.
func (*ProposeRequest) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{1}
}

func (x *ProposeRequest) GetSource() uint32 {
	if x != nil {
		return x.Source
	}
	return 0
}

func (x *ProposeRequest) GetTarget() uint32 {
	if x != nil {
		return x.Target
	}
	return 0
}

func (x *ProposeRequest) GetKeylist() []byte {
	if x != nil {
		return x.Keylist
	}
	return nil
}

func (x *ProposeRequest) GetCommitType() CommitType {
	if x != nil {
		return x.CommitType
	}
	return CommitType_TWO_PHASE_COMMIT
}

func (x *ProposeRequest) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

type PrecommitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *PrecommitRequest) Reset() {
	*x = PrecommitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrecommitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrecommitRequest) ProtoMessage() {}

func (x *PrecommitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrecommitRequest.ProtoReflect.Descriptor instead.
func (*PrecommitRequest) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{2}
}

func (x *PrecommitRequest) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

type CommitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index      uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	IsRollback bool   `protobuf:"varint,2,opt,name=isRollback,proto3" json:"isRollback,omitempty"`
}

func (x *CommitRequest) Reset() {
	*x = CommitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitRequest) ProtoMessage() {}

func (x *CommitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitRequest.ProtoReflect.Descriptor instead.
func (*CommitRequest) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{3}
}

func (x *CommitRequest) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *CommitRequest) GetIsRollback() bool {
	if x != nil {
		return x.IsRollback
	}
	return false
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  AckType      `protobuf:"varint,1,opt,name=type,proto3,enum=protocol.AckType" json:"type,omitempty"`
	Index uint64       `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
	Ftype FollowerType `protobuf:"varint,3,opt,name=ftype,proto3,enum=protocol.FollowerType" json:"ftype,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{4}
}

func (x *Response) GetType() AckType {
	if x != nil {
		return x.Type
	}
	return AckType_ACK
}

func (x *Response) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Response) GetFtype() FollowerType {
	if x != nil {
		return x.Ftype
	}
	return FollowerType_SOURCE
}

// 根据address获取当前状态数据(key=address, value=state)
type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{5}
}

func (x *Msg) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{6}
}

func (x *Value) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// 根据当前height生成事务编号index
type Info struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Height uint64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *Info) Reset() {
	*x = Info{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stateTransfer_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Info) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Info) ProtoMessage() {}

func (x *Info) ProtoReflect() protoreflect.Message {
	mi := &file_stateTransfer_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Info.ProtoReflect.Descriptor instead.
func (*Info) Descriptor() ([]byte, []int) {
	return file_stateTransfer_proto_rawDescGZIP(), []int{7}
}

func (x *Info) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

var File_stateTransfer_proto protoreflect.FileDescriptor

var file_stateTransfer_proto_rawDesc = []byte{
	0x0a, 0x13, 0x73, 0x74, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x71, 0x0a, 0x0f,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x4b, 0x65, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x22,
	0xa6, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x54, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x07, 0x4b, 0x65, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0a,
	0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x28, 0x0a, 0x10, 0x50, 0x72, 0x65, 0x63,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x22, 0x45, 0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x73, 0x52,
	0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69,
	0x73, 0x52, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x22, 0x75, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x41,
	0x63, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x2c, 0x0a, 0x05, 0x66, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x05, 0x66, 0x74, 0x79, 0x70, 0x65,
	0x22, 0x17, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x1d, 0x0a, 0x05, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x1e, 0x0a, 0x04, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x2a, 0x3a, 0x0a, 0x0a, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x57, 0x4f, 0x5f, 0x50, 0x48,
	0x41, 0x53, 0x45, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x49, 0x54, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12,
	0x54, 0x48, 0x52, 0x45, 0x45, 0x5f, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x43, 0x4f, 0x4d, 0x4d,
	0x49, 0x54, 0x10, 0x01, 0x2a, 0x1c, 0x0a, 0x07, 0x41, 0x63, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x07, 0x0a, 0x03, 0x41, 0x43, 0x4b, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x41, 0x43, 0x4b,
	0x10, 0x01, 0x2a, 0x26, 0x0a, 0x0c, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x54, 0x41, 0x52, 0x47, 0x45, 0x54, 0x10, 0x01, 0x32, 0xd0, 0x02, 0x0a, 0x06, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x37, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65,
	0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x70,
	0x6f, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b,
	0x0a, 0x09, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x06, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x25, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x4d, 0x73, 0x67, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x4e, 0x6f, 0x64,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x04, 0x5a,
	0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stateTransfer_proto_rawDescOnce sync.Once
	file_stateTransfer_proto_rawDescData = file_stateTransfer_proto_rawDesc
)

func file_stateTransfer_proto_rawDescGZIP() []byte {
	file_stateTransfer_proto_rawDescOnce.Do(func() {
		file_stateTransfer_proto_rawDescData = protoimpl.X.CompressGZIP(file_stateTransfer_proto_rawDescData)
	})
	return file_stateTransfer_proto_rawDescData
}

var file_stateTransfer_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_stateTransfer_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_stateTransfer_proto_goTypes = []interface{}{
	(CommitType)(0),          // 0: protocol.CommitType
	(AckType)(0),             // 1: protocol.AckType
	(FollowerType)(0),        // 2: protocol.FollowerType
	(*TransferRequest)(nil),  // 3: protocol.TransferRequest
	(*ProposeRequest)(nil),   // 4: protocol.ProposeRequest
	(*PrecommitRequest)(nil), // 5: protocol.PrecommitRequest
	(*CommitRequest)(nil),    // 6: protocol.CommitRequest
	(*Response)(nil),         // 7: protocol.Response
	(*Msg)(nil),              // 8: protocol.Msg
	(*Value)(nil),            // 9: protocol.Value
	(*Info)(nil),             // 10: protocol.Info
	(*emptypb.Empty)(nil),    // 11: google.protobuf.Empty
}
var file_stateTransfer_proto_depIdxs = []int32{
	0,  // 0: protocol.ProposeRequest.CommitType:type_name -> protocol.CommitType
	1,  // 1: protocol.Response.type:type_name -> protocol.AckType
	2,  // 2: protocol.Response.ftype:type_name -> protocol.FollowerType
	4,  // 3: protocol.Commit.Propose:input_type -> protocol.ProposeRequest
	5,  // 4: protocol.Commit.Precommit:input_type -> protocol.PrecommitRequest
	6,  // 5: protocol.Commit.Commit:input_type -> protocol.CommitRequest
	3,  // 6: protocol.Commit.StateTransfer:input_type -> protocol.TransferRequest
	8,  // 7: protocol.Commit.Get:input_type -> protocol.Msg
	11, // 8: protocol.Commit.NodeInfo:input_type -> google.protobuf.Empty
	7,  // 9: protocol.Commit.Propose:output_type -> protocol.Response
	7,  // 10: protocol.Commit.Precommit:output_type -> protocol.Response
	7,  // 11: protocol.Commit.Commit:output_type -> protocol.Response
	7,  // 12: protocol.Commit.StateTransfer:output_type -> protocol.Response
	9,  // 13: protocol.Commit.Get:output_type -> protocol.Value
	10, // 14: protocol.Commit.NodeInfo:output_type -> protocol.Info
	9,  // [9:15] is the sub-list for method output_type
	3,  // [3:9] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_stateTransfer_proto_init() }
func file_stateTransfer_proto_init() {
	if File_stateTransfer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stateTransfer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransferRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stateTransfer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProposeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stateTransfer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrecommitRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stateTransfer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stateTransfer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stateTransfer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stateTransfer_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stateTransfer_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Info); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stateTransfer_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stateTransfer_proto_goTypes,
		DependencyIndexes: file_stateTransfer_proto_depIdxs,
		EnumInfos:         file_stateTransfer_proto_enumTypes,
		MessageInfos:      file_stateTransfer_proto_msgTypes,
	}.Build()
	File_stateTransfer_proto = out.File
	file_stateTransfer_proto_rawDesc = nil
	file_stateTransfer_proto_goTypes = nil
	file_stateTransfer_proto_depIdxs = nil
}