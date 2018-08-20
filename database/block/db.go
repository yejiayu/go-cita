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
	"context"
	"encoding/binary"

	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/pb"
)

var (
	blockHeaderTable = []byte("block-header-")
	blockBodyTable   = []byte("block-body-")
	blockLatestTable = []byte("block-latest")
)

type Interface interface {
	GetHeaderByHeight(ctx context.Context, height uint64) (*pb.BlockHeader, error)
	GetBodyByHeight(ctx context.Context, height uint64) (*pb.BlockBody, error)
	GetHeaderByLatest(ctx context.Context) (*pb.BlockHeader, error)

	AddBlock(ctx context.Context, b *pb.Block) error
}

func New(raw raw.Interface) Interface {
	return &blockDB{raw: raw}
}

type blockDB struct {
	raw raw.Interface
}

func (db *blockDB) GetHeaderByHeight(ctx context.Context, height uint64) (*pb.BlockHeader, error) {
	data, err := db.raw.Get(ctx, blockHeaderTable, uint64ToBytes(height))
	if err != nil {
		return nil, err
	}

	var header pb.BlockHeader
	if err := proto.Unmarshal(data, &header); err != nil {
		return nil, err
	}

	return &header, nil
}

func (db *blockDB) GetBodyByHeight(ctx context.Context, height uint64) (*pb.BlockBody, error) {
	data, err := db.raw.Get(ctx, blockBodyTable, uint64ToBytes(height))
	if err != nil {
		return nil, err
	}

	var body pb.BlockBody
	if err := proto.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	return &body, nil
}

func (db *blockDB) GetHeaderByLatest(ctx context.Context) (*pb.BlockHeader, error) {
	height, err := db.getLatest(ctx)
	if err != nil {
		return nil, err
	}

	return db.GetHeaderByHeight(ctx, height)
}

func (db *blockDB) AddBlock(ctx context.Context, b *pb.Block) error {
	header := b.GetHeader()
	hData, err := proto.Marshal(header)
	if err != nil {
		return err
	}

	// tx, err := db.raw.BeginTransaction(ctx)
	// if err != nil {
	// 	return err
	// }

	// log.Info("put block header")
	if err = db.raw.Put(ctx, blockHeaderTable, uint64ToBytes(header.GetHeight()), hData); err != nil {
		// tx.Rollback()
		return err
	}

	body := b.GetBody()
	if body != nil && len(body.GetTxHashes()) > 0 {
		bData, err := proto.Marshal(body)
		if err != nil {
			// tx.Rollback()
			return err
		}
		if err := db.raw.Put(ctx, blockBodyTable, uint64ToBytes(header.GetHeight()), bData); err != nil {
			// tx.Rollback()
			return err
		}
	}

	if err := db.raw.Put(ctx, blockLatestTable, nil, uint64ToBytes(header.GetHeight())); err != nil {
		// tx.Rollback()
		return err
	}
	return nil
	// err = tx.Commit()
	// if err != nil {
	// 	log.Errorf("add block error %s", err)
	// }
	// log.Info("add block commited")
	// return err
}

func (db *blockDB) getLatest(ctx context.Context) (uint64, error) {
	data, err := db.raw.Get(ctx, blockLatestTable, nil)
	if err != nil {
		return 0, err
	}

	return bytesToUint64(data), nil
}

func uint64ToBytes(num uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, num)
	return buf
}

func bytesToUint64(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}
