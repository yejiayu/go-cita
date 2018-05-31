package base

import (
	dbConfig "github.com/yejiayu/go-cita/config/db"
)

type Interface interface {
	// // Helper to create a new transaction.
	// Transaction() *gorocksdb.TransactionDB
	// Get(col uint32, key []byte) (common.Hash, error)
	// GetByPrefix(col uint32, prefix []byte)([]byte, error)
	// WriteBuffered(tx gorocksdb.TransactionDB)
	// Write()
	// Flush()
	// Iter()
	// IterFromPrefix()
	// Restore()
	//   /// Get a value by key.
	//   fn get(&self, col: Option<u32>, key: &[u8]) -> Result<Option<DBValue>, String>;
	//
	//   /// Get a value by partial key. Only works for flushed data.
	//   fn get_by_prefix(&self, col: Option<u32>, prefix: &[u8]) -> Option<Box<[u8]>>;
	//
	//   /// Write a transaction of changes to the buffer.
	//   fn write_buffered(&self, transaction: DBTransaction);
	//
	//   /// Write a transaction of changes to the backing store.
	//   fn write(&self, transaction: DBTransaction) -> Result<(), String> {
	//       self.write_buffered(transaction);
	//       self.flush()
	//   }
	//
	//   /// Flush all buffered data.
	//   fn flush(&self) -> Result<(), String>;
	//
	//   /// Iterate over flushed data for a given column.
	//   fn iter<'a>(&'a self, col: Option<u32>) -> Box<Iterator<Item = (Box<[u8]>, Box<[u8]>)> + 'a>;
	//
	//   /// Iterate over flushed data for a given column, starting from a given prefix.
	//   fn iter_from_prefix<'a>(&'a self, col: Option<u32>, prefix: &'a [u8]) -> Box<Iterator<Item = (Box<[u8]>, Box<[u8]>)> + 'a>;
	//
	//   /// Attempt to replace this database with a new one located at the given path.
	//   fn restore(&self, new_db: &str) -> Result<(), UtilError>;
}

func New(config dbConfig.Config) Interface {
	return &baseDB{config: config}
}

type baseDB struct {
	config dbConfig.Config
}
