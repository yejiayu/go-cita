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
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/golang/glog"
	"github.com/yejiayu/go-cita/api"
)

var (
	serverPort = kingpin.Flag("port", "server port of api").Default("8080").String()

	authServer = kingpin.Flag("auth-server", "url of auth server").Default("127.0.0.1:9001").String()
)

func main() {
	kingpin.Parse()
	defer glog.Flush()

	if err := api.NewServer(
		*serverPort,
		*authServer,
	); err != nil {
		glog.Error(err)
	}
}
