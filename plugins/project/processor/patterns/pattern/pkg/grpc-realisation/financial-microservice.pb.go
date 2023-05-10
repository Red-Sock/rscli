// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: financial-microservice/financial-microservice.proto

package grpc_realisation

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientTimestamp *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=client_timestamp,json=clientTimestamp,proto3" json:"client_timestamp,omitempty"`
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_financial_microservice_financial_microservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_financial_microservice_financial_microservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_financial_microservice_financial_microservice_proto_rawDescGZIP(), []int{0}
}

func (x *PingRequest) GetClientTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.ClientTimestamp
	}
	return nil
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Took uint32 `protobuf:"varint,1,opt,name=took,proto3" json:"took,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_financial_microservice_financial_microservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_financial_microservice_financial_microservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_financial_microservice_financial_microservice_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetTook() uint32 {
	if x != nil {
		return x.Took
	}
	return 0
}

var File_financial_microservice_financial_microservice_proto protoreflect.FileDescriptor

var file_financial_microservice_financial_microservice_proto_rawDesc = []byte{
	0x0a, 0x33, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c, 0x2d, 0x6d, 0x69, 0x63, 0x72,
	0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69,
	0x61, 0x6c, 0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c,
	0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x0b,
	0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x45, 0x0a, 0x10, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x22, 0x22, 0x0a, 0x0c, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x6f, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x74, 0x6f, 0x6f, 0x6b, 0x32, 0x7c, 0x0a, 0x0c, 0x46, 0x69, 0x6e, 0x61, 0x6e, 0x63,
	0x69, 0x61, 0x6c, 0x41, 0x50, 0x49, 0x12, 0x6c, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x23, 0x2e, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69, 0x61, 0x6c, 0x5f, 0x6d, 0x69,
	0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69,
	0x61, 0x6c, 0x5f, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x10, 0x3a, 0x01, 0x2a, 0x22, 0x0b, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x42, 0x2a, 0x5a, 0x28, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69,
	0x61, 0x6c, 0x2d, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_financial_microservice_financial_microservice_proto_rawDescOnce sync.Once
	file_financial_microservice_financial_microservice_proto_rawDescData = file_financial_microservice_financial_microservice_proto_rawDesc
)

func file_financial_microservice_financial_microservice_proto_rawDescGZIP() []byte {
	file_financial_microservice_financial_microservice_proto_rawDescOnce.Do(func() {
		file_financial_microservice_financial_microservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_financial_microservice_financial_microservice_proto_rawDescData)
	})
	return file_financial_microservice_financial_microservice_proto_rawDescData
}

var file_financial_microservice_financial_microservice_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_financial_microservice_financial_microservice_proto_goTypes = []interface{}{
	(*PingRequest)(nil),           // 0: financial_microservice.PingRequest
	(*PingResponse)(nil),          // 1: financial_microservice.PingResponse
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_financial_microservice_financial_microservice_proto_depIdxs = []int32{
	2, // 0: financial_microservice.PingRequest.client_timestamp:type_name -> google.protobuf.Timestamp
	0, // 1: financial_microservice.FinancialAPI.Version:input_type -> financial_microservice.PingRequest
	1, // 2: financial_microservice.FinancialAPI.Version:output_type -> financial_microservice.PingResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_financial_microservice_financial_microservice_proto_init() }
func file_financial_microservice_financial_microservice_proto_init() {
	if File_financial_microservice_financial_microservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_financial_microservice_financial_microservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
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
		file_financial_microservice_financial_microservice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
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
			RawDescriptor: file_financial_microservice_financial_microservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_financial_microservice_financial_microservice_proto_goTypes,
		DependencyIndexes: file_financial_microservice_financial_microservice_proto_depIdxs,
		MessageInfos:      file_financial_microservice_financial_microservice_proto_msgTypes,
	}.Build()
	File_financial_microservice_financial_microservice_proto = out.File
	file_financial_microservice_financial_microservice_proto_rawDesc = nil
	file_financial_microservice_financial_microservice_proto_goTypes = nil
	file_financial_microservice_financial_microservice_proto_depIdxs = nil
}
