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

package consensus

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/clients"
	cfg "github.com/yejiayu/go-cita/config/consensus"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/consensus/tendermint"
)

type Server interface {
	Run()
}

func New() Server {
	return &server{
		grpcS: grpc.NewServer(),
		tendermint: tendermint.New(
			clients.NewAuthClient(cfg.GetAuthURL()),
			clients.NewChainClient(cfg.GetChainURL()),
			clients.NewNetworkClient(cfg.GetNetworkURL()),
		),
	}
}

type server struct {
	grpcS      *grpc.Server
	tendermint tendermint.Interface
}

func (s *server) Run() {
	port := cfg.GetPort()
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("The consensus server listens on port %s", port)
	pb.RegisterConsensusServer(s.grpcS, s)

	if err := s.grpcS.Serve(lis); err != nil {
		log.Panic(err)
	}
}

func (s *server) SetProposal(ctx context.Context, in *pb.SetProposalReq) (*pb.Empty, error) {
	if err := s.tendermint.SetProposal(ctx, in.GetProposal(), in.GetSignature()); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) AddVote(ctx context.Context, in *pb.AddVoteReq) (*pb.Empty, error) {
	if err := s.tendermint.SetVote(ctx, in.GetVote(), in.GetSignature()); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
