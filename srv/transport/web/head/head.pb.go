// Copyright (c) 2024 The horm-database Authors (such as CaoHao <18500482693@163.com>). All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.17.3
// source: head.proto

package head

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// RequestHeader 请求头
type WebReqHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version     string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`                             // 客户端版本
	RequestType int32  `protobuf:"varint,2,opt,name=request_type,json=requestType,proto3" json:"request_type,omitempty"` // 请求类型 0-rpc 请求 1-http 请求 2-web 请求
	RequestId   uint64 `protobuf:"varint,3,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`       // 请求唯一id
	Timestamp   uint64 `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                        // 请求时间戳（精确到毫秒）
	Timeout     uint32 `protobuf:"varint,5,opt,name=timeout,proto3" json:"timeout,omitempty"`                            // 请求超时时间，单位ms
	Userid      uint64 `protobuf:"varint,6,opt,name=userid,proto3" json:"userid,omitempty"`                              // 用户id
	WorkspaceId uint32 `protobuf:"varint,7,opt,name=workspace_id,json=workspaceId,proto3" json:"workspace_id,omitempty"` // workspace id
	Ip          string `protobuf:"bytes,8,opt,name=ip,proto3" json:"ip,omitempty"`                                       // ip地址
	Caller      string `protobuf:"bytes,9,opt,name=caller,proto3" json:"caller,omitempty"`                               // 调用方
	Callee      string `protobuf:"bytes,10,opt,name=callee,proto3" json:"callee,omitempty"`                              // 被调方
	AuthRand    uint32 `protobuf:"varint,11,opt,name=auth_rand,json=authRand,proto3" json:"auth_rand,omitempty"`         // 随机生成 0-9999999 的数字，相同 timestamp 不允许出现同样的 ip、auth_rand。为了避免碰撞，0-9999999，单机理论最大支持 100 亿/秒的并发。
	Sign        string `protobuf:"bytes,12,opt,name=sign,proto3" json:"sign,omitempty"`                                  // 签名，为 md5(workspace_id+userid+token+version+request_id+timestamp+timeout+caller+auth_rand)
}

func (x *WebReqHeader) Reset() {
	*x = WebReqHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_head_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebReqHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebReqHeader) ProtoMessage() {}

func (x *WebReqHeader) ProtoReflect() protoreflect.Message {
	mi := &file_head_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebReqHeader.ProtoReflect.Descriptor instead.
func (*WebReqHeader) Descriptor() ([]byte, []int) {
	return file_head_proto_rawDescGZIP(), []int{0}
}

func (x *WebReqHeader) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *WebReqHeader) GetRequestType() int32 {
	if x != nil {
		return x.RequestType
	}
	return 0
}

func (x *WebReqHeader) GetRequestId() uint64 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *WebReqHeader) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *WebReqHeader) GetTimeout() uint32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

func (x *WebReqHeader) GetUserid() uint64 {
	if x != nil {
		return x.Userid
	}
	return 0
}

func (x *WebReqHeader) GetWorkspaceId() uint32 {
	if x != nil {
		return x.WorkspaceId
	}
	return 0
}

func (x *WebReqHeader) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *WebReqHeader) GetCaller() string {
	if x != nil {
		return x.Caller
	}
	return ""
}

func (x *WebReqHeader) GetCallee() string {
	if x != nil {
		return x.Callee
	}
	return ""
}

func (x *WebReqHeader) GetAuthRand() uint32 {
	if x != nil {
		return x.AuthRand
	}
	return 0
}

func (x *WebReqHeader) GetSign() string {
	if x != nil {
		return x.Sign
	}
	return ""
}

// ResponseHeader 响应头
type WebRespHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version   string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`                       // 客户端版本
	RequestId uint64 `protobuf:"varint,2,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"` // 请求唯一id
}

func (x *WebRespHeader) Reset() {
	*x = WebRespHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_head_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebRespHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebRespHeader) ProtoMessage() {}

func (x *WebRespHeader) ProtoReflect() protoreflect.Message {
	mi := &file_head_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebRespHeader.ProtoReflect.Descriptor instead.
func (*WebRespHeader) Descriptor() ([]byte, []int) {
	return file_head_proto_rawDescGZIP(), []int{1}
}

func (x *WebRespHeader) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *WebRespHeader) GetRequestId() uint64 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

var File_head_proto protoreflect.FileDescriptor

var file_head_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x68, 0x65, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xce, 0x02, 0x0a,
	0x0c, 0x57, 0x65, 0x62, 0x52, 0x65, 0x71, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x77, 0x6f, 0x72,
	0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x16, 0x0a, 0x06,
	0x63, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61,
	0x6c, 0x6c, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x6c, 0x6c, 0x65, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x6c, 0x6c, 0x65, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x75, 0x74, 0x68, 0x5f, 0x72, 0x61, 0x6e, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x61, 0x75, 0x74, 0x68, 0x52, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x67,
	0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x69, 0x67, 0x6e, 0x22, 0x48, 0x0a,
	0x0d, 0x57, 0x65, 0x62, 0x52, 0x65, 0x73, 0x70, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_head_proto_rawDescOnce sync.Once
	file_head_proto_rawDescData = file_head_proto_rawDesc
)

func file_head_proto_rawDescGZIP() []byte {
	file_head_proto_rawDescOnce.Do(func() {
		file_head_proto_rawDescData = protoimpl.X.CompressGZIP(file_head_proto_rawDescData)
	})
	return file_head_proto_rawDescData
}

var file_head_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_head_proto_goTypes = []interface{}{
	(*WebReqHeader)(nil),  // 0: WebReqHeader
	(*WebRespHeader)(nil), // 1: WebRespHeader
}
var file_head_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_head_proto_init() }
func file_head_proto_init() {
	if File_head_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_head_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebReqHeader); i {
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
		file_head_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebRespHeader); i {
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
			RawDescriptor: file_head_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_head_proto_goTypes,
		DependencyIndexes: file_head_proto_depIdxs,
		MessageInfos:      file_head_proto_msgTypes,
	}.Build()
	File_head_proto = out.File
	file_head_proto_rawDesc = nil
	file_head_proto_goTypes = nil
	file_head_proto_depIdxs = nil
}
