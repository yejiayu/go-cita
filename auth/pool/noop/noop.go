package noop

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/database/tx"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/auth/pool"
)

const poolKey = "tx.noop.pool"

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
	return &noopPool{
		txDB:      txDB,
		connPool:  connPool,
		poolCount: poolCount,
	}
}

type noopPool struct {
	txDB      tx.Interface
	connPool  *redis.Pool
	poolCount uint32
}

func (p *noopPool) Add(ctx context.Context, signedTx *pb.SignedTransaction) error {
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

	_, err = redis.Int64(conn.Do("HSET", poolKey, txHash.String(), data))
	return err
}

func (p *noopPool) Pull(ctx context.Context, count uint32) ([]*pb.SignedTransaction, error) {
	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := redis.ByteSlices(conn.Do("HVALS", poolKey))
	if err != nil {
		return nil, err
	}

	if len(reply) > int(count) {
		reply = reply[:count]
	}

	txs := make([]*pb.SignedTransaction, len(reply))
	for i, data := range reply {
		var signeTx pb.SignedTransaction
		if err := proto.Unmarshal(data, &signeTx); err != nil {
			log.Error(err)
			continue
		}

		txs[i] = &signeTx
	}

	return txs, nil
}

func (p *noopPool) Get(ctx context.Context, hashes []hash.Hash) ([]*pb.SignedTransaction, error) {
	keys := make([]interface{}, len(hashes)+1)
	keys[0] = poolKey
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

func (p *noopPool) Del(ctx context.Context, hashes []hash.Hash) (int64, error) {
	keys := make([]interface{}, len(hashes)+1)
	keys[0] = poolKey
	for i, hash := range hashes {
		keys[i+1] = hash.String()
	}

	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return redis.Int64(conn.Do("HDEL", keys...))
}

func (p *noopPool) Flush(ctx context.Context, hashes []hash.Hash) (int64, error) {
	keys := make([]interface{}, len(hashes)+1)
	keys[0] = poolKey
	for i, hash := range hashes {
		keys[i+1] = hash.String()
	}

	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return redis.Int64(conn.Do("HDEL", keys...))
}

func (p *noopPool) Len(ctx context.Context) (int64, error) {
	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("HLEN", poolKey))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

func (p *noopPool) exists(ctx context.Context, hash hash.Hash) (bool, error) {
	conn, err := p.connPool.GetContext(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("HEXISTS", poolKey, hash.String()))
	if err != nil && err != redis.ErrNil {
		return false, errors.WithStack(err)
	}
	if exists {
		return true, nil
	}

	return p.txDB.Exists(ctx, hash)
}
