package clients

import (
	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"
)

func NewAuthClient(url string) pb.AuthClient {
	conn, err := create(url)
	if err != nil {
		log.Panic(err)
	}

	return pb.NewAuthClient(conn)
}

func NewChainClient(url string) pb.ChainClient {
	conn, err := create(url)
	if err != nil {
		log.Panic(err)
	}

	return pb.NewChainClient(conn)
}

func NewNetworkClient(url string) pb.NetworkClient {
	conn, err := create(url)
	if err != nil {
		log.Panic(err)
	}

	return pb.NewNetworkClient(conn)
}

func NewConsensusClient(url string) pb.ConsensusClient {
	conn, err := create(url)
	if err != nil {
		log.Panic(err)
	}

	return pb.NewConsensusClient(conn)
}

func create(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithInsecure())
}
