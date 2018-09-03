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

type Receipt struct {
	QuotaUsed       uint64      `json:"quota_used" mapstructure:"quota_used"`
	Quota           uint64      `json:"quota" mapstructure:"quota"`
	LogBloom        string      `json:"log_bloom" mapstructure:"log_bloom"`
	Error           string      `json:"error" mapstructure:"error"`
	StateRoot       string      `json:"state_root" mapstructure:"state_root"`
	TransactionHash string      `json:"transaction_hash" mapstructure:"transaction_hash"`
	ContractAddress string      `json:"contract_address" mapstructure:"contract_address"`
	Logs            []*LogEntry `json:"logs" mapstructure:"logs"`
}

type LogEntry struct {
	Address string   `json:"address" mapstructure:"address"`
	Topics  []string `json:"topics" mapstructure:"topics"`
	Data    string   `json:"data" mapstructure:"data"`
}

var ReceiptObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Receipt",
	Fields: graphql.Fields{
		"quota_used":       {Type: Uint64},
		"quota":            {Type: Uint64},
		"log_bloom":        {Type: Hex},
		"error":            {Type: graphql.String},
		"state_root":       {Type: graphql.String},
		"transaction_hash": {Type: Hex},
		"contract_address": {Type: Hex},
		"logs":             {Type: graphql.NewList(LogEntryObject)},
	},
})

var LogEntryObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "LogEntry",
	Fields: graphql.Fields{
		"address": {Type: Address},
		"topics":  {Type: graphql.NewList(Hex)},
		"data":    {Type: Hex},
	},
})
