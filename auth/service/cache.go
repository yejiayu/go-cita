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
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gomodule/redigo/redis"
	ot "github.com/opentracing/opentracing-go"
)

func newCache(redisURL string) (*cache, error) {
	client, err := redis.DialURL("redis://" + redisURL)
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

func (c *cache) getPublicKey(ctx context.Context, hash common.Hash) (*ecdsa.PublicKey, error) {
	span, _ := ot.StartSpanFromContext(ctx, "publick-key-cache")
	defer span.Finish()
	span.SetTag("method", "HGET")
	span.SetTag("tx_hash", hash.String())

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

func (c *cache) setPublicKey(ctx context.Context, hash common.Hash, pk *ecdsa.PublicKey) error {
	span, _ := ot.StartSpanFromContext(ctx, "publick-key-cache")
	defer span.Finish()
	span.SetTag("method", "HSET")
	span.SetTag("tx_hash", hash.String())

	_, err := c.client.Do("HSET", "tx.pk", hash.String(), ethcrypto.CompressPubkey(pk))
	return err
}
