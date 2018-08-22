// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"context"
	"crypto/ecdsa"
	"errors"

	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/auth/cache"
	"github.com/yejiayu/go-cita/auth/pool"
	"github.com/yejiayu/go-cita/database"
	blockdb "github.com/yejiayu/go-cita/database/block"
	txdb "github.com/yejiayu/go-cita/database/tx"
)

var (
	errBadChainID    = errors.New("bad chain id")
	errInvalidNonce  = errors.New("invalid nonce")
	errInvalidSig    = errors.New("invalid sig")
	errMoreThanQuota = errors.New("more than the quota")
	errInvalidHashes = errors.New("invalid tx hashes")
)

type Interface interface {
	AddUnverifyTx(ctx context.Context, untx *pb.UnverifiedTransaction) (hash.Hash, error)
	EnsureFromPool(ctx context.Context, nodeID uint32, quotaUsed uint64, hashes []hash.Hash) error
	GetHashFromPool(ctx context.Context, count uint32, quotaLimit uint64) ([]hash.Hash, error)
	ClearPool(ctx context.Context, height uint64) error

	EnsureHashes(ctx context.Context, nodeID uint32, hashes []hash.Hash) ([]hash.Hash, error)
}

const (
	nonceLenLimit = 128
)

type service struct {
	chainID    uint32
	quotaLimit uint64

	cache cache.Interface
	pool  pool.Interface

	networkClient pb.NetworkClient

	txDB    txdb.Interface
	blockDB blockdb.Interface
}

func New(
	dbFactory database.Factory,
	networkClient pb.NetworkClient,
	cache cache.Interface,
	pool pool.Interface,
) Interface {
	s := &service{
		cache: cache,
		pool:  pool,

		networkClient: networkClient,

		txDB:    dbFactory.TxDB(),
		blockDB: dbFactory.BlockDB(),
	}

	return s
}

func (s *service) AddUnverifyTx(ctx context.Context, untx *pb.UnverifiedTransaction) (hash.Hash, error) {
	tx := untx.GetTransaction()
	data, err := proto.Marshal(untx.GetTransaction())
	if err != nil {
		return hash.Hash{}, err
	}
	txHash := hash.BytesToSha3(data)

	//TODO: black verify
	pubkey, err := s.verifySig(ctx, txHash, untx.GetSignature(), untx.GetCrypto())
	if err != nil || pubkey == nil {
		return hash.Hash{}, errInvalidSig
	}
	if err := s.verifyTransaction(ctx, tx); err != nil {
		return hash.Hash{}, err
	}

	signedTx := &pb.SignedTransaction{
		TransactionWithSig: untx,
		TxHash:             txHash.Bytes(),
		Signer:             crypto.PubkeyToAddress(*pubkey).Bytes(),
	}

	if err := s.pool.Add(ctx, signedTx); err != nil {
		return hash.Hash{}, err
	}

	go s.networkClient.BroadcastTransaction(ctx, &pb.BroadcastTransactionReq{Untx: untx})
	return txHash, nil
}

func (s *service) EnsureFromPool(ctx context.Context, nodeID uint32, quotaUsed uint64, hashes []hash.Hash) error {
	if len(hashes) == 0 {
		return nil
	}

	signedTxs, err := s.pool.Get(ctx, hashes)
	if err != nil {
		return err
	}

	// If these transactions do not exist in the pool, from the target node to pull
	nonexistentHashMap := make(map[hash.Hash]bool)
	var quotaCount uint64
	for i, signedTx := range signedTxs {
		if signedTx == nil {
			// The pool return value is the same as the incoming order
			nonexistentHashMap[hashes[i]] = false
			continue
		}

		quotaUsed += signedTx.GetTransactionWithSig().GetTransaction().GetQuota()
		if quotaCount > quotaUsed {
			return errMoreThanQuota
		}
	}

	res, err := s.networkClient.GetUnverifyTxs(ctx, &pb.GetUnverifyTxsReq{
		NodeID: nodeID,
		TxHashes: func() [][]byte {
			var hashes [][]byte
			for txHash := range nonexistentHashMap {
				hashes = append(hashes, txHash.Bytes())
			}
			return hashes
		}(),
	})
	if err != nil {
		return err
	}

	// Verify these transactions
	untxs := res.GetTxs()
	if len(untxs) != len(nonexistentHashMap) {
		return errInvalidHashes
	}
	for _, untx := range untxs {
		txHash, err := s.AddUnverifyTx(ctx, untx)
		if err != nil {
			return err
		}

		delete(nonexistentHashMap, txHash)
	}
	if len(nonexistentHashMap) != 0 {
		return errInvalidHashes
	}
	return nil
}

