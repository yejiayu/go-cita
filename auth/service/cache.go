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
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gomodule/redigo/redis"
)

func newCache() (*cache, error) {
	client, err := redis.DialURL("redis://127.0.0.1:6379")
	if err != nil {
		return nil, err
	}

	return &cache{
		client: client,
	}, nil
}

type cache struct {
	client redis.Conn
}

func (c *cache) getPublicKey(hash common.Hash) (*ecdsa.PublicKey, error) {
	reply, err := c.client.Do("HGET", "tx.pk", hash.String())
	if err != nil {
		return nil, err
	}

	if reply == nil {
		return nil, nil
	}

	data := reply.([]byte)
	return ethcrypto.DecompressPubkey(data)
}

func (c *cache) setPublicKey(hash common.Hash, pk *ecdsa.PublicKey) error {
	_, err := c.client.Do("HSET", "tx.pk", hash.String(), ethcrypto.CompressPubkey(pk))
	return err
}

func (c *cache) setHistoryTx(height uint64, txHashes []common.Hash) error {
	keys := []string{}
	for _, hashes := range txHashes {
		keys = append(keys, hashes.Hex())
	}

	_, err := c.client.Do("SADD", fmt.Sprintf("history.txs.%d", height), keys)
	return err
}
