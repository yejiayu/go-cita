package vm

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/tx"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	cfg "github.com/yejiayu/go-cita/config/vm"
)

type Server interface {
	Run()
}

func NewServer(factory database.Factory) Server {
	return &server{
		grpcS:    grpc.NewServer(),
		executor: NewExecutor(factory),

		txDB: factory.TxDB(),
	}
}

type server struct {
	grpcS    *grpc.Server
	executor Executor

	txDB tx.Interface
}

func (s *server) Run() {
	port := cfg.GetPort()
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("The vm server listens on port %s", port)
	pb.RegisterVMServer(s.grpcS, s)
	if err := s.grpcS.Serve(lis); err != nil {
		log.Panic(err)
	}
}

func (s *server) Call(ctx context.Context, req *pb.CallReq) (*pb.CallRes, error) {
	header := req.GetHeader()
	signedTxs, err := s.txDB.Get(ctx, req.GetTxHashes())
	if err != nil {
		return nil, err
	}

	if len(signedTxs) == 0 {
		return &pb.CallRes{StateRoot: header.GetStateRoot()}, nil
	}

	repiects, root, err := s.executor.Call(ctx, header, signedTxs)
	if err != nil {
		return nil, err
	}

	return &pb.CallRes{Receipts: repiects, StateRoot: root}, nil
}

func (s *server) StaticCall(ctx context.Context, req *pb.StaticCallReq) (*pb.StaticCallRes, error) {
	ret, err := s.executor.StaticCall(ctx, req.GetHeight(), req.GetFrom(), req.GetTo(), req.GetData())
	if err != nil {
		return nil, err
	}

	return &pb.StaticCallRes{Result: ret}, nil
}
