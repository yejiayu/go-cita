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
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/log"

	cfg "github.com/yejiayu/go-cita/config/vm"
	"github.com/yejiayu/go-cita/vm"
)

func main() {
	file, err := os.Open(cfg.GetGenesisPath())
	if err != nil {
		log.Panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}

	var genesis vm.Genesis
	if err = json.Unmarshal(data, &genesis); err != nil {
		log.Panic(err)
	}

	factory, err := database.NewFactory(cfg.GetDBType(), cfg.GetDBUrl())
	if err != nil {
		log.Panic(err)
	}

	if err := vm.SetupGenesis(factory, &genesis, false); err != nil {
		log.Panic(err)
	}

	vm.NewServer(factory).Run()
}
