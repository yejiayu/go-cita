package wal

import (
	"context"
	"testing"
)

func TestFileWal(t *testing.T) {
	wal, err := NewFileWAL("file_wal")
	if err != nil {
		t.Fatal(err)
	}

	if err := wal.SetHeight(context.Background(), 10); err != nil {
		t.Fatal(err)
	}

	if err := wal.Save(context.Background(), LogTypeVote, []byte("string")); err != nil {
		t.Fatal(err)
	}

	logType, data, err := wal.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%d", logType)
	t.Log(string(data))
}
