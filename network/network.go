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

package network

import (
	"github.com/streadway/amqp"

	networkConfig "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/mq"
	"github.com/yejiayu/go-cita/network/connection"
	"github.com/yejiayu/go-cita/network/server"
)

type Interface interface {
	Run(quit chan<- error)
}

func New(config networkConfig.Config, queue mq.Queue) (Interface, error) {
	serve, err := server.New(config.Port)
	if err != nil {
		return nil, err
	}

	cm := connection.NewManager(config)
	return &network{
		config:      config,
		connManager: cm,
		server:      serve,
		syncHandler: newSynchronizer(cm),

		queue: queue,
	}, nil
}

type network struct {
	config      networkConfig.Config
	connManager connection.Manager
	server      server.Interface
	syncHandler *synchronizer

	queue mq.Queue
}

func (n *network) Run(quit chan<- error) {
	go n.connManager.Run(quit)
	go n.server.Run(quit)

	go n.handleServer(quit)
	go n.subQueue(quit)
}

func (n *network) handleServer(quit chan<- error) {
	for {
		m := n.server.Message()

		switch mq.RoutingKey(m.Key) {
		case mq.SyncUnverifiedTx:
			n.queue.Pub(mq.NetworkUnverifiedTx, m.Message.Payload())
		}
	}
}

func (n *network) subQueue(quit chan<- error) {
	delivery, err := n.queue.Sub()
	if err != nil {
		quit <- err
		return
	}

	for msg := range delivery {
		go n.handleMQ(&msg)
	}
}

func (n *network) handleMQ(msg *amqp.Delivery) {
	key := mq.RoutingKey(msg.RoutingKey)
	data := msg.Body

	switch key {
	case mq.AuthUnverifiedTx:
		n.connManager.Broadcast(mq.SyncUnverifiedTx, data)
	}
}
