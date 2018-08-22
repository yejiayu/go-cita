package pool

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/database/tx"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"
)

func NewNoop(txDB tx.Interface, conn redis.Conn, poolCount uint32) Interface {
	return &noopPool{
		txDB:      txDB,
		conn:      conn,
		poolCount: poolCount,
	}
}

type noopPool struct {
	txDB      tx.Interface
	conn      redis.Conn
	poolCount uint32
}

func (p *noopPool) Add(ctx context.Context, signedTx *pb.SignedTransaction) error {
	count, err := p.Len(ctx)
	if err != nil {
		return err
	}

	if count >= int64(p.poolCount) {
		return errPoolReachLimit
	}

	txHash := hash.BytesToHash(signedTx.GetTxHash())
	data, err := proto.Marshal(signedTx)
	if err != nil {
		return err
	}

	// first, add to tx DB.
	exists, err := p.txDB.Exists(ctx, txHash)
	if err != nil {
		return err
	}
	if exists {
		return errTxDup
	}

	if err = p.txDB.Add(ctx, signedTx); err != nil {
		return err
	}

	// then, add to pool
	reply, err := redis.Int64(p.conn.Do("HSET", poolKey, txHash.String(), data))
	if err != nil {
		return err
	}

	if reply == 0 {
		return errTxDup
	}
	return nil
}

func (p *noopPool) GetAll(ctx context.Context, count uint32) ([]*pb.SignedTransaction, error) {
	reply, err := redis.ByteSlices(p.conn.Do("HVALS", poolKey))
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

	reply, err := redis.ByteSlices(p.conn.Do("HMGET", keys))
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
	count, err := p.del(hashes)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *noopPool) del(hashes []hash.Hash) (int64, error) {
	keys := make([]interface{}, len(hashes)+1)
	keys[0] = poolKey
	for i, hash := range hashes {
		keys[i+1] = hash.String()
	}

	return redis.Int64(p.conn.Do("HDEL", keys...))
}

func (p *noopPool) Len(ctx context.Context) (int64, error) {
	reply, err := redis.Int64(p.conn.Do("HLEN", poolKey))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

func (p *noopPool) Exists(ctx context.Context, hash hash.Hash) (bool, error) {
	exists, err := redis.Bool(p.conn.Do("HEXISTS", poolKey, hash.String()))
	log.Error(err)
	log.Info(exists)
	log.Info(redis.ErrNil == err)
	if err != nil && err != redis.ErrNil {
		return false, errors.WithStack(err)
	}

	return exists, nil
}
