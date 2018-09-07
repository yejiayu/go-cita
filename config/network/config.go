package network

import (
	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	Name string `env:"NAME" envDefault:"network"`
	Port string `env:"PORT" envDefault:"7002"`

	Address       string   `env:"ADDRESS" envDefault:"0xd92f2df9ab6bae68541e83cd38f22808f202363a"`
	NodeAddresses []string `env:"NODE_ADDRESSES" envSeparator:","`
	NodeURLs      []string `env:"NODE_URLS" envSeparator:","`

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

func GetName() string {
	return cfg.Name
}

func GetPort() string {
	return cfg.Port
}

func GetAddress() string {
	return cfg.Address
}

func GetNodeAddresses() []string {
	return cfg.NodeAddresses
}

func GetNodeURLs() []string {
	return cfg.NodeURLs
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
