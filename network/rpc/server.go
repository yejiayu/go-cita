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

package rpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	cfg "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/network/protocol"
	"github.com/yejiayu/go-cita/network/rpc/connection"
)

type Server interface {
	Run()
}

func New() Server {
	return &server{
		grpcS: grpc.NewServer(),
		m:     connection.NewManager(),
	}
}

type server struct {
	grpcS *grpc.Server

	m connection.Manager
}

func (s *server) Run() {
	port := cfg.GetRPCPort()
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("The network server listens on port %s", port)
	pb.RegisterNetworkServer(s.grpcS, s)

	if err := s.grpcS.Serve(lis); err != nil {
		log.Panic(err)
	}
}

func (s *server) BroadcastProposal(ctx context.Context, req *pb.BroadcastProposalReq) (*pb.Empty, error) {
	err := s.m.Broadcast(protocol.KeyBroadcastProposal, req)
	return &pb.Empty{}, err
}

func (s *server) BroadcastVote(ctx context.Context, req *pb.BroadcastVotelReq) (*pb.Empty, error) {
	err := s.m.Broadcast(protocol.KeyBroadcastVote, req)
	return &pb.Empty{}, err
}

func (s *server) BroadcastTransaction(ctx context.Context, req *pb.BroadcastTransactionReq) (*pb.Empty, error) {
	err := s.m.Broadcast(protocol.KeyBroadcastTransaction, req)
	return &pb.Empty{}, err
}

func (s *server) GetUnverifyTxs(ctx context.Context, req *pb.GetUnverifyTxsReq) (*pb.GetUnverifyTxsRes, error) {
	return nil, nil
}
