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

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/auth/service"
	"github.com/yejiayu/go-cita/types"
)

func newServer(srv service.Service) *grpc.Server {
	s := grpc.NewServer()
	types.RegisterAuthServer(s, &server{srv: srv})
	return s
}

// AuthServer is the server API for Auth service.
type server struct {
	srv service.Service
}

func (s *server) SendTransaction(ctx context.Context, untx *types.UnverifiedTransaction) (*types.SendTransactionReply, error) {
	hash, err := s.srv.Untx(untx)
	if err != nil {
		return nil, err
	}

	return &types.SendTransactionReply{
		TxHash: hash.Bytes(),
	}, nil
}
