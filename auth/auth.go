package auth

import (
	"github.com/golang/glog"
	"github.com/streadway/amqp"

	"github.com/yejiayu/go-cita/auth/service"
	"github.com/yejiayu/go-cita/mq"
)

type Interface interface {
	Run(quit chan<- error)
}

func New(queue mq.Queue) Interface {
	return &auth{
		queue:      queue,
		mqService:  service.NewMQ(queue),
		rpcService: service.NewRPC(queue),
	}
}

type auth struct {
	queue mq.Queue

	mqService  service.MQ
	rpcService service.RPC
}

func (a *auth) Run(quit chan<- error) {
	go a.loopMQ(quit)
}

func (a *auth) loopMQ(quit chan<- error) {
	delivery, err := a.queue.Sub()
	if err != nil {
		quit <- err
		return
	}

	for msg := range delivery {
		go func(msg *amqp.Delivery) {
			key := mq.RoutingKey(msg.RoutingKey)
			data := msg.Body

			if err := a.mqService.Call(key, data); err != nil {
				glog.Error(err)
				msg.Ack(false)
			} else {
				msg.Ack(true)
			}
		}(&msg)
	}
}
