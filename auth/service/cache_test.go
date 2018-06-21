package service

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestAdd(t *testing.T) {
	cache, err := newCache()
	if err != nil {
		t.Fatal(err)
	}

	if err := cache.setHistoryTx(1, []common.Hash{common.Hash{}}); err != nil {
		t.Fatal(err)
	}
}
