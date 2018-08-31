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

package database

import (
	"log"
	"strings"

	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/database/raw/redis"
	"github.com/yejiayu/go-cita/database/raw/tikv"

	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/database/ethdb"
	"github.com/yejiayu/go-cita/database/tx"
)

type Factory interface {
	RawDB() raw.Interface
	BlockDB() block.Interface
	TxDB() tx.Interface
	// StateDB(root hash.Hash) state.Interface
	EthDB() ethdb.Database
}

func NewFactory(t string, urls []string) (Factory, error) {
	var raw raw.Interface
	var err error

	switch strings.ToLower(t) {
	case "tikv":
		raw, err = tikv.New(urls)
	case "redis":
		raw = redis.New(urls)
	default:
		log.Panicf("Can't match type %s", t)
	}

	if err != nil {
		return nil, err
	}
	return &factory{raw: raw}, nil
}

type factory struct {
	raw raw.Interface
}

func (f *factory) RawDB() raw.Interface {
	return f.raw
}

func (f *factory) BlockDB() block.Interface {
	return block.New(f.raw)
}

func (f *factory) TxDB() tx.Interface {
	return tx.New(f.raw)
}

// func (f *factory) StateDB(root hash.Hash) state.Interface {
// 	return state.New(root, f.raw)
// }

func (f *factory) EthDB() ethdb.Database {
	return ethdb.New(f.raw)
}
