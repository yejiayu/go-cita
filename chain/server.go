package chain

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/chain/service"
	"github.com/yejiayu/go-cita/database"
	blockdb "github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/types"
)

func New(port string, dbFactory database.Factory) error {
	s := grpc.NewServer()

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}

	svc, err := service.New(dbFactory)
	if err != nil {
		return err
	}

	types.RegisterChainServer(s, &server{
		blockDB: dbFactory.BlockDB(),
		svc:     svc,
	})

	log.Infof("The chain server listens on port %s", port)
	return s.Serve(lis)
}

type server struct {
	blockDB blockdb.Interface
	svc     service.Interface
}

func (s *server) NewBlock(ctx context.Context, req *types.NewBlockReq) (*types.NewBlockRes, error) {
	block := req.GetBlock()

	if err := s.blockDB.AddBlock(ctx, block); err != nil {
		return nil, err
	}

	return &types.NewBlockRes{Height: block.GetHeader().GetHeight()}, nil
}

func (s *server) LatestHeight(ctx context.Context, req *types.LatestHeightReq) (*types.LatestHeightRes, error) {
	height, err := s.svc.GetLatestHeight(ctx)
	if err != nil {
		return nil, err
	}

	return &types.LatestHeightRes{Height: height}, nil
}
