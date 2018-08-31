package ethdb

import (
	"context"
	"sync"

	eth "github.com/ethereum/go-ethereum/ethdb"

	"github.com/yejiayu/go-cita/log"

	"github.com/yejiayu/go-cita/database/raw"
)

// Code using batches should try to add this much data to the batch.
// The value was determined empirically.
const (
	ethDBNamespace = "eth.db"
	IdealBatchSize = 100 * 1024
)

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(key []byte, value []byte) error
}

// Deleter wraps the database delete operation supported by both batches and regular databases.
type Deleter interface {
	Delete(key []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	Putter
	Deleter
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Close()
	NewBatch() eth.Batch
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type Batch interface {
	Putter
	Deleter
	ValueSize() int // amount of data in the batch
	Write() error
	// Reset resets the batch for reuse
	Reset()
}

func New(raw raw.Interface) Database {
	return &ethDB{raw: raw}
}

type ethDB struct {
	raw raw.Interface
}

func (db *ethDB) Get(key []byte) ([]byte, error) {
	return db.raw.Get(context.TODO(), ethDBNamespace, key)
}

func (db *ethDB) Has(key []byte) (bool, error) {
	data, err := db.raw.Get(context.TODO(), ethDBNamespace, key)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func (db *ethDB) Put(key []byte, value []byte) error {
	return db.raw.Put(context.TODO(), ethDBNamespace, key, value)
}

func (db *ethDB) Delete(key []byte) error {
	return db.raw.Delete(context.TODO(), ethDBNamespace, key)
}

func (db *ethDB) Close() {
	if err := db.raw.Close(); err != nil {
		log.Errorf("close ethdb %s", err)
	}
}

func (db *ethDB) NewBatch() eth.Batch {
	return &batch{raw: db.raw}
}

type batch struct {
	raw raw.Interface

	mu   sync.RWMutex
	size int

	putKeys    [][]byte
	putValues  [][]byte
	deleteKeys [][]byte
}

func (b *batch) Put(key []byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.putKeys = append(b.putKeys, key)
	b.putValues = append(b.putValues, value)
	b.size += len(value)
	return nil
}

func (b *batch) Delete(key []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.deleteKeys = append(b.deleteKeys, key)
	b.size++
	return nil
}

func (b *batch) ValueSize() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.size
}

func (b *batch) Write() error {
	if len(b.putKeys) > 0 && len(b.putKeys) == len(b.putValues) {
		if err := b.raw.BatchPut(context.TODO(), ethDBNamespace, b.putKeys, b.putValues); err != nil {
			return err
		}
	}

	if len(b.deleteKeys) > 0 {
		return b.raw.BatchDelete(context.TODO(), ethDBNamespace, b.deleteKeys)
	}

	return nil
}

func (b *batch) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.putKeys = nil
	b.putValues = nil
	b.deleteKeys = nil
	b.size = 0
}
