package redis

import (
	"context"
	"errors"
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

func (db *db) BatchGet(ctx context.Context, namespace string, keys [][]byte) ([][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	conn, err := db.connPool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	redisKeys := make([]interface{}, len(keys)+1)
	redisKeys[0] = namespace
	for i, key := range keys {
		redisKeys[i+1] = key
	}
	return redis.ByteSlices(conn.Do("HMGET", redisKeys...))
}

func (db *db) BatchPut(ctx context.Context, namespace string, keys, values [][]byte) error {
	if len(keys) == 0 || len(values) == 0 {
		return nil
	}
	if len(keys) != len(values) {
		return errors.New("batch put The length of the key must be the same as the value")
	}

	conn, err := db.connPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	lenth := len(keys)
	redisKeys := make([]interface{}, lenth*2+1)
	redisKeys[0] = namespace

	for i := 1; i <= lenth; i++ {
		redisKeys[2*i-1] = keys[i-1]
		redisKeys[2*i] = values[i-1]
	}

	_, err = conn.Do("HMSET", redisKeys...)
	return err
}

func (db *db) BatchDelete(ctx context.Context, namespace string, keys [][]byte) error {
	conn, err := db.connPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	redisKeys := make([]interface{}, len(keys)+1)
	redisKeys[0] = namespace
	for i, key := range keys {
		redisKeys[i+1] = key
	}
	_, err = conn.Do("HDEL", redisKeys...)
	return err
}

func (db *db) Scan(ctx context.Context, namespace string, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return nil, nil, nil
}

func (db *db) Close() error {
	return db.connPool.Close()
}
