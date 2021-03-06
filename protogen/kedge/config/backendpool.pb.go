// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kedge/config/backendpool.proto

/*
Package kedge_config is a generated protocol buffer package.

It is generated from these files:
	kedge/config/backendpool.proto
	kedge/config/director.proto

It has these top-level messages:
	BackendPoolConfig
	TlsServerConfig
	DirectorConfig
*/
package kedge_config

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/mwitkow/go-proto-validators"
import kedge_config_grpc_backends "github.com/improbable-eng/kedge/protogen/kedge/config/grpc/backends"
import kedge_config_http_backends "github.com/improbable-eng/kedge/protogen/kedge/config/http/backends"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// / Config is the top level configuration message for a backend pool.
type BackendPoolConfig struct {
	TlsServerConfigs []*TlsServerConfig      `protobuf:"bytes,1,rep,name=tls_server_configs,json=tlsServerConfigs" json:"tls_server_configs,omitempty"`
	Grpc             *BackendPoolConfig_Grpc `protobuf:"bytes,2,opt,name=grpc" json:"grpc,omitempty"`
	Http             *BackendPoolConfig_Http `protobuf:"bytes,3,opt,name=http" json:"http,omitempty"`
}

func (m *BackendPoolConfig) Reset()                    { *m = BackendPoolConfig{} }
func (m *BackendPoolConfig) String() string            { return proto.CompactTextString(m) }
func (*BackendPoolConfig) ProtoMessage()               {}
func (*BackendPoolConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *BackendPoolConfig) GetTlsServerConfigs() []*TlsServerConfig {
	if m != nil {
		return m.TlsServerConfigs
	}
	return nil
}

func (m *BackendPoolConfig) GetGrpc() *BackendPoolConfig_Grpc {
	if m != nil {
		return m.Grpc
	}
	return nil
}

func (m *BackendPoolConfig) GetHttp() *BackendPoolConfig_Http {
	if m != nil {
		return m.Http
	}
	return nil
}

type BackendPoolConfig_Grpc struct {
	Backends []*kedge_config_grpc_backends.Backend `protobuf:"bytes,1,rep,name=backends" json:"backends,omitempty"`
}

func (m *BackendPoolConfig_Grpc) Reset()                    { *m = BackendPoolConfig_Grpc{} }
func (m *BackendPoolConfig_Grpc) String() string            { return proto.CompactTextString(m) }
func (*BackendPoolConfig_Grpc) ProtoMessage()               {}
func (*BackendPoolConfig_Grpc) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *BackendPoolConfig_Grpc) GetBackends() []*kedge_config_grpc_backends.Backend {
	if m != nil {
		return m.Backends
	}
	return nil
}

type BackendPoolConfig_Http struct {
	Backends []*kedge_config_http_backends.Backend `protobuf:"bytes,1,rep,name=backends" json:"backends,omitempty"`
}

func (m *BackendPoolConfig_Http) Reset()                    { *m = BackendPoolConfig_Http{} }
func (m *BackendPoolConfig_Http) String() string            { return proto.CompactTextString(m) }
func (*BackendPoolConfig_Http) ProtoMessage()               {}
func (*BackendPoolConfig_Http) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *BackendPoolConfig_Http) GetBackends() []*kedge_config_http_backends.Backend {
	if m != nil {
		return m.Backends
	}
	return nil
}

type TlsServerConfig struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *TlsServerConfig) Reset()                    { *m = TlsServerConfig{} }
func (m *TlsServerConfig) String() string            { return proto.CompactTextString(m) }
func (*TlsServerConfig) ProtoMessage()               {}
func (*TlsServerConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TlsServerConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*BackendPoolConfig)(nil), "kedge.config.BackendPoolConfig")
	proto.RegisterType((*BackendPoolConfig_Grpc)(nil), "kedge.config.BackendPoolConfig.Grpc")
	proto.RegisterType((*BackendPoolConfig_Http)(nil), "kedge.config.BackendPoolConfig.Http")
	proto.RegisterType((*TlsServerConfig)(nil), "kedge.config.TlsServerConfig")
}

