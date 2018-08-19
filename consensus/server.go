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

	"github.com/yejiayu/go-cita/common/hash"
	cfg "github.com/yejiayu/go-cita/config/consensus"
	"github.com/yejiayu/go-cita/consensus/tendermint"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"
)

type Server interface {
	Run() error
}

func New(dbFactory database.Factory) (Server, error) {
	// err = tendermint.New(dbFactory, authClient, chainClient).Run()

	return &server{
		grpcS: grpc.NewServer(),
		// tendermint: tendermint.New(, chainClient types.ChainClient, privKey *ecdsa.PrivateKey)
	}, nil
}

type server struct {
	grpcS      *grpc.Server
	tendermint tendermint.Interface
}

func (s *server) Run() error {
	port := cfg.GetPort()
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}

	log.Infof("The consensus server listens on port %s", port)
	pb.RegisterConsensusServer(s, &server{})
	return s.grpcS.Serve(lis)
}

func (s *server) SignedProposal(ctx context.Context, in *pb.SignedProposalReq) (*pb.Empty, error) {
	if err := s.tendermint.SignedProposal(ctx, in.GetProposal(), in.GetSignature()); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) VotingBlock(ctx context.Context, in *pb.VotingBlockReq) (*pb.Empty, error) {
	if err := s.tendermint.VotingBlock(ctx, hash.BytesToHash(in.GetHash()), in.GetSignature()); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
