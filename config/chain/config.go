package chain

import (
	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/log"
)

var cfg *config

type config struct {
	DbType string   `env:"DB_TYPE" envDefault:"redis"`
	DbURL  []string `env:"DB_URL" envSeparator:"," envDefault:"127.0.0.1:6379"`

	Name       string   `env:"NAME" envDefault:"chain"`
	Port       string   `env:"PORT" envDefault:"9002"`
	Validators []string `env:"VALIDATORS" envSeparator:","`
	VMURL      string   `env:"VM_URL" envDefault:"127.0.0.1:9003"`

	TracingURL string `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func init() {
	cfg = &config{}
	if err := env.Parse(cfg); err != nil {
		log.Panic(err)
	}

	log.Infof("The chain config is %+v", cfg)
}

func GetChain() string {
	return cfg.Name
}

func GetDbType() string {
	return cfg.DbType
}

func GetDbURL() []string {
	return cfg.DbURL
}

func GetPort() string {
	return cfg.Port
}

func GetValidators() []string {
	return cfg.Validators
}

func GetVMURL() string {
	return cfg.VMURL
}

func GetTracingURL() string {
	return cfg.TracingURL
}
