package clients

import (
	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"
)

func NewAuthClient(url string) pb.AuthClient {
	conn, err := grpc.Dial(url)
	if err != nil {
		log.Panic(err)
	}

	return pb.NewAuthClient(conn)
}

func NewChainClient(url string) pb.ChainClient {
	conn, err := grpc.Dial(url)
	if err != nil {
		log.Panic(err)
	}

	return pb.NewChainClient(conn)
}

func NewNetworkClient(url string) pb.NetworkClient {
	conn, err := grpc.Dial(url)
	if err != nil {
		log.Panic(err)
	}

	return pb.NewNetworkClient(conn)
}
