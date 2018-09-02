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

package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/graphql-go/graphql"
)

type Transaction struct {
	To              common.Address `json:"to" mapstructure:"to"`
	Nonce           string         `json:"nonce" mapstructure:"nonce"`
	Quota           uint64         `json:"quota" mapstructure:"quota"`
	ValidUntilBlock uint64         `json:"valid_until_block" mapstructure:"valid_until_block"`
	Data            string         `json:"data" mapstructure:"data"`
	Value           string         `json:"value" mapstructure:"value"`
	ChainID         uint32         `json:"chain_id" mapstructure:"chain_id"`
	Version         uint32         `json:"version" mapstructure:"version"`
}

var TransactionObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Transaction",
	Fields: graphql.Fields{
		"to":                {Type: Address},
		"nonce":             {Type: graphql.String},
		"quota":             {Type: graphql.NewNonNull(Uint64)},
		"valid_until_block": {Type: Uint64},
		"data":              {Type: Hex},
		"value":             {Type: Hex},
		"chain_id":          {Type: graphql.NewNonNull(Uint32)},
		"version":           {Type: Uint32},
	},
})

var CryptoEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "Crypto",
	Values: graphql.EnumValueConfigMap{
		"SECP": &graphql.EnumValueConfig{Value: 0},
		"SM2":  &graphql.EnumValueConfig{Value: 1},
	},
})

type UnverifiedTransaction struct {
	Transaction Transaction `json:"transaction" mapstructure:"transaction"`
	Crypto      int32       `json:"crypto" mapstructure:"crypto"`
	Signature   string      `json:"signature" mapstructure:"signature"`
}

type UnsafeTransaction struct {
	Transaction Transaction `json:"transaction" mapstructure:"transaction"`
	Crypto      int32       `json:"crypto" mapstructure:"crypto"`
	PrivateKey  string      `json:"privateKey" mapstructure:"privateKey"`
}

var TransactionInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "TransactionInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"to":                {Type: Address},
		"nonce":             {Type: graphql.String},
		"quota":             {Type: graphql.NewNonNull(Uint64)},
		"valid_until_block": {Type: Uint64},
		"data":              {Type: Hex},
		"value":             {Type: Hex},
		"chain_id":          {Type: graphql.NewNonNull(Uint32)},
		"version":           {Type: Uint32},
	},
})
