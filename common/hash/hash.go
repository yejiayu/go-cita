package hash

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/golang/protobuf/proto"
)

// Lengths of hashes in bytes.
const (
	HashLength = 32
)

// Hash represents the 32 byte sha3-256 hash of arbitrary data.
type Hash [HashLength]byte

func ToHex(data []byte) string {
	return common.ToHex(data)
}

func FromHex(data string) []byte {
	return common.FromHex(data)
}

func HashesToBytesS(hashes []Hash) [][]byte {
	bytesS := make([][]byte, len(hashes))
	for i, h := range hashes {
		bytesS[i] = h.Bytes()
	}

	return bytesS
}

func BytesSToHashes(bytesS [][]byte) []Hash {
	hashes := make([]Hash, len(bytesS))
	for i, b := range bytesS {
		hashes[i] = BytesToHash(b)
	}

	return hashes
}

func IsZeroHash(h Hash) bool {
	return bytes.Equal(h.Bytes(), Hash{}.Bytes())
}

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BytesToSha3 creates a new SHA3-256 hash.
func BytesToSha3(data []byte) Hash {
	h := sha3.New256()
	h.Write(data)
	return BytesToHash(h.Sum(nil))
}

func ProtoToSha3(in proto.Message) (Hash, error) {
	data, err := proto.Marshal(in)
	if err != nil {
		return Hash{}, nil
	}

	return BytesToSha3(data), nil
}

// SetBytes Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) String() string {
	return common.ToHex(h.Bytes())
}
