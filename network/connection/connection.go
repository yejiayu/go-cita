package connection

import (
	"fmt"
	"net"

	networkConfig "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/network/protocol"
)

func newConnection(id uint32, peer networkConfig.Peer) (*connection, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.IP, peer.Port))
	if err != nil {
		return nil, err
	}

	return &connection{id: id, conn: conn, peer: peer}, nil
}

type connection struct {
	id   uint32
	conn net.Conn
	peer networkConfig.Peer

	closeCh chan struct{}
}

func (c *connection) Loop() error {
	codec := protocol.NewCodec()

	for {
		key, data, err := codec.Decode(c.conn)
		if err != nil {
			return err
		}

		if err := c.handler(key, data); err != nil {
			return err
		}
	}
}

func (c *connection) close() error {
	defer close(c.closeCh)
	return nil
}

func (c *connection) handler(key string, data []byte) error {
	return nil
}
