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
