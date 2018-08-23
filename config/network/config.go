package network

import (
	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	RPCPort      string `env:"RPC_PORT" envDefault:"7001"`
	ExternalPort string `env:"EXTERNAL_PORT" envDefault:"7000"`
	ID           uint32 `env:"ID"`

	Peers []string `env:"PEERS" envSeparator:"," envDefault:"1-127.0.0.1:7101"`

	ConsensusURL string `env:"CONSENSUS_URL" envDefault:"127.0.0.1:8001"`
	AuthURL      string `env:"AUTH_URL" envDefault:"127.0.0.1:9001"`
	ChainURL     string `env:"CHAIN_URL" envDefault:"127.0.0.1:9002"`

	TracingURL string `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	log.Infof("The network config %+v", cfg)
}

func GetRPCPort() string {
	return cfg.RPCPort
}

func GetExternalPort() string {
	return cfg.ExternalPort
}

func GetID() uint32 {
	return cfg.ID
}

func GetPeers() []string {
	return cfg.Peers
}

func GetConsensusURL() string {
	return cfg.ConsensusURL
}

func GetAuthURL() string {
	return cfg.AuthURL
}

func GetChainURL() string {
	return cfg.ChainURL
}

func GetTracingURL() string {
	return cfg.TracingURL
}
