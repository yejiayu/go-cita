package api

import (
	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	Name string `env:"NAME" envDefault:"api-graphql"`
	Port string `env:"PORT" envDefault:"8000"`

	AuthURL  string `env:"AUTH_URL" envDefault:"127.0.0.1:9001"`
	ChainURL string `env:"CHAIN_URL" envDefault:"127.0.0.1:9002"`
	VMURL    string `env:"VM_URL" envDefault:"127.0.0.1:9003"`

	TracingURL string `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		log.Panic(err)
	}

	log.Infof("The api config is %+v", cfg)
}

func GetName() string {
	return cfg.Name
}

func GetPort() string {
	return cfg.Port
}

func GetAuthURL() string {
	return cfg.AuthURL
}

func GetChainURL() string {
	return cfg.ChainURL
}

func GetVMURL() string {
	return cfg.VMURL
}

func GetTracingURL() string {
	return cfg.TracingURL
}
