// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: .proto/pm0.proto

package pb

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

type Unit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Pid           *int32 `protobuf:"varint,3,opt,name=pid,proto3,oneof" json:"pid,omitempty"`
	Status        uint32 `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
	RestartsCount uint32 `protobuf:"varint,5,opt,name=restarts_count,json=restartsCount,proto3" json:"restarts_count,omitempty"`
	StartedAt     int64  `protobuf:"varint,6,opt,name=started_at,json=startedAt,proto3" json:"started_at,omitempty"` // TODO: add cpu/memory
}

func (x *Unit) Reset() {
	*x = Unit{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Unit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Unit) ProtoMessage() {}

func (x *Unit) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Unit.ProtoReflect.Descriptor instead.
func (*Unit) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{0}
}

func (x *Unit) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Unit) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Unit) GetPid() int32 {
	if x != nil && x.Pid != nil {
		return *x.Pid
	}
	return 0
}

func (x *Unit) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Unit) GetRestartsCount() uint32 {
	if x != nil {
		return x.RestartsCount
	}
	return 0
}

func (x *Unit) GetStartedAt() int64 {
	if x != nil {
		return x.StartedAt
	}
	return 0
}

type StartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cwd  string   `protobuf:"bytes,1,opt,name=cwd,proto3" json:"cwd,omitempty"`
	Bin  string   `protobuf:"bytes,2,opt,name=bin,proto3" json:"bin,omitempty"`
	Name string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Args []string `protobuf:"bytes,4,rep,name=args,proto3" json:"args,omitempty"`
}

func (x *StartRequest) Reset() {
	*x = StartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartRequest) ProtoMessage() {}

func (x *StartRequest) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartRequest.ProtoReflect.Descriptor instead.
func (*StartRequest) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{1}
}

func (x *StartRequest) GetCwd() string {
	if x != nil {
		return x.Cwd
	}
	return ""
}

func (x *StartRequest) GetBin() string {
	if x != nil {
		return x.Bin
	}
	return ""
}

func (x *StartRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *StartRequest) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

type StartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Unit *Unit `protobuf:"bytes,1,opt,name=unit,proto3" json:"unit,omitempty"`
}

func (x *StartResponse) Reset() {
	*x = StartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartResponse) ProtoMessage() {}

func (x *StartResponse) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartResponse.ProtoReflect.Descriptor instead.
func (*StartResponse) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{2}
}

func (x *StartResponse) GetUnit() *Unit {
	if x != nil {
		return x.Unit
	}
	return nil
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Units []*Unit `protobuf:"bytes,1,rep,name=units,proto3" json:"units,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{3}
}

func (x *ListResponse) GetUnits() []*Unit {
	if x != nil {
		return x.Units
	}
	return nil
}

type StopRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnitIds []uint32 `protobuf:"varint,1,rep,packed,name=unit_ids,json=unitIds,proto3" json:"unit_ids,omitempty"`
}

func (x *StopRequest) Reset() {
	*x = StopRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopRequest) ProtoMessage() {}

func (x *StopRequest) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopRequest.ProtoReflect.Descriptor instead.
func (*StopRequest) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{4}
}

func (x *StopRequest) GetUnitIds() []uint32 {
	if x != nil {
		return x.UnitIds
	}
	return nil
}

type StopResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnitId uint32 `protobuf:"varint,1,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	Unit   *Unit  `protobuf:"bytes,2,opt,name=unit,proto3,oneof" json:"unit,omitempty"`
	Error  string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *StopResponse) Reset() {
	*x = StopResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopResponse) ProtoMessage() {}

func (x *StopResponse) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopResponse.ProtoReflect.Descriptor instead.
func (*StopResponse) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{5}
}

func (x *StopResponse) GetUnitId() uint32 {
	if x != nil {
		return x.UnitId
	}
	return 0
}

func (x *StopResponse) GetUnit() *Unit {
	if x != nil {
		return x.Unit
	}
	return nil
}

func (x *StopResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type LogsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnitId uint32 `protobuf:"varint,1,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
	Follow bool   `protobuf:"varint,2,opt,name=follow,proto3" json:"follow,omitempty"`
	Lines  uint64 `protobuf:"varint,3,opt,name=lines,proto3" json:"lines,omitempty"`
}

func (x *LogsRequest) Reset() {
	*x = LogsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsRequest) ProtoMessage() {}

func (x *LogsRequest) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsRequest.ProtoReflect.Descriptor instead.
func (*LogsRequest) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{6}
}

func (x *LogsRequest) GetUnitId() uint32 {
	if x != nil {
		return x.UnitId
	}
	return 0
}

func (x *LogsRequest) GetFollow() bool {
	if x != nil {
		return x.Follow
	}
	return false
}

func (x *LogsRequest) GetLines() uint64 {
	if x != nil {
		return x.Lines
	}
	return 0
}

type LogsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lines []string `protobuf:"bytes,1,rep,name=lines,proto3" json:"lines,omitempty"`
}