func init() { proto.RegisterFile("kedge/config/backendpool.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 318 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x4d, 0x4b, 0xc3, 0x30,
	0x18, 0xc7, 0xe9, 0x56, 0x44, 0x33, 0x61, 0x2e, 0x20, 0x94, 0x82, 0x5a, 0xe6, 0x0e, 0x15, 0x6c,
	0x0a, 0x55, 0x86, 0x07, 0x41, 0x98, 0x87, 0x09, 0x5e, 0xa4, 0x7a, 0x13, 0x2d, 0x7d, 0x89, 0x5d,
	0xe9, 0x4b, 0x42, 0x12, 0x37, 0x50, 0xfc, 0xac, 0x82, 0x07, 0x3f, 0x87, 0x24, 0xb5, 0xc3, 0xfa,
	0x7e, 0x0b, 0xc9, 0xff, 0xf7, 0xfb, 0x3f, 0x4f, 0x0b, 0xb6, 0x73, 0x9c, 0xa4, 0xd8, 0x8d, 0x49,
	0x75, 0x97, 0xa5, 0x6e, 0x14, 0xc6, 0x39, 0xae, 0x12, 0x4a, 0x48, 0x81, 0x28, 0x23, 0x82, 0xc0,
	0x75, 0xf5, 0x8e, 0xea, 0x77, 0x73, 0x9c, 0x66, 0x62, 0x76, 0x1f, 0xa1, 0x98, 0x94, 0x6e, 0xb9,
	0xc8, 0x44, 0x4e, 0x16, 0x6e, 0x4a, 0x1c, 0x15, 0x75, 0xe6, 0x61, 0x91, 0x25, 0xa1, 0x20, 0x8c,
	0xbb, 0xcb, 0x63, 0x6d, 0x31, 0xed, 0x56, 0x4b, 0xca, 0x68, 0xdc, 0x54, 0xf1, 0xe6, 0xf0, 0x6d,
	0x72, 0x26, 0x04, 0xfd, 0x21, 0x39, 0x7c, 0xed, 0x80, 0xc1, 0xa4, 0xbe, 0xb9, 0x20, 0xa4, 0x38,
	0x55, 0x04, 0x3c, 0x07, 0x50, 0x14, 0x3c, 0xe0, 0x98, 0xcd, 0x31, 0x0b, 0x6a, 0x0d, 0x37, 0x34,
	0xab, 0x6b, 0xf7, 0xbc, 0x2d, 0xf4, 0x71, 0x19, 0x74, 0x55, 0xf0, 0x4b, 0x15, 0xab, 0x51, 0x7f,
	0x43, 0xb4, 0x2f, 0x38, 0x3c, 0x02, 0xba, 0x9c, 0xd5, 0xe8, 0x58, 0x9a, 0xdd, 0xf3, 0x46, 0x6d,
	0xfc, 0x4b, 0x37, 0x9a, 0x32, 0x1a, 0xfb, 0x8a, 0x90, 0xa4, 0x9c, 0xdd, 0xe8, 0xfe, 0x8f, 0x3c,
	0x13, 0x82, 0xfa, 0x8a, 0x30, 0xa7, 0x40, 0x97, 0x1e, 0x78, 0x02, 0x56, 0x9b, 0xc5, 0xdf, 0xc7,
	0xdf, 0x6d, 0x5b, 0x64, 0x0f, 0x6a, 0x22, 0x8d, 0xd3, 0x5f, 0x42, 0x52, 0x24, 0xb5, 0x7f, 0x8b,
	0x64, 0xed, 0x2f, 0xa2, 0xe1, 0x31, 0xe8, 0x7f, 0xfa, 0x54, 0x70, 0x0f, 0xe8, 0x55, 0x58, 0x62,
	0x43, 0xb3, 0x34, 0x7b, 0x6d, 0xb2, 0xf9, 0xf2, 0xbc, 0x33, 0x00, 0xfd, 0xdb, 0xeb, 0xd0, 0x79,
	0x08, 0xd0, 0xcd, 0xa3, 0xb7, 0x3f, 0x3e, 0x7c, 0x1a, 0xf9, 0x2a, 0x12, 0xad, 0xa8, 0xbf, 0x75,
	0xf0, 0x16, 0x00, 0x00, 0xff, 0xff, 0x85, 0x9f, 0x6a, 0xa2, 0x69, 0x02, 0x00, 0x00,
}