func (s *service) GetHashFromPool(ctx context.Context, count uint32, quotaLimit uint64) ([]hash.Hash, error) {
	signedTxs, err := s.pool.GetAll(ctx, count)
	if err != nil {
		return nil, err
	}

	var quotaCount uint64
	var invalidHashes []hash.Hash
	var validHashes []hash.Hash
	for _, signedTx := range signedTxs {
		tx := signedTx.GetTransactionWithSig().GetTransaction()
		txHash := hash.BytesToHash(signedTx.GetTxHash())

		// // check valid_until_block
		// if tx.GetValidUntilBlock() < header.GetHeight() {
		// 	invalidHashes = append(invalidHashes, txHash)
		// 	continue
		// }

		quotaCount += tx.GetQuota()
		// check quota limit
		if quotaCount > quotaLimit {
			// restore quota
			quotaCount -= tx.GetQuota()
			continue
		}

		validHashes = append(validHashes, txHash)
	}

	// delete invali tx
	if _, err := s.pool.Del(ctx, invalidHashes); err != nil {
		return nil, err
	}

	return validHashes, nil
}

func (s *service) ClearPool(ctx context.Context, height uint64) error {
	body, err := s.blockDB.GetBodyByHeight(ctx, height)
	if err != nil {
		return err
	}

	hashes := make([]hash.Hash, len(body.GetTxHashes()))
	for i, h := range body.GetTxHashes() {
		hashes[i] = hash.BytesToHash(h)
	}

	_, err = s.pool.Del(ctx, hashes)
	return err
}

func (s *service) EnsureHashes(ctx context.Context, nodeID uint32, hashes []hash.Hash) ([]hash.Hash, error) {
	var a []hash.Hash
	for _, hash := range hashes {
		exists, err := s.txDB.Exists(ctx, hash)
		if err != nil {
			return nil, err
		}

		if !exists {
			a = append(a, hash)
		}
	}

	return nil, nil
}

func (s *service) verifySig(ctx context.Context, txHash hash.Hash, signature []byte, cryptoType pb.Crypto) (*ecdsa.PublicKey, error) {
	pubkey, err := s.cache.GetPublicKey(ctx, txHash)
	if err != nil {
		log.Error(err)
	}
	if pubkey == nil {
		switch cryptoType {
		case pb.Crypto_SECP:
			pubkey, err = crypto.SigToPub(txHash, signature)
			if err != nil {
				return nil, err
			}
		}

		if err = s.cache.SetPublicKey(ctx, txHash, pubkey); err != nil {
			return nil, err
		}
	}

	return pubkey, err
}

func (s *service) verifyTransaction(ctx context.Context, tx *pb.Transaction) error {
	if tx.GetChainId() != s.chainID {
		return errBadChainID
	}
	if len(tx.GetNonce()) > nonceLenLimit {
		return errInvalidNonce
	}
	return nil
}

// // TODO: check quota
// func (s *service) checkQuota(quota uint64, address common.Address) error {
// 	// if quota > s.blockQuotaLimit {
// 	// 	return errors.New("quota not enough")
// 	// }
// 	return nil
// }
