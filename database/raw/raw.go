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
