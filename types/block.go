package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Block is encoded as it is on the block chain.
type Block struct {
	Header *BlockHeader
	Body   *BlockBody
}

type BlockHeader struct {
	ParentHash   common.Hash
	Timestamp    uint64
	Number       uint64
	TxRoot       common.Hash
	StateRoot    common.Hash
	ReceiptsRoot common.Hash
	LogBloom     LogBloom
	GasUsed      *big.Int
	GasLimit     *big.Int
	Hash         []common.Hash
	Proof        Proof
	Hashes       []common.Hash
	Version      uint32
	Proposer     common.Address
}

type BlockBody struct {
	Transactions []*SignedTransaction
}

// LogBoom is A record of execution for a `LOG` operation.
type LogBloom struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}
