package rpc

import (
	"context"
	"net"

	"github.com/sourcegraph/jsonrpc2"
)

func New(port string) error {
	lis, err := net.Listen("TCP", port)
	if err != nil {
		return err
	}

	return serve(lis, nil)
}

func serve(lis net.Listener, opt ...jsonrpc2.ConnOpt) error {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}

		ctx := context.Background()
		jsonrpc2.NewConn(
			ctx,
			jsonrpc2.NewBufferedStream(conn, jsonrpc2.VarintObjectCodec{}),
			&handler{},
			opt...,
		)
	}
}
