package service

import (
	"crypto/ecdsa"

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
	// publicKey map[common.Hash]*ecdsa.PublicKey
	client redis.Conn
}

func (c *cache) getPublicKey(hash common.Hash) (*ecdsa.PublicKey, error) {
	reply, err := c.client.Do("HGET", hash.String())
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}

		return nil, err
	}

	data := reply.([]byte)
	return ethcrypto.DecompressPubkey(data)
}

func (c *cache) setPublicKey(hash common.Hash, pk *ecdsa.PublicKey) error {
	_, err := c.client.Do("HSET", hash.String(), ethcrypto.CompressPubkey(pk))
	return err
}
