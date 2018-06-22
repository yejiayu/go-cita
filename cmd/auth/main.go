// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
