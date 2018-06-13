package service

import (
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/auth/service/logic"
	"github.com/yejiayu/go-cita/mq"
	"github.com/yejiayu/go-cita/types"
)

type MQ interface {
	Call(key mq.RoutingKey, data []byte) error
}

func NewMQ(queue mq.Queue) (MQ, error) {
	authLogic, err := logic.NewAuth()
	if err != nil {
		return nil, err
	}

	return &mqService{
		queue:     queue,
		authLogic: authLogic,
	}
}

type mqService struct {
	queue     mq.Queue
	authLogic logic.Auth
}

func (h *mqService) Call(key mq.RoutingKey, data []byte) error {
	switch key {
	case mq.NetworkUnverifiedTx:
		return h.untx(data)
	}
	return nil
}

func (h *mqService) untx(data []byte) error {
	var untx types.UnverifiedTransaction
	if err := proto.Unmarshal(data, &untx); err != nil {
		return err
	}

	if err := h.authLogic.Untx(&untx); err != nil {
		return err
	}

	return h.queue.Pub(mq.AuthUnverifiedTx, data)
}
