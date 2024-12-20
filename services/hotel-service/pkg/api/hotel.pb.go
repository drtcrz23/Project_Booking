// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: hotel.proto

package api

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

// Запрос для получения отеля по ID
type GetHotelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HotelId int32 `protobuf:"varint,1,opt,name=hotel_id,json=hotelId,proto3" json:"hotel_id,omitempty"`
}

func (x *GetHotelRequest) Reset() {
	*x = GetHotelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotel_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetHotelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetHotelRequest) ProtoMessage() {}

func (x *GetHotelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_hotel_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetHotelRequest.ProtoReflect.Descriptor instead.
func (*GetHotelRequest) Descriptor() ([]byte, []int) {
	return file_hotel_proto_rawDescGZIP(), []int{0}
}

func (x *GetHotelRequest) GetHotelId() int32 {
	if x != nil {
		return x.HotelId
	}
	return 0
}

// Структура отеля
type Hotel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Price      string  `protobuf:"bytes,3,opt,name=price,proto3" json:"price,omitempty"`
	HotelierId int32   `protobuf:"varint,4,opt,name=hotelier_id,json=hotelierId,proto3" json:"hotelier_id,omitempty"`
	Rooms      []*Room `protobuf:"bytes,5,rep,name=rooms,proto3" json:"rooms,omitempty"`
}

func (x *Hotel) Reset() {
	*x = Hotel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotel_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Hotel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hotel) ProtoMessage() {}

func (x *Hotel) ProtoReflect() protoreflect.Message {
	mi := &file_hotel_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hotel.ProtoReflect.Descriptor instead.
func (*Hotel) Descriptor() ([]byte, []int) {
	return file_hotel_proto_rawDescGZIP(), []int{1}
}

func (x *Hotel) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Hotel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Hotel) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *Hotel) GetHotelierId() int32 {
	if x != nil {
		return x.HotelierId
	}
	return 0
}

func (x *Hotel) GetRooms() []*Room {
	if x != nil {
		return x.Rooms
	}
	return nil
}

// Структура комнаты
type Room struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	HotelId    int32   `protobuf:"varint,2,opt,name=hotel_id,json=hotelId,proto3" json:"hotel_id,omitempty"`
	RoomNumber string  `protobuf:"bytes,3,opt,name=room_number,json=roomNumber,proto3" json:"room_number,omitempty"`
	Type       string  `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Price      float32 `protobuf:"fixed32,5,opt,name=price,proto3" json:"price,omitempty"`
	Status     string  `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"` // Статус комнаты (available, booked, maintenance)
}

func (x *Room) Reset() {
	*x = Room{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotel_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Room) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Room) ProtoMessage() {}

func (x *Room) ProtoReflect() protoreflect.Message {
	mi := &file_hotel_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Room.ProtoReflect.Descriptor instead.
func (*Room) Descriptor() ([]byte, []int) {
	return file_hotel_proto_rawDescGZIP(), []int{2}
}

func (x *Room) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Room) GetHotelId() int32 {
	if x != nil {
		return x.HotelId
	}
	return 0
}

func (x *Room) GetRoomNumber() string {
	if x != nil {
		return x.RoomNumber
	}
	return ""
}

func (x *Room) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Room) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Room) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_hotel_proto protoreflect.FileDescriptor

var file_hotel_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x68, 0x6f, 0x74, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x68,
	0x6f, 0x74, 0x65, 0x6c, 0x22, 0x2c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x48, 0x6f, 0x74, 0x65, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x6f, 0x74, 0x65, 0x6c,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x68, 0x6f, 0x74, 0x65, 0x6c,
	0x49, 0x64, 0x22, 0x85, 0x01, 0x0a, 0x05, 0x48, 0x6f, 0x74, 0x65, 0x6c, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x6f, 0x74, 0x65, 0x6c, 0x69,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x68, 0x6f, 0x74,
	0x65, 0x6c, 0x69, 0x65, 0x72, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x68, 0x6f, 0x74, 0x65, 0x6c, 0x2e, 0x52,
	0x6f, 0x6f, 0x6d, 0x52, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73, 0x22, 0x94, 0x01, 0x0a, 0x04, 0x52,
	0x6f, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x6f, 0x74, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x68, 0x6f, 0x74, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x6f, 0x6f, 0x6d, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x32, 0x44, 0x0a, 0x0c, 0x48, 0x6f, 0x74, 0x65, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x34, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x48, 0x6f, 0x74, 0x65, 0x6c, 0x42, 0x79, 0x49,
	0x64, 0x12, 0x16, 0x2e, 0x68, 0x6f, 0x74, 0x65, 0x6c, 0x2e, 0x47, 0x65, 0x74, 0x48, 0x6f, 0x74,
	0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x68, 0x6f, 0x74, 0x65,
	0x6c, 0x2e, 0x48, 0x6f, 0x74, 0x65, 0x6c, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2e, 0x2f, 0x2e, 0x2e,
	0x2f, 0x2e, 0x2e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hotel_proto_rawDescOnce sync.Once
	file_hotel_proto_rawDescData = file_hotel_proto_rawDesc
)

func file_hotel_proto_rawDescGZIP() []byte {
	file_hotel_proto_rawDescOnce.Do(func() {
		file_hotel_proto_rawDescData = protoimpl.X.CompressGZIP(file_hotel_proto_rawDescData)
	})
	return file_hotel_proto_rawDescData
}

var file_hotel_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_hotel_proto_goTypes = []interface{}{
	(*GetHotelRequest)(nil), // 0: hotel.GetHotelRequest
	(*Hotel)(nil),           // 1: hotel.Hotel
	(*Room)(nil),            // 2: hotel.Room
}
var file_hotel_proto_depIdxs = []int32{
	2, // 0: hotel.Hotel.rooms:type_name -> hotel.Room
	0, // 1: hotel.HotelService.GetHotelById:input_type -> hotel.GetHotelRequest
	1, // 2: hotel.HotelService.GetHotelById:output_type -> hotel.Hotel
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_hotel_proto_init() }
func file_hotel_proto_init() {
	if File_hotel_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hotel_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetHotelRequest); i {
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
		file_hotel_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Hotel); i {
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
		file_hotel_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Room); i {
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
			RawDescriptor: file_hotel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_hotel_proto_goTypes,
		DependencyIndexes: file_hotel_proto_depIdxs,
		MessageInfos:      file_hotel_proto_msgTypes,
	}.Build()
	File_hotel_proto = out.File
	file_hotel_proto_rawDesc = nil
	file_hotel_proto_goTypes = nil
	file_hotel_proto_depIdxs = nil
}