func (x *LogsResponse) Reset() {
	*x = LogsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsResponse) ProtoMessage() {}

func (x *LogsResponse) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsResponse.ProtoReflect.Descriptor instead.
func (*LogsResponse) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{7}
}

func (x *LogsResponse) GetLines() []string {
	if x != nil {
		return x.Lines
	}
	return nil
}

type ShowRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnitId uint32 `protobuf:"varint,1,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
}

func (x *ShowRequest) Reset() {
	*x = ShowRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShowRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShowRequest) ProtoMessage() {}

func (x *ShowRequest) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShowRequest.ProtoReflect.Descriptor instead.
func (*ShowRequest) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{8}
}

func (x *ShowRequest) GetUnitId() uint32 {
	if x != nil {
		return x.UnitId
	}
	return 0
}

type ShowResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Cwd     string `protobuf:"bytes,3,opt,name=cwd,proto3" json:"cwd,omitempty"`
	Command string `protobuf:"bytes,4,opt,name=command,proto3" json:"command,omitempty"`
}

func (x *ShowResponse) Reset() {
	*x = ShowResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShowResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShowResponse) ProtoMessage() {}

func (x *ShowResponse) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShowResponse.ProtoReflect.Descriptor instead.
func (*ShowResponse) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{9}
}

func (x *ShowResponse) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ShowResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ShowResponse) GetCwd() string {
	if x != nil {
		return x.Cwd
	}
	return ""
}

func (x *ShowResponse) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

type LogsClearRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnitIds []uint32 `protobuf:"varint,1,rep,packed,name=unit_ids,json=unitIds,proto3" json:"unit_ids,omitempty"`
}

func (x *LogsClearRequest) Reset() {
	*x = LogsClearRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_pm0_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogsClearRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsClearRequest) ProtoMessage() {}

func (x *LogsClearRequest) ProtoReflect() protoreflect.Message {
	mi := &file___proto_pm0_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsClearRequest.ProtoReflect.Descriptor instead.
func (*LogsClearRequest) Descriptor() ([]byte, []int) {
	return file___proto_pm0_proto_rawDescGZIP(), []int{10}
}

func (x *LogsClearRequest) GetUnitIds() []uint32 {
	if x != nil {
		return x.UnitIds
	}
	return nil
}

var File___proto_pm0_proto protoreflect.FileDescriptor

var file___proto_pm0_proto_rawDesc = []byte{
	0x0a, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6d, 0x30, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x70, 0x6d, 0x30, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7, 0x01, 0x0a, 0x04, 0x55, 0x6e, 0x69, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x15, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00,
	0x52, 0x03, 0x70, 0x69, 0x64, 0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x73, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x70, 0x69, 0x64, 0x22, 0x5a,
	0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x63, 0x77, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x77, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x62, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x62,
	0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x22, 0x2e, 0x0a, 0x0d, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x04, 0x75,
	0x6e, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x6d, 0x30, 0x2e,
	0x55, 0x6e, 0x69, 0x74, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x22, 0x2f, 0x0a, 0x0c, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x05, 0x75, 0x6e,
	0x69, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70, 0x6d, 0x30, 0x2e,
	0x55, 0x6e, 0x69, 0x74, 0x52, 0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x22, 0x28, 0x0a, 0x0b, 0x53,
	0x74, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x6e,
	0x69, 0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x75, 0x6e,
	0x69, 0x74, 0x49, 0x64, 0x73, 0x22, 0x6a, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49, 0x64, 0x12, 0x22,
	0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70,
	0x6d, 0x30, 0x2e, 0x55, 0x6e, 0x69, 0x74, 0x48, 0x00, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x75, 0x6e, 0x69,
	0x74, 0x22, 0x54, 0x0a, 0x0b, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
	0x77, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x05, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x22, 0x24, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6e, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x22, 0x26, 0x0a,
	0x0b, 0x53, 0x68, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75,
	0x6e, 0x69, 0x74, 0x49, 0x64, 0x22, 0x5e, 0x0a, 0x0c, 0x53, 0x68, 0x6f, 0x77, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x77, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x63, 0x77, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x2d, 0x0a, 0x10, 0x4c, 0x6f, 0x67, 0x73, 0x43, 0x6c, 0x65,
	0x61, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x6e, 0x69,
	0x74, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x75, 0x6e, 0x69,
	0x74, 0x49, 0x64, 0x73, 0x32, 0x9d, 0x03, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2e, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x12, 0x11, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x11, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x04, 0x53, 0x74,
	0x6f, 0x70, 0x12, 0x10, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x30, 0x0a, 0x07, 0x52, 0x65, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x6f,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x2d, 0x0a, 0x04, 0x4c,
	0x6f, 0x67, 0x73, 0x12, 0x10, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x4c, 0x6f, 0x67, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x2f, 0x0a, 0x06, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x12, 0x10, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x74, 0x6f,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x2b, 0x0a, 0x04, 0x53,
	0x68, 0x6f, 0x77, 0x12, 0x10, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x68, 0x6f, 0x77, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x53, 0x68, 0x6f, 0x77,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x73,
	0x43, 0x6c, 0x65, 0x61, 0x72, 0x12, 0x15, 0x2e, 0x70, 0x6d, 0x30, 0x2e, 0x4c, 0x6f, 0x67, 0x73,
	0x43, 0x6c, 0x65, 0x61, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file___proto_pm0_proto_rawDescOnce sync.Once
	file___proto_pm0_proto_rawDescData = file___proto_pm0_proto_rawDesc
)

func file___proto_pm0_proto_rawDescGZIP() []byte {
	file___proto_pm0_proto_rawDescOnce.Do(func() {
		file___proto_pm0_proto_rawDescData = protoimpl.X.CompressGZIP(file___proto_pm0_proto_rawDescData)
	})
	return file___proto_pm0_proto_rawDescData
}

var file___proto_pm0_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file___proto_pm0_proto_goTypes = []interface{}{
	(*Unit)(nil),             // 0: pm0.Unit
	(*StartRequest)(nil),     // 1: pm0.StartRequest
	(*StartResponse)(nil),    // 2: pm0.StartResponse
	(*ListResponse)(nil),     // 3: pm0.ListResponse
	(*StopRequest)(nil),      // 4: pm0.StopRequest
	(*StopResponse)(nil),     // 5: pm0.StopResponse
	(*LogsRequest)(nil),      // 6: pm0.LogsRequest
	(*LogsResponse)(nil),     // 7: pm0.LogsResponse
	(*ShowRequest)(nil),      // 8: pm0.ShowRequest
	(*ShowResponse)(nil),     // 9: pm0.ShowResponse
	(*LogsClearRequest)(nil), // 10: pm0.LogsClearRequest
	(*emptypb.Empty)(nil),    // 11: google.protobuf.Empty
}
var file___proto_pm0_proto_depIdxs = []int32{
	0,  // 0: pm0.StartResponse.unit:type_name -> pm0.Unit
	0,  // 1: pm0.ListResponse.units:type_name -> pm0.Unit
	0,  // 2: pm0.StopResponse.unit:type_name -> pm0.Unit
	1,  // 3: pm0.ProcessService.Start:input_type -> pm0.StartRequest
	11, // 4: pm0.ProcessService.List:input_type -> google.protobuf.Empty
	4,  // 5: pm0.ProcessService.Stop:input_type -> pm0.StopRequest
	4,  // 6: pm0.ProcessService.Restart:input_type -> pm0.StopRequest
	6,  // 7: pm0.ProcessService.Logs:input_type -> pm0.LogsRequest
	4,  // 8: pm0.ProcessService.Delete:input_type -> pm0.StopRequest
	8,  // 9: pm0.ProcessService.Show:input_type -> pm0.ShowRequest
	10, // 10: pm0.ProcessService.LogsClear:input_type -> pm0.LogsClearRequest
	2,  // 11: pm0.ProcessService.Start:output_type -> pm0.StartResponse
	3,  // 12: pm0.ProcessService.List:output_type -> pm0.ListResponse
	5,  // 13: pm0.ProcessService.Stop:output_type -> pm0.StopResponse
	5,  // 14: pm0.ProcessService.Restart:output_type -> pm0.StopResponse
	7,  // 15: pm0.ProcessService.Logs:output_type -> pm0.LogsResponse
	5,  // 16: pm0.ProcessService.Delete:output_type -> pm0.StopResponse
	9,  // 17: pm0.ProcessService.Show:output_type -> pm0.ShowResponse
	11, // 18: pm0.ProcessService.LogsClear:output_type -> google.protobuf.Empty
	11, // [11:19] is the sub-list for method output_type
	3,  // [3:11] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file___proto_pm0_proto_init() }
func file___proto_pm0_proto_init() {
	if File___proto_pm0_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file___proto_pm0_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Unit); i {
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
		file___proto_pm0_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartRequest); i {
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
		file___proto_pm0_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartResponse); i {
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
		file___proto_pm0_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResponse); i {
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
		file___proto_pm0_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopRequest); i {
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
		file___proto_pm0_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopResponse); i {
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
		file___proto_pm0_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogsRequest); i {
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
		file___proto_pm0_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogsResponse); i {
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
		file___proto_pm0_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShowRequest); i {
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
		file___proto_pm0_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShowResponse); i {
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
		file___proto_pm0_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogsClearRequest); i {
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
	file___proto_pm0_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file___proto_pm0_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file___proto_pm0_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file___proto_pm0_proto_goTypes,
		DependencyIndexes: file___proto_pm0_proto_depIdxs,
		MessageInfos:      file___proto_pm0_proto_msgTypes,
	}.Build()
	File___proto_pm0_proto = out.File
	file___proto_pm0_proto_rawDesc = nil
	file___proto_pm0_proto_goTypes = nil
	file___proto_pm0_proto_depIdxs = nil
}
