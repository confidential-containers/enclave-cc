// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/confidential-containers/enclave-cc/src/shim/runtime/v2/rune/rune/agent/grpc/image.proto

package grpc

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"

	github_com_containerd_ttrpc "github.com/containerd/ttrpc"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type PullImageRequest struct {
	// Image name (e.g. docker.io/library/busybox:latest).
	Image string `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	// Unique image identifier, used to avoid duplication when unpacking the image layers.
	ContainerId          string   `protobuf:"bytes,2,opt,name=container_id,json=containerId,proto3" json:"container_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PullImageRequest) Reset()      { *m = PullImageRequest{} }
func (*PullImageRequest) ProtoMessage() {}
func (*PullImageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9624c68e2b547544, []int{0}
}
func (m *PullImageRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PullImageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PullImageRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PullImageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PullImageRequest.Merge(m, src)
}
func (m *PullImageRequest) XXX_Size() int {
	return m.Size()
}
func (m *PullImageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PullImageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PullImageRequest proto.InternalMessageInfo

type PullImageResponse struct {
	// Reference to the image in use. For most runtimes, this should be an
	// image ID or digest.
	ImageRef             string   `protobuf:"bytes,1,opt,name=image_ref,json=imageRef,proto3" json:"image_ref,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PullImageResponse) Reset()      { *m = PullImageResponse{} }
func (*PullImageResponse) ProtoMessage() {}
func (*PullImageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9624c68e2b547544, []int{1}
}
func (m *PullImageResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PullImageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PullImageResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PullImageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PullImageResponse.Merge(m, src)
}
func (m *PullImageResponse) XXX_Size() int {
	return m.Size()
}
func (m *PullImageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PullImageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PullImageResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*PullImageRequest)(nil), "grpc.PullImageRequest")
	proto.RegisterType((*PullImageResponse)(nil), "grpc.PullImageResponse")
}

func init() { proto.RegisterFile("image.proto", fileDescriptor_9624c68e2b547544) }

var fileDescriptor_9624c68e2b547544 = []byte{
	// 252 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xcc, 0x4d, 0x4c,
	0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x49, 0x2f, 0x2a, 0x48, 0x56, 0xf2, 0xe6,
	0x12, 0x08, 0x28, 0xcd, 0xc9, 0xf1, 0x04, 0x49, 0x04, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08,
	0x89, 0x70, 0xb1, 0x82, 0x15, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x38, 0x42, 0x8a,
	0x5c, 0x3c, 0xc9, 0xf9, 0x79, 0x25, 0x89, 0x99, 0x79, 0xa9, 0x45, 0xf1, 0x99, 0x29, 0x12, 0x4c,
	0x60, 0x49, 0x6e, 0xb8, 0x98, 0x67, 0x8a, 0x92, 0x01, 0x97, 0x20, 0x92, 0x61, 0xc5, 0x05, 0xf9,
	0x79, 0xc5, 0xa9, 0x42, 0xd2, 0x5c, 0x9c, 0x60, 0x03, 0xe2, 0x8b, 0x52, 0xd3, 0xa0, 0x26, 0x72,
	0x64, 0x42, 0x54, 0xa4, 0x19, 0xb9, 0x73, 0xb1, 0x82, 0x55, 0x0b, 0xd9, 0x71, 0x71, 0xc2, 0xb5,
	0x0a, 0x89, 0xe9, 0x81, 0xdc, 0xa6, 0x87, 0xee, 0x30, 0x29, 0x71, 0x0c, 0x71, 0x88, 0x1d, 0x4a,
	0x0c, 0x4e, 0x15, 0x27, 0x1e, 0xca, 0x31, 0xdc, 0x78, 0x28, 0xc7, 0xd0, 0xf0, 0x48, 0x8e, 0xf1,
	0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x8c, 0x8a, 0x4b, 0xcf,
	0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0xcf, 0x4e, 0x2c, 0x49, 0xd4, 0x85, 0xbb,
	0xb8, 0x18, 0x83, 0x5f, 0x5c, 0x94, 0xac, 0x5f, 0x54, 0x9a, 0x57, 0x92, 0x99, 0x9b, 0xaa, 0x5f,
	0x96, 0x59, 0x54, 0x82, 0x24, 0x55, 0x90, 0x9d, 0xae, 0x9f, 0x98, 0x9e, 0x9a, 0x57, 0xa2, 0x0f,
	0x0e, 0xbe, 0xe4, 0xfc, 0x9c, 0x62, 0x7d, 0x90, 0x6b, 0x92, 0xd8, 0xc0, 0x7c, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xca, 0x2d, 0x8f, 0xdf, 0x5d, 0x01, 0x00, 0x00,
}

func (m *PullImageRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PullImageRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PullImageRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ContainerId) > 0 {
		i -= len(m.ContainerId)
		copy(dAtA[i:], m.ContainerId)
		i = encodeVarintImage(dAtA, i, uint64(len(m.ContainerId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Image) > 0 {
		i -= len(m.Image)
		copy(dAtA[i:], m.Image)
		i = encodeVarintImage(dAtA, i, uint64(len(m.Image)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PullImageResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PullImageResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PullImageResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ImageRef) > 0 {
		i -= len(m.ImageRef)
		copy(dAtA[i:], m.ImageRef)
		i = encodeVarintImage(dAtA, i, uint64(len(m.ImageRef)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintImage(dAtA []byte, offset int, v uint64) int {
	offset -= sovImage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PullImageRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Image)
	if l > 0 {
		n += 1 + l + sovImage(uint64(l))
	}
	l = len(m.ContainerId)
	if l > 0 {
		n += 1 + l + sovImage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *PullImageResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ImageRef)
	if l > 0 {
		n += 1 + l + sovImage(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovImage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozImage(x uint64) (n int) {
	return sovImage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *PullImageRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PullImageRequest{`,
		`Image:` + fmt.Sprintf("%v", this.Image) + `,`,
		`ContainerId:` + fmt.Sprintf("%v", this.ContainerId) + `,`,
		`XXX_unrecognized:` + fmt.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func (this *PullImageResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&PullImageResponse{`,
		`ImageRef:` + fmt.Sprintf("%v", this.ImageRef) + `,`,
		`XXX_unrecognized:` + fmt.Sprintf("%v", this.XXX_unrecognized) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringImage(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}

type ImageService interface {
	PullImage(ctx context.Context, req *PullImageRequest) (*PullImageResponse, error)
}

func RegisterImageService(srv *github_com_containerd_ttrpc.Server, svc ImageService) {
	srv.Register("grpc.Image", map[string]github_com_containerd_ttrpc.Method{
		"PullImage": func(ctx context.Context, unmarshal func(interface{}) error) (interface{}, error) {
			var req PullImageRequest
			if err := unmarshal(&req); err != nil {
				return nil, err
			}
			return svc.PullImage(ctx, &req)
		},
	})
}

type imageClient struct {
	client *github_com_containerd_ttrpc.Client
}

func NewImageClient(client *github_com_containerd_ttrpc.Client) ImageService {
	return &imageClient{
		client: client,
	}
}

func (c *imageClient) PullImage(ctx context.Context, req *PullImageRequest) (*PullImageResponse, error) {
	var resp PullImageResponse
	if err := c.client.Call(ctx, "grpc.Image", "PullImage", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (m *PullImageRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowImage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PullImageRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PullImageRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Image", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowImage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthImage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthImage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Image = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContainerId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowImage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthImage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthImage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ContainerId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipImage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthImage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PullImageResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowImage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PullImageResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PullImageResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageRef", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowImage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthImage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthImage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageRef = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipImage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthImage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipImage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowImage
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowImage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowImage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthImage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupImage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthImage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthImage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowImage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupImage = fmt.Errorf("proto: unexpected end of group")
)
