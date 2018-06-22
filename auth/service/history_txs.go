// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

func (ht *historyTx) AddTxs(height uint64, txs []*types.SignedTransaction) error {
	if ht.maxHeight > height {
		return nil
	}

	hashes := []common.Hash{}
	for _, signtx := range txs {
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
