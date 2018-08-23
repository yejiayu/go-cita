package service

import (
	"context"
	"math"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/pb"
)

type Interface interface {
	GetBlockHeader(ctx context.Context, height uint64) (*pb.BlockHeader, error)
	GetValidators(ctx context.Context, height uint64) ([][]byte, error)

	NewBlock(ctx context.Context, block *pb.Block) error
}

func New(dbFactory database.Factory) Interface {
	return &service{
		blockDB: dbFactory.BlockDB(),
	}
}

type service struct {
	blockDB block.Interface
}

func (svc *service) GetBlockHeader(ctx context.Context, height uint64) (*pb.BlockHeader, error) {
	if height == math.MaxUint64 {
		return svc.blockDB.GetHeaderByLatest(ctx)
	}

	return svc.blockDB.GetHeaderByHeight(ctx, height)
}

func (svc *service) GetValidators(ctx context.Context, height uint64) ([][]byte, error) {
	priv1, err := crypto.HexToECDSA("add757cf60afa08fc54376db9cd1f313f2d20d907f3ac984f227ea0835fc0111")
	if err != nil {
		return nil, err
	}

	return [][]byte{crypto.CompressPubkey(&priv1.PublicKey)}, nil
}

func (svc *service) NewBlock(ctx context.Context, block *pb.Block) error {
	return svc.blockDB.AddBlock(ctx, block)
}
