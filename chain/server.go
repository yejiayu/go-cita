package chain

import (
	"context"
	"math"
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/grpc/middleware/logger"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/chain/service"
	cfg "github.com/yejiayu/go-cita/config/chain"
)

type Server interface {
	Run()
}

func New(dbFactory database.Factory, vmClient pb.VMClient) Server {
	return &server{
		grpcS: grpc.NewServer(
			grpc.UnaryInterceptor(loggger.NewServer()),
		),

		blockDB: dbFactory.BlockDB(),

		svc: service.New(dbFactory, vmClient),
	}
}

type server struct {
	grpcS *grpc.Server

	blockDB block.Interface

	svc service.Interface
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

func (s *server) NewBlock(ctx context.Context, req *pb.NewBlockReq) (*pb.Empty, error) {
	err := s.svc.NewBlock(ctx, req.GetBlock())
	return &pb.Empty{}, err
}

func (s *server) GetValidators(ctx context.Context, req *pb.GetValidatorsReq) (*pb.GetValidatorsRes, error) {
	vals, err := s.svc.GetValidators(ctx, req.GetHeight())
	if err != nil {
		return nil, err
	}

	return &pb.GetValidatorsRes{Vals: vals}, nil
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
