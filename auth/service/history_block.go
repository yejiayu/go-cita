package service

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yejiayu/go-cita/types"
)

type historyTx struct {
	maxHeight uint64
	mu        sync.RWMutex
	txMap     map[uint64]map[common.Hash]bool
}

func newHistoryTx() *historyTx {
	return &historyTx{
		txMap: make(map[uint64]map[common.Hash]bool),
	}
}

func (ht *historyTx) AddBlock(b *types.Block) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	txs, ok := ht.txMap[b.Header.Height]
	if ok {
		return
	}

	txs = map[common.Hash]bool{}
	for _, signtx := range b.Body.Transactions {
		hash := signtx.GetTxHash()
		txs[common.BytesToHash(hash)] = true
	}

	ht.txMap[b.Header.Height] = txs

	if len(ht.txMap) > 100 {
		for height := range ht.txMap {
			if height < b.Header.Height-100 {
				delete(ht.txMap, height)
			}
		}
	}

	ht.maxHeight = b.Header.Height
}

func (ht *historyTx) Contains(txHash common.Hash) bool {
	for _, txMap := range ht.txMap {
		_, ok := txMap[txHash]
		if ok {
			return true
		}
	}

	return false
}
