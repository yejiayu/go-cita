package tikv

import (
	"context"
	"fmt"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"

	"github.com/yejiayu/go-cita/database/raw"
)

func New(urls []string) (raw.Interface, error) {
	client, err := tikv.NewRawKVClient(urls, config.Security{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create the RawKVClient %s", err)
	}

	// driver := tikv.Driver{}
	// storage, err := driver.Open("tikv://" + strings.Join(urls, ","))
	// if err != nil {
	// 	return nil, err
	// }
	// storage.Begin()

	return &db{client: client}, nil
}

type db struct {
	client *tikv.RawKVClient
	// storage kv.Storage
}

func (db *db) Put(ctx context.Context, namespace string, key, value []byte) error {
	return db.client.Put(buildKey(namespace, key), value)
}

func (db *db) Get(ctx context.Context, namespace string, key []byte) ([]byte, error) {
	return db.client.Get(buildKey(namespace, key))
}

func (db *db) Delete(ctx context.Context, namespace string, key []byte) error {
	return db.client.Delete(buildKey(namespace, key))
}

func (db *db) Scan(ctx context.Context, namespace string, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return db.client.Scan(buildKey(namespace, prefix), limit)
}

func (db *db) BatchGet(ctx context.Context, namespace string, keys [][]byte) ([][]byte, error) {
	return nil, nil
}

func (db *db) BatchPut(ctx context.Context, namespace string, keys, values [][]byte) error {
	return nil
}

func (db *db) BatchDelete(ctx context.Context, namespace string, keys [][]byte) error {
	return nil
}

func (db *db) Close() error {
	return db.client.Close()
}

func buildKey(namespace string, key []byte) []byte {
	if key != nil {
		return append([]byte(namespace), key...)
	}

	return []byte(namespace)
}
