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

package resolvers

import (
	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/types"
)

type Resolver struct {
	Mutation *Mutation
	Query    *Query
}

type clients struct {
	auth  types.AuthClient
	chain types.ChainClient
}

func New(authClient, chainClient string) (*Resolver, error) {
	conn, err := grpc.Dial(authClient, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	auth := types.NewAuthClient(conn)

	conn, err = grpc.Dial(chainClient, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	chain := types.NewChainClient(conn)

	cs := &clients{
		auth:  auth,
		chain: chain,
	}

	return &Resolver{
		Mutation: &Mutation{clients: cs},
		Query:    &Query{clients: cs},
	}, nil
}
