package block

import (
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/types"
)

var (
	blockPrefix = byte(10)
)

type Interface interface {
	GetByHeight(height uint64) (*types.Block, error)
	Scan(startHeight uint64, limit int) ([]*types.Block, error)
}

func New(rawDB raw.Interface) Interface {
	return &blockDB{rawDB: rawDB}
}

type blockDB struct {
	rawDB raw.Interface
}

func (db *blockDB) GetByHeight(height uint64) (*types.Block, error) {
	data, err := db.rawDB.Get(blockKey(height))
	if err != nil {
		return nil, err
	}

	var block types.Block
	if err := proto.Unmarshal(data, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func (db *blockDB) Scan(startHeight uint64, limit int) ([]*types.Block, error) {
	_, values, err := db.rawDB.Scan(blockKey(startHeight), limit)
	if err != nil {
		return nil, err
	}

	var bs []*types.Block
	for _, value := range values {
		var b types.Block
		if err := proto.Unmarshal(value, &b); err != nil {
			return nil, err
		}

		bs = append(bs, &b)
	}

	return bs, nil
}

func blockKey(height uint64) []byte {
	var key []byte
	binary.BigEndian.PutUint64(key, height)
	return append([]byte{blockPrefix}, key...)
}
