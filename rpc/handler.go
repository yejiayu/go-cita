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
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
)

var methodMap = map[string]handleMethod{
	"net_peerCount":             nil,
	"cita_blockNumber":          nil,
	"cita_sendTransaction":      nil,
	"cita_getBlockByHash":       nil,
	"cita_getBlockByNumber":     nil,
	"eth_getTransactionReceipt": nil,
	"eth_getLogs":               nil,
	"eth_call":                  nil,
	"cita_getTransaction":       nil,
	"eth_getTransactionCount":   nil,
	"eth_getCode":               nil,
	"eth_getAbi":                nil,
	"eth_getBalance":            nil,
	"eth_newFilter":             nil,
	"eth_newBlockFilter":        nil,
	"eth_uninstallFilter":       nil,
	"eth_getFilterChanges":      nil,
	"eth_getFilterLogs":         nil,
	"cita_getTransactionProof":  nil,
	"cita_getMetaData":          nil,
}

type handleMethod func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) error
type handler struct{}

func (h *handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Notif {
		return // handler notification
	}

	m, ok := methodMap[req.Method]
	if !ok {
		conn.Reply(ctx, req.ID, fmt.Sprintf("method %s is not found", req.Method))
		return
	}

	if err := m(ctx, conn, req); err != nil {
		conn.Reply(ctx, req.ID, err.Error())
	}
}
