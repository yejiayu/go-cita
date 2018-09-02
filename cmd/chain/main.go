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
	"github.com/yejiayu/go-cita/clients"
	"github.com/yejiayu/go-cita/common/tracing"
	"github.com/yejiayu/go-cita/log"

	"github.com/yejiayu/go-cita/chain"
	cfg "github.com/yejiayu/go-cita/config/chain"
	"github.com/yejiayu/go-cita/database"
)

func main() {
	otClose, err := tracing.Configure("cita-chain", cfg.GetTracingURL())
	if err != nil {
		log.Error(err)
	} else {
		defer otClose.Close()
	}

	factory, err := database.NewFactory(cfg.GetDbType(), cfg.GetDbURL())
	if err != nil {
		log.Panic(err)
	}
	vmClient := clients.NewVMClient(cfg.GetVMURL())

	chain.New(factory, vmClient).Run()
}
