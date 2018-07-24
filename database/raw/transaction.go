package raw

import (
	"context"

	"github.com/pingcap/tidb/kv"
)

type Transaction interface {
	Commit() error
	Rollback() error

	Delete(tableName, key []byte) error
	Put(tableName, key, data []byte) error
}

func beginTransaction(ctx context.Context, tx kv.Transaction) Transaction {
	return &transaction{
		ctx: ctx,
		tx:  tx,
	}
}

type transaction struct {
	tx  kv.Transaction
	ctx context.Context
}

func (t *transaction) Commit() error {
	return t.tx.Commit(t.ctx)
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *transaction) Delete(tableName, key []byte) error {
	key = kv.Key(buildKey(tableName, key))
	return t.tx.Delete(key)
}

func (t *transaction) Put(tableName, key, data []byte) error {
	key = kv.Key(buildKey(tableName, key))
	return t.tx.Set(key, data)
}
