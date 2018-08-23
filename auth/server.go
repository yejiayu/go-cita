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

package auth

import (
	"context"
	"net"

	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/auth/cache"
	poolquota "github.com/yejiayu/go-cita/auth/pool/quota"
	"github.com/yejiayu/go-cita/auth/service"
	cfg "github.com/yejiayu/go-cita/config/auth"
	"github.com/yejiayu/go-cita/database"
)

type Server interface {
	Run()
}

func New(dbFactory database.Factory, networkClient pb.NetworkClient) Server {
	conn, err := redis.DialURL("redis://" + cfg.GetRedisURL())
	if err != nil {
		log.Panic(err)
	}

	pool := poolquota.New(dbFactory.TxDB(), cfg.GetRedisURL(), cfg.GetPoolCount())
	cache := cache.New(conn)
	svc := service.New(dbFactory, networkClient, cache, pool)

	return &server{
		grpcS: grpc.NewServer(),
		svc:   svc,
	}
}

// AuthServer is the server API for Auth service.
type server struct {
	grpcS *grpc.Server
	svc   service.Interface
}

func (s *server) Run() {
	port := cfg.GetPort()
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("The auth server listens on port %s", port)
	pb.RegisterAuthServer(s.grpcS, s)

	if err := s.grpcS.Serve(lis); err != nil {
		log.Panic(err)
	}
}

func (s *server) AddUnverifyTx(ctx context.Context, req *pb.AddUnverifyTxReq) (*pb.AddUnverifyTxRes, error) {
	hash, err := s.svc.AddUnverifyTx(ctx, req.GetUntx())
	if err != nil {
		return nil, err
	}

	return &pb.AddUnverifyTxRes{
		TxHash: hash.Bytes(),
	}, nil
}

func (s *server) GetTxFromPool(ctx context.Context, req *pb.GetTxFromPoolReq) (*pb.GetTxFromPoolRes, error) {
	hashes, err := s.svc.GetHashFromPool(ctx, req.GetTxCount(), req.GetQuotaLimit())
	if err != nil {
		return nil, err
	}

	return &pb.GetTxFromPoolRes{
		TxHashes: hash.HashesToBytesS(hashes),
	}, nil
}

func (s *server) EnsureFromPool(ctx context.Context, req *pb.EnsureFromPoolReq) (*pb.Empty, error) {
	if err := s.svc.EnsureFromPool(ctx, req.GetNodeId(), req.GetQuotaUsed(), hash.BytesSToHashes(req.GetTxHashes())); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) ClearPool(ctx context.Context, req *pb.ClearPoolReq) (*pb.Empty, error) {
	if err := s.svc.ClearPool(ctx, req.GetHeight()); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) EnsureHashes(ctx context.Context, req *pb.EnsureHashesReq) (*pb.EnsureHashesRes, error) {
	return nil, nil
}
