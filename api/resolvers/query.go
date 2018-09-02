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
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/graphql-go/graphql"
	ot "github.com/opentracing/opentracing-go"

	"github.com/yejiayu/go-cita/pb"
)

type Query struct {
	clients *clients
}

func (q *Query) Hello(p graphql.ResolveParams) (interface{}, error) {
	span, _ := ot.StartSpanFromContext(p.Context, "hello")
	defer span.Finish()

	return "hello", nil
}

func (q *Query) LatestHeight(p graphql.ResolveParams) (interface{}, error) {
	res, err := q.clients.chain.GetBlockHeader(p.Context, &pb.GetBlockHeaderReq{
		Height: math.MaxUint64,
	})
	if err != nil {
		return nil, err
	}
	return res.GetHeader().GetHeight(), nil
}

func (q *Query) Call(p graphql.ResolveParams) (interface{}, error) {
	height, ok := p.Args["height"].(uint64)
	if !ok {
		height = math.MaxUint64
	}
	from := p.Args["from"].(string)
	to := p.Args["to"].(string)
	data := p.Args["data"].(string)

	res, err := q.clients.vm.StaticCall(p.Context, &pb.StaticCallReq{
		Height: height,
		From:   common.FromHex(from),
		To:     common.FromHex(to),
		Data:   common.FromHex(data),
	})
	if err != nil {
		return nil, err
	}

	return common.ToHex(res.GetResult()), nil
}
