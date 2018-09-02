// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chain.proto

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

type ProofType int32

const (
	ProofType_AuthorityRound ProofType = 0
	ProofType_Raft           ProofType = 1
	ProofType_Tendermint     ProofType = 2
)

var ProofType_name = map[int32]string{
	0: "AuthorityRound",
	1: "Raft",
	2: "Tendermint",
}
var ProofType_value = map[string]int32{
	"AuthorityRound": 0,
	"Raft":           1,
	"Tendermint":     2,
}

func (x ProofType) String() string {
	return proto.EnumName(ProofType_name, int32(x))
}
func (ProofType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{0}
}

type Proof struct {
	Content              []byte    `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Type                 ProofType `protobuf:"varint,2,opt,name=type,enum=ProofType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Proof) Reset()         { *m = Proof{} }
func (m *Proof) String() string { return proto.CompactTextString(m) }
func (*Proof) ProtoMessage()    {}
func (*Proof) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{0}
}
func (m *Proof) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Proof.Unmarshal(m, b)
}
func (m *Proof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Proof.Marshal(b, m, deterministic)
}
func (dst *Proof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proof.Merge(dst, src)
}
func (m *Proof) XXX_Size() int {
	return xxx_messageInfo_Proof.Size(m)
}
func (m *Proof) XXX_DiscardUnknown() {
	xxx_messageInfo_Proof.DiscardUnknown(m)
}

var xxx_messageInfo_Proof proto.InternalMessageInfo

func (m *Proof) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *Proof) GetType() ProofType {
	if m != nil {
		return m.Type
	}
	return ProofType_AuthorityRound
}

type Block struct {
	Version              uint32       `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Header               *BlockHeader `protobuf:"bytes,2,opt,name=header" json:"header,omitempty"`
	Body                 *BlockBody   `protobuf:"bytes,3,opt,name=body" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{1}
}
func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (dst *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(dst, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Block) GetHeader() *BlockHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Block) GetBody() *BlockBody {
	if m != nil {
		return m.Body
	}
	return nil
}

type BlockHeader struct {
	Prevhash             []byte   `protobuf:"bytes,1,opt,name=prevhash,proto3" json:"prevhash,omitempty"`
	Timestamp            uint64   `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
	Height               uint64   `protobuf:"varint,3,opt,name=height" json:"height,omitempty"`
	StateRoot            []byte   `protobuf:"bytes,4,opt,name=state_root,json=stateRoot,proto3" json:"state_root,omitempty"`
	TransactionsRoot     []byte   `protobuf:"bytes,5,opt,name=transactions_root,json=transactionsRoot,proto3" json:"transactions_root,omitempty"`
	ReceiptsRoot         []byte   `protobuf:"bytes,6,opt,name=receipts_root,json=receiptsRoot,proto3" json:"receipts_root,omitempty"`
	QuotaUsed            uint64   `protobuf:"varint,7,opt,name=quota_used,json=quotaUsed" json:"quota_used,omitempty"`
	QuotaLimit           uint64   `protobuf:"varint,8,opt,name=quota_limit,json=quotaLimit" json:"quota_limit,omitempty"`
	Proof                *Proof   `protobuf:"bytes,9,opt,name=proof" json:"proof,omitempty"`
	Proposer             []byte   `protobuf:"bytes,10,opt,name=proposer,proto3" json:"proposer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockHeader) Reset()         { *m = BlockHeader{} }
func (m *BlockHeader) String() string { return proto.CompactTextString(m) }
func (*BlockHeader) ProtoMessage()    {}
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{2}
}
func (m *BlockHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockHeader.Unmarshal(m, b)
}
func (m *BlockHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockHeader.Marshal(b, m, deterministic)
}
func (dst *BlockHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHeader.Merge(dst, src)
}
func (m *BlockHeader) XXX_Size() int {
	return xxx_messageInfo_BlockHeader.Size(m)
}
func (m *BlockHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHeader proto.InternalMessageInfo

func (m *BlockHeader) GetPrevhash() []byte {
	if m != nil {
		return m.Prevhash
	}
	return nil
}

func (m *BlockHeader) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *BlockHeader) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockHeader) GetStateRoot() []byte {
	if m != nil {
		return m.StateRoot
	}
	return nil
}

