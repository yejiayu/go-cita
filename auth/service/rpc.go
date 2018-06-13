package service

import (
	"github.com/yejiayu/go-cita/auth/service/logic"
	"github.com/yejiayu/go-cita/mq"
)

type RPC interface {
	Call(key mq.RoutingKey, data []byte, out interface{}) error
}

func NewRPC(queue mq.Queue) (RPC, error) {
	authLogic, err := logic.NewAuth()
	if err != nil {
		return nil, err
	}

	return &rpcService{
		queue:     queue,
		authLogic: authLogic,
	}
}

type rpcService struct {
	queue     mq.Queue
	authLogic logic.Auth
}

func (h *rpcService) Call(key mq.RoutingKey, data []byte, out interface{}) error {
	return nil
}
