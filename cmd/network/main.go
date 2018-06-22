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
