package benchmark

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/auth/pool"
	"github.com/yejiayu/go-cita/auth/pool/noop"
	"github.com/yejiayu/go-cita/auth/pool/quota"
)

func TestPool(t *testing.T) {
	noopPool := noop.New(&mockTxDB{}, "127.0.0.1:6379", 200000)
	quotaPool := quota.New(&mockTxDB{}, "127.0.0.1:6379", 200000)

	t.Log("-----------exec noop start---------")
	exec(noopPool, t)
	t.Log("-----------exec noop end---------")

	t.Log("-----------exec quota start---------")
	exec(quotaPool, t)
	t.Log("-----------exec quota end---------")
}

func exec(pool pool.Interface, t *testing.T) {
	count := 100000
	mockData := make([]*pb.SignedTransaction, count)
	t.Logf("building mock data, count %d", count)
	mockStart := time.Now()

	for i := 0; i < count; i++ {
		tx := &pb.UnverifiedTransaction{
			Transaction: &pb.Transaction{
				To:              "0xffffffffffffffffffffffffffffff",
				Nonce:           fmt.Sprintf("%d-%d", i, time.Now().UnixNano()),
				Quota:           10,
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
		mockData[i] = signedTx
	}

	// exec noop pool
	ctx := context.Background()
	var wg sync.WaitGroup
	for _, tx := range mockData {
		wg.Add(1)
		go func(wg *sync.WaitGroup, tx *pb.SignedTransaction) {
			if err := pool.Add(ctx, tx); err != nil {
				t.Fatal(err)
			}
			wg.Done()
		}(&wg, tx)
	}
	wg.Wait()

	mockEnd := time.Now()
	t.Logf("mock data completion, %s", mockEnd.Sub(mockStart))
	startTime := time.Now()

	t.Logf("Pull %d transaction", 100000)
	poolStart := time.Now()
	txs, err := pool.Pull(context.Background(), 50000)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Pull end, %s", time.Now().Sub(poolStart))
	hashes := make([]hash.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = hash.BytesToHash(tx.GetTxHash())
	}

	t.Log("Flush transaction")
	flushStart := time.Now()
	delCount, err := pool.Flush(context.Background(), hashes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Flush end %s, remove %d", time.Now().Sub(flushStart), delCount)

	endTime := time.Now()
	t.Logf("Total execution %s", endTime.Sub(startTime).String())
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

func buildData(t *testing.T, count int) []*pb.SignedTransaction {
	txs := make([]*pb.SignedTransaction, count)

	for i := 0; i < count; i++ {
		tx := &pb.UnverifiedTransaction{
			Transaction: &pb.Transaction{
				To:              "0xffffffffffffffffffffffffffffff",
				Nonce:           fmt.Sprintf("%d-%d", i, time.Now().UnixNano()),
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

		txs[i] = signedTx
	}

	return txs
}
