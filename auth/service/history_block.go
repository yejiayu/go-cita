package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/yejiayu/go-cita/types"
)

type historyTx struct {
	maxHeight uint64
	cache     *cache
}

func newHistoryTx(cache *cache) *historyTx {
	return &historyTx{
		cache: cache,
	}
}

func (ht *historyTx) AddBlock(b *types.Block) error {
	height := b.Header.Height
	if ht.maxHeight > height {
		return nil
	}

	hashes := []common.Hash{}
	for _, signtx := range b.Body.Transactions {
		hashes = append(hashes, common.BytesToHash(signtx.GetTxHash()))
		if len(hashes) > 100 {
			if err := ht.cache.setHistoryTx(height, hashes); err != nil {
				return err
			}
		}
	}

	if len(hashes) > 0 {
		if err := ht.cache.setHistoryTx(height, hashes); err != nil {
			return err
		}
	}

	ht.maxHeight = height
	if ht.maxHeight > 100 {
		return ht.cache.deleteHistoryTx(ht.maxHeight - 100)
	}
	return nil
}

func (ht *historyTx) Contains(txHash common.Hash) (bool, error) {
	var count uint64
	for count < 100 {
		height := ht.maxHeight - count
		if height < 0 {
			return false, nil
		}

		exists, err := ht.cache.historyTxExists(uint64(height), txHash)
		if err != nil {
			return false, err
		}
		if exists {
			return true, nil
		}
	}

	return false, nil
}
