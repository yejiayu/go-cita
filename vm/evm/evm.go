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

package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// EVM provides a basic interface for the EVM calling conventions. The EVM EVM
// depends on this context being implemented for doing subcalls and initialising new EVM contracts.
type EVM interface {
	// Call another contract
	Call(caller vm.ContractRef, addr common.Address, data []byte, quota uint64, value *big.Int) ([]byte, uint64, error)
	// Take another's contract code and execute within our own context
	CallCode(caller vm.ContractRef, addr common.Address, data []byte, quota uint64, value *big.Int) ([]byte, uint64, error)
	// Same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(caller vm.ContractRef, addr common.Address, data []byte, quota uint64) ([]byte, uint64, error)
	// Create a new contract
	Create(caller vm.ContractRef, data []byte, quota uint64, value *big.Int) ([]byte, common.Address, uint64, error)
}

func New(ctx vm.Context, stateDB *state.StateDB, chainConfig *params.ChainConfig, vmConfig vm.Config) EVM {
	return &evm{
		vm: vm.NewEVM(ctx, stateDB, chainConfig, vmConfig),
	}
}

type evm struct {
	vm *vm.EVM
}

// Call another contract
func (e *evm) Call(caller vm.ContractRef, addr common.Address, data []byte, quota uint64, value *big.Int) ([]byte, uint64, error) {
	return e.vm.Call(caller, addr, data, quota, value)
}

// Take another's contract code and execute within our own context
func (e *evm) CallCode(caller vm.ContractRef, addr common.Address, data []byte, quota uint64, value *big.Int) ([]byte, uint64, error) {
	return e.vm.CallCode(caller, addr, data, quota, value)
}

// Same as CallCode except sender and value is propagated from parent to child scope
func (e *evm) DelegateCall(caller vm.ContractRef, addr common.Address, data []byte, quota uint64) ([]byte, uint64, error) {
	return e.vm.DelegateCall(caller, addr, data, quota)
}

// Create a new contract
func (e *evm) Create(caller vm.ContractRef, data []byte, quota uint64, value *big.Int) ([]byte, common.Address, uint64, error) {
	return e.vm.Create(caller, data, quota, value)
}
