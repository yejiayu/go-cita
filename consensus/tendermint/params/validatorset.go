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
	"sort"

	"github.com/yejiayu/go-cita/common/hash"
)

// ValidatorSet represent a set of *Validator at a given height.
// The validators can be fetched by address or index.
// The index is in order of .Address, so the indices are fixed
// for all rounds of a given blockchain height.
// On the other hand, the .AccumPower of each validator and
// the designated .GetProposer() of a set changes every round,
// upon calling .IncrementAccum().
// NOTE: Not goroutine-safe.
// NOTE: All get/set to validators should copy the value for safety.
type ValidatorSet struct {
	// NOTE: persisted via reflect, must be exported.
	Validators []*Validator `json:"validators"`
	Proposer   *Validator   `json:"proposer"`

	// cached (unexported)
	totalVotingPower int64
}

func NewValidatorSet(vals []*Validator) *ValidatorSet {
	validators := make([]*Validator, len(vals))
	for i, val := range vals {
		validators[i] = val.Copy()
	}
	sort.Sort(ValidatorsSort(validators))
	vs := &ValidatorSet{
		Validators: validators,
	}

	if len(vs.Validators) > 0 {
		// vs.IncrementAccum(1)
	}

	return vs
}

// GetProposer returns the current proposer. If the validator set is empty, nil
// is returned.
func (vs *ValidatorSet) GetProposer(height, round uint64) *Validator {
	length := uint64(len(vs.Validators))
	if length == 0 {
		return nil
	}

	index := (height + round) % length
	return vs.Validators[index]
}

func (vs *ValidatorSet) GetByAddress(address hash.Address) *Validator {
	for _, val := range vs.Validators {
		if val.Address == address {
			return val
		}
	}

	return nil
}

//-------------------------------------
// Implements sort for sorting validators by address.

// Sort validators by address
type ValidatorsSort []*Validator

func (v ValidatorsSort) Len() int {
	return len(v)
}

func (v ValidatorsSort) Less(i, j int) bool {
	return v[i].ID < v[j].ID
}

func (v ValidatorsSort) Swap(i, j int) {
	it := v[i]
	v[i] = v[j]
	v[j] = it
}
