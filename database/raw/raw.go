package raw

import (
	"fmt"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
)

type Interface interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error

	Scan(prefix []byte, limit int) ([][]byte, [][]byte, error)
}

func New(urls []string) (Interface, error) {
	// "47.75.129.215:2379", "47.75.129.215:2380", "47.75.129.215:2381"
	client, err := tikv.NewRawKVClient(urls, config.Security{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create the RawKVClient %s", err)
	}

	return &rawDB{
		client: client,
	}, nil
}

type rawDB struct {
	client *tikv.RawKVClient
}

func (db *rawDB) Put(key, value []byte) error {
	return db.client.Put(key, value)
}

func (db *rawDB) Get(key []byte) ([]byte, error) {
	return db.client.Get(key)
}

func (db *rawDB) Delete(key []byte) error {
	return db.client.Delete(key)
}

func (db *rawDB) Scan(prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return db.client.Scan(prefix, limit)
}
