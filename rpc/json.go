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

	"github.com/sourcegraph/jsonrpc2"
)

func New(port string) error {
	lis, err := net.Listen("TCP", port)
	if err != nil {
		return err
	}

	return serve(lis, nil)
}

func serve(lis net.Listener, opt ...jsonrpc2.ConnOpt) error {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}

		ctx := context.Background()
		jsonrpc2.NewConn(
			ctx,
			jsonrpc2.NewBufferedStream(conn, jsonrpc2.VarintObjectCodec{}),
			&handler{},
			opt...,
		)
	}
}
