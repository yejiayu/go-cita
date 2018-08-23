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
	"net"

	"github.com/yejiayu/go-cita/network/protocol"
)

func newConnection(id uint32, address string) (*connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &connection{
		id:      id,
		address: address,

		conn:  conn,
		codec: protocol.NewCodec(),
	}, nil
}

type connection struct {
	id      uint32
	address string

	conn  net.Conn
	codec protocol.Codec
}

func (c *connection) send(key string, data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *connection) close() error {
	return c.conn.Close()
}
