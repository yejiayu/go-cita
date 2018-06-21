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
