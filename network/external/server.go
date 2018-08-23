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

package external

import (
	"context"
	"io"
	"net"

	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	cfg "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/network/protocol"
)

type Server interface {
	Run()
}

func New(
	consensusClient pb.ConsensusClient,
	authClient pb.AuthClient,
	chainClient pb.ChainClient,
) Server {
	return &server{
		codec: protocol.NewCodec(),

		consensusClient: consensusClient,
		authClient:      authClient,
		chainClient:     chainClient,
	}
}

type server struct {
	codec protocol.Codec

	consensusClient pb.ConsensusClient
	authClient      pb.AuthClient
	chainClient     pb.ChainClient
}

func (s *server) Run() {
	port := cfg.GetExternalPort()
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Panic(err)
	}

	log.Infof("The network external server listens on port %s", port)
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Panic(err)
			}

			go func(conn net.Conn) {
				err := s.handler(conn)
				if err != io.EOF {
					log.Errorf("external handler err, %s", err)
					conn.Close()
				}
			}(conn)
		}
	}()
}

func (s *server) handler(conn net.Conn) error {
	for {
		key, data, err := s.codec.Decode(conn)
		if err != nil {
			return err
		}

		ctx := context.Background()
		err = s.dispatch(ctx, key, protocol.NewMessageWithRaw(data))
		if err != nil {
			return err
		}
	}
}

func (s *server) dispatch(ctx context.Context, key string, msg protocol.Message) error {
	origin := msg.Origin()
	payload := msg.Payload()

	switch key {
	case protocol.KeyBroadcastProposal:
		var req pb.BroadcastProposalReq
		if err := proto.Unmarshal(payload, &req); err != nil {
			return err
		}
		_, err := s.consensusClient.SetProposal(ctx, &pb.SetProposalReq{
			Proposal:  req.GetProposal(),
			Signature: req.GetSignature(),
		})
		return err
	case protocol.KeyBroadcastVote:
		var req pb.BroadcastVotelReq
		if err := proto.Unmarshal(payload, &req); err != nil {
			return err
		}
		_, err := s.consensusClient.AddVote(ctx, &pb.AddVoteReq{
			Vote:      req.GetVote(),
			Signature: req.GetSignature(),
		})
		return err
	case protocol.KeyBroadcastTransaction:
		var req pb.BroadcastTransactionReq
		if err := proto.Unmarshal(payload, &req); err != nil {
			return err
		}
		_, err := s.authClient.AddUnverifyTx(ctx, &pb.AddUnverifyTxReq{
			Untx: req.GetUntx(),
		})
		return err
	default:
		log.Panicf("unknown key %s, from %d", key, origin)
	}

	return nil
}
