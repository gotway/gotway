// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.7.1
// source: route.proto

package route

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  int32 `protobuf:"varint,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude int32 `protobuf:"varint,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *Point) Reset() {
	*x = Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_route_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Point) ProtoMessage() {}

func (x *Point) ProtoReflect() protoreflect.Message {
	mi := &file_route_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Point.ProtoReflect.Descriptor instead.
func (*Point) Descriptor() ([]byte, []int) {
	return file_route_proto_rawDescGZIP(), []int{0}
}

func (x *Point) GetLatitude() int32 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Point) GetLongitude() int32 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type Feature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Location *Point `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Feature) Reset() {
	*x = Feature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_route_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Feature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Feature) ProtoMessage() {}

func (x *Feature) ProtoReflect() protoreflect.Message {
	mi := &file_route_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Feature.ProtoReflect.Descriptor instead.
func (*Feature) Descriptor() ([]byte, []int) {
	return file_route_proto_rawDescGZIP(), []int{1}
}

func (x *Feature) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Feature) GetLocation() *Point {
	if x != nil {
		return x.Location
	}
	return nil
}

type Rectangle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lo *Point `protobuf:"bytes,1,opt,name=lo,proto3" json:"lo,omitempty"`
	Hi *Point `protobuf:"bytes,2,opt,name=hi,proto3" json:"hi,omitempty"`
}

func (x *Rectangle) Reset() {
	*x = Rectangle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_route_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rectangle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rectangle) ProtoMessage() {}

func (x *Rectangle) ProtoReflect() protoreflect.Message {
	mi := &file_route_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rectangle.ProtoReflect.Descriptor instead.
func (*Rectangle) Descriptor() ([]byte, []int) {
	return file_route_proto_rawDescGZIP(), []int{2}
}

func (x *Rectangle) GetLo() *Point {
	if x != nil {
		return x.Lo
	}
	return nil
}

func (x *Rectangle) GetHi() *Point {
	if x != nil {
		return x.Hi
	}
	return nil
}

type RouteSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PointCount   int32 `protobuf:"varint,1,opt,name=point_count,json=pointCount,proto3" json:"point_count,omitempty"`
	FeatureCount int32 `protobuf:"varint,2,opt,name=feature_count,json=featureCount,proto3" json:"feature_count,omitempty"`
	Distance     int32 `protobuf:"varint,3,opt,name=distance,proto3" json:"distance,omitempty"`
	ElapsedTime  int32 `protobuf:"varint,4,opt,name=elapsed_time,json=elapsedTime,proto3" json:"elapsed_time,omitempty"`
}

func (x *RouteSummary) Reset() {
	*x = RouteSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_route_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteSummary) ProtoMessage() {}

func (x *RouteSummary) ProtoReflect() protoreflect.Message {
	mi := &file_route_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteSummary.ProtoReflect.Descriptor instead.
func (*RouteSummary) Descriptor() ([]byte, []int) {
	return file_route_proto_rawDescGZIP(), []int{3}
}

func (x *RouteSummary) GetPointCount() int32 {
	if x != nil {
		return x.PointCount
	}
	return 0
}

func (x *RouteSummary) GetFeatureCount() int32 {
	if x != nil {
		return x.FeatureCount
	}
	return 0
}

func (x *RouteSummary) GetDistance() int32 {
	if x != nil {
		return x.Distance
	}
	return 0
}

func (x *RouteSummary) GetElapsedTime() int32 {
	if x != nil {
		return x.ElapsedTime
	}
	return 0
}

type RouteNote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Location *Point `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	Message  string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RouteNote) Reset() {
	*x = RouteNote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_route_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouteNote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouteNote) ProtoMessage() {}

