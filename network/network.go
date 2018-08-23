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
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/network/external"
	"github.com/yejiayu/go-cita/network/rpc"
)

type Interface interface {
	Run()
}

func New(
	consensusClient pb.ConsensusClient,
	authClient pb.AuthClient,
	chainClient pb.ChainClient,
) Interface {
	return &network{
		externalServer: external.New(consensusClient, authClient, chainClient),
		rpcServer:      rpc.New(),
	}
}

type network struct {
	externalServer external.Server
	rpcServer      external.Server
}

func (n *network) Run() {
	go n.externalServer.Run()
	n.rpcServer.Run()
}
