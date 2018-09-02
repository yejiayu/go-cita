package loggger

import (
	"context"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/log"
)

func NewServer() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Infof("grpc: request method %s", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			log.Error(err)
		}
		return resp, err
	}
}
