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

package connection

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"

	"github.com/yejiayu/go-cita/log"

	cfg "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/network/protocol"
)

type Manager interface {
	// Run()
	// UpdateConfig(config networkConfig.Config)

	// Available() []networkConfig.Peer

	Broadcast(key string, protoMsg proto.Message) error
	// Single(key mq.RoutingKey, origin uint32, data []byte)
	// Subtract(key mq.RoutingKey, origin uint32, data []byte)
}

func NewManager() Manager {
	peers := make(map[uint32]string)
	conns := make(map[uint32]*connection)

	for _, peer := range cfg.GetPeers() {
		ss := strings.Split(peer, "-")
		if len(ss) != 2 {
			log.Panic("config peer invalid")
		}
		idStr := ss[0]
		address := ss[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Panic(err)
		}

		peers[uint32(id)] = address
		conn, err := newConnection(uint32(id), address)
		if err != nil {
			log.Error(err)
			continue
		}

		conns[uint32(id)] = conn
	}

	m := &manager{
		id:    cfg.GetID(),
		peers: peers,
		conns: conns,
	}
	go func() {
		for {
			m.retry()
		}
	}()
	return m
}

type manager struct {
	mu sync.RWMutex

	id    uint32
	peers map[uint32]string
	conns map[uint32]*connection
}

func (m *manager) Broadcast(key string, protoMsg proto.Message) error {
	data, err := proto.Marshal(protoMsg)
	if err != nil {
		return err
	}

	msg := protocol.NewMessage(protocol.OpTypeBroadcast, m.id, data)

	for peer, conn := range m.conns {
		err := conn.send(key, msg.Raw())
		if err != nil {
			glog.Errorf("network: the message send to %-v failed, %s\n", peer, err.Error())
		}
	}

	return nil
}

func (m *manager) retry() {
	m.mu.Lock()
	for id, address := range m.peers {
		conn, ok := m.conns[id]
		if !ok || conn == nil {
			conn, err := newConnection(id, address)
			if err != nil {
				log.Errorf("%s, retry", err)
				continue
			}

			m.conns[id] = conn
		}
	}
	m.mu.Unlock()
	time.Sleep(time.Second * 2)
}

// func (m *manager) Single(key mq.RoutingKey, origin uint32, data []byte) {
// 	sendMsg := protocol.NewMessage(protocol.OpTypeSingle, m.id, data)
//
// 	for peer, conn := range m.conns {
// 		if peer.ID == origin {
// 			err := conn.send(key, sendMsg.Raw())
// 			if err != nil {
// 				glog.Errorf("network: the message send to %-v failed, %s\n", peer, err.Error())
// 			}
//
// 			return
// 		}
// 	}
// }
//
// func (m *manager) Subtract(key mq.RoutingKey, origin uint32, data []byte) {
// 	sendMsg := protocol.NewMessage(protocol.OpTypeSubtract, m.id, data)
//
// 	for peer, conn := range m.conns {
// 		if peer.ID != origin {
// 			err := conn.send(key, sendMsg.Raw())
// 			if err != nil {
// 				glog.Errorf("network: the message send to %-v failed, %s\n", peer, err.Error())
// 			}
//
// 			return
// 		}
// 	}
// }
//
// func (m *manager) checkReadied() {
// 	for unReadiedPeer := range m.unReadyCh {
// 		time.Sleep(1 * time.Second)
//
// 		m.mu.Lock()
// 		_, ok := m.peers[unReadiedPeer]
// 		if !ok {
// 			m.mu.Unlock()
// 			continue
// 		}
//
// 		conn, err := newConnection(m.id, unReadiedPeer)
// 		if err != nil {
// 			m.mu.Unlock()
// 			m.unReadyCh <- unReadiedPeer
// 			continue
// 		}
//
// 		m.conns[unReadiedPeer] = conn
// 		m.mu.Unlock()
// 	}
// }
