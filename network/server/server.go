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
