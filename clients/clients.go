package clients

import (
	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/types"
)

type Factory interface {
	Auth(url string) (types.AuthClient, error)
	Chain(url string) (types.ChainClient, error)
}

func New() Factory {
	return &factory{}
}

type factory struct{}

func (f *factory) Auth(url string) (types.AuthClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return types.NewAuthClient(conn), nil
}

func (f *factory) Chain(url string) (types.ChainClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return types.NewChainClient(conn), nil
}
