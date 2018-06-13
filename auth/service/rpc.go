package service

import (
	"github.com/yejiayu/go-cita/auth/service/logic"
	"github.com/yejiayu/go-cita/mq"
)

type RPC interface {
	Call(key mq.RoutingKey, data []byte, out interface{}) error
}

func NewRPC(queue mq.Queue) RPC {
	return &rpcService{
		queue: queue,
		logic: logic.NewAuth(),
	}
}

type rpcService struct {
	queue mq.Queue
	logic logic.Auth
}

func (h *rpcService) Call(key mq.RoutingKey, data []byte, out interface{}) error {
	return nil
}
