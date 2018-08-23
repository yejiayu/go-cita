package quota

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gomodule/redigo/redis"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/database/tx"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/auth/pool"
)

const (
	poolKey     = "tx.quota.pool"
	poolHashKey = "tx.quota.hash.pool"
	poolTempKey = "tx.quota.temp.pool"
)

func New(txDB tx.Interface, url string, poolCount uint32) pool.Interface {
	connPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   40,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				url,
				redis.DialConnectTimeout(time.Second*1),
				redis.DialReadTimeout(time.Second*3),
				redis.DialWriteTimeout(time.Second*3),
			)
		},
	}
	return &quotaPool{
		txDB:      txDB,
		connPool:  connPool,
		poolCount: poolCount,
	}
}

type quotaPool struct {
	txDB      tx.Interface
	connPool  *redis.Pool
	poolCount uint32
}

func (p *quotaPool) Add(ctx context.Context, signedTx *pb.SignedTransaction) error {
	len, err := p.Len(ctx)
	if err != nil {
		return err
	}

	if len >= int64(p.poolCount) {
		return pool.ErrPoolReachLimit
	}

	txHash := hash.BytesToHash(signedTx.GetTxHash())
	exists, err := p.exists(ctx, txHash)
	if err != nil {
		return err
	}
	if exists {
		return pool.ErrTxDup
	}

	data, err := proto.Marshal(signedTx)
	if err != nil {
		return err
	}

	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	quota := signedTx.GetTransactionWithSig().GetTransaction().GetQuota()
	conn.Send("HSET", poolHashKey, txHash.String(), data)
	conn.Send("ZADD", poolKey, quota, txHash.String())
	conn.Flush()
	conn.Receive()
	conn.Receive()

	return nil
}

func (p *quotaPool) Pull(ctx context.Context, count uint32) ([]*pb.SignedTransaction, error) {
	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	keys := []interface{}{poolKey, 0, count - 1}
	replyKeys, err := redis.Strings(conn.Do("ZREVRANGE", keys...))
	if err != nil {
		return nil, err
	}

	keys = make([]interface{}, len(replyKeys)+1)
	keys[0] = poolHashKey
	for i, key := range replyKeys {
		keys[i+1] = key
	}

	replyValues, err := redis.ByteSlices(conn.Do("HMGET", keys...))
	if err != nil {
		return nil, err
	}

	txs := make([]*pb.SignedTransaction, len(replyValues))
	for i, data := range replyValues {
		var signeTx pb.SignedTransaction
		if err := proto.Unmarshal(data, &signeTx); err != nil {
			log.Error(err)
			continue
		}

		txs[i] = &signeTx
	}

	return txs, nil
}

func (p *quotaPool) Get(ctx context.Context, hashes []hash.Hash) ([]*pb.SignedTransaction, error) {
	keys := make([]interface{}, len(hashes)+1)
	keys[0] = poolHashKey
	for i, hash := range hashes {
		keys[i+1] = hash.String()
	}

	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := redis.ByteSlices(conn.Do("HMGET", keys))
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, nil
	}

	txs := make([]*pb.SignedTransaction, len(reply))
	for i, data := range reply {
		if data == nil {
			txs[i] = nil
			continue
		}

		var signedTx pb.SignedTransaction
		if err := proto.Unmarshal(data, &signedTx); err != nil {
			log.Error(err)
			continue
		}

		txs[i] = &signedTx
	}

	return txs, nil
}

func (p *quotaPool) Flush(ctx context.Context, hashes []hash.Hash) (int64, error) {
	listKeys := make([]interface{}, len(hashes)+1)
	hashKeys := make([]interface{}, len(hashes)+1)
	listKeys[0] = poolKey
	hashKeys[0] = poolHashKey

	for i, hash := range hashes {
		key := hash.String()
		listKeys[i+1] = key
		hashKeys[i+1] = key
	}

	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	conn.Send("DEL", poolTempKey)
	conn.Send("ZREM", listKeys...)
	conn.Send("HDEL", hashKeys...)
	conn.Flush()
	conn.Receive()
	conn.Receive()
	return redis.Int64(conn.Receive())
}

func (p *quotaPool) Del(ctx context.Context, hashes []hash.Hash) (int64, error) {
	listKeys := make([]interface{}, len(hashes)+1)
	hashKeys := make([]interface{}, len(hashes)+1)
	listKeys[0] = poolKey
	hashKeys[0] = poolHashKey

	for i, hash := range hashes {
		key := hash.String()
		listKeys[i+1] = key
		hashKeys[i+1] = key
	}

	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	conn.Send("ZREM", listKeys...)
	conn.Send("HDEL", hashKeys...)
	conn.Flush()
	conn.Receive()
	return redis.Int64(conn.Receive())
}

func (p *quotaPool) Len(ctx context.Context) (int64, error) {
	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("HLEN", poolHashKey))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

func (p *quotaPool) exists(ctx context.Context, hash hash.Hash) (bool, error) {
	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("HEXISTS", poolHashKey, hash.String()))
	if err != nil && err != redis.ErrNil {
		return false, err
	}
	if exists {
		return true, nil
	}

	return p.txDB.Exists(ctx, hash)
}

func (p *quotaPool) copyToTemp(ctx context.Context, hashes []string) error {
	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	keys := make([]interface{}, len(hashes))
	keys[0] = poolTempKey
	for i, hash := range hashes {
		keys[i+1] = hash
	}

	_, err = conn.Do("SADD", keys...)
	return err
}
