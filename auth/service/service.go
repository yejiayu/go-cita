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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/database"
	blockdb "github.com/yejiayu/go-cita/database/block"
	txdb "github.com/yejiayu/go-cita/database/tx"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/types"
)

type Interface interface {
	Auth(ctx context.Context, untx *types.UnverifiedTransaction) ([]byte, error)
	PackTransactions(ctx context.Context, height uint64) ([][]byte, error)
}

const (
	nonceLenLimit        = 128
	validUntilBlockLimit = 100

	packTransactionsLimit int = 10000
)

type service struct {
	chainID         uint32
	blockQuotaLimit uint64

	cache   *cache
	txDB    txdb.Interface
	blockDB blockdb.Interface
}

func New(redisURL string, dbFactory database.Factory) (Interface, error) {
	cache, err := newCache(redisURL)
	if err != nil {
		return nil, err
	}

	s := &service{
		cache:   cache,
		txDB:    dbFactory.TxDB(),
		blockDB: dbFactory.BlockDB(),
	}

	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *service) init() error {
	return nil
}

func (s *service) Auth(ctx context.Context, untx *types.UnverifiedTransaction) ([]byte, error) {
	tx := untx.GetTransaction()

	data, err := proto.Marshal(untx.GetTransaction())
	if err != nil {
		return nil, err
	}

	h := sha3.New256()
	h.Write(data)
	txHash := common.BytesToHash(h.Sum(nil))

	log.Infof("tx hash is", txHash.String())
	pk, err := s.cache.getPublicKey(ctx, txHash)
	if err != nil {
		log.Error(err)
	}
	if pk == nil {
		pk, err = s.verifyTxSig(txHash.Bytes(), untx.GetSignature(), untx.GetCrypto())
		if err != nil {
			return nil, err
		}

		s.cache.setPublicKey(ctx, txHash, pk)
	}

	//TODO: black verify

	if err := s.checkTxParams(ctx, tx, pk, txHash); err != nil {
		return nil, err
	}

	signTx := &types.SignedTransaction{
		TransactionWithSig: untx,
		TxHash:             txHash.Bytes(),
		Signer:             ethcrypto.CompressPubkey(pk),
	}

	if err := s.txDB.AddPool(ctx, signTx); err != nil {
		return nil, err
	}

	return txHash.Bytes(), nil
}

func (s *service) PackTransactions(ctx context.Context, height uint64) ([][]byte, error) {
	body, err := s.blockDB.GetBodyByHeight(ctx, height)
	if err != nil {
		return nil, err
	}

	hashes := body.GetTxHashes()
	if len(hashes) > 0 {
		if err := s.txDB.UpdateTx(ctx, hashes); err != nil {
			return nil, err
		}
	}

	return s.txDB.GetTxhashesFromPool(ctx, packTransactionsLimit)
}

func (s *service) checkTxParams(ctx context.Context, tx *types.Transaction, signer *ecdsa.PublicKey, txHash common.Hash) error {
	if tx.ChainId != s.chainID {
		return errors.New("bad chain id")
	}

	if len(tx.Nonce) > nonceLenLimit {
		return errors.New("invalid nonce")
	}

	if err := s.checkValidUntilBlock(tx.ValidUntilBlock); err != nil {
		return err
	}

	if err := s.checkHistoryTxs(ctx, txHash); err != nil {
		return err
	}

	if err := s.checkQuota(tx.Quota, ethcrypto.PubkeyToAddress(*signer)); err != nil {
		return err
	}

	return nil
}

func (s *service) checkValidUntilBlock(validUntilBlock uint64) error {
	latestHeader, err := s.blockDB.GetHeaderByLatest(context.Background())
	if err != nil {
		return err
	}

	latestHeight := latestHeader.GetHeight()
	if validUntilBlock < latestHeight+1 || validUntilBlock >= (latestHeight+1+validUntilBlockLimit) {
		return errors.New("invalid until block")
	}

	return nil
}

func (s *service) checkHistoryTxs(ctx context.Context, hash common.Hash) error {
	exists, err := s.txDB.Exists(ctx, hash)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("tx dup")
	}
	return nil
}

// TODO: check quota
func (s *service) checkQuota(quota uint64, address common.Address) error {
	// if quota > s.blockQuotaLimit {
	// 	return errors.New("quota not enough")
	// }
	return nil
}

func (s *service) verifyTxSig(hash, signature []byte, crypto types.Crypto) (*ecdsa.PublicKey, error) {
	switch crypto {
	case types.Crypto_SECP:
		return ethcrypto.SigToPub(hash, signature)
	}

	return nil, fmt.Errorf("%s is Unexpected crypto", crypto.String())
}
