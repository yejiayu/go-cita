// Code generated by protoc-gen-go. DO NOT EDIT.
// source: consensus.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type VoteType int32

const (
	VoteType_Prevote   VoteType = 0
	VoteType_precommit VoteType = 1
)

var VoteType_name = map[int32]string{
	0: "Prevote",
	1: "precommit",
}
var VoteType_value = map[string]int32{
	"Prevote":   0,
	"precommit": 1,
}

func (x VoteType) String() string {
	return proto.EnumName(VoteType_name, int32(x))
}
func (VoteType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_consensus_f7d7a63bbbcf5fe0, []int{0}
}

type Proposal struct {
	Block                *Block   `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
	Islock               bool     `protobuf:"varint,2,opt,name=islock" json:"islock,omitempty"`
	LockRound            uint64   `protobuf:"varint,3,opt,name=lock_round,json=lockRound" json:"lock_round,omitempty"`
	LockVotes            []*Vote  `protobuf:"bytes,4,rep,name=lock_votes,json=lockVotes" json:"lock_votes,omitempty"`
	Round                uint64   `protobuf:"varint,5,opt,name=round" json:"round,omitempty"`
	Height               uint64   `protobuf:"varint,6,opt,name=height" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Proposal) Reset()         { *m = Proposal{} }
func (m *Proposal) String() string { return proto.CompactTextString(m) }
func (*Proposal) ProtoMessage()    {}
func (*Proposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_f7d7a63bbbcf5fe0, []int{0}
}
func (m *Proposal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Proposal.Unmarshal(m, b)
}
func (m *Proposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Proposal.Marshal(b, m, deterministic)
}
func (dst *Proposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proposal.Merge(dst, src)
}
func (m *Proposal) XXX_Size() int {
	return xxx_messageInfo_Proposal.Size(m)
}
func (m *Proposal) XXX_DiscardUnknown() {
	xxx_messageInfo_Proposal.DiscardUnknown(m)
}

var xxx_messageInfo_Proposal proto.InternalMessageInfo

func (m *Proposal) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

func (m *Proposal) GetIslock() bool {
	if m != nil {
		return m.Islock
	}
	return false
}

func (m *Proposal) GetLockRound() uint64 {
	if m != nil {
		return m.LockRound
	}
	return 0
}

func (m *Proposal) GetLockVotes() []*Vote {
	if m != nil {
		return m.LockVotes
	}
	return nil
}

