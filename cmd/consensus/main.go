package main

import (
	"flag"

	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/consensus"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/log"
)

type config struct {
	DbURL []string `env:"DB_URL" envSeparator:"," envDefault:"47.75.129.215:2379,47.75.129.215:2380,47.75.129.215:2381"`

	Port     string `env:"PORT" envDefault:"8004"`
	AuthURL  string `env:"AUTH_URL" envDefault:"127.0.0.1:8001"`
	ChainURL string `env:"CHAIN_URL" envDefault:"127.0.0.1:8003"`
}

func main() {
	flag.Parse()

	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("consensus config %+v", cfg)

	dbFactory, err := database.NewFactory(cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := consensus.New(
		cfg.Port,
		cfg.AuthURL,
		cfg.ChainURL,
		dbFactory,
	); err != nil {
		log.Fatal(err)
	}
}
