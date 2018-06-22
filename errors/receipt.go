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

package errors

var (
	// ExecutionError
	NotEnoughBaseGas        = new("NotEnoughBaseGas")
	BlockGasLimitReached    = new("BlockGasLimitReached")
	AccountGasLimitReached  = new("AccountGasLimitReached")
	InvalidTransactionNonce = new("InvalidTransactionNonce")
	NotEnoughCash           = new("NotEnoughCash")
	NoTransactionPermission = new("NoTransactionPermission")
	NoContractPermission    = new("NoContractPermission")
	NoCallPermission        = new("NoCallPermission")
	ExecutionInternal       = new("ExecutionInternal")
	TransactionMalformed    = new("TransactionMalformed")

	// EvmError
	OutOfGas                   = new("OutOfGas")
	BadJumpDestination         = new("BadJumpDestination")
	BadInstruction             = new("BadInstruction")
	StackUnderflow             = new("StackUnderflow")
	OutOfStack                 = new("OutOfStack")
	Internal                   = new("Internal")
	MutableCallInStaticContext = new("MutableCallInStaticContext")
	OutOfBounds                = new("OutOfBounds")
	Reverted                   = new("Reverted")
)
