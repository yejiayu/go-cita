package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/yejiayu/go-cita/database/raw"
)

func New(urls []string) raw.Interface {
	connPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   20,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				urls[0],
				redis.DialConnectTimeout(time.Second*1),
				redis.DialReadTimeout(time.Second*3),
				redis.DialWriteTimeout(time.Second*3),
			)
		},
	}
	return &db{connPool: connPool}
}

type db struct {
	connPool *redis.Pool
}

func (db *db) Put(ctx context.Context, namespace string, key, value []byte) error {
	conn, err := db.connPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("HSET", namespace, key, value)
	return err
}

func (db *db) Get(ctx context.Context, namespace string, key []byte) ([]byte, error) {
	conn, err := db.connPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("HGET", namespace, key))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	return data, nil
}

func (db *db) Delete(ctx context.Context, namespace string, key []byte) error {
	conn, err := db.connPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("HDEL", namespace, key)
	return err
}

func (db *db) Scan(ctx context.Context, namespace string, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return nil, nil, nil
}
