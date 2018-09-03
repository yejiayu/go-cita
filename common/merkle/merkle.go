package merkle

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

type DerivableList interface {
	Len() int
	GetRlp(i int) []byte
}

type transactions struct {
	txHashes []hash.Hash
}

func (list *transactions) Len() int {
	return len(list.txHashes)
}

func (list *transactions) GetRlp(i int) []byte {
	txHash := list.txHashes[i]
	enc, _ := rlp.EncodeToBytes(txHash)
	return enc
}

func TransactionsToRoot(txHashes []hash.Hash) hash.Hash {
	if len(txHashes) == 0 {
		return hash.Hash{}
	}

	list := &transactions{txHashes: txHashes}
	root := types.DeriveSha(list)
	return hash.BytesToHash(root.Bytes())
}

type receipts struct {
	receipts []*pb.Receipt
}

func (list *receipts) Len() int {
	return len(list.receipts)
}

func (list *receipts) GetRlp(i int) []byte {
	receipt := list.receipts[i]
	enc, _ := rlp.EncodeToBytes(receipt)
	return enc
}

func ReceiptsToRoot(rs []*pb.Receipt) hash.Hash {
	if len(rs) == 0 {
		return hash.Hash{}
	}

	list := &receipts{receipts: rs}
	root := types.DeriveSha(list)
	return hash.BytesToHash(root.Bytes())
}
