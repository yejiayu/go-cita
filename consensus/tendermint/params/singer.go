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

package params

import (
	"crypto/ecdsa"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

type Singer struct {
	privKey *ecdsa.PrivateKey
	address hash.Address
}

func NewSinger(privKey *ecdsa.PrivateKey) *Singer {
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	return &Singer{
		privKey: privKey,
		address: address,
	}
}

func (s *Singer) SignVote(vote *pb.Vote) ([]byte, error) {
	data, err := hash.ProtoToSha3(vote)
	if err != nil {
		return nil, err
	}

	return s.sign(data)
}

func (s *Singer) SignBlock(block *pb.Block) ([]byte, error) {
	data, err := hash.ProtoToSha3(block)
	if err != nil {
		return nil, err
	}

	return s.sign(data)
}

func (s *Singer) Address() hash.Address {
	return s.address
}

func (s *Singer) sign(data hash.Hash) ([]byte, error) {
	return crypto.Sign(data, s.privKey)
}
