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

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/auth/service"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/types"
)

func New(port, redisURL string, dbFactory database.Factory) error {
	s := grpc.NewServer()
	svc, err := service.New(redisURL, dbFactory)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}

	log.Infof("The auth server listens on port %s", port)
	types.RegisterAuthServer(s, &server{svc: svc})
	return s.Serve(lis)
}

// AuthServer is the server API for Auth service.
type server struct {
	svc service.Interface
}

func (s *server) SendTransaction(ctx context.Context, req *types.SendTransactionReq) (*types.SendTransactionRes, error) {
	hash, err := s.svc.Auth(ctx, req.GetUntx())
	if err != nil {
		return nil, err
	}

	return &types.SendTransactionRes{
		TxHash: hash,
	}, nil
}

func (s *server) PackTransactions(ctx context.Context, req *types.PackTransactionsReq) (*types.PackTransactionsRes, error) {
	hashes, err := s.svc.PackTransactions(ctx, req.GetHeight())
	if err != nil {
		return nil, err
	}

	return &types.PackTransactionsRes{
		TxHashes: hashes,
	}, nil
}
