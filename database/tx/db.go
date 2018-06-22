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

package tx

import (
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/types"
)

var (
	txPoolPrefix = byte(20)
)

type Interface interface {
	AddPool(signedTx *types.SignedTransaction) error
}

func New(rawDB raw.Interface) Interface {
	return &txDB{rawDB: rawDB}
}

type txDB struct {
	rawDB raw.Interface
}

func (db *txDB) AddPool(signedTx *types.SignedTransaction) error {
	data, err := proto.Marshal(signedTx)
	if err != nil {
		return err
	}

	key := txPoolKey(signedTx.GetTxHash())
	return db.rawDB.Put(key, data)
}

func txPoolKey(key []byte) []byte {
	return joinKey(txPoolPrefix, key)
}

func joinKey(prefix byte, key []byte) []byte {
	return append(append([]byte{prefix}, []byte(".")...), key...)
}