func (m *BlockHeader) GetTransactionsRoot() []byte {
	if m != nil {
		return m.TransactionsRoot
	}
	return nil
}

func (m *BlockHeader) GetReceiptsRoot() []byte {
	if m != nil {
		return m.ReceiptsRoot
	}
	return nil
}

func (m *BlockHeader) GetQuotaUsed() uint64 {
	if m != nil {
		return m.QuotaUsed
	}
	return 0
}

func (m *BlockHeader) GetQuotaLimit() uint64 {
	if m != nil {
		return m.QuotaLimit
	}
	return 0
}

func (m *BlockHeader) GetProof() *Proof {
	if m != nil {
		return m.Proof
	}
	return nil
}

func (m *BlockHeader) GetProposer() []byte {
	if m != nil {
		return m.Proposer
	}
	return nil
}

// data precompile API
type BlockBody struct {
	TxHashes             [][]byte `protobuf:"bytes,1,rep,name=tx_hashes,json=txHashes,proto3" json:"tx_hashes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockBody) Reset()         { *m = BlockBody{} }
func (m *BlockBody) String() string { return proto.CompactTextString(m) }
func (*BlockBody) ProtoMessage()    {}
func (*BlockBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{3}
}
func (m *BlockBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockBody.Unmarshal(m, b)
}
func (m *BlockBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockBody.Marshal(b, m, deterministic)
}
func (dst *BlockBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockBody.Merge(dst, src)
}
func (m *BlockBody) XXX_Size() int {
	return xxx_messageInfo_BlockBody.Size(m)
}
func (m *BlockBody) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockBody.DiscardUnknown(m)
}

var xxx_messageInfo_BlockBody proto.InternalMessageInfo

func (m *BlockBody) GetTxHashes() [][]byte {
	if m != nil {
		return m.TxHashes
	}
	return nil
}

type Receipt struct {
	StateRoot            []byte      `protobuf:"bytes,1,opt,name=state_root,json=stateRoot,proto3" json:"state_root,omitempty"`
	QuotaUsed            uint64      `protobuf:"varint,2,opt,name=quota_used,json=quotaUsed" json:"quota_used,omitempty"`
	Quota                uint64      `protobuf:"varint,3,opt,name=quota" json:"quota,omitempty"`
	LogBloom             []byte      `protobuf:"bytes,4,opt,name=log_bloom,json=logBloom,proto3" json:"log_bloom,omitempty"`
	Logs                 []*LogEntry `protobuf:"bytes,5,rep,name=logs" json:"logs,omitempty"`
	Error                string      `protobuf:"bytes,6,opt,name=error" json:"error,omitempty"`
	TransactionHash      []byte      `protobuf:"bytes,7,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"`
	ContractAddress      []byte      `protobuf:"bytes,8,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return proto.CompactTextString(m) }
func (*Receipt) ProtoMessage()    {}
func (*Receipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{4}
}
func (m *Receipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Receipt.Unmarshal(m, b)
}
func (m *Receipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Receipt.Marshal(b, m, deterministic)
}
func (dst *Receipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Receipt.Merge(dst, src)
}
func (m *Receipt) XXX_Size() int {
	return xxx_messageInfo_Receipt.Size(m)
}
func (m *Receipt) XXX_DiscardUnknown() {
	xxx_messageInfo_Receipt.DiscardUnknown(m)
}

var xxx_messageInfo_Receipt proto.InternalMessageInfo

func (m *Receipt) GetStateRoot() []byte {
	if m != nil {
		return m.StateRoot
	}
	return nil
}

func (m *Receipt) GetQuotaUsed() uint64 {
	if m != nil {
		return m.QuotaUsed
	}
	return 0
}

func (m *Receipt) GetQuota() uint64 {
	if m != nil {
		return m.Quota
	}
	return 0
}

func (m *Receipt) GetLogBloom() []byte {
	if m != nil {
		return m.LogBloom
	}
	return nil
}

func (m *Receipt) GetLogs() []*LogEntry {
	if m != nil {
		return m.Logs
	}
	return nil
}

func (m *Receipt) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *Receipt) GetTransactionHash() []byte {
	if m != nil {
		return m.TransactionHash
	}
	return nil
}

