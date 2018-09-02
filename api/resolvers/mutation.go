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
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/protobuf/proto"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"

	types "github.com/yejiayu/go-cita/api/schema/types"
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

	req := &pb.AddUnverifyTxReq{
		Untx: &pb.UnverifiedTransaction{
			Signature: common.Hex2Bytes(untx.Signature[2:]),
			Crypto:    pb.Crypto(untx.Crypto),
			Transaction: &pb.Transaction{
				To:              to,
				Nonce:           untx.Transaction.Nonce,
				Quota:           untx.Transaction.Quota,
				ValidUntilBlock: untx.Transaction.ValidUntilBlock,
				Data:            data,
				Value:           common.FromHex(untx.Transaction.Value),
				ChainId:         untx.Transaction.ChainID,
				Version:         untx.Transaction.Version,
			},
		},
	}
	reply, err := m.clients.auth.AddUnverifyTx(p.Context, req)

	if err != nil {
		return nil, err
	}

	return common.BytesToHash(reply.GetTxHash()).Hex(), nil
}

func (m *Mutation) SendRawTransaction(p graphql.ResolveParams) (interface{}, error) {
	data := p.Args["data"]
	protoData := common.FromHex(data.(string))
	var untx pb.UnverifiedTransaction
	if err := proto.Unmarshal(protoData, &untx); err != nil {
		return nil, err
	}

	reply, err := m.clients.auth.AddUnverifyTx(p.Context, &pb.AddUnverifyTxReq{
		Untx: &untx,
	})
	if err != nil {
		return nil, err
	}

	return common.BytesToHash(reply.GetTxHash()).Hex(), nil
}

func (m *Mutation) SendUnsafeTransaction(p graphql.ResolveParams) (interface{}, error) {
	var unsafeTx types.UnsafeTransaction
	if err := mapstructure.Decode(p.Args, &unsafeTx); err != nil {
		return nil, err
	}

	pbTx := &pb.Transaction{
		To:              unsafeTx.Transaction.To.String(),
		Nonce:           unsafeTx.Transaction.Nonce,
		Quota:           unsafeTx.Transaction.Quota,
		ValidUntilBlock: unsafeTx.Transaction.ValidUntilBlock,
		Data:            common.FromHex(unsafeTx.Transaction.Data),
		Value:           common.FromHex(unsafeTx.Transaction.Value),
		ChainId:         unsafeTx.Transaction.ChainID,
		Version:         unsafeTx.Transaction.Version,
	}

	if pbTx.Nonce == "" {
		pbTx.Nonce = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if pbTx.ValidUntilBlock == 0 {
		res, err := m.clients.chain.GetBlockHeader(p.Context, &pb.GetBlockHeaderReq{
			Height: math.MaxUint64,
		})
		if err != nil {
			return nil, err
		}
		pbTx.ValidUntilBlock = res.GetHeader().GetHeight() + 88
	}
	if pbTx.Version == 0 {
		pbTx.Version = 2
	}

	sha3Data, err := hash.ProtoToSha3(pbTx)
	if err != nil {
		return nil, err
	}

	privatekey, err := crypto.HexToECDSA(unsafeTx.PrivateKey)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(sha3Data, privatekey)
	if err != nil {
		return nil, err
	}

	untx := &pb.UnverifiedTransaction{
		Transaction: pbTx,
		Signature:   signature,
		Crypto:      pb.Crypto(unsafeTx.Crypto),
	}

	res, err := m.clients.auth.AddUnverifyTx(p.Context, &pb.AddUnverifyTxReq{
		Untx: untx,
	})
	if err != nil {
		return nil, err
	}

	return common.ToHex(res.GetTxHash()), nil
}