func (m *Proposal) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *Proposal) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type SetProposalReq struct {
	Proposal             *Proposal `protobuf:"bytes,1,opt,name=proposal" json:"proposal,omitempty"`
	Signature            []byte    `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *SetProposalReq) Reset()         { *m = SetProposalReq{} }
func (m *SetProposalReq) String() string { return proto.CompactTextString(m) }
func (*SetProposalReq) ProtoMessage()    {}
func (*SetProposalReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_f7d7a63bbbcf5fe0, []int{1}
}
func (m *SetProposalReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetProposalReq.Unmarshal(m, b)
}
func (m *SetProposalReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetProposalReq.Marshal(b, m, deterministic)
}
func (dst *SetProposalReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetProposalReq.Merge(dst, src)
}
func (m *SetProposalReq) XXX_Size() int {
	return xxx_messageInfo_SetProposalReq.Size(m)
}
func (m *SetProposalReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SetProposalReq.DiscardUnknown(m)
}

var xxx_messageInfo_SetProposalReq proto.InternalMessageInfo

func (m *SetProposalReq) GetProposal() *Proposal {
	if m != nil {
		return m.Proposal
	}
	return nil
}

func (m *SetProposalReq) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type Vote struct {
	Height               uint64   `protobuf:"varint,1,opt,name=height" json:"height,omitempty"`
	Round                uint64   `protobuf:"varint,2,opt,name=round" json:"round,omitempty"`
	VoteType             VoteType `protobuf:"varint,3,opt,name=voteType,enum=VoteType" json:"voteType,omitempty"`
	Address              []byte   `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Hash                 []byte   `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vote) Reset()         { *m = Vote{} }
func (m *Vote) String() string { return proto.CompactTextString(m) }
func (*Vote) ProtoMessage()    {}
func (*Vote) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_f7d7a63bbbcf5fe0, []int{2}
}
func (m *Vote) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vote.Unmarshal(m, b)
}
func (m *Vote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vote.Marshal(b, m, deterministic)
}
func (dst *Vote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vote.Merge(dst, src)
}
func (m *Vote) XXX_Size() int {
	return xxx_messageInfo_Vote.Size(m)
}
func (m *Vote) XXX_DiscardUnknown() {
	xxx_messageInfo_Vote.DiscardUnknown(m)
}

var xxx_messageInfo_Vote proto.InternalMessageInfo

func (m *Vote) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Vote) GetRound() uint64 {
	if m != nil {
		return m.Round
	}
	return 0
}

func (m *Vote) GetVoteType() VoteType {
	if m != nil {
		return m.VoteType
	}
	return VoteType_Prevote
}

func (m *Vote) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Vote) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

type AddVoteReq struct {
	Vote                 *Vote    `protobuf:"bytes,1,opt,name=vote" json:"vote,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddVoteReq) Reset()         { *m = AddVoteReq{} }
func (m *AddVoteReq) String() string { return proto.CompactTextString(m) }
func (*AddVoteReq) ProtoMessage()    {}
func (*AddVoteReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_consensus_f7d7a63bbbcf5fe0, []int{3}
}
func (m *AddVoteReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddVoteReq.Unmarshal(m, b)
}
func (m *AddVoteReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddVoteReq.Marshal(b, m, deterministic)
}
func (dst *AddVoteReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddVoteReq.Merge(dst, src)
}
func (m *AddVoteReq) XXX_Size() int {
	return xxx_messageInfo_AddVoteReq.Size(m)
}
func (m *AddVoteReq) XXX_DiscardUnknown() {
	xxx_messageInfo_AddVoteReq.DiscardUnknown(m)
}

var xxx_messageInfo_AddVoteReq proto.InternalMessageInfo

func (m *AddVoteReq) GetVote() *Vote {
	if m != nil {
		return m.Vote
	}
	return nil
}

func (m *AddVoteReq) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*Proposal)(nil), "Proposal")
	proto.RegisterType((*SetProposalReq)(nil), "SetProposalReq")
	proto.RegisterType((*Vote)(nil), "Vote")
	proto.RegisterType((*AddVoteReq)(nil), "AddVoteReq")
	proto.RegisterEnum("VoteType", VoteType_name, VoteType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConsensusClient is the client API for Consensus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConsensusClient interface {
	SetProposal(ctx context.Context, in *SetProposalReq, opts ...grpc.CallOption) (*Empty, error)
	AddVote(ctx context.Context, in *AddVoteReq, opts ...grpc.CallOption) (*Empty, error)
}

type consensusClient struct {
	cc *grpc.ClientConn
}

func NewConsensusClient(cc *grpc.ClientConn) ConsensusClient {
	return &consensusClient{cc}
}

func (c *consensusClient) SetProposal(ctx context.Context, in *SetProposalReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Consensus/SetProposal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consensusClient) AddVote(ctx context.Context, in *AddVoteReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Consensus/AddVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConsensusServer is the server API for Consensus service.
type ConsensusServer interface {
	SetProposal(context.Context, *SetProposalReq) (*Empty, error)
	AddVote(context.Context, *AddVoteReq) (*Empty, error)
}

func RegisterConsensusServer(s *grpc.Server, srv ConsensusServer) {
	s.RegisterService(&_Consensus_serviceDesc, srv)
}

func _Consensus_SetProposal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetProposalReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).SetProposal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Consensus/SetProposal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).SetProposal(ctx, req.(*SetProposalReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Consensus_AddVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddVoteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsensusServer).AddVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Consensus/AddVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsensusServer).AddVote(ctx, req.(*AddVoteReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Consensus_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Consensus",
	HandlerType: (*ConsensusServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetProposal",
			Handler:    _Consensus_SetProposal_Handler,
		},
		{
			MethodName: "AddVote",
			Handler:    _Consensus_AddVote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "consensus.proto",
}

func init() { proto.RegisterFile("consensus.proto", fileDescriptor_consensus_f7d7a63bbbcf5fe0) }

var fileDescriptor_consensus_f7d7a63bbbcf5fe0 = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0xc1, 0x6e, 0x9b, 0x40,
	0x14, 0xf4, 0xda, 0x80, 0xe1, 0xe1, 0xda, 0xd6, 0x53, 0x55, 0x51, 0xcb, 0x95, 0x10, 0xaa, 0x2b,
	0xd4, 0x03, 0x07, 0xfa, 0x05, 0x6d, 0xe5, 0xbb, 0x45, 0x9b, 0xe4, 0x18, 0x61, 0x58, 0x19, 0x14,
	0x9b, 0xc5, 0xec, 0x3a, 0x92, 0x3f, 0x21, 0xdf, 0x93, 0x1f, 0x8c, 0x76, 0x59, 0x0c, 0xb9, 0xe4,
	0xc4, 0xce, 0xcc, 0x3e, 0xde, 0xcc, 0x68, 0x61, 0x91, 0xb1, 0x8a, 0xd3, 0x8a, 0x5f, 0x78, 0x54,
	0x37, 0x4c, 0xb0, 0x95, 0x9b, 0x15, 0x69, 0x59, 0x69, 0x30, 0xcb, 0xd8, 0xe9, 0xc4, 0x34, 0x0a,
	0x5e, 0x09, 0xd8, 0xbb, 0x86, 0xd5, 0x8c, 0xa7, 0x47, 0x5c, 0x83, 0xb9, 0x3f, 0xb2, 0xec, 0xc9,
	0x23, 0x3e, 0x09, 0xdd, 0xd8, 0x8a, 0xfe, 0x48, 0x94, 0xb4, 0x24, 0x7e, 0x01, 0xab, 0xe4, 0x4a,
	0x1e, 0xfb, 0x24, 0xb4, 0x13, 0x8d, 0xf0, 0x1b, 0x80, 0xfc, 0x3e, 0x36, 0xec, 0x52, 0xe5, 0xde,
	0xc4, 0x27, 0xa1, 0x91, 0x38, 0x6a, 0x50, 0x12, 0xf8, 0x5d, 0xcb, 0xcf, 0x4c, 0x50, 0xee, 0x19,
	0xfe, 0x24, 0x74, 0x63, 0x33, 0xba, 0x67, 0x82, 0xb6, 0xb7, 0xe4, 0x89, 0xe3, 0x67, 0x30, 0xdb,
	0x79, 0x53, 0xcd, 0xb7, 0x40, 0xae, 0x2c, 0x68, 0x79, 0x28, 0x84, 0x67, 0x29, 0x5a, 0xa3, 0xe0,
	0x0e, 0xe6, 0xff, 0xa8, 0xe8, 0x7c, 0x27, 0xf4, 0x8c, 0x1b, 0xb0, 0x6b, 0x0d, 0xb5, 0x7b, 0x27,
	0xba, 0xe9, 0x37, 0x09, 0xd7, 0xe0, 0xf0, 0xf2, 0x50, 0xa5, 0xe2, 0xd2, 0x50, 0x15, 0x63, 0x96,
	0xf4, 0x44, 0xf0, 0x42, 0xc0, 0x90, 0x76, 0x06, 0x7b, 0xc9, 0x70, 0x6f, 0xef, 0x72, 0x3c, 0x74,
	0xb9, 0x01, 0x5b, 0x86, 0xfb, 0x7f, 0xad, 0xa9, 0x8a, 0x3f, 0x8f, 0x1d, 0x95, 0x4f, 0x12, 0xc9,
	0x4d, 0x42, 0x0f, 0xa6, 0x69, 0x9e, 0x37, 0x94, 0xcb, 0x16, 0xe4, 0xe6, 0x0e, 0x22, 0x82, 0x51,
	0xa4, 0xbc, 0x50, 0xd9, 0x67, 0x89, 0x3a, 0x07, 0x5b, 0x80, 0xdf, 0x79, 0xae, 0x6a, 0xa2, 0x67,
	0xfc, 0x0a, 0x86, 0xfc, 0x8f, 0x8e, 0xa6, 0xeb, 0x53, 0xd4, 0xc7, 0x91, 0x7e, 0xfe, 0x00, 0xbb,
	0xb3, 0x82, 0x2e, 0x4c, 0x77, 0x0d, 0x95, 0x43, 0xcb, 0x11, 0x7e, 0x02, 0xa7, 0x6e, 0xa8, 0x7c,
	0x0b, 0xa5, 0x58, 0x92, 0xf8, 0x01, 0x9c, 0xbf, 0xdd, 0xab, 0xc1, 0x10, 0xdc, 0x41, 0xbd, 0xb8,
	0x88, 0xde, 0x97, 0xbd, 0xb2, 0xa2, 0xed, 0xa9, 0x16, 0xd7, 0x60, 0x84, 0x3e, 0x4c, 0xb5, 0x4b,
	0x74, 0xa3, 0xde, 0x6f, 0x7f, 0x63, 0x6f, 0xa9, 0x77, 0xf6, 0xeb, 0x2d, 0x00, 0x00, 0xff, 0xff,
	0x35, 0x6b, 0x75, 0x71, 0x95, 0x02, 0x00, 0x00,
}
