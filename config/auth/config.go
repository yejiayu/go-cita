package auth

import (
	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	DbType   string   `env:"DB_TYPE" envDefault:"redis"`
	DbURL    []string `env:"DB_URL" envSeparator:"," envDefault:"127.0.0.1:6379"`
	RedisURL string   `env:"REDIS_URL" envDefault:"127.0.0.1:6379"`

	Name      string `env:"NAME" envDefault:"auth"`
	Port      string `env:"PORT" envDefault:"9001"`
	PoolCount int    `env:"POOL_COUNT" envDefault:"50000"`
	ChainID   int    `env:"CHAIN_ID" envDefault:"1"`

	NetworkURL string `env:"NETWORK_URL" envDefault:"127.0.0.1:7002"`

	TracingURL string `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		log.Panic(err)
	}

	log.Infof("The auth config is %+v", cfg)
}

func GetName() string {
	return cfg.Name
}

func GetDbType() string {
	return cfg.DbType
}

func GetDbURL() []string {
	return cfg.DbURL
}

func GetRedisURL() string {
	return cfg.RedisURL
}

func GetPort() string {
	return cfg.Port
}

func GetNetworkURL() string {
	return cfg.NetworkURL
}

func GetPoolCount() uint32 {
	return uint32(cfg.PoolCount)
}

func GetChainID() uint32 {
	return uint32(cfg.ChainID)
}

func GetTracingURL() string {
	return cfg.TracingURL
}
