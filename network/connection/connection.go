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
	"fmt"
	"net"

	networkConfig "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/mq"
	"github.com/yejiayu/go-cita/network/protocol"
)

func newConnection(id uint32, peer networkConfig.Peer) (*connection, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.IP, peer.Port))
	if err != nil {
		return nil, err
	}

	return &connection{
		id:    id,
		conn:  conn,
		peer:  peer,
		codec: protocol.NewCodec(),
	}, nil
}

type connection struct {
	id   uint32
	conn net.Conn
	peer networkConfig.Peer

	codec protocol.Codec
}

func (c *connection) send(key mq.RoutingKey, data []byte) error {
	encodedDdata, err := c.codec.Encode(key, data)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(encodedDdata)
	return err
}

func (c *connection) close() error {
	return c.conn.Close()
}
