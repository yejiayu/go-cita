package types

type ProofType uint8

const (
	ProofAuthorityRound ProofType = 0
	ProofRaft           ProofType = 1
	ProofTendermint     ProofType = 2
)

type Proof struct {
	Content      []byte
	FieldType    ProofType
	UnKonwFields UnknownFields
	CachedSize   uint32
}

type UnknownFields struct {
	Fields map[uint32]UnknownValues
}

type UnknownValues struct {
	Fixed32         uint32
	Fixed64         uint64
	Varint          uint64
	LengthDelimited [][]byte
}
