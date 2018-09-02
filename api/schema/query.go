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

package schema

import (
	"github.com/graphql-go/graphql"

	"github.com/yejiayu/go-cita/api/resolvers"
	"github.com/yejiayu/go-cita/api/schema/types"
)

func NewQuery(r *resolvers.Resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"Hello": &graphql.Field{
				Type:    graphql.String,
				Resolve: r.Query.Hello,
			},
			"LatestHeight": &graphql.Field{
				Type:    types.Uint64,
				Resolve: r.Query.LatestHeight,
			},
			"Call": &graphql.Field{
				Type: types.Hex,
				Args: graphql.FieldConfigArgument{
					"height": &graphql.ArgumentConfig{Type: types.Uint64},
					"from":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.Hex)},
					"to":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.Hex)},
					"data":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.Hex)},
				},
				Resolve: r.Query.Call,
			},
			"GetReceipt": &graphql.Field{
				Type: types.ReceiptObject,
				Args: graphql.FieldConfigArgument{
					"txHash": &graphql.ArgumentConfig{Type: types.Hex},
				},
				Resolve: r.Query.GetReceipt,
			},
		},
	})
}
