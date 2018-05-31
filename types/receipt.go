package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yejiayu/go-cita/errors"
)

// Receipt is Information describing execution of a transaction.
type Receipt struct {
	// The state root after executing the transaction. Optional since EIP98
	StateRoot common.Hash
	// The total gas used in the block following execution of the transaction.
	GasUsed *big.Int
	// The OR-wide combination of all logs' blooms for this transaction.
	// pub log_bloom: LogBloom,
	// The logs stemming from this transaction.
	// pub logs: Vec<LogEntry>,
	// For calculating contract address
	Contract common.Address
	// Transaction hash.
	TransactionHash common.Hash
	// Transaction transact error
	ErrorMsg errors.ReceiptError
}
