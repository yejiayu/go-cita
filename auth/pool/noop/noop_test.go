package noop

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/yejiayu/go-cita/auth/pool"
	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

func TestAdd(t *testing.T) {
	p := getTestNoop(t)
	err := p.Add(context.Background(), &pb.SignedTransaction{
		TxHash: hash.BytesToSha3([]byte(fmt.Sprintf("test-%d", time.Now().Unix()))).Bytes(),
	})

	if err == pool.ErrTxDup {
		t.Fatal(err)
	}
	t.Log("add success")
}

func TestGet(t *testing.T) {
	pool := getTestNoop(t)
	if _, err := pool.Pull(context.Background(), 5); err != nil {
		t.Fatal(err)
	}
}

func TestDel(t *testing.T) {
	pool := getTestNoop(t)
	txs, err := pool.Pull(context.Background(), 50)
	if err != nil {
		t.Fatal(err)
	}

	hashes := make([]hash.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = hash.BytesToHash(tx.GetTxHash())
	}

	count, err := pool.Flush(context.Background(), hashes)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("delete %d tx", count)
}

func TestLen(t *testing.T) {
	pool := getTestNoop(t)
	len, err := pool.Len(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("tx pool length is %d", len)
}

func getTestNoop(t *testing.T) pool.Interface {
	return New(&mockTxDB{}, "127.0.0.1:6379", math.MaxUint32)
}

type mockTxDB struct {
}

func (mock *mockTxDB) Add(ctx context.Context, signedTx *pb.SignedTransaction) error {
	return nil
}

func (mock *mockTxDB) GetByHash(ctx context.Context, hash hash.Hash) (*pb.SignedTransaction, error) {
	return nil, nil
}

func (mock *mockTxDB) Exists(ctx context.Context, hash hash.Hash) (bool, error) {
	return false, nil
}
