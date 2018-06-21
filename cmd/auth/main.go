package main

import (
	"github.com/golang/glog"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/yejiayu/go-cita/auth"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/mq"
)

var (
	dbURLs = kingpin.Flag("db-url", "url of tikv").Required().Strings()
	mqURL  = kingpin.Flag("mq-url", "url of rabbitmq").Default("amqp://guest:guest@localhost:5672").String()
)

func main() {
	kingpin.Parse()
	defer glog.Flush()

	subKeys := []mq.RoutingKey{
		mq.NetworkUnverifiedTx,
	}

	queue, err := mq.New(*mqURL, "auth", subKeys)
	if err != nil {
		glog.Fatal(err)
	}

	dbFactory, err := database.NewFactory(*dbURLs)
	if err != nil {
		glog.Fatal(err)
	}

	quit := make(chan error)
	a, err := auth.New(queue, dbFactory)
	if err != nil {
		glog.Fatal(err)
	}

	a.Run(quit)
	glog.Info("auth start")
	if err := <-quit; err != nil {
		glog.Fatal(err)
	}
}