func (m *Receipt) GetContractAddress() []byte {
	if m != nil {
		return m.ContractAddress
	}
	return nil
}

type LogEntry struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Topics               [][]byte `protobuf:"bytes,2,rep,name=topics,proto3" json:"topics,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}
func (*LogEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{5}
}
func (m *LogEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEntry.Unmarshal(m, b)
}
func (m *LogEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEntry.Marshal(b, m, deterministic)
}
func (dst *LogEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEntry.Merge(dst, src)
}
func (m *LogEntry) XXX_Size() int {
	return xxx_messageInfo_LogEntry.Size(m)
}
func (m *LogEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEntry.DiscardUnknown(m)
}

var xxx_messageInfo_LogEntry proto.InternalMessageInfo

func (m *LogEntry) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *LogEntry) GetTopics() [][]byte {
	if m != nil {
		return m.Topics
	}
	return nil
}

func (m *LogEntry) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type NewBlockReq struct {
	Block                *Block   `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewBlockReq) Reset()         { *m = NewBlockReq{} }
func (m *NewBlockReq) String() string { return proto.CompactTextString(m) }
func (*NewBlockReq) ProtoMessage()    {}
func (*NewBlockReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{6}
}
func (m *NewBlockReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewBlockReq.Unmarshal(m, b)
}
func (m *NewBlockReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewBlockReq.Marshal(b, m, deterministic)
}
func (dst *NewBlockReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewBlockReq.Merge(dst, src)
}
func (m *NewBlockReq) XXX_Size() int {
	return xxx_messageInfo_NewBlockReq.Size(m)
}
func (m *NewBlockReq) XXX_DiscardUnknown() {
	xxx_messageInfo_NewBlockReq.DiscardUnknown(m)
}

var xxx_messageInfo_NewBlockReq proto.InternalMessageInfo

func (m *NewBlockReq) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type NewBlockRes struct {
	Height               uint64   `protobuf:"varint,1,opt,name=height" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewBlockRes) Reset()         { *m = NewBlockRes{} }
func (m *NewBlockRes) String() string { return proto.CompactTextString(m) }
func (*NewBlockRes) ProtoMessage()    {}
func (*NewBlockRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{7}
}
func (m *NewBlockRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewBlockRes.Unmarshal(m, b)
}
func (m *NewBlockRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewBlockRes.Marshal(b, m, deterministic)
}
func (dst *NewBlockRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewBlockRes.Merge(dst, src)
}
func (m *NewBlockRes) XXX_Size() int {
	return xxx_messageInfo_NewBlockRes.Size(m)
}
func (m *NewBlockRes) XXX_DiscardUnknown() {
	xxx_messageInfo_NewBlockRes.DiscardUnknown(m)
}

var xxx_messageInfo_NewBlockRes proto.InternalMessageInfo

func (m *NewBlockRes) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type GetValidatorsReq struct {
	Height               uint64   `protobuf:"varint,1,opt,name=height" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetValidatorsReq) Reset()         { *m = GetValidatorsReq{} }
