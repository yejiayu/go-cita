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

package cache

import (
	"context"
	"crypto/ecdsa"

	"github.com/gomodule/redigo/redis"
	ot "github.com/opentracing/opentracing-go"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/common/hash"
)

type Interface interface {
	SetPublicKey(ctx context.Context, hash hash.Hash, pk *ecdsa.PublicKey) error
	GetPublicKey(ctx context.Context, hash hash.Hash) (*ecdsa.PublicKey, error)
}

func New(conn redis.Conn) Interface {
	return &cache{conn: conn}
}

type cache struct {
	conn redis.Conn
}

func (c *cache) GetPublicKey(ctx context.Context, hash hash.Hash) (*ecdsa.PublicKey, error) {
	span, _ := ot.StartSpanFromContext(ctx, "publick-key-cache")
	defer span.Finish()
	span.SetTag("method", "HGET")
	span.SetTag("tx_hash", hash.String())

	reply, err := redis.Bytes(c.conn.Do("HGET", "tx.pk", hash.String()))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	if reply == nil {
		return nil, nil
	}

	return crypto.DecompressPubkey(reply)
}

func (c *cache) SetPublicKey(ctx context.Context, hash hash.Hash, pk *ecdsa.PublicKey) error {
	span, _ := ot.StartSpanFromContext(ctx, "publick-key-cache")
	defer span.Finish()
	span.SetTag("method", "HSET")
	span.SetTag("tx_hash", hash.String())

	_, err := c.conn.Do("HSET", "tx.pk", hash.String(), crypto.CompressPubkey(pk))
	return err
}
