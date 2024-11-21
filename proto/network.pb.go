// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: proto/network.proto

package network

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RaftNetworkDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RaftNetworkDataRequest) Reset() {
	*x = RaftNetworkDataRequest{}
	mi := &file_proto_network_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RaftNetworkDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftNetworkDataRequest) ProtoMessage() {}

func (x *RaftNetworkDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_network_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftNetworkDataRequest.ProtoReflect.Descriptor instead.
func (*RaftNetworkDataRequest) Descriptor() ([]byte, []int) {
	return file_proto_network_proto_rawDescGZIP(), []int{0}
}

type PeerData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Port     int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *PeerData) Reset() {
	*x = PeerData{}
	mi := &file_proto_network_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PeerData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerData) ProtoMessage() {}

func (x *PeerData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_network_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerData.ProtoReflect.Descriptor instead.
func (*PeerData) Descriptor() ([]byte, []int) {
	return file_proto_network_proto_rawDescGZIP(), []int{1}
}

func (x *PeerData) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *PeerData) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type RaftNetworkDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId            string               `protobuf:"bytes,1,opt,name=peerId,proto3" json:"peerId,omitempty"`
	Peers             map[string]*PeerData `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CurrentLeaderId   string               `protobuf:"bytes,3,opt,name=currentLeaderId,proto3" json:"currentLeaderId,omitempty"`
	CurrentLeaderPort int32                `protobuf:"varint,4,opt,name=currentLeaderPort,proto3" json:"currentLeaderPort,omitempty"`
	IsBootNode        bool                 `protobuf:"varint,5,opt,name=isBootNode,proto3" json:"isBootNode,omitempty"`
	BootNodeEndpoint  string               `protobuf:"bytes,6,opt,name=bootNodeEndpoint,proto3" json:"bootNodeEndpoint,omitempty"`
	BootNodePort      int32                `protobuf:"varint,7,opt,name=bootNodePort,proto3" json:"bootNodePort,omitempty"`
}

func (x *RaftNetworkDataResponse) Reset() {
	*x = RaftNetworkDataResponse{}
	mi := &file_proto_network_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RaftNetworkDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftNetworkDataResponse) ProtoMessage() {}

func (x *RaftNetworkDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_network_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftNetworkDataResponse.ProtoReflect.Descriptor instead.
func (*RaftNetworkDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_network_proto_rawDescGZIP(), []int{2}
}

func (x *RaftNetworkDataResponse) GetPeerId() string {
	if x != nil {
		return x.PeerId
	}
	return ""
}

func (x *RaftNetworkDataResponse) GetPeers() map[string]*PeerData {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *RaftNetworkDataResponse) GetCurrentLeaderId() string {
	if x != nil {
		return x.CurrentLeaderId
	}
	return ""
}

func (x *RaftNetworkDataResponse) GetCurrentLeaderPort() int32 {
	if x != nil {
		return x.CurrentLeaderPort
	}
	return 0
}

func (x *RaftNetworkDataResponse) GetIsBootNode() bool {
	if x != nil {
		return x.IsBootNode
	}
	return false
}

func (x *RaftNetworkDataResponse) GetBootNodeEndpoint() string {
	if x != nil {
		return x.BootNodeEndpoint
	}
	return ""
}

func (x *RaftNetworkDataResponse) GetBootNodePort() int32 {
	if x != nil {
		return x.BootNodePort
	}
	return 0
}

var File_proto_network_proto protoreflect.FileDescriptor

var file_proto_network_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x22, 0x18,
	0x0a, 0x16, 0x52, 0x61, 0x66, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x0a, 0x08, 0x50, 0x65, 0x65, 0x72,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x22, 0x89, 0x03, 0x0a, 0x17, 0x52, 0x61, 0x66, 0x74, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x2e, 0x52, 0x61, 0x66, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x12, 0x28, 0x0a, 0x0f, 0x63,
	0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x50,
	0x6f, 0x72, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x73, 0x42, 0x6f, 0x6f, 0x74, 0x4e, 0x6f, 0x64,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x42, 0x6f, 0x6f, 0x74, 0x4e,
	0x6f, 0x64, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x62, 0x6f, 0x6f, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x45,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x62,
	0x6f, 0x6f, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12,
	0x22, 0x0a, 0x0c, 0x62, 0x6f, 0x6f, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x62, 0x6f, 0x6f, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x50,
	0x6f, 0x72, 0x74, 0x1a, 0x4b, 0x0a, 0x0a, 0x50, 0x65, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x50, 0x65, 0x65,
	0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x32, 0x67, 0x0a, 0x0a, 0x52, 0x61, 0x66, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x59,
	0x0a, 0x12, 0x47, 0x65, 0x74, 0x52, 0x61, 0x66, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x1f, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x52,
	0x61, 0x66, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e,
	0x52, 0x61, 0x66, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x65, 0x74, 0x2d, 0x6d, 0x6f, 0x64,
	0x69, 0x2f, 0x52, 0x61, 0x66, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_network_proto_rawDescOnce sync.Once
	file_proto_network_proto_rawDescData = file_proto_network_proto_rawDesc
)

func file_proto_network_proto_rawDescGZIP() []byte {
	file_proto_network_proto_rawDescOnce.Do(func() {
		file_proto_network_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_network_proto_rawDescData)
	})
	return file_proto_network_proto_rawDescData
}

var file_proto_network_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_network_proto_goTypes = []any{
	(*RaftNetworkDataRequest)(nil),  // 0: network.RaftNetworkDataRequest
	(*PeerData)(nil),                // 1: network.PeerData
	(*RaftNetworkDataResponse)(nil), // 2: network.RaftNetworkDataResponse
	nil,                             // 3: network.RaftNetworkDataResponse.PeersEntry
}
var file_proto_network_proto_depIdxs = []int32{
	3, // 0: network.RaftNetworkDataResponse.peers:type_name -> network.RaftNetworkDataResponse.PeersEntry
	1, // 1: network.RaftNetworkDataResponse.PeersEntry.value:type_name -> network.PeerData
	0, // 2: network.RaftServer.GetRaftNetworkData:input_type -> network.RaftNetworkDataRequest
	2, // 3: network.RaftServer.GetRaftNetworkData:output_type -> network.RaftNetworkDataResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_network_proto_init() }
func file_proto_network_proto_init() {
	if File_proto_network_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_network_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_network_proto_goTypes,
		DependencyIndexes: file_proto_network_proto_depIdxs,
		MessageInfos:      file_proto_network_proto_msgTypes,
	}.Build()
	File_proto_network_proto = out.File
	file_proto_network_proto_rawDesc = nil
	file_proto_network_proto_goTypes = nil
	file_proto_network_proto_depIdxs = nil
}