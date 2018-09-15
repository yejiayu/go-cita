package service

import (
	"context"
	"math"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/common/merkle"
	"github.com/yejiayu/go-cita/pb"

	cfg "github.com/yejiayu/go-cita/config/chain"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
)

type Interface interface {
	GetBlockHeader(ctx context.Context, height uint64) (*pb.BlockHeader, error)
	GetBlockBody(ctx context.Context, height uint64) (*pb.BlockBody, error)
	GetValidators(ctx context.Context, height uint64) ([][]byte, error)

	NewBlock(ctx context.Context, block *pb.Block) error

	GetReceipt(ctx context.Context, txHash hash.Hash) (*pb.Receipt, error)
}

func New(dbFactory database.Factory, vmClient pb.VMClient) Interface {
	return &service{
		blockDB:  dbFactory.BlockDB(),
		vmClient: vmClient,
	}
}

type service struct {
	blockDB block.Interface

	vmClient pb.VMClient
}

func (svc *service) GetBlockHeader(ctx context.Context, height uint64) (*pb.BlockHeader, error) {
	if height == math.MaxUint64 {
		return svc.blockDB.GetHeaderByLatest(ctx)
	}

	return svc.blockDB.GetHeaderByHeight(ctx, height)
}

func (svc *service) GetBlockBody(ctx context.Context, height uint64) (*pb.BlockBody, error) {
	return svc.blockDB.GetBodyByHeight(ctx, height)
}

func (svc *service) GetValidators(ctx context.Context, height uint64) ([][]byte, error) {
	validators := make([][]byte, len(cfg.GetValidators()))
	for i, v := range cfg.GetValidators() {
		validators[i] = hash.FromHex(v)
	}

	return validators, nil
}

func (svc *service) NewBlock(ctx context.Context, block *pb.Block) error {
	preHeader, err := svc.GetBlockHeader(ctx, block.GetHeader().GetHeight()-1)
	if err != nil {
		return err
	}

	res, err := svc.vmClient.Call(ctx, &pb.CallReq{
		Header:   preHeader,
		TxHashes: block.GetBody().GetTxHashes(),
	})
	if err != nil {
		return err
	}

	var quotaUsed uint64
	for _, receipt := range res.GetReceipts() {
		receipt.StateRoot = res.GetStateRoot()
		quotaUsed += receipt.GetQuotaUsed()
		receipt.BlockHeight = block.GetHeader().GetHeight()
	}
	receiptsRoot := merkle.ReceiptsToRoot(res.GetReceipts())

	block.GetHeader().ReceiptsRoot = receiptsRoot.Bytes()
	block.GetHeader().StateRoot = res.GetStateRoot()
	block.GetHeader().QuotaUsed = quotaUsed

	return svc.blockDB.AddBlock(ctx, block, res.GetReceipts())
}

func (svc *service) GetReceipt(ctx context.Context, txHash hash.Hash) (*pb.Receipt, error) {
	return svc.blockDB.GetReceipt(ctx, txHash)
}
