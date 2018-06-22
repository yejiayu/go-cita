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
	"github.com/yejiayu/go-cita/auth/service"
	"github.com/yejiayu/go-cita/mq"
)

func newRPC(queue mq.Queue, srv service.Service) *rpcHandler {
	return &rpcHandler{
		queue: queue,
		srv:   srv,
	}
}

type rpcHandler struct {
	queue mq.Queue
	srv   service.Service
}

func (h *rpcHandler) Call(key mq.RoutingKey, data []byte, out interface{}) error {
	return nil
}
