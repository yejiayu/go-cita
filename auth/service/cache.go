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
		if err == redis.ErrNil {
			return nil, nil
		}

		return nil, err
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

func (c *cache) historyTxExists(height uint64, txHash common.Hash) (bool, error) {
	_, err := c.client.Do("SMEMBERS", fmt.Sprintf("history.txs.%d", height), txHash.Hex())
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func (c *cache) deleteHistoryTx(height uint64) error {
	_, err := c.client.Do("DEL", fmt.Sprintf("history.txs.%d", height))
	if err != nil {
		if err == redis.ErrNil {
			return nil
		}

		return err
	}

	return nil
}
