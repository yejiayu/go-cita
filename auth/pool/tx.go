package pool

import (
	"github.com/yejiayu/go-cita/types"
)

type TxPool interface {
	Add(signedTx *types.SignedTransaction) error
}

func NewTxPool() TxPool {
	return &txPool{}
}

type txPool struct{}

func (pool *txPool) Add(signedTx *types.SignedTransaction) error {
	return nil
}
