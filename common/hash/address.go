package hash

import "github.com/ethereum/go-ethereum/common"

// Lengths of addresses in bytes.
const (
	AddressLength = 20
)

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func (a Address) Bytes() []byte {
	return a[:]
}

// SetBytes Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}

	copy(a[AddressLength-len(b):], b)
}

func (a Address) String() string {
	return common.ToHex(a.Bytes())
}
