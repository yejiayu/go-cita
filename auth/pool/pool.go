package pool

import (
	"context"
	"errors"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

var (
	ErrTxDup          = errors.New("tx dup")
	ErrPoolReachLimit = errors.New("tx pool reach the limit")
)

type Interface interface {
	// Add a signed transaction, and the transaction does not exist on the pool and block
	Add(ctx context.Context, signedTx *pb.SignedTransaction) error

	// Pulling the "count" number of signed transactions from the pool, but does not remove the transaction from the pool
	Pull(ctx context.Context, count uint32) ([]*pb.SignedTransaction, error)

	// Returns the transaction associated with hash.
	// For every transaction that does not exist in the pool, a nil value is returned
	Get(ctx context.Context, hashes []hash.Hash) ([]*pb.SignedTransaction, error)

	// Remove the transaction from the pool and no extra action
	Del(ctx context.Context, hashes []hash.Hash) (int64, error)

	// Flush pool, in addition to removing transactions from the pool may do something extra
	Flush(ctx context.Context, hashes []hash.Hash) (int64, error)

	Len(ctx context.Context) (int64, error)
}
