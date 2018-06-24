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

package block

import (
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/types"
)

var (
	headerPrefix = byte(10)
	bodyPrefix   = byte(11)

	latestHeader = []byte("block.header.latest")
)

type Interface interface {
	GetHeaderByHeight(height uint64) (*types.BlockHeader, error)
	GetBodyByHeight(height uint64) (*types.BlockBody, error)
	GetHeaderByLatest() (*types.BlockHeader, error)
}

func New(rawDB raw.Interface) Interface {
	return &blockDB{rawDB: rawDB}
}

type blockDB struct {
	rawDB raw.Interface
}

func (db *blockDB) GetHeaderByHeight(height uint64) (*types.BlockHeader, error) {
	data, err := db.rawDB.Get(headerKey(height))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var h types.BlockHeader
	if err := proto.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

func (db *blockDB) GetBodyByHeight(height uint64) (*types.BlockBody, error) {
	data, err := db.rawDB.Get(bodyKey(height))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var body types.BlockBody
	if err := proto.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	return &body, nil
}

func (db *blockDB) GetHeaderByLatest() (*types.BlockHeader, error) {
	data, err := db.rawDB.Get(latestHeader)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var header types.BlockHeader
	if err := proto.Unmarshal(data, &header); err != nil {
		return nil, err
	}

	return &header, nil
}

// func (db *blockDB) Scan(startHeight uint64, limit int) ([]*types.Block, error) {
// 	_, values, err := db.rawDB.Scan(blockKey(startHeight), limit)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var bs []*types.Block
// 	for _, value := range values {
// 		var b types.Block
// 		if err := proto.Unmarshal(value, &b); err != nil {
// 			return nil, err
// 		}
//
// 		bs = append(bs, &b)
// 	}
//
// 	return bs, nil
// }

func headerKey(height uint64) []byte {
	var key []byte
	binary.BigEndian.PutUint64(key, height)
	return joinKey(headerPrefix, key)
}

func bodyKey(height uint64) []byte {
	var key []byte
	binary.BigEndian.PutUint64(key, height)
	return joinKey(bodyPrefix, key)
}

func joinKey(prefix byte, key []byte) []byte {
	return append(append([]byte{prefix}, []byte(".")...), key...)
}
