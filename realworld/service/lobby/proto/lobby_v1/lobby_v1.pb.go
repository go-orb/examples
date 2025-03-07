// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.2
// source: lobby_v1/lobby_v1.proto

package lobby_v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListGamesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Games         []*Game                `protobuf:"bytes,1,rep,name=games,proto3" json:"games,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListGamesResponse) Reset() {
	*x = ListGamesResponse{}
	mi := &file_lobby_v1_lobby_v1_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListGamesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListGamesResponse) ProtoMessage() {}

func (x *ListGamesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lobby_v1_lobby_v1_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListGamesResponse.ProtoReflect.Descriptor instead.
func (*ListGamesResponse) Descriptor() ([]byte, []int) {
	return file_lobby_v1_lobby_v1_proto_rawDescGZIP(), []int{0}
}

func (x *ListGamesResponse) GetGames() []*Game {
	if x != nil {
		return x.Games
	}
	return nil
}

type Game struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Game) Reset() {
	*x = Game{}
	mi := &file_lobby_v1_lobby_v1_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Game) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Game) ProtoMessage() {}

func (x *Game) ProtoReflect() protoreflect.Message {
	mi := &file_lobby_v1_lobby_v1_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Game.ProtoReflect.Descriptor instead.
func (*Game) Descriptor() ([]byte, []int) {
	return file_lobby_v1_lobby_v1_proto_rawDescGZIP(), []int{1}
}

func (x *Game) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Game) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_lobby_v1_lobby_v1_proto protoreflect.FileDescriptor

var file_lobby_v1_lobby_v1_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x5f, 0x76, 0x31, 0x2f, 0x6c, 0x6f, 0x62, 0x62, 0x79,
	0x5f, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6c, 0x6f, 0x62, 0x62, 0x79,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x39, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x05, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x52, 0x05, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x2a, 0x0a, 0x04, 0x47,
	0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0x50, 0x0a, 0x0c, 0x4c, 0x6f, 0x62, 0x62, 0x79,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x47,
	0x61, 0x6d, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1b, 0x2e, 0x6c,
	0x6f, 0x62, 0x62, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x47, 0x61, 0x6d, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x12, 0x5a, 0x10, 0x2e, 0x2f, 0x6c,
	0x6f, 0x62, 0x62, 0x79, 0x3b, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_lobby_v1_lobby_v1_proto_rawDescOnce sync.Once
	file_lobby_v1_lobby_v1_proto_rawDescData []byte
)

func file_lobby_v1_lobby_v1_proto_rawDescGZIP() []byte {
	file_lobby_v1_lobby_v1_proto_rawDescOnce.Do(func() {
		file_lobby_v1_lobby_v1_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_lobby_v1_lobby_v1_proto_rawDesc), len(file_lobby_v1_lobby_v1_proto_rawDesc)))
	})
	return file_lobby_v1_lobby_v1_proto_rawDescData
}

var file_lobby_v1_lobby_v1_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_lobby_v1_lobby_v1_proto_goTypes = []any{
	(*ListGamesResponse)(nil), // 0: lobby.v1.ListGamesResponse
	(*Game)(nil),              // 1: lobby.v1.Game
	(*emptypb.Empty)(nil),     // 2: google.protobuf.Empty
}
var file_lobby_v1_lobby_v1_proto_depIdxs = []int32{
	1, // 0: lobby.v1.ListGamesResponse.games:type_name -> lobby.v1.Game
	2, // 1: lobby.v1.LobbyService.ListGames:input_type -> google.protobuf.Empty
	0, // 2: lobby.v1.LobbyService.ListGames:output_type -> lobby.v1.ListGamesResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_lobby_v1_lobby_v1_proto_init() }
func file_lobby_v1_lobby_v1_proto_init() {
	if File_lobby_v1_lobby_v1_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_lobby_v1_lobby_v1_proto_rawDesc), len(file_lobby_v1_lobby_v1_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lobby_v1_lobby_v1_proto_goTypes,
		DependencyIndexes: file_lobby_v1_lobby_v1_proto_depIdxs,
		MessageInfos:      file_lobby_v1_lobby_v1_proto_msgTypes,
	}.Build()
	File_lobby_v1_lobby_v1_proto = out.File
	file_lobby_v1_lobby_v1_proto_goTypes = nil
	file_lobby_v1_lobby_v1_proto_depIdxs = nil
}
