package consensus

import (
	"net"

	"google.golang.org/grpc"

	"github.com/yejiayu/go-cita/clients"
	"github.com/yejiayu/go-cita/consensus/pbft"
	"github.com/yejiayu/go-cita/database"
	blockdb "github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/log"
)

func New(
	port, authURL, chainURL string,
	dbFactory database.Factory,
) error {
	cFactory := clients.New()
	authClient, err := cFactory.Auth(authURL)
	if err != nil {
		return err
	}
	chainClient, err := cFactory.Chain(chainURL)
	if err != nil {
		return err
	}

	go func() {
		err = pbft.New(dbFactory, authClient, chainClient).Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	s := grpc.NewServer()

	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return err
	}

	log.Infof("The consensus server listens on port %s", port)
	// types.RegisterConsensusServer(s, &server{})
	grpc.Dial(authURL, grpc.WithInsecure())
	return s.Serve(lis)
}

type server struct {
	blockDB blockdb.Interface
}
