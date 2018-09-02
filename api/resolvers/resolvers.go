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
	"github.com/yejiayu/go-cita/pb"
)

type Resolver struct {
	Mutation *Mutation
	Query    *Query
}

type clients struct {
	auth  pb.AuthClient
	chain pb.ChainClient
	vm    pb.VMClient
}

func New(
	authClient pb.AuthClient,
	chainClient pb.ChainClient,
	vmClient pb.VMClient,
) *Resolver {
	cs := &clients{
		auth:  authClient,
		chain: chainClient,
		vm:    vmClient,
	}

	return &Resolver{
		Mutation: &Mutation{clients: cs},
		Query:    &Query{clients: cs},
	}
}
