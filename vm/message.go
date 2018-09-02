package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func NewMessage(
	from common.Address, to *common.Address,
	data []byte, gasLimit uint64, value *big.Int,
	txHash []byte,
) *Message {
	msg := &Message{
		from:       from,
		nonce:      1,
		amount:     value,
		gasLimit:   gasLimit,
		gasPrice:   big.NewInt(0),
		data:       data,
		checkNonce: false,
		txHash:     txHash,
	}

	if to != nil {
		msg.to = to
	}

	return msg
}

type Message struct {
	to         *common.Address
	from       common.Address
	nonce      uint64
	amount     *big.Int
	gasLimit   uint64
	gasPrice   *big.Int
	data       []byte
	checkNonce bool
	txHash     []byte
}

func (m Message) From() common.Address { return m.from }
func (m Message) To() *common.Address  { return m.to }
func (m Message) GasPrice() *big.Int   { return m.gasPrice }
func (m Message) Value() *big.Int      { return m.amount }
func (m Message) Gas() uint64          { return m.gasLimit }
func (m Message) Nonce() uint64        { return m.nonce }
func (m Message) Data() []byte         { return m.data }
func (m Message) CheckNonce() bool     { return m.checkNonce }
func (m Message) TxHash() []byte       { return m.txHash }
