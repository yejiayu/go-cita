package chain

import (
	"context"
	"math"
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	cfg "github.com/yejiayu/go-cita/config/chain"
)

type Server interface {
	Run()
}

func New(dbFactory database.Factory) Server {
	return &server{
		blockDB: dbFactory.BlockDB(),
		grpcS:   grpc.NewServer(),
	}
}

type server struct {
	blockDB block.Interface
	grpcS   *grpc.Server
}

func (s *server) Run() {
	port := cfg.GetPort()
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("The chain server listens on port %s", port)
	pb.RegisterChainServer(s.grpcS, s)
	if err := s.grpcS.Serve(lis); err != nil {
		log.Panic(err)
	}
}

func (s *server) NewBlock(ctx context.Context, req *pb.NewBlockReq) (*pb.NewBlockRes, error) {
	if err := s.blockDB.AddBlock(ctx, req.GetBlock()); err != nil {
		return nil, err
	}
	return &pb.NewBlockRes{Height: req.GetBlock().GetHeader().GetHeight()}, nil
}

func (s *server) NodeList(ctx context.Context, req *pb.NodeListReq) (*pb.NodeListRes, error) {
	priv1, err := crypto.HexToECDSA("add757cf60afa08fc54376db9cd1f313f2d20d907f3ac984f227ea0835fc0111")
	if err != nil {
		return nil, err
	}

	return &pb.NodeListRes{
		Nodes: [][]byte{crypto.CompressPubkey(&priv1.PublicKey)},
	}, nil
}

func (s *server) GetBlockHeader(ctx context.Context, req *pb.GetBlockHeaderReq) (*pb.GetBlockHeaderRes, error) {
	if req.GetHeight() == math.MaxUint64 {
		header, err := s.blockDB.GetHeaderByLatest(ctx)
		if err != nil {
			return nil, err
		}
		return &pb.GetBlockHeaderRes{Header: header}, nil
	}

	header, err := s.blockDB.GetHeaderByHeight(ctx, req.GetHeight())
	if err != nil {
		return nil, err
	}
	return &pb.GetBlockHeaderRes{Header: header}, nil
}
