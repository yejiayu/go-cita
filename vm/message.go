package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func NewMessage(
	from common.Address, to *common.Address,
	data []byte, quota uint64, value *big.Int,
	txHash []byte,
) *Message {
	msg := &Message{
		from:       from,
		nonce:      1,
		amount:     value,
		quota:      quota,
		quotaPrice: big.NewInt(1),
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
	quota      uint64
	quotaPrice *big.Int
	data       []byte
	checkNonce bool
	txHash     []byte
}

func (m Message) From() common.Address { return m.from }
func (m Message) To() *common.Address  { return m.to }
func (m Message) QuotaPrice() *big.Int { return m.quotaPrice }
func (m Message) Value() *big.Int      { return m.amount }
func (m Message) Quota() uint64        { return m.quota }
func (m Message) Nonce() uint64        { return m.nonce }
func (m Message) Data() []byte         { return m.data }
func (m Message) CheckNonce() bool     { return m.checkNonce }
func (m Message) TxHash() []byte       { return m.txHash }
