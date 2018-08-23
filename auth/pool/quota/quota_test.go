package quota

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/auth/pool"
)

func TestListAdd(t *testing.T) {
	p := getTestQuota(t)
	err := p.Add(context.Background(), &pb.SignedTransaction{
		TxHash: hash.BytesToSha3([]byte(fmt.Sprintf("test-%d", time.Now().UnixNano()))).Bytes(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("add success")
}

func TestListGetAll(t *testing.T) {
	p := getTestQuota(t)
	txs, err := p.Pull(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(txs))
}

func TestListDel(t *testing.T) {
	p := getTestQuota(t)
	txs, err := p.Pull(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}

	hashes := make([]hash.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = hash.BytesToHash(tx.GetTxHash())
	}

	count, err := p.Flush(context.Background(), hashes)
	t.Log(count, err)
}

func getTestQuota(t *testing.T) pool.Interface {
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