func (m *GetValidatorsReq) String() string { return proto.CompactTextString(m) }
func (*GetValidatorsReq) ProtoMessage()    {}
func (*GetValidatorsReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{8}
}
func (m *GetValidatorsReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetValidatorsReq.Unmarshal(m, b)
}
func (m *GetValidatorsReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetValidatorsReq.Marshal(b, m, deterministic)
}
func (dst *GetValidatorsReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetValidatorsReq.Merge(dst, src)
}
func (m *GetValidatorsReq) XXX_Size() int {
	return xxx_messageInfo_GetValidatorsReq.Size(m)
}
func (m *GetValidatorsReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetValidatorsReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetValidatorsReq proto.InternalMessageInfo

func (m *GetValidatorsReq) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type GetValidatorsRes struct {
	Vals                 [][]byte `protobuf:"bytes,3,rep,name=vals,proto3" json:"vals,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetValidatorsRes) Reset()         { *m = GetValidatorsRes{} }
func (m *GetValidatorsRes) String() string { return proto.CompactTextString(m) }
func (*GetValidatorsRes) ProtoMessage()    {}
func (*GetValidatorsRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{9}
}
func (m *GetValidatorsRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetValidatorsRes.Unmarshal(m, b)
}
func (m *GetValidatorsRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetValidatorsRes.Marshal(b, m, deterministic)
}
func (dst *GetValidatorsRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetValidatorsRes.Merge(dst, src)
}
func (m *GetValidatorsRes) XXX_Size() int {
	return xxx_messageInfo_GetValidatorsRes.Size(m)
}
func (m *GetValidatorsRes) XXX_DiscardUnknown() {
	xxx_messageInfo_GetValidatorsRes.DiscardUnknown(m)
}

var xxx_messageInfo_GetValidatorsRes proto.InternalMessageInfo

func (m *GetValidatorsRes) GetVals() [][]byte {
	if m != nil {
		return m.Vals
	}
	return nil
}

type GetBlockHeaderReq struct {
	Height               uint64   `protobuf:"varint,1,opt,name=height" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBlockHeaderReq) Reset()         { *m = GetBlockHeaderReq{} }
func (m *GetBlockHeaderReq) String() string { return proto.CompactTextString(m) }
func (*GetBlockHeaderReq) ProtoMessage()    {}
func (*GetBlockHeaderReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{10}
}
func (m *GetBlockHeaderReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetBlockHeaderReq.Unmarshal(m, b)
}
func (m *GetBlockHeaderReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetBlockHeaderReq.Marshal(b, m, deterministic)
}
func (dst *GetBlockHeaderReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetBlockHeaderReq.Merge(dst, src)
}
func (m *GetBlockHeaderReq) XXX_Size() int {
	return xxx_messageInfo_GetBlockHeaderReq.Size(m)
}
func (m *GetBlockHeaderReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetBlockHeaderReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetBlockHeaderReq proto.InternalMessageInfo

func (m *GetBlockHeaderReq) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type GetBlockHeaderRes struct {
	Header               *BlockHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GetBlockHeaderRes) Reset()         { *m = GetBlockHeaderRes{} }
func (m *GetBlockHeaderRes) String() string { return proto.CompactTextString(m) }
func (*GetBlockHeaderRes) ProtoMessage()    {}
func (*GetBlockHeaderRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{11}
}
func (m *GetBlockHeaderRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetBlockHeaderRes.Unmarshal(m, b)
}
func (m *GetBlockHeaderRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetBlockHeaderRes.Marshal(b, m, deterministic)
}
func (dst *GetBlockHeaderRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetBlockHeaderRes.Merge(dst, src)
}
func (m *GetBlockHeaderRes) XXX_Size() int {
	return xxx_messageInfo_GetBlockHeaderRes.Size(m)
}
func (m *GetBlockHeaderRes) XXX_DiscardUnknown() {
	xxx_messageInfo_GetBlockHeaderRes.DiscardUnknown(m)
}

var xxx_messageInfo_GetBlockHeaderRes proto.InternalMessageInfo

func (m *GetBlockHeaderRes) GetHeader() *BlockHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type GetReceiptReq struct {
	TxHash               []byte   `protobuf:"bytes,1,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetReceiptReq) Reset()         { *m = GetReceiptReq{} }
func (m *GetReceiptReq) String() string { return proto.CompactTextString(m) }
func (*GetReceiptReq) ProtoMessage()    {}
func (*GetReceiptReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{12}
}
func (m *GetReceiptReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetReceiptReq.Unmarshal(m, b)
}
func (m *GetReceiptReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetReceiptReq.Marshal(b, m, deterministic)
}
func (dst *GetReceiptReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetReceiptReq.Merge(dst, src)
}
func (m *GetReceiptReq) XXX_Size() int {
	return xxx_messageInfo_GetReceiptReq.Size(m)
}
func (m *GetReceiptReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetReceiptReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetReceiptReq proto.InternalMessageInfo

func (m *GetReceiptReq) GetTxHash() []byte {
	if m != nil {
		return m.TxHash
	}
	return nil
}

type GetReceiptRes struct {
	Receipt              *Receipt `protobuf:"bytes,1,opt,name=receipt" json:"receipt,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetReceiptRes) Reset()         { *m = GetReceiptRes{} }
func (m *GetReceiptRes) String() string { return proto.CompactTextString(m) }
func (*GetReceiptRes) ProtoMessage()    {}
func (*GetReceiptRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_chain_0aaaa1caeb9f0a22, []int{13}
}
func (m *GetReceiptRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetReceiptRes.Unmarshal(m, b)
}
func (m *GetReceiptRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetReceiptRes.Marshal(b, m, deterministic)
}
func (dst *GetReceiptRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetReceiptRes.Merge(dst, src)
}
func (m *GetReceiptRes) XXX_Size() int {
	return xxx_messageInfo_GetReceiptRes.Size(m)
}
func (m *GetReceiptRes) XXX_DiscardUnknown() {
	xxx_messageInfo_GetReceiptRes.DiscardUnknown(m)
}

var xxx_messageInfo_GetReceiptRes proto.InternalMessageInfo

func (m *GetReceiptRes) GetReceipt() *Receipt {
	if m != nil {
		return m.Receipt
	}
	return nil
}

func init() {
	proto.RegisterType((*Proof)(nil), "Proof")
	proto.RegisterType((*Block)(nil), "Block")
	proto.RegisterType((*BlockHeader)(nil), "BlockHeader")
	proto.RegisterType((*BlockBody)(nil), "BlockBody")
	proto.RegisterType((*Receipt)(nil), "Receipt")
	proto.RegisterType((*LogEntry)(nil), "LogEntry")
	proto.RegisterType((*NewBlockReq)(nil), "NewBlockReq")
	proto.RegisterType((*NewBlockRes)(nil), "NewBlockRes")
	proto.RegisterType((*GetValidatorsReq)(nil), "GetValidatorsReq")
	proto.RegisterType((*GetValidatorsRes)(nil), "GetValidatorsRes")
	proto.RegisterType((*GetBlockHeaderReq)(nil), "GetBlockHeaderReq")
	proto.RegisterType((*GetBlockHeaderRes)(nil), "GetBlockHeaderRes")
	proto.RegisterType((*GetReceiptReq)(nil), "GetReceiptReq")
	proto.RegisterType((*GetReceiptRes)(nil), "GetReceiptRes")
	proto.RegisterEnum("ProofType", ProofType_name, ProofType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ChainClient is the client API for Chain service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ChainClient interface {
	NewBlock(ctx context.Context, in *NewBlockReq, opts ...grpc.CallOption) (*Empty, error)
	GetValidators(ctx context.Context, in *GetValidatorsReq, opts ...grpc.CallOption) (*GetValidatorsRes, error)
	GetBlockHeader(ctx context.Context, in *GetBlockHeaderReq, opts ...grpc.CallOption) (*GetBlockHeaderRes, error)
	GetReceipt(ctx context.Context, in *GetReceiptReq, opts ...grpc.CallOption) (*GetReceiptRes, error)
}

type chainClient struct {
	cc *grpc.ClientConn
}

func NewChainClient(cc *grpc.ClientConn) ChainClient {
	return &chainClient{cc}
}

func (c *chainClient) NewBlock(ctx context.Context, in *NewBlockReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/Chain/NewBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chainClient) GetValidators(ctx context.Context, in *GetValidatorsReq, opts ...grpc.CallOption) (*GetValidatorsRes, error) {
	out := new(GetValidatorsRes)
	err := c.cc.Invoke(ctx, "/Chain/GetValidators", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chainClient) GetBlockHeader(ctx context.Context, in *GetBlockHeaderReq, opts ...grpc.CallOption) (*GetBlockHeaderRes, error) {
	out := new(GetBlockHeaderRes)
	err := c.cc.Invoke(ctx, "/Chain/GetBlockHeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chainClient) GetReceipt(ctx context.Context, in *GetReceiptReq, opts ...grpc.CallOption) (*GetReceiptRes, error) {
	out := new(GetReceiptRes)
	err := c.cc.Invoke(ctx, "/Chain/GetReceipt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChainServer is the server API for Chain service.
type ChainServer interface {
	NewBlock(context.Context, *NewBlockReq) (*Empty, error)
	GetValidators(context.Context, *GetValidatorsReq) (*GetValidatorsRes, error)
	GetBlockHeader(context.Context, *GetBlockHeaderReq) (*GetBlockHeaderRes, error)
	GetReceipt(context.Context, *GetReceiptReq) (*GetReceiptRes, error)
}

func RegisterChainServer(s *grpc.Server, srv ChainServer) {
	s.RegisterService(&_Chain_serviceDesc, srv)
}

func _Chain_NewBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewBlockReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainServer).NewBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Chain/NewBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainServer).NewBlock(ctx, req.(*NewBlockReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chain_GetValidators_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetValidatorsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainServer).GetValidators(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Chain/GetValidators",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainServer).GetValidators(ctx, req.(*GetValidatorsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chain_GetBlockHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockHeaderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainServer).GetBlockHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Chain/GetBlockHeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainServer).GetBlockHeader(ctx, req.(*GetBlockHeaderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chain_GetReceipt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReceiptReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChainServer).GetReceipt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Chain/GetReceipt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChainServer).GetReceipt(ctx, req.(*GetReceiptReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Chain_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Chain",
	HandlerType: (*ChainServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewBlock",
			Handler:    _Chain_NewBlock_Handler,
		},
		{
			MethodName: "GetValidators",
			Handler:    _Chain_GetValidators_Handler,
		},
		{
			MethodName: "GetBlockHeader",
			Handler:    _Chain_GetBlockHeader_Handler,
		},
		{
			MethodName: "GetReceipt",
			Handler:    _Chain_GetReceipt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chain.proto",
}

func init() { proto.RegisterFile("chain.proto", fileDescriptor_chain_0aaaa1caeb9f0a22) }

var fileDescriptor_chain_0aaaa1caeb9f0a22 = []byte{
	// 744 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0xcb, 0x6e, 0xeb, 0x36,
	0x10, 0xb5, 0xfc, 0xd6, 0x58, 0xf1, 0x75, 0x88, 0xa2, 0x15, 0xdc, 0x9b, 0xd6, 0x60, 0x1f, 0x70,
	0x6f, 0x00, 0x2d, 0x7c, 0x17, 0x45, 0xba, 0x4b, 0x8a, 0x20, 0x59, 0x04, 0x45, 0x40, 0xa4, 0xdd,
	0x1a, 0xb2, 0xc4, 0xd8, 0x42, 0x25, 0x51, 0x21, 0xe9, 0x34, 0xfe, 0x80, 0xfe, 0x5c, 0x57, 0xfd,
	0xa4, 0x82, 0x43, 0x2a, 0x91, 0x6d, 0xa4, 0x3b, 0x9f, 0x33, 0x87, 0x33, 0xc3, 0x33, 0x23, 0x1a,
	0x46, 0xc9, 0x26, 0xce, 0xca, 0xa8, 0x92, 0x42, 0x8b, 0x69, 0x90, 0x88, 0xa2, 0x10, 0x0e, 0xd1,
	0x4b, 0xe8, 0xdd, 0x4b, 0x21, 0x1e, 0x49, 0x08, 0x83, 0x44, 0x94, 0x9a, 0x97, 0x3a, 0xf4, 0x66,
	0xde, 0x3c, 0x60, 0x35, 0x24, 0xdf, 0x40, 0x57, 0xef, 0x2a, 0x1e, 0xb6, 0x67, 0xde, 0x7c, 0xbc,
	0x80, 0x08, 0xf5, 0x0f, 0xbb, 0x8a, 0x33, 0xe4, 0xe9, 0x1a, 0x7a, 0x57, 0xb9, 0x48, 0xfe, 0x34,
	0x29, 0x9e, 0xb9, 0x54, 0x99, 0x28, 0x31, 0xc5, 0x09, 0xab, 0x21, 0xf9, 0x1e, 0xfa, 0x1b, 0x1e,
	0xa7, 0x5c, 0x62, 0x92, 0xd1, 0x22, 0x88, 0xf0, 0xc4, 0x2d, 0x72, 0xcc, 0xc5, 0x4c, 0xa1, 0x95,
	0x48, 0x77, 0x61, 0x07, 0x35, 0x60, 0x35, 0x57, 0x22, 0xdd, 0x31, 0xe4, 0xe9, 0x3f, 0x6d, 0x18,
	0x35, 0xce, 0x91, 0x29, 0x0c, 0x2b, 0xc9, 0x9f, 0x37, 0xb1, 0xda, 0xb8, 0x9e, 0x5f, 0x31, 0xf9,
	0x08, 0xbe, 0xce, 0x0a, 0xae, 0x74, 0x5c, 0x54, 0x58, 0xb4, 0xcb, 0xde, 0x08, 0xf2, 0xa5, 0xe9,
	0x27, 0x5b, 0x6f, 0x34, 0xd6, 0xea, 0x32, 0x87, 0xc8, 0x19, 0x80, 0xd2, 0xb1, 0xe6, 0x4b, 0x29,
	0x84, 0x0e, 0xbb, 0x98, 0xd3, 0x47, 0x86, 0x09, 0xa1, 0xc9, 0x39, 0x9c, 0x6a, 0x19, 0x97, 0x2a,
	0x4e, 0x74, 0x26, 0x4a, 0x65, 0x55, 0x3d, 0x54, 0x4d, 0x9a, 0x01, 0x14, 0x7f, 0x07, 0x27, 0x92,
	0x27, 0x3c, 0xab, 0xb4, 0x13, 0xf6, 0x51, 0x18, 0xd4, 0x24, 0x8a, 0xce, 0x00, 0x9e, 0xb6, 0x42,
	0xc7, 0xcb, 0xad, 0xe2, 0x69, 0x38, 0xb0, 0x7d, 0x22, 0xf3, 0xbb, 0xe2, 0x29, 0xf9, 0x16, 0x46,
	0x36, 0x9c, 0x67, 0x45, 0xa6, 0xc3, 0x21, 0xc6, 0xed, 0x89, 0x3b, 0xc3, 0x90, 0x8f, 0xd0, 0xab,
	0xcc, 0x38, 0x42, 0x1f, 0x3d, 0xeb, 0xdb, 0xe1, 0x30, 0x4b, 0x5a, 0x83, 0x44, 0x25, 0x14, 0x97,
	0x21, 0xd4, 0x06, 0x59, 0x4c, 0xe7, 0xe0, 0xbf, 0xfa, 0x4b, 0xbe, 0x06, 0x5f, 0xbf, 0x2c, 0x8d,
	0x71, 0x5c, 0x85, 0xde, 0xac, 0x63, 0x94, 0xfa, 0xe5, 0x16, 0x31, 0xfd, 0xbb, 0x0d, 0x03, 0x66,
	0x9b, 0x3e, 0x30, 0xc8, 0x3b, 0x34, 0x68, 0xff, 0x3a, 0xed, 0xc3, 0xeb, 0x7c, 0x01, 0x3d, 0x04,
	0xce, 0x75, 0x0b, 0x4c, 0xf1, 0x5c, 0xac, 0x97, 0xab, 0x5c, 0x88, 0xc2, 0x79, 0x3e, 0xcc, 0xc5,
	0xfa, 0xca, 0x60, 0x72, 0x06, 0xdd, 0x5c, 0xac, 0x55, 0xd8, 0x9b, 0x75, 0xe6, 0xa3, 0x85, 0x1f,
	0xdd, 0x89, 0xf5, 0x75, 0xa9, 0xe5, 0x8e, 0x21, 0x6d, 0x32, 0x72, 0x29, 0x85, 0x44, 0x73, 0x7d,
	0x66, 0x01, 0xf9, 0x09, 0x9a, 0xe3, 0xc0, 0x7b, 0xa1, 0xb7, 0x01, 0xfb, 0xd0, 0xe0, 0xcd, 0xf5,
	0x8c, 0xd4, 0xec, 0xb9, 0x8c, 0x13, 0xbd, 0x8c, 0xd3, 0x54, 0x72, 0xa5, 0xd0, 0xe6, 0x80, 0x7d,
	0xa8, 0xf9, 0x4b, 0x4b, 0xd3, 0x7b, 0x18, 0xd6, 0xd5, 0xcd, 0xaa, 0xd7, 0x6a, 0xf7, 0xb5, 0x38,
	0x68, 0x56, 0x4b, 0x8b, 0x2a, 0x4b, 0x54, 0xd8, 0x46, 0x1f, 0x1d, 0x22, 0x04, 0xba, 0x69, 0xec,
	0xae, 0x1e, 0x30, 0xfc, 0x4d, 0xcf, 0x61, 0xf4, 0x1b, 0xff, 0x0b, 0xc7, 0xc0, 0xf8, 0x93, 0x19,
	0xe6, 0xca, 0xfc, 0xc6, 0x94, 0x66, 0x98, 0x36, 0x62, 0x49, 0xfa, 0x43, 0x53, 0xac, 0x1a, 0x2b,
	0xec, 0x35, 0x57, 0x98, 0x7e, 0x82, 0xc9, 0x0d, 0xd7, 0x7f, 0xc4, 0x79, 0x96, 0xc6, 0x5a, 0x48,
	0x65, 0x12, 0xbf, 0xa7, 0xfd, 0xf1, 0x48, 0x8b, 0x7d, 0x3e, 0xc7, 0xb9, 0x0a, 0x3b, 0xd8, 0x3d,
	0xfe, 0xa6, 0xe7, 0x70, 0x7a, 0xc3, 0x75, 0xf3, 0x93, 0xfd, 0x9f, 0xa4, 0x17, 0xc7, 0x62, 0xd5,
	0x78, 0x00, 0xbc, 0xf7, 0x1f, 0x00, 0x3a, 0x87, 0x93, 0x1b, 0xae, 0xdd, 0xae, 0x99, 0x1a, 0x5f,
	0xc1, 0xc0, 0xed, 0xa5, 0xb3, 0xb9, 0x6f, 0xb7, 0x92, 0x7e, 0xde, 0x57, 0x2a, 0x42, 0x61, 0xe0,
	0x3e, 0x2c, 0x57, 0x61, 0x18, 0xd5, 0xd1, 0x3a, 0xf0, 0xe9, 0x02, 0xfc, 0xd7, 0xb7, 0x8b, 0x10,
	0x18, 0x5f, 0x6e, 0xf5, 0x46, 0xc8, 0x4c, 0xef, 0x98, 0xd8, 0x96, 0xe9, 0xa4, 0x45, 0x86, 0xd0,
	0x65, 0xf1, 0xa3, 0x9e, 0x78, 0x64, 0x0c, 0xf0, 0xc0, 0xcb, 0x94, 0xcb, 0x22, 0x2b, 0xf5, 0xa4,
	0xbd, 0xf8, 0xd7, 0x83, 0xde, 0xaf, 0xe6, 0x11, 0x25, 0x14, 0x86, 0xf5, 0x18, 0x48, 0x10, 0x35,
	0xc6, 0x37, 0xed, 0x47, 0xd7, 0x45, 0xa5, 0x77, 0xb4, 0x45, 0x7e, 0xc6, 0xee, 0xde, 0x7c, 0x25,
	0xa7, 0xd1, 0xe1, 0x4c, 0xa6, 0x47, 0x94, 0xa2, 0x2d, 0xf2, 0x0b, 0x8c, 0xf7, 0xbd, 0x23, 0x24,
	0x3a, 0x72, 0x7e, 0x7a, 0xcc, 0x99, 0xb3, 0x11, 0xc0, 0x9b, 0x25, 0x64, 0x1c, 0xed, 0x39, 0x39,
	0xdd, 0xc7, 0x8a, 0xb6, 0x56, 0x7d, 0xfc, 0x03, 0xf8, 0xfc, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x9e, 0x35, 0x8c, 0x5e, 0x1d, 0x06, 0x00, 0x00,
}