func (x *RouteNote) ProtoReflect() protoreflect.Message {
	mi := &file_route_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouteNote.ProtoReflect.Descriptor instead.
func (*RouteNote) Descriptor() ([]byte, []int) {
	return file_route_proto_rawDescGZIP(), []int{4}
}

func (x *RouteNote) GetLocation() *Point {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *RouteNote) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_route_proto protoreflect.FileDescriptor

var file_route_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x72,
	0x6f, 0x75, 0x74, 0x65, 0x22, 0x41, 0x0a, 0x05, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e,
	0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6c, 0x6f,
	0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x47, 0x0a, 0x07, 0x46, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x47, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x12, 0x1c, 0x0a,
	0x02, 0x6c, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x6f, 0x75, 0x74,
	0x65, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x02, 0x6c, 0x6f, 0x12, 0x1c, 0x0a, 0x02, 0x68,
	0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x02, 0x68, 0x69, 0x22, 0x93, 0x01, 0x0a, 0x0c, 0x52, 0x6f,
	0x75, 0x74, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0b, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22,
	0x4f, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x08,
	0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0xd8, 0x01, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x2c, 0x0a, 0x0a, 0x47, 0x65,
	0x74, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x0c, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x1a, 0x0e, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x46,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74,
	0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x10, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65,
	0x2e, 0x52, 0x65, 0x63, 0x74, 0x61, 0x6e, 0x67, 0x6c, 0x65, 0x1a, 0x0e, 0x2e, 0x72, 0x6f, 0x75,
	0x74, 0x65, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x34,
	0x0a, 0x0b, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x0c, 0x2e,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x1a, 0x13, 0x2e, 0x72, 0x6f,
	0x75, 0x74, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x22, 0x00, 0x28, 0x01, 0x12, 0x35, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x43, 0x68, 0x61,
	0x74, 0x12, 0x10, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4e,
	0x6f, 0x74, 0x65, 0x1a, 0x10, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x4e, 0x6f, 0x74, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x2b, 0x5a, 0x29, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6d, 0x6f, 0x6e, 0x74, 0x65,
	0x73, 0x31, 0x31, 0x2f, 0x67, 0x6f, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x72, 0x6f, 0x75, 0x74,
	0x65, 0x73, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_route_proto_rawDescOnce sync.Once
	file_route_proto_rawDescData = file_route_proto_rawDesc
)

func file_route_proto_rawDescGZIP() []byte {
	file_route_proto_rawDescOnce.Do(func() {
		file_route_proto_rawDescData = protoimpl.X.CompressGZIP(file_route_proto_rawDescData)
	})
	return file_route_proto_rawDescData
}

var file_route_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_route_proto_goTypes = []interface{}{
	(*Point)(nil),        // 0: route.Point
	(*Feature)(nil),      // 1: route.Feature
	(*Rectangle)(nil),    // 2: route.Rectangle
	(*RouteSummary)(nil), // 3: route.RouteSummary
	(*RouteNote)(nil),    // 4: route.RouteNote
}
var file_route_proto_depIdxs = []int32{
	0, // 0: route.Feature.location:type_name -> route.Point
	0, // 1: route.Rectangle.lo:type_name -> route.Point
	0, // 2: route.Rectangle.hi:type_name -> route.Point
	0, // 3: route.RouteNote.location:type_name -> route.Point
	0, // 4: route.Route.GetFeature:input_type -> route.Point
	2, // 5: route.Route.ListFeatures:input_type -> route.Rectangle
	0, // 6: route.Route.RecordRoute:input_type -> route.Point
	4, // 7: route.Route.RouteChat:input_type -> route.RouteNote
	1, // 8: route.Route.GetFeature:output_type -> route.Feature
	1, // 9: route.Route.ListFeatures:output_type -> route.Feature
	3, // 10: route.Route.RecordRoute:output_type -> route.RouteSummary
	4, // 11: route.Route.RouteChat:output_type -> route.RouteNote
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_route_proto_init() }
func file_route_proto_init() {
	if File_route_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_route_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Point); i {
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
		file_route_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Feature); i {
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
		file_route_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rectangle); i {
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
		file_route_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteSummary); i {
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
		file_route_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouteNote); i {
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
			RawDescriptor: file_route_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_route_proto_goTypes,
		DependencyIndexes: file_route_proto_depIdxs,
		MessageInfos:      file_route_proto_msgTypes,
	}.Build()
	File_route_proto = out.File
	file_route_proto_rawDesc = nil
	file_route_proto_goTypes = nil
	file_route_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RouteClient is the client API for Route service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RouteClient interface {
	GetFeature(ctx context.Context, in *Point, opts ...grpc.CallOption) (*Feature, error)
	ListFeatures(ctx context.Context, in *Rectangle, opts ...grpc.CallOption) (Route_ListFeaturesClient, error)
	RecordRoute(ctx context.Context, opts ...grpc.CallOption) (Route_RecordRouteClient, error)
	RouteChat(ctx context.Context, opts ...grpc.CallOption) (Route_RouteChatClient, error)
}

type routeClient struct {
	cc grpc.ClientConnInterface
}

func NewRouteClient(cc grpc.ClientConnInterface) RouteClient {
	return &routeClient{cc}
}

func (c *routeClient) GetFeature(ctx context.Context, in *Point, opts ...grpc.CallOption) (*Feature, error) {
	out := new(Feature)
	err := c.cc.Invoke(ctx, "/route.Route/GetFeature", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routeClient) ListFeatures(ctx context.Context, in *Rectangle, opts ...grpc.CallOption) (Route_ListFeaturesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Route_serviceDesc.Streams[0], "/route.Route/ListFeatures", opts...)
	if err != nil {
		return nil, err
	}
	x := &routeListFeaturesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Route_ListFeaturesClient interface {
	Recv() (*Feature, error)
	grpc.ClientStream
}

type routeListFeaturesClient struct {
	grpc.ClientStream
}

func (x *routeListFeaturesClient) Recv() (*Feature, error) {
	m := new(Feature)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *routeClient) RecordRoute(ctx context.Context, opts ...grpc.CallOption) (Route_RecordRouteClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Route_serviceDesc.Streams[1], "/route.Route/RecordRoute", opts...)
	if err != nil {
		return nil, err
	}
	x := &routeRecordRouteClient{stream}
	return x, nil
}

