package pool

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

func TestAdd(t *testing.T) {
	pool := getTestNoop(t)
	err := pool.Add(context.Background(), &pb.SignedTransaction{
		TxHash: hash.BytesToSha3([]byte(fmt.Sprintf("test-%d", time.Now().Unix()))).Bytes(),
	})

	if err == errTxDup {
		t.Fatal(err)
	}
	t.Log("add success")
}

func TestGet(t *testing.T) {
	pool := getTestNoop(t)
	if _, err := pool.GetAll(context.Background(), 5); err != nil {
		t.Fatal(err)
	}
}

func TestDel(t *testing.T) {
	pool := getTestNoop(t)
	txs, err := pool.GetAll(context.Background(), 50000)
	if err != nil {
		t.Fatal(err)
	}

	hashes := make([]hash.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = hash.BytesToHash(tx.GetTxHash())
	}

	count, err := pool.Del(context.Background(), hashes)
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

func TestNoopPool(t *testing.T) {
	pool := getTestNoop(t)

	count := 100000
	for i := 0; i < count; i++ {
		tx := &pb.UnverifiedTransaction{
			Transaction: &pb.Transaction{
				To:              "0xffffffffffffffffffffffffffffff",
				Nonce:           fmt.Sprintf("%d", i),
				Quota:           uint64(i * 1),
				ValidUntilBlock: 100,
				Data:            []byte("0xffffffffffffffffffffffffffffffffffffffffffffffffff"),
				Value:           10,
				ChainId:         10,
				Version:         2,
			},
			Signature: []byte("0xf39b9183786e3347b48c1647fd6058c1832ddc9e3c909743fead6f4092dc47d825e6b20be4eb7e2a272b73c2e8a670719113a94f4e4989beb4e74492aa5b332000"),
			Crypto:    pb.Crypto_SECP,
		}
		txHash, err := hash.ProtoToSha3(tx)
		if err != nil {
			t.Fatal(err)
		}
		signedTx := &pb.SignedTransaction{
			TransactionWithSig: tx,
			TxHash:             txHash.Bytes(),
			Signer:             []byte("0xf39b9183786e3347b48c1647fd6058c1832ddc9e3c909743fead6f4092dc47d825e6b20be4eb7e2a272b73c2e8a670719113a94f4e4989beb4e74492aa5b332000"),
		}

		if err := pool.Add(context.Background(), signedTx); err != nil {
			t.Fatal(err)
		}
	}

	startTime := time.Now()
	txs, err := pool.GetAll(context.Background(), 200000)
	if err != nil {
		t.Fatal(err)
	}

	hashes := make([]hash.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = hash.BytesToHash(tx.GetTxHash())
	}

	delCount, err := pool.Del(context.Background(), hashes)
	if err != nil {
		t.Fatal(err)
	}
	endTime := time.Now()
	t.Logf("delete %d tx, exec %s", delCount, endTime.Sub(startTime).String())
}

func getTestNoop(t *testing.T) Interface {
	conn, err := redis.DialURL("redis://127.0.0.1:6379")
	if err != nil {
		t.Fatal(err)
	}

	return NewNoop(&mockTxDB{}, conn, math.MaxUint32)
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
