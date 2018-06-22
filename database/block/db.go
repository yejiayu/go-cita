package block

import (
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/yejiayu/go-cita/database/raw"
	"github.com/yejiayu/go-cita/types"
)

var (
	headerPrefix = byte(10)
	bodyPrefix   = byte(11)
)

type Interface interface {
	GetHeaderByHeight(height uint64) (*types.BlockHeader, error)
	GetBodyByHeight(height uint64) (*types.BlockBody, error)
}

func New(rawDB raw.Interface) Interface {
	return &blockDB{rawDB: rawDB}
}

type blockDB struct {
	rawDB raw.Interface
}

func (db *blockDB) GetHeaderByHeight(height uint64) (*types.BlockHeader, error) {
	data, err := db.rawDB.Get(headerKey(height))
	if err != nil {
		return nil, err
	}

	var h types.BlockHeader
	if err := proto.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

func (db *blockDB) GetBodyByHeight(height uint64) (*types.BlockBody, error) {
	data, err := db.rawDB.Get(bodyKey(height))
	if err != nil {
		return nil, err
	}

	var body types.BlockBody
	if err := proto.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	return &body, nil
}

// func (db *blockDB) Scan(startHeight uint64, limit int) ([]*types.Block, error) {
// 	_, values, err := db.rawDB.Scan(blockKey(startHeight), limit)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var bs []*types.Block
// 	for _, value := range values {
// 		var b types.Block
// 		if err := proto.Unmarshal(value, &b); err != nil {
// 			return nil, err
// 		}
//
// 		bs = append(bs, &b)
// 	}
//
// 	return bs, nil
// }

func headerKey(height uint64) []byte {
	var key []byte
	binary.BigEndian.PutUint64(key, height)
	return joinKey(headerPrefix, key)
}

func bodyKey(height uint64) []byte {
	var key []byte
	binary.BigEndian.PutUint64(key, height)
	return joinKey(bodyPrefix, key)
}

func joinKey(prefix byte, key []byte) []byte {
	return append(append([]byte{prefix}, []byte(".")...), key...)
}
