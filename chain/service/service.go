package service

import (
	"context"
	"time"

	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/types"
)

type Interface interface {
}

func New(dbFactory database.Factory) (Interface, error) {
	svc := &service{
		blockDB: dbFactory.BlockDB(),
	}

	if err := svc.init(); err != nil {
		return nil, err
	}
	return svc, nil
}

type service struct {
	blockDB block.Interface
}

func (s *service) init() error {
	// h, err := s.blockDB.GetHeaderByHeight(context.Background(), 0)
	// if err != nil {
	// 	return err
	// }

	// if h == nil {
	// 	return s.blockDB.AddBlock(context.Background(), &types.Block{
	// 		Header: &types.BlockHeader{
	// 			Height:    0,
	// 			Timestamp: uint64(time.Now().Unix()),
	// 		},
	// 		Body: &types.BlockBody{
	// 			TxHashes: [][]byte{},
	// 		},
	// 	})
	// }

	return s.blockDB.AddBlock(context.Background(), &types.Block{
		Header: &types.BlockHeader{
			Height:    0,
			Timestamp: uint64(time.Now().Unix()),
		},
	})

	// return nil
}
