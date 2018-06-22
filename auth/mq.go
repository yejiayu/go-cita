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
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/auth/service"
	"github.com/yejiayu/go-cita/mq"
	"github.com/yejiayu/go-cita/types"
)

func newMQ(queue mq.Queue, srv service.Service) *mqHandler {
	return &mqHandler{
		queue: queue,
		srv:   srv,
	}
}

type mqHandler struct {
	queue mq.Queue
	srv   service.Service
}

func (h *mqHandler) Call(key mq.RoutingKey, data []byte) error {
	switch key {
	case mq.NetworkUnverifiedTx:
		return h.untx(data)
	case mq.ChainSyncStatus:
		return h.syncStatus(data)
	}
	return nil
}

func (h *mqHandler) untx(data []byte) error {
	var untx types.UnverifiedTransaction
	if err := proto.Unmarshal(data, &untx); err != nil {
		return err
	}

	if _, err := h.srv.Untx(&untx); err != nil {
		return err
	}

	return h.queue.Pub(mq.AuthUnverifiedTx, data)
}

func (h *mqHandler) syncStatus(data []byte) error {
	var status types.Status
	if err := proto.Unmarshal(data, &status); err != nil {
		return err
	}

	return h.srv.AddBlock(status.GetHeight())
}
