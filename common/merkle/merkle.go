package merkle

import (
	"reflect"

	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/rlp"
)

type Content struct {
	In []byte
}

func (c Content) CalculateHash() ([]byte, error) {
	return rlp.EncodeToBytes(c.In)
}

func (c Content) Equals(other merkletree.Content) (bool, error) {
	return reflect.DeepEqual(c.In, other.(Content).In), nil
}

func New(data [][]byte) (*merkletree.MerkleTree, error) {
	var list []merkletree.Content
	for _, d := range data {
		list = append(list, Content{In: d})
	}

	return merkletree.NewTree(list)
}
