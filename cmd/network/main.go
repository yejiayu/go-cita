package main

import (
	"flag"

	"github.com/golang/glog"
	networkConfig "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/mq"
	"github.com/yejiayu/go-cita/network"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	config := networkConfig.Config{
		ID:   1,
		Port: 8001,
	}

	subKeys := []mq.RoutingKey{
		mq.AuthUntx,
	}
	queue, err := mq.New("amqp://guest:guest@localhost:5672", "network", subKeys)
	if err != nil {
		glog.Fatal(err)
	}

	n, err := network.New(config, queue)
	if err != nil {
		glog.Fatal(err)
	}

	quit := make(chan error)
	n.Run(quit)

	glog.Info("network start")
	if err := <-quit; err != nil {
		glog.Fatal(err)
	}
}
