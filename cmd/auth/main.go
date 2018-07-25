// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"

	"github.com/caarlos0/env"

	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/tools/tracing"

	"github.com/yejiayu/go-cita/auth"
	"github.com/yejiayu/go-cita/database"
)

type config struct {
	Port       string   `env:"PORT" envDefault:"8001"`
	DbURL      []string `env:"DB_URL" envSeparator:"," envDefault:"47.75.129.215:2379,47.75.129.215:2380,47.75.129.215:2381"`
	RedisURL   string   `env:"REDIS_URL" envDefault:"127.0.0.1:6379"`
	TracingURL string   `env:"TRACING_URL" envDefault:"zipkin.istio-system:9411"`
}

func main() {
	flag.Parse()

	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("env config %+v", cfg)

	otClose, err := tracing.Configure("cita-auth", cfg.TracingURL)
	if err != nil {
		log.Error(err)
	} else {
		defer otClose.Close()
	}

	dbFactory, err := database.NewFactory(cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := auth.New(cfg.Port, cfg.RedisURL, dbFactory); err != nil {
		log.Fatal(err)
	}
}
