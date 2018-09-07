package consensus

import (
	"crypto/ecdsa"

	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	PrivKeyHex string `env:"PRIVATE_KEY_HEX" envDefault:"0x3a91bade4b0b370f3958196ee7bdbdc86207f8470872a8159e30004a282c41c8"`
	QuotaLimit int    `env:"QUOTA_LIMIT" envDefault:"99999999"`
	TxCount    int    `env:"TX_COUNT" envDefault:"20000"`

	Name string `env:"NAME" envDefault:"consensus"`
	Port string `env:"PORT" envDefault:"8001"`

	AuthURL    string `env:"AUTH_URL" envDefault:"127.0.0.1:9001"`
	ChainURL   string `env:"CHAIN_URL" envDefault:"127.0.0.1:9002"`
	NetworkURL string `env:"NETWORK_URL" envDefault:"127.0.0.1:7002"`

	TracingURL string `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	log.Infof("consensus config %+v", cfg)
}

func GetName() string {
	return cfg.Name
}

func GetQuotaLimit() uint64 {
	return uint64(cfg.QuotaLimit)
}

func GetTxCount() uint32 {
	return uint32(cfg.TxCount)
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

func GetNetworkURL() string {
	return cfg.NetworkURL
}

func GetTracingURL() string {
	return cfg.TracingURL
}

func GetPrivKey() (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(cfg.PrivKeyHex)
}
