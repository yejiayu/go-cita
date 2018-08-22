package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"

	"github.com/yejiayu/go-cita/database/raw"
)

func New(urls []string) (raw.Interface, error) {
	conn, err := redis.DialURL("redis://" + urls[0])
	if err != nil {
		return nil, err
	}

	return &db{conn: conn}, nil
}

type db struct {
	conn redis.Conn
}

func (db *db) Put(ctx context.Context, namespace string, key, value []byte) error {
	_, err := db.conn.Do("HSET", namespace, key, value)
	return err
}

func (db *db) Get(ctx context.Context, namespace string, key []byte) ([]byte, error) {
	data, err := redis.Bytes(db.conn.Do("HGET", namespace, key))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	return data, nil
}

func (db *db) Delete(ctx context.Context, namespace string, key []byte) error {
	_, err := db.conn.Do("HDEL", namespace, key)
	return err
}

func (db *db) Scan(ctx context.Context, namespace string, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return nil, nil, nil
}
