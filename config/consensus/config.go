package consensus

import (
	"crypto/ecdsa"

	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	DbURL []string `env:"DB_URL" envSeparator:"," envDefault:"47.75.129.215:2379,47.75.129.215:2380,47.75.129.215:2381"`

	PrivKeyHex string `env:"PRIVATE_KEY_HEX"`

	Port string `env:"PORT" envDefault:"8004"`

	AuthURL    string `env:"AUTH_URL" envDefault:"127.0.0.1:8001"`
	ChainURL   string `env:"CHAIN_URL" envDefault:"127.0.0.1:8003"`
	TracingURL string `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	log.Infof("consensus config %+v", cfg)
}

func GetDbURL() []string {
	return cfg.DbURL
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

func GetTracingURL() string {
	return cfg.TracingURL
}

func GetPrivKey() (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(cfg.PrivKeyHex)
}
