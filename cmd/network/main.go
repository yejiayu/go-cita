package main

import (
	"flag"

	"github.com/golang/glog"
	networkConfig "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/network"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	config := networkConfig.Config{
		ID:   1,
		Port: 8001,
	}

	n, err := network.New(config)
	if err != nil {
		glog.Fatal(err)
	}

	quit := make(chan error)
	n.Run(quit)

	if err := <-quit; err != nil {
		glog.Fatal(err)
	}
}
