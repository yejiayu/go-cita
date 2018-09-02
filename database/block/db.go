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
	"math"

	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/pb"
)

var (
	nsBlockHeader  = "block.header"
	nsBlockBody    = "block.body"
	nsLatestHeight = "block.latest.height"

	nsReceipts = "block.receipts"
)

type Interface interface {
	GetHeaderByHeight(ctx context.Context, height uint64) (*pb.BlockHeader, error)
	GetBodyByHeight(ctx context.Context, height uint64) (*pb.BlockBody, error)
	GetHeaderByLatest(ctx context.Context) (*pb.BlockHeader, error)

	AddBlock(ctx context.Context, b *pb.Block, receipts []*pb.Receipt) error
}

func New(raw raw.Interface) Interface {
	return &blockDB{raw: raw}
}

type blockDB struct {
	raw raw.Interface
}

func (db *blockDB) GetHeaderByHeight(ctx context.Context, height uint64) (*pb.BlockHeader, error) {
	if height == math.MaxUint64 {
		var err error
		height, err = db.getLatest(ctx)
		if err != nil {
			return nil, err
		}
	}

	data, err := db.raw.Get(ctx, nsBlockHeader, uint64ToBytes(height))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	var header pb.BlockHeader
	if err := proto.Unmarshal(data, &header); err != nil {
		return nil, err
	}

	return &header, nil
}

func (db *blockDB) GetBodyByHeight(ctx context.Context, height uint64) (*pb.BlockBody, error) {
	data, err := db.raw.Get(ctx, nsBlockBody, uint64ToBytes(height))
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
	if height == 0 {
		return db.GetHeaderByHeight(ctx, height)
	}

	return db.GetHeaderByHeight(ctx, height)
}

// TODO: transaction
func (db *blockDB) AddBlock(ctx context.Context, b *pb.Block, receipts []*pb.Receipt) error {
	header := b.GetHeader()
	hData, err := proto.Marshal(header)
	if err != nil {
		return err
	}
	if err = db.raw.Put(ctx, nsBlockHeader, uint64ToBytes(header.GetHeight()), hData); err != nil {
		return err
	}

	body := b.GetBody()
	if body != nil && len(body.GetTxHashes()) > 0 {
		bData, err := proto.Marshal(body)
		if err != nil {
			return err
		}
		if err := db.raw.Put(ctx, nsBlockBody, uint64ToBytes(header.GetHeight()), bData); err != nil {
			return err
		}
	}

	if err := db.raw.Put(ctx, nsLatestHeight, []byte("latest"), uint64ToBytes(header.GetHeight())); err != nil {
		return err
	}

	if len(receipts) > 0 {
		receiptKeys := make([][]byte, len(receipts))
		receiptValues := make([][]byte, len(receipts))

		for i, receipt := range receipts {
			value, err := proto.Marshal(receipt)
			if err != nil {
				return err
			}
			key := hash.BytesToSha3(value).Bytes()
			receiptKeys[i] = key
			receiptValues[i] = value
		}

		return db.raw.BatchPut(ctx, nsReceipts, receiptKeys, receiptValues)
	}
	return nil
}

func (db *blockDB) getLatest(ctx context.Context) (uint64, error) {
	data, err := db.raw.Get(ctx, nsLatestHeight, []byte("latest"))
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, nil
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
