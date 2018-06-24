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

package resolvers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"

	types "github.com/yejiayu/go-cita/api/schema/types"
	pbTypes "github.com/yejiayu/go-cita/types"
)

type Mutation struct {
	clients *clients
}

func (m *Mutation) SendTransaction(p graphql.ResolveParams) (interface{}, error) {
	var untx types.UnverifiedTransaction
	if err := mapstructure.Decode(p.Args, &untx); err != nil {
		return nil, err
	}

	to := ""
	if (untx.Transaction.To != common.Address{}) {
		to = untx.Transaction.To.Hex()
	}

	data := []byte{}
	if len(untx.Transaction.Data) > 2 {
		data = common.Hex2Bytes(untx.Transaction.Data[2:])
	}

	reply, err := m.clients.auth.SendTransaction(p.Context, &pbTypes.UnverifiedTransaction{
		Signature: common.Hex2Bytes(untx.Signature[2:]),
		Crypto:    pbTypes.Crypto(untx.Crypto),
		Transaction: &pbTypes.Transaction{
			To:              to,
			Nonce:           untx.Transaction.Nonce,
			Quota:           untx.Transaction.Quota,
			ValidUntilBlock: untx.Transaction.ValidUntilBlock,
			Data:            data,
			Value:           untx.Transaction.Value,
			ChainId:         untx.Transaction.ChainID,
			Version:         untx.Transaction.Version,
		},
	})

	if err != nil {
		return nil, err
	}

	return common.BytesToHash(reply.GetTxHash()).Hex(), nil
}
