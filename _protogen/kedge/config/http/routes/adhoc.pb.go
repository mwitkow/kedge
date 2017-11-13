// Code generated by protoc-gen-go.
// source: kedge/config/http/routes/adhoc.proto
// DO NOT EDIT!

/*
Package kedge_config_http_routes is a generated protocol buffer package.

It is generated from these files:
	kedge/config/http/routes/adhoc.proto
	kedge/config/http/routes/routes.proto

It has these top-level messages:
	Adhoc
	Route
*/
package kedge_config_http_routes

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// / Adhoc describes an adhoc proxying method that is not backed by a backend, but dials a "free form" DNS record.
type Adhoc struct {
	// / dns_name_matcher matches the hostname that will be resolved using A records.
	// / The names are matched with a * prefix. For example:
	// / - *.pod.cluster.local
	// / - *.my_namespace.svc.cluster.local
	// / - *.local
	// / The first rule that matches a DNS name will be used, and its ports will be checked.
	DnsNameMatcher string `protobuf:"bytes,1,opt,name=dns_name_matcher,json=dnsNameMatcher" json:"dns_name_matcher,omitempty"`
	// / Port controls the :port behaviour of the URI requested.
	Port *Adhoc_Port `protobuf:"bytes,2,opt,name=port" json:"port,omitempty"`
}

func (m *Adhoc) Reset()                    { *m = Adhoc{} }
func (m *Adhoc) String() string            { return proto.CompactTextString(m) }
func (*Adhoc) ProtoMessage()               {}
func (*Adhoc) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Adhoc) GetDnsNameMatcher() string {
	if m != nil {
		return m.DnsNameMatcher
	}
	return ""
}

func (m *Adhoc) GetPort() *Adhoc_Port {
	if m != nil {
		return m.Port
	}
	return nil
}

// / Port controls how the :port part of the URI is processed.
type Adhoc_Port struct {
	// / default is the default port used if no entry is present.
	// / This defaults to 80.
	Default uint32 `protobuf:"varint,1,opt,name=default" json:"default,omitempty"`
	// / allowed ports is a list of whitelisted ports that this Adhoc rule will allow.
	Allowed []uint32 `protobuf:"varint,3,rep,packed,name=allowed" json:"allowed,omitempty"`
	// / allowed_ranges is a list of whitelisted port ranges that this Adhoc rule will allow.
	AllowedRanges []*Adhoc_Port_Range `protobuf:"bytes,4,rep,name=allowed_ranges,json=allowedRanges" json:"allowed_ranges,omitempty"`
}

func (m *Adhoc_Port) Reset()                    { *m = Adhoc_Port{} }
func (m *Adhoc_Port) String() string            { return proto.CompactTextString(m) }
func (*Adhoc_Port) ProtoMessage()               {}
func (*Adhoc_Port) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Adhoc_Port) GetDefault() uint32 {
	if m != nil {
		return m.Default
	}
	return 0
}

func (m *Adhoc_Port) GetAllowed() []uint32 {
	if m != nil {
		return m.Allowed
	}
	return nil
}

func (m *Adhoc_Port) GetAllowedRanges() []*Adhoc_Port_Range {
	if m != nil {
		return m.AllowedRanges
	}
	return nil
}

type Adhoc_Port_Range struct {
	// / from is an inclusive lower bound for the port range
	From uint32 `protobuf:"varint,1,opt,name=from" json:"from,omitempty"`
	// / to is an inclusive upper bound for the port range
	To uint32 `protobuf:"varint,2,opt,name=to" json:"to,omitempty"`
}

func (m *Adhoc_Port_Range) Reset()                    { *m = Adhoc_Port_Range{} }
func (m *Adhoc_Port_Range) String() string            { return proto.CompactTextString(m) }
func (*Adhoc_Port_Range) ProtoMessage()               {}
func (*Adhoc_Port_Range) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0, 0} }

func (m *Adhoc_Port_Range) GetFrom() uint32 {
	if m != nil {
		return m.From
	}
	return 0
}

