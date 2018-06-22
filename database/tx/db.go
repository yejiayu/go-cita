package tx

import (
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/types"
)

var (
	txPoolPrefix = byte(20)
)

type Interface interface {
	AddPool(signedTx *types.SignedTransaction) error
}

func New(rawDB raw.Interface) Interface {
	return &txDB{rawDB: rawDB}
}

type txDB struct {
	rawDB raw.Interface
}

func (db *txDB) AddPool(signedTx *types.SignedTransaction) error {
	data, err := proto.Marshal(signedTx)
	if err != nil {
		return err
	}

	key := txPoolKey(signedTx.GetTxHash())
	return db.rawDB.Put(key, data)
}

func txPoolKey(key []byte) []byte {
	return joinKey(txPoolPrefix, key)
}

func joinKey(prefix byte, key []byte) []byte {
	return append(append([]byte{prefix}, []byte(".")...), key...)
}
