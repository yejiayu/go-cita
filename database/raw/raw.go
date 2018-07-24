// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package raw

import (
	"context"
	"fmt"
	"strings"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/kv"
	"github.com/pingcap/tidb/store/tikv"
)

type Interface interface {
	Put(ctx context.Context, tableName, key, value []byte) error
	Get(ctx context.Context, tableName, key []byte) ([]byte, error)
	Delete(ctx context.Context, tableName, key []byte) error

	Scan(ctx context.Context, tableName, prefix []byte, limit int) ([][]byte, [][]byte, error)

	BeginTransaction(ctx context.Context) (Transaction, error)
}

func New(urls []string) (Interface, error) {
	client, err := tikv.NewRawKVClient(urls, config.Security{})
	if err != nil {
		return nil, fmt.Errorf("Failed to create the RawKVClient %s", err)
	}

	driver := tikv.Driver{}
	storage, err := driver.Open("tikv://" + strings.Join(urls, ","))
	if err != nil {
		return nil, err
	}
	storage.Begin()

	return &rawDB{
		client:  client,
		storage: storage,
	}, nil
}

type rawDB struct {
	client  *tikv.RawKVClient
	storage kv.Storage
}

func (db *rawDB) Put(ctx context.Context, tableName, key, value []byte) error {
	return db.client.Put(buildKey(tableName, key), value)
}

func (db *rawDB) Get(ctx context.Context, tableName, key []byte) ([]byte, error) {
	return db.client.Get(buildKey(tableName, key))
	// data, err := db.client.Get(buildKey(tableName, key))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if len(data) > 0 {
	// 	return nil, errors.DatabaseNotFound
	// }
	//
	// return data, nil
}

func (db *rawDB) Delete(ctx context.Context, tableName, key []byte) error {
	return db.client.Delete(buildKey(tableName, key))
}

func (db *rawDB) Scan(ctx context.Context, tableName, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	return db.client.Scan(buildKey(tableName, prefix), limit)
}

func (db *rawDB) BeginTransaction(ctx context.Context) (Transaction, error) {
	tx, err := db.storage.Begin()
	if err != nil {
		return nil, err
	}

	return beginTransaction(ctx, tx), nil
}

func buildKey(tableName, key []byte) []byte {
	if key != nil {
		return append(tableName, key...)
	}

	return tableName
}
