package database

import (
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/database/tx"
)

type Factory interface {
	BlockDB() block.Interface
	TxDB() tx.Interface
}

func NewFactory(urls []string) (Factory, error) {
	rawDB, err := raw.New(urls)
	if err != nil {
		return nil, err
	}

	return &factory{rawDB: rawDB}, nil
}

type factory struct {
	rawDB raw.Interface
}

func (f *factory) BlockDB() block.Interface {
	return block.New(f.rawDB)
}

func (f *factory) TxDB() tx.Interface {
	return tx.New(f.rawDB)
}
