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

package server

import (
	"fmt"
	"io"
	"net"

	"github.com/golang/glog"

	"github.com/yejiayu/go-cita/network/protocol"
)

type PeerMessage struct {
	Key     string
	Message protocol.Message
}

type Interface interface {
	Run(quit chan<- error)
	Message() *PeerMessage
}

// New returns tcp listener.
func New(port uint32) (Interface, error) {
	address := fmt.Sprintf("0.0.0.0:%d", port)
	glog.Infof("network server listen address is %s", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &server{
		listener: listener,

		MessageCh: make(chan *PeerMessage, 10),
	}, nil
}

type server struct {
	listener net.Listener

	MessageCh chan *PeerMessage
}

func (s *server) Run(quit chan<- error) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			quit <- err
			return
		}

		glog.Infof("accept conn from %s", conn.RemoteAddr().String())
		go func() {
			if err := s.handler(conn); err != nil {
				if err != io.EOF {
					glog.Error(err)
					conn.Close()
				}
			}
		}()
	}
}

func (s *server) Message() *PeerMessage {
	return <-s.MessageCh
}

func (s *server) handler(conn net.Conn) error {
	codec := protocol.NewCodec()

	for {
		key, data, err := codec.Decode(conn)
		if err != nil {
			return err
		}

		s.MessageCh <- &PeerMessage{
			Key:     key,
			Message: protocol.NewMessageWithRaw(data),
		}
	}
}
