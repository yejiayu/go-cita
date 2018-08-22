package pool

import (
	"context"
	"errors"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

var (
	errTxDup          = errors.New("tx dup")
	errPoolReachLimit = errors.New("tx pool reach the limit")
)

const (
	poolKey = "tx.pool"
)

type Interface interface {
	Add(ctx context.Context, signedTx *pb.SignedTransaction) error

	GetAll(ctx context.Context, count uint32) ([]*pb.SignedTransaction, error)
	Get(ctx context.Context, hashes []hash.Hash) ([]*pb.SignedTransaction, error)

	Del(ctx context.Context, hashes []hash.Hash) (int64, error)

	Len(ctx context.Context) (int64, error)
	Exists(ctx context.Context, hash hash.Hash) (bool, error)
}
