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

package network

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	cfg "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/grpc/middleware/logger"
)

type Server interface {
	Run()
}

func NewServer(
	address string,
	nodeClientMap map[string]pb.NetworkClient,
	consensusClient pb.ConsensusClient,
	authClient pb.AuthClient,
	chainClient pb.ChainClient,
) Server {
	return &server{
		grpcS: grpc.NewServer(
			grpc.UnaryInterceptor(loggger.NewServer()),
		),

		address: address,

		nodeClientMap:   nodeClientMap,
		consensusClient: consensusClient,
		authClient:      authClient,
		chainClient:     chainClient,
	}
}

type server struct {
	grpcS *grpc.Server

	address string

	nodeClientMap   map[string]pb.NetworkClient
	consensusClient pb.ConsensusClient
	authClient      pb.AuthClient
	chainClient     pb.ChainClient
}

func (s *server) Run() {
	port := cfg.GetPort()
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
	for address, node := range s.nodeClientMap {
		go func(address string, node pb.NetworkClient) {
			_, err := node.SetProposal(context.Background(), &pb.SetProposalReq{
				Proposal:  req.GetProposal(),
				Signature: req.GetSignature(),
			})
			if err != nil {
				log.Errorf("broadcast proposal %s, node %s", err.Error(), address)
			}
		}(address, node)
	}

	return &pb.Empty{}, nil
}

func (s *server) BroadcastVote(ctx context.Context, req *pb.BroadcastVotelReq) (*pb.Empty, error) {
	for address, node := range s.nodeClientMap {
		go func(address string, node pb.NetworkClient) {
			_, err := node.AddVote(context.Background(), &pb.AddVoteReq{
				Vote:      req.GetVote(),
				Signature: req.GetSignature(),
			})
			if err != nil {
				log.Errorf("broadcast vote %s, node %s", err.Error(), address)
			}
		}(address, node)
	}

	return &pb.Empty{}, nil
}

func (s *server) BroadcastTransaction(ctx context.Context, req *pb.BroadcastTransactionReq) (*pb.Empty, error) {
	for address, node := range s.nodeClientMap {
		go func(address string, node pb.NetworkClient) {
			_, err := node.AddUnverifyTx(context.Background(), &pb.AddUnverifyTxReq{
				Untx: req.GetUntx(),
			})
			if err != nil {
				log.Errorf("broadcast transaction %s, node %s", err.Error(), address)
			}
		}(address, node)
	}

	return &pb.Empty{}, nil
}

func (s *server) BroadcastHeight(ctx context.Context, req *pb.BroadcastHeightReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *server) SetProposal(ctx context.Context, req *pb.SetProposalReq) (*pb.Empty, error) {
	return s.consensusClient.SetProposal(ctx, req)
}

func (s *server) AddVote(ctx context.Context, req *pb.AddVoteReq) (*pb.Empty, error) {
	return s.consensusClient.AddVote(ctx, req)
}

func (s *server) AddUnverifyTx(ctx context.Context, req *pb.AddUnverifyTxReq) (*pb.AddUnverifyTxRes, error) {
	return s.authClient.AddUnverifyTx(ctx, req)
}

func (s *server) GetUnverifyTxs(ctx context.Context, req *pb.GetUnverifyTxsReq) (*pb.GetUnverifyTxsRes, error) {
	address := hash.ToHex(req.GetNode())
	if s.address == address {
		return s.authClient.GetUnverifyTxs(ctx, &pb.GetUnverifyTxsReq{
			TxHashes: req.GetTxHashes(),
			Node:     req.GetNode(),
		})
	}

	node, ok := s.nodeClientMap[address]
	if !ok {
		return &pb.GetUnverifyTxsRes{}, nil
	}

	return node.GetUnverifyTxs(ctx, &pb.GetUnverifyTxsReq{
		TxHashes: req.GetTxHashes(),
		Node:     req.GetNode(),
	})
}
