package server

import (
	"net"
)

type Server interface {
}

// NewServer returns tcp listener.
func NewServer(address string) (Server, error) {
	listener, err := net.Listen("TCP", address)
	if err != nil {
		return nil, err
	}

	return server{
		listener: listener,
	}, nil
}

type server struct {
	listener net.Listener
}
