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

package api

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/api/resolvers"
	"github.com/yejiayu/go-cita/api/schema"
	cfg "github.com/yejiayu/go-cita/config/api"
)

type Server interface {
	Run()
}

func NewServer(
	authClient pb.AuthClient,
	chainClient pb.ChainClient,
	vmClient pb.VMClient,
) Server {
	r := resolvers.New(authClient, chainClient, vmClient)

	schema, err := schema.NewSchema(r)
	if err != nil {
		log.Panic(err)
	}
	return &server{schema: &schema}
}

type server struct {
	schema *graphql.Schema
}

func (s *server) Run() {
	port := cfg.GetPort()

	h := handler.New(&handler.Config{
		Schema:   s.schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// serve HTTP
	http.Handle("/", h)

	log.Infof("The api server listens on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Panic(err)
	}
}