type Route_RecordRouteClient interface {
	Send(*Point) error
	CloseAndRecv() (*RouteSummary, error)
	grpc.ClientStream
}

type routeRecordRouteClient struct {
	grpc.ClientStream
}

func (x *routeRecordRouteClient) Send(m *Point) error {
	return x.ClientStream.SendMsg(m)
}

func (x *routeRecordRouteClient) CloseAndRecv() (*RouteSummary, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RouteSummary)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *routeClient) RouteChat(ctx context.Context, opts ...grpc.CallOption) (Route_RouteChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Route_serviceDesc.Streams[2], "/route.Route/RouteChat", opts...)
	if err != nil {
		return nil, err
	}
	x := &routeRouteChatClient{stream}
	return x, nil
}

type Route_RouteChatClient interface {
	Send(*RouteNote) error
	Recv() (*RouteNote, error)
	grpc.ClientStream
}

type routeRouteChatClient struct {
	grpc.ClientStream
}

func (x *routeRouteChatClient) Send(m *RouteNote) error {
	return x.ClientStream.SendMsg(m)
}

func (x *routeRouteChatClient) Recv() (*RouteNote, error) {
	m := new(RouteNote)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouteServer is the server API for Route service.
type RouteServer interface {
	GetFeature(context.Context, *Point) (*Feature, error)
	ListFeatures(*Rectangle, Route_ListFeaturesServer) error
	RecordRoute(Route_RecordRouteServer) error
	RouteChat(Route_RouteChatServer) error
}

// UnimplementedRouteServer can be embedded to have forward compatible implementations.
type UnimplementedRouteServer struct {
}

func (*UnimplementedRouteServer) GetFeature(context.Context, *Point) (*Feature, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeature not implemented")
}
func (*UnimplementedRouteServer) ListFeatures(*Rectangle, Route_ListFeaturesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListFeatures not implemented")
}
func (*UnimplementedRouteServer) RecordRoute(Route_RecordRouteServer) error {
	return status.Errorf(codes.Unimplemented, "method RecordRoute not implemented")
}
func (*UnimplementedRouteServer) RouteChat(Route_RouteChatServer) error {
	return status.Errorf(codes.Unimplemented, "method RouteChat not implemented")
}

func RegisterRouteServer(s *grpc.Server, srv RouteServer) {
	s.RegisterService(&_Route_serviceDesc, srv)
}

func _Route_GetFeature_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Point)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouteServer).GetFeature(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/route.Route/GetFeature",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouteServer).GetFeature(ctx, req.(*Point))
	}
	return interceptor(ctx, in, info, handler)
}

func _Route_ListFeatures_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Rectangle)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RouteServer).ListFeatures(m, &routeListFeaturesServer{stream})
}

type Route_ListFeaturesServer interface {
	Send(*Feature) error
	grpc.ServerStream
}

type routeListFeaturesServer struct {
	grpc.ServerStream
}

func (x *routeListFeaturesServer) Send(m *Feature) error {
	return x.ServerStream.SendMsg(m)
}

func _Route_RecordRoute_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RouteServer).RecordRoute(&routeRecordRouteServer{stream})
}

type Route_RecordRouteServer interface {
	SendAndClose(*RouteSummary) error
	Recv() (*Point, error)
	grpc.ServerStream
}

type routeRecordRouteServer struct {
	grpc.ServerStream
}

func (x *routeRecordRouteServer) SendAndClose(m *RouteSummary) error {
	return x.ServerStream.SendMsg(m)
}

func (x *routeRecordRouteServer) Recv() (*Point, error) {
	m := new(Point)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Route_RouteChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RouteServer).RouteChat(&routeRouteChatServer{stream})
}

type Route_RouteChatServer interface {
	Send(*RouteNote) error
	Recv() (*RouteNote, error)
	grpc.ServerStream
}

type routeRouteChatServer struct {
	grpc.ServerStream
}

func (x *routeRouteChatServer) Send(m *RouteNote) error {
	return x.ServerStream.SendMsg(m)
}

func (x *routeRouteChatServer) Recv() (*RouteNote, error) {
	m := new(RouteNote)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Route_serviceDesc = grpc.ServiceDesc{
	ServiceName: "route.Route",
	HandlerType: (*RouteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFeature",
			Handler:    _Route_GetFeature_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListFeatures",
			Handler:       _Route_ListFeatures_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "RecordRoute",
			Handler:       _Route_RecordRoute_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RouteChat",
			Handler:       _Route_RouteChat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "route.proto",
}
