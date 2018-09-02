package service

import (
	"context"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"
)

type Interface interface {
	GetBlockHeader(ctx context.Context, height uint64) (*pb.BlockHeader, error)
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

func (svc *service) GetValidators(ctx context.Context, height uint64) ([][]byte, error) {
	priv1, err := crypto.HexToECDSA("add757cf60afa08fc54376db9cd1f313f2d20d907f3ac984f227ea0835fc0111")
	if err != nil {
		return nil, err
	}

	return [][]byte{crypto.CompressPubkey(&priv1.PublicKey)}, nil
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

	// TODO: use MPT?
	var receiptsData []byte
	for _, receipt := range res.GetReceipts() {
		log.Info(common.BytesToAddress(receipt.GetContractAddress()).String())
		data, err := proto.Marshal(receipt)
		if err != nil {
			return err
		}
		receiptsData = append(receiptsData, data...)
	}
	receiptsRoot := hash.BytesToSha3(receiptsData).Bytes()

	block.GetHeader().ReceiptsRoot = receiptsRoot
	block.GetHeader().StateRoot = res.GetStateRoot()

	return svc.blockDB.AddBlock(ctx, block, res.GetReceipts())
}

func (svc *service) GetReceipt(ctx context.Context, txHash hash.Hash) (*pb.Receipt, error) {
	return svc.blockDB.GetReceipt(ctx, txHash)
}
