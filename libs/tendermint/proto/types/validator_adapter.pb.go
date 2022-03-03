// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tendermint/types/validator.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	pc "github.com/okex/exchain/libs/tendermint/proto/crypto/keys"
	io "io"
	math "math"
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


type SimpleValidator struct {
	PubKey      *pc.PublicKey `protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	VotingPower int64             `protobuf:"varint,2,opt,name=voting_power,json=votingPower,proto3" json:"voting_power,omitempty"`
}

func (m *SimpleValidator) Reset()         { *m = SimpleValidator{} }
func (m *SimpleValidator) String() string { return proto.CompactTextString(m) }
func (*SimpleValidator) ProtoMessage()    {}
func (*SimpleValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e92274df03d3088, []int{2}
}
func (m *SimpleValidator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SimpleValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SimpleValidator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SimpleValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimpleValidator.Merge(m, src)
}
func (m *SimpleValidator) XXX_Size() int {
	return m.Size()
}
func (m *SimpleValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_SimpleValidator.DiscardUnknown(m)
}

var xxx_messageInfo_SimpleValidator proto.InternalMessageInfo

func (m *SimpleValidator) GetPubKey() *pc.PublicKey {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *SimpleValidator) GetVotingPower() int64 {
	if m != nil {
		return m.VotingPower
	}
	return 0
}

func init() {
	proto.RegisterType((*ValidatorSet)(nil), "tendermint.types.ValidatorSet")
	proto.RegisterType((*Validator)(nil), "tendermint.types.Validator")
	proto.RegisterType((*SimpleValidator)(nil), "tendermint.types.SimpleValidator")
}

func init() { proto.RegisterFile("tendermint/types/validator.proto", fileDescriptor_4e92274df03d3088) }

var fileDescriptor_4e92274df03d3088 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcf, 0x4e, 0xc2, 0x40,
	0x10, 0xc6, 0xbb, 0x40, 0x40, 0x17, 0x12, 0x71, 0xe3, 0xa1, 0x41, 0x52, 0x2b, 0x27, 0x12, 0x4d,
	0x9b, 0x68, 0x0c, 0x07, 0x6e, 0x5c, 0xb9, 0x60, 0x49, 0x38, 0x78, 0x69, 0x5a, 0xba, 0xa9, 0x1b,
	0x0a, 0xbb, 0xd9, 0x6e, 0x31, 0xfb, 0x16, 0x3e, 0x8b, 0x4f, 0xc1, 0x91, 0xa3, 0x27, 0x63, 0xe0,
	0x45, 0x4c, 0x5b, 0xfa, 0x27, 0xa8, 0xe1, 0x36, 0x9d, 0xef, 0x9b, 0x99, 0x5f, 0x37, 0x1f, 0xd4,
	0x05, 0x5e, 0x79, 0x98, 0x2f, 0xc9, 0x4a, 0x98, 0x42, 0x32, 0x1c, 0x9a, 0x6b, 0x27, 0x20, 0x9e,
	0x23, 0x28, 0x37, 0x18, 0xa7, 0x82, 0xa2, 0x76, 0xe1, 0x30, 0x12, 0x47, 0xe7, 0xca, 0xa7, 0x3e,
	0x4d, 0x44, 0x33, 0xae, 0x52, 0x5f, 0xa7, 0x5b, 0xda, 0x34, 0xe7, 0x92, 0x09, 0x6a, 0x2e, 0xb0,
	0x0c, 0x53, 0xb5, 0xf7, 0x01, 0x60, 0x6b, 0x96, 0x6d, 0x9e, 0x62, 0x81, 0x86, 0x10, 0xe6, 0x97,
	0x42, 0x15, 0xe8, 0xd5, 0x7e, 0xf3, 0xe1, 0xda, 0x38, 0xbe, 0x65, 0xe4, 0x33, 0x56, 0xc9, 0x8e,
	0x06, 0xf0, 0x8c, 0x71, 0xca, 0x68, 0x88, 0xb9, 0x5a, 0xd1, 0xc1, 0xa9, 0xd1, 0xdc, 0x8c, 0xee,
	0x21, 0x12, 0x54, 0x38, 0x81, 0xbd, 0xa6, 0x82, 0xac, 0x7c, 0x9b, 0xd1, 0x37, 0xcc, 0xd5, 0xaa,
	0x0e, 0xfa, 0x55, 0xab, 0x9d, 0x28, 0xb3, 0x44, 0x98, 0xc4, 0xfd, 0x18, 0xfa, 0x3c, 0xdf, 0x82,
	0x54, 0xd8, 0x70, 0x3c, 0x8f, 0xe3, 0x30, 0xc6, 0x05, 0xfd, 0x96, 0x95, 0x7d, 0xa2, 0x21, 0x6c,
	0xb0, 0xc8, 0xb5, 0x17, 0x58, 0x1e, 0x68, 0xba, 0x65, 0x9a, 0xf4, 0x31, 0x8c, 0x49, 0xe4, 0x06,
	0x64, 0x3e, 0xc6, 0x72, 0x54, 0xdb, 0x7c, 0xdd, 0x28, 0x56, 0x9d, 0x45, 0xee, 0x18, 0x4b, 0x74,
	0x0b, 0x5b, 0x7f, 0xc0, 0x34, 0xd7, 0x05, 0x07, 0xba, 0x83, 0x97, 0xd9, 0x1f, 0xd8, 0x8c, 0x13,
	0xca, 0x89, 0x90, 0x6a, 0x2d, 0x85, 0xce, 0x84, 0xc9, 0xa1, 0xdf, 0x5b, 0xc0, 0x8b, 0x29, 0x59,
	0xb2, 0x00, 0x17, 0xe4, 0x4f, 0x05, 0x1f, 0x38, 0xcd, 0xf7, 0x2f, 0x59, 0xe5, 0x17, 0xd9, 0xe8,
	0x79, 0xb3, 0xd3, 0xc0, 0x76, 0xa7, 0x81, 0xef, 0x9d, 0x06, 0xde, 0xf7, 0x9a, 0xb2, 0xdd, 0x6b,
	0xca, 0xe7, 0x5e, 0x53, 0x5e, 0x06, 0x3e, 0x11, 0xaf, 0x91, 0x6b, 0xcc, 0xe9, 0xd2, 0x2c, 0x67,
	0xac, 0x28, 0xd3, 0x04, 0x1d, 0xe7, 0xcf, 0xad, 0x27, 0xfd, 0xc7, 0x9f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x48, 0xbf, 0x34, 0x35, 0x9a, 0x02, 0x00, 0x00,
}


func (m *SimpleValidator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SimpleValidator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SimpleValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.VotingPower != 0 {
		i = encodeVarintValidator(dAtA, i, uint64(m.VotingPower))
		i--
		dAtA[i] = 0x10
	}
	if m.PubKey != nil {
		{
			size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintValidator(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}


func (m *SimpleValidator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PubKey != nil {
		l = m.PubKey.Size()
		n += 1 + l + sovValidator(uint64(l))
	}
	if m.VotingPower != 0 {
		n += 1 + sovValidator(uint64(m.VotingPower))
	}
	return n
}

func (m *SimpleValidator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowValidator
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
			return fmt.Errorf("proto: SimpleValidator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SimpleValidator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthValidator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthValidator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PubKey == nil {
				m.PubKey = &pc.PublicKey{}
			}
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VotingPower", wireType)
			}
			m.VotingPower = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowValidator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VotingPower |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipValidator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthValidator
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}