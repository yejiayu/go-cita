package vm

import (
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/yejiayu/go-cita/pb"
)

func NewMessage(gasLimit uint64, signedTx *pb.SignedTransaction) *Message {
	tx := signedTx.GetTransactionWithSig().GetTransaction()

	nonce, err := strconv.Atoi(tx.GetNonce())
	if err != nil {
		log.Fatal(err)
	}
	amount := &big.Int{}
	amount.SetBytes(tx.GetValue())

	msg := &Message{
		from:       common.BytesToAddress(signedTx.GetSigner()),
		nonce:      uint64(nonce),
		amount:     amount,
		gasLimit:   gasLimit,
		gasPrice:   big.NewInt(0),
		data:       tx.GetData(),
		checkNonce: false,
		txHash:     signedTx.GetTxHash(),
	}

	if tx.GetTo() != "" {
		to := common.HexToAddress(tx.GetTo())
		msg.to = &to
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
