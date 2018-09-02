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

func NewMutation(r *resolvers.Resolver) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"SendTransaction": &graphql.Field{
				Type: types.Hash,
				Args: graphql.FieldConfigArgument{
					"transaction": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.TransactionInput)},
					"signature":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.Hex)},
					"crypto":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.CryptoEnum)},
				},
				Resolve: r.Mutation.SendTransaction,
			},
			"SendRawTransaction": &graphql.Field{
				Type: types.Hash,
				Args: graphql.FieldConfigArgument{
					"data": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.Hex)},
				},
				Resolve: r.Mutation.SendRawTransaction,
			},
			"SendUnsafeTransaction": &graphql.Field{
				Type:              types.Hash,
				DeprecationReason: "Warning!! This will give away your private key and should never be used in a production environment",
				Args: graphql.FieldConfigArgument{
					"transaction": &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.TransactionInput)},
					"privateKey":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.Hex)},
					"crypto":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(types.CryptoEnum)},
				},
				Resolve: r.Mutation.SendUnsafeTransaction,
			},
		},
	})
}
