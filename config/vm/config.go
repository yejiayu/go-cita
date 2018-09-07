package vm

import (
	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	DbType string   `env:"DB_TYPE" envDefault:"redis"`
	DbURL  []string `env:"DB_URL" envSeparator:"," envDefault:"127.0.0.1:6379"`

	ChainID     int    `env:"CHAIN_ID" envDefault:"1"`
	Name        string `env:"NAME" envDefault:"vm"`
	Port        string `env:"Port" envDefault:"9003"`
	GenesisPath string `env:"GENESIS_PATH" envDefault:"genesis.json"`

	TracingURL string `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	log.Infof("The network config %+v", cfg)
}

func GetChainID() uint32 {
	return uint32(cfg.ChainID)
}

func GetName() string {
	return cfg.Name
}

func GetDBType() string {
	return cfg.DbType
}

func GetDBUrl() []string {
	return cfg.DbURL
}

func GetPort() string {
	return cfg.Port
}

func GetGenesisPath() string {
	return cfg.GenesisPath
}

func GetTracingURL() string {
	return cfg.TracingURL
}
