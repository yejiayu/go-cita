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
)

// Volatile state for each Validator
// NOTE: The Accum is not included in Validator.Hash();
// make sure to update that method if changes are made here
type Validator struct {
	ID      uint32           `json:"id"`
	Address hash.Address     `json:"address"`
	PubKey  *ecdsa.PublicKey `json:"pub_key"`

	Accum int64 `json:"accum"`
}

func NewValidator(id uint32, pubKey *ecdsa.PublicKey) *Validator {
	return &Validator{
		ID:      id,
		Address: crypto.PubkeyToAddress(*pubKey),
		PubKey:  pubKey,
		Accum:   0,
	}
}

// Creates a new copy of the validator so we can mutate accum.
// Panics if the validator is nil.
func (v *Validator) Copy() *Validator {
	vCopy := *v
	return &vCopy
}

func (v *Validator) VerifySignature(hash hash.Hash, signature []byte) bool {
	return crypto.VerifySignature(v.PubKey, hash, signature[:len(signature)-1])
}
