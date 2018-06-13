package main

import (
	"github.com/golang/glog"
	"github.com/yejiayu/go-cita/mq"
)

func main() {
	subKeys := []mq.RoutingKey{
		mq.NetworkUntx,
	}

	queue, err := mq.New("amqp://guest:guest@localhost:5672", "auth", subKeys)
	if err != nil {
		glog.Fatal(err)
	}

}
