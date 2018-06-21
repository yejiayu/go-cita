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

	if err := h.srv.Untx(&untx); err != nil {
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
