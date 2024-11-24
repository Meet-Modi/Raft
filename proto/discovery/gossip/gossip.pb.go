// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.3
// source: proto/discovery/gossip/gossip.proto

package gossip

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

type PeerData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uri string `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
}

func (x *PeerData) Reset() {
	*x = PeerData{}
	mi := &file_proto_discovery_gossip_gossip_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PeerData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerData) ProtoMessage() {}

func (x *PeerData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_discovery_gossip_gossip_proto_msgTypes[0]
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
	return file_proto_discovery_gossip_gossip_proto_rawDescGZIP(), []int{0}
}

func (x *PeerData) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

type DiscoveryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId  string               `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	MyPeers map[string]*PeerData `protobuf:"bytes,2,rep,name=my_peers,json=myPeers,proto3" json:"my_peers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *DiscoveryRequest) Reset() {
	*x = DiscoveryRequest{}
	mi := &file_proto_discovery_gossip_gossip_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DiscoveryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryRequest) ProtoMessage() {}

func (x *DiscoveryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_discovery_gossip_gossip_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryRequest.ProtoReflect.Descriptor instead.
func (*DiscoveryRequest) Descriptor() ([]byte, []int) {
	return file_proto_discovery_gossip_gossip_proto_rawDescGZIP(), []int{1}
}

func (x *DiscoveryRequest) GetPeerId() string {
	if x != nil {
		return x.PeerId
	}
	return ""
}

func (x *DiscoveryRequest) GetMyPeers() map[string]*PeerData {
	if x != nil {
		return x.MyPeers
	}
	return nil
}

type DiscoveryDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId          string               `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	Peers           map[string]*PeerData `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CurrentLeaderId string               `protobuf:"bytes,3,opt,name=current_leader_id,json=currentLeaderId,proto3" json:"current_leader_id,omitempty"`
	BootNodeId      string               `protobuf:"bytes,4,opt,name=boot_node_id,json=bootNodeId,proto3" json:"boot_node_id,omitempty"`
}

func (x *DiscoveryDataResponse) Reset() {
	*x = DiscoveryDataResponse{}
	mi := &file_proto_discovery_gossip_gossip_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DiscoveryDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryDataResponse) ProtoMessage() {}

func (x *DiscoveryDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_discovery_gossip_gossip_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryDataResponse.ProtoReflect.Descriptor instead.
func (*DiscoveryDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_discovery_gossip_gossip_proto_rawDescGZIP(), []int{2}
}

func (x *DiscoveryDataResponse) GetPeerId() string {
	if x != nil {
		return x.PeerId
	}
	return ""
}

func (x *DiscoveryDataResponse) GetPeers() map[string]*PeerData {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *DiscoveryDataResponse) GetCurrentLeaderId() string {
	if x != nil {
		return x.CurrentLeaderId
	}
	return ""
}

func (x *DiscoveryDataResponse) GetBootNodeId() string {
	if x != nil {
		return x.BootNodeId
	}
	return ""
}

var File_proto_discovery_gossip_gossip_proto protoreflect.FileDescriptor

var file_proto_discovery_gossip_gossip_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x2f, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2f, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x22, 0x1c, 0x0a,
	0x08, 0x50, 0x65, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x69,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x22, 0xbb, 0x01, 0x0a, 0x10,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x40, 0x0a, 0x08, 0x6d, 0x79, 0x5f,
	0x70, 0x65, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x67, 0x6f,
	0x73, 0x73, 0x69, 0x70, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x79, 0x50, 0x65, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x6d, 0x79, 0x50, 0x65, 0x65, 0x72, 0x73, 0x1a, 0x4c, 0x0a, 0x0c, 0x4d,
	0x79, 0x50, 0x65, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67,
	0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8a, 0x02, 0x0a, 0x15, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x3e, 0x0a, 0x05,
	0x70, 0x65, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x67, 0x6f,
	0x73, 0x73, 0x69, 0x70, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x12, 0x2a, 0x0a, 0x11,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0c, 0x62, 0x6f, 0x6f, 0x74,
	0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x62, 0x6f, 0x6f, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x1a, 0x4a, 0x0a, 0x0a, 0x50, 0x65,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x6f, 0x73, 0x73,
	0x69, 0x70, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x63, 0x0a, 0x10, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x12, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x50, 0x65, 0x65, 0x72, 0x73,
	0x12, 0x18, 0x2e, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x6f, 0x73,
	0x73, 0x69, 0x70, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x32, 0x5a, 0x30, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x65, 0x65, 0x74, 0x2d, 0x6d,
	0x6f, 0x64, 0x69, 0x2f, 0x52, 0x61, 0x66, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_discovery_gossip_gossip_proto_rawDescOnce sync.Once
	file_proto_discovery_gossip_gossip_proto_rawDescData = file_proto_discovery_gossip_gossip_proto_rawDesc
)

func file_proto_discovery_gossip_gossip_proto_rawDescGZIP() []byte {
	file_proto_discovery_gossip_gossip_proto_rawDescOnce.Do(func() {
		file_proto_discovery_gossip_gossip_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_discovery_gossip_gossip_proto_rawDescData)
	})
	return file_proto_discovery_gossip_gossip_proto_rawDescData
}

var file_proto_discovery_gossip_gossip_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_discovery_gossip_gossip_proto_goTypes = []any{
	(*PeerData)(nil),              // 0: gossip.PeerData
	(*DiscoveryRequest)(nil),      // 1: gossip.DiscoveryRequest
	(*DiscoveryDataResponse)(nil), // 2: gossip.DiscoveryDataResponse
	nil,                           // 3: gossip.DiscoveryRequest.MyPeersEntry
	nil,                           // 4: gossip.DiscoveryDataResponse.PeersEntry
}
var file_proto_discovery_gossip_gossip_proto_depIdxs = []int32{
	3, // 0: gossip.DiscoveryRequest.my_peers:type_name -> gossip.DiscoveryRequest.MyPeersEntry
	4, // 1: gossip.DiscoveryDataResponse.peers:type_name -> gossip.DiscoveryDataResponse.PeersEntry
	0, // 2: gossip.DiscoveryRequest.MyPeersEntry.value:type_name -> gossip.PeerData
	0, // 3: gossip.DiscoveryDataResponse.PeersEntry.value:type_name -> gossip.PeerData
	1, // 4: gossip.DiscoveryService.ServeDiscoverPeers:input_type -> gossip.DiscoveryRequest
	2, // 5: gossip.DiscoveryService.ServeDiscoverPeers:output_type -> gossip.DiscoveryDataResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_discovery_gossip_gossip_proto_init() }
func file_proto_discovery_gossip_gossip_proto_init() {
	if File_proto_discovery_gossip_gossip_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_discovery_gossip_gossip_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_discovery_gossip_gossip_proto_goTypes,
		DependencyIndexes: file_proto_discovery_gossip_gossip_proto_depIdxs,
		MessageInfos:      file_proto_discovery_gossip_gossip_proto_msgTypes,
	}.Build()
	File_proto_discovery_gossip_gossip_proto = out.File
	file_proto_discovery_gossip_gossip_proto_rawDesc = nil
	file_proto_discovery_gossip_gossip_proto_goTypes = nil
	file_proto_discovery_gossip_gossip_proto_depIdxs = nil
}
