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
	"github.com/yejiayu/go-cita/common/tracing"
	cfg "github.com/yejiayu/go-cita/config/consensus"
	"github.com/yejiayu/go-cita/log"

	"github.com/yejiayu/go-cita/consensus"
)

func main() {
	otClose, err := tracing.Configure("cita-consensus", cfg.GetTracingURL())
	if err != nil {
		log.Error(err)
	} else {
		defer otClose.Close()
	}

	server, err := consensus.New()
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