func (m *Adhoc_Port_Range) GetTo() uint32 {
	if m != nil {
		return m.To
	}
	return 0
}

func init() {
	proto.RegisterType((*Adhoc)(nil), "kedge.config.http.routes.Adhoc")
	proto.RegisterType((*Adhoc_Port)(nil), "kedge.config.http.routes.Adhoc.Port")
	proto.RegisterType((*Adhoc_Port_Range)(nil), "kedge.config.http.routes.Adhoc.Port.Range")
}

func init() { proto.RegisterFile("kedge/config/http/routes/adhoc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xcd, 0x4a, 0xc4, 0x30,
	0x14, 0x85, 0xe9, 0xcf, 0x28, 0xde, 0xa1, 0x45, 0xb2, 0x0a, 0xb3, 0x2a, 0x32, 0x8b, 0xa2, 0x90,
	0xc2, 0xb8, 0x71, 0xeb, 0x03, 0x28, 0x9a, 0x17, 0x28, 0xb1, 0x49, 0x5b, 0xb1, 0xcd, 0x2d, 0xc9,
	0x1d, 0x7c, 0x25, 0x9f, 0xc1, 0xa7, 0x93, 0xa4, 0x75, 0x29, 0xb8, 0x3b, 0x27, 0xf7, 0xbb, 0x27,
	0x27, 0x81, 0xe3, 0x87, 0xd1, 0x83, 0x69, 0x3a, 0xb4, 0xfd, 0xfb, 0xd0, 0x8c, 0x44, 0x4b, 0xe3,
	0xf0, 0x4c, 0xc6, 0x37, 0x4a, 0x8f, 0xd8, 0x89, 0xc5, 0x21, 0x21, 0xe3, 0x91, 0x12, 0x2b, 0x25,
	0x02, 0x25, 0x56, 0xea, 0xe6, 0x2b, 0x85, 0xdd, 0x63, 0x20, 0x59, 0x0d, 0xd7, 0xda, 0xfa, 0xd6,
	0xaa, 0xd9, 0xb4, 0xb3, 0xa2, 0x6e, 0x34, 0x8e, 0x27, 0x55, 0x52, 0x5f, 0xc9, 0x52, 0x5b, 0xff,
	0xac, 0x66, 0xf3, 0xb4, 0x9e, 0xb2, 0x07, 0xc8, 0x17, 0x74, 0xc4, 0xd3, 0x2a, 0xa9, 0xf7, 0xa7,
	0xa3, 0xf8, 0x2b, 0x5c, 0xc4, 0x60, 0xf1, 0x82, 0x8e, 0x64, 0xdc, 0x38, 0x7c, 0x27, 0x90, 0x07,
	0xcb, 0x38, 0x5c, 0x6a, 0xd3, 0xab, 0xf3, 0x44, 0xf1, 0x8e, 0x42, 0xfe, 0xda, 0x30, 0x51, 0xd3,
	0x84, 0x9f, 0x46, 0xf3, 0xac, 0xca, 0xc2, 0x64, 0xb3, 0xec, 0x15, 0xca, 0x4d, 0xb6, 0x4e, 0xd9,
	0xc1, 0x78, 0x9e, 0x57, 0x59, 0xbd, 0x3f, 0xdd, 0xfe, 0xa7, 0x80, 0x90, 0x61, 0x45, 0x16, 0x5b,
	0x42, 0x74, 0xfe, 0x70, 0x07, 0xbb, 0xa8, 0x18, 0x83, 0xbc, 0x77, 0x38, 0x6f, 0x65, 0xa2, 0x66,
	0x25, 0xa4, 0x84, 0xf1, 0x91, 0x85, 0x4c, 0x09, 0xdf, 0x2e, 0xe2, 0x5f, 0xde, 0xff, 0x04, 0x00,
	0x00, 0xff, 0xff, 0xc7, 0x58, 0xa5, 0xc7, 0x73, 0x01, 0x00, 0x00,
}
