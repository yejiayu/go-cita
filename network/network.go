package network

import (
	networkConfig "github.com/yejiayu/go-cita/config/network"
	"github.com/yejiayu/go-cita/network/connection"
	"github.com/yejiayu/go-cita/network/server"
)

type Interface interface {
	Run(quit chan<- error)
}

func New(config networkConfig.Config) (Interface, error) {
	serve, err := server.New(config.Port)
	if err != nil {
		return nil, err
	}

	return &network{
		config:      config,
		connManager: connection.NewManager(config),
		server:      serve,
	}, nil
}

type network struct {
	config      networkConfig.Config
	connManager connection.Manager
	server      server.Interface
}

func (n *network) Run(quit chan<- error) {
	go n.connManager.Run(quit)
	go n.server.Run(quit)
}
