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

	"github.com/golang/glog"
	"github.com/graphql-go/handler"
	"github.com/yejiayu/go-cita/log"

	"github.com/yejiayu/go-cita/api/resolvers"
	"github.com/yejiayu/go-cita/api/schema"
)

func New(port, authServer string) error {
	r, err := resolvers.New(authServer)
	if err != nil {
		return err
	}

	schema, err := schema.NewSchema(r)
	if err != nil {
		glog.Fatal(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// serve HTTP
	http.Handle("/", h)

	log.Infof("The api server listens on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}
