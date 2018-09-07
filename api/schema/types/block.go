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
	"github.com/graphql-go/graphql"
)

type BlockHeader struct {
	Prevhash        string `json:"prevhash" mapstructure:"prevhash"`
	Timestamp       uint64 `json:"timestamp" mapstructure:"timestamp"`
	Height          uint64 `json:"height" mapstructure:"height"`
	StateRoot       string `json:"state_root" mapstructure:"state_root"`
	TransactionRoot string `json:"transactions_root" mapstructure:"transactions_root"`
	ReceiptsRoot    string `json:"receipts_root" mapstructure:"receipts_root"`
	QuotaUsed       uint64 `json:"quota_used" mapstructure:"quota_used"`
	QuotaLimit      uint64 `json:"quota_limit" mapstructure:"quota_limit"`
	Proposer        string `json:"proposer" mapstructure:"proposer"`
}

type BlockBody struct {
	TxHashes []string `json:"tx_hashes"`
}

var BlockHeaderObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "BlockHeader",
	Fields: graphql.Fields{
		"Prevhash":        {Type: graphql.String},
		"Timestamp":       {Type: Uint64},
		"Height":          {Type: Uint64},
		"StateRoot":       {Type: graphql.String},
		"TransactionRoot": {Type: graphql.String},
		"ReceiptsRoot":    {Type: graphql.String},
		"QuotaUsed":       {Type: Uint64},
		"QuotaLimit":      {Type: Uint64},
		"Proposer":        {Type: graphql.String},
	},
})

var BlockBodyObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "BlockBody",
	Fields: graphql.Fields{
		"TxHashes": {Type: graphql.NewList(graphql.String)},
	},
})
