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
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/pb"
)

var (
	txTable     = []byte("tx-")
	txPoolTable = []byte("tx-pool-")
)

type Interface interface {
	AddPool(ctx context.Context, signedTx *pb.SignedTransaction) error

	GetByHash(ctx context.Context, hash common.Hash) (*pb.SignedTransaction, error)
	GetTxhashesFromPool(ctx context.Context, limit int) ([][]byte, error)
	Exists(ctx context.Context, hash common.Hash) (bool, error)

	UpdateTx(ctx context.Context, hashes [][]byte) error
}

func New(raw raw.Interface) Interface {
	return &txDB{raw: raw}
}

type txDB struct {
	raw raw.Interface
}

func (db *txDB) AddPool(ctx context.Context, signedTx *pb.SignedTransaction) error {
	data, err := proto.Marshal(signedTx)
	if err != nil {
		return err
	}

	// TODO: begin transaction
	if err := db.raw.Put(ctx, txTable, signedTx.GetTxHash(), data); err != nil {
		return err
	}
	return db.raw.Put(ctx, txPoolTable, signedTx.GetTxHash(), signedTx.GetTxHash())
}

func (db *txDB) GetByHash(ctx context.Context, hash common.Hash) (*pb.SignedTransaction, error) {
	data, err := db.raw.Get(ctx, txTable, hash.Bytes())
	if err != nil {
		return nil, err
	}

	var signedTx pb.SignedTransaction
	if err := proto.Unmarshal(data, &signedTx); err != nil {
		return nil, err
	}

	return &signedTx, nil
}

func (db *txDB) UpdateTx(ctx context.Context, hashes [][]byte) error {
	for _, hash := range hashes {
		if err := db.raw.Delete(ctx, txPoolTable, hash); err != nil {
			return err
		}
	}

	return nil
}

func (db *txDB) Exists(ctx context.Context, hash common.Hash) (bool, error) {
	data, err := db.raw.Get(ctx, txTable, hash.Bytes())
	if err != nil {
		return false, err
	}
	return len(data) > 0, nil
}

func (db *txDB) GetTxhashesFromPool(ctx context.Context, limit int) ([][]byte, error) {
	_, hashes, err := db.raw.Scan(ctx, txPoolTable, nil, limit)
	if err != nil {
		return nil, err
	}

	return hashes, nil
}
