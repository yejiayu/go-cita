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
	"crypto/ecdsa"
	"errors"
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/yejiayu/go-cita/database"
	blockdb "github.com/yejiayu/go-cita/database/block"
	txdb "github.com/yejiayu/go-cita/database/tx"
	"github.com/yejiayu/go-cita/types"
)

type Service interface {
	Untx(untx *types.UnverifiedTransaction) (common.Hash, error)
	AddBlock(height uint64) error
}

const (
	nonceLenLimit        = 128
	validUntilBlockLimit = 100
)

type service struct {
	chainID         uint32
	blockQuotaLimit uint64

	historyTx *historyTx

	cache   *cache
	txDB    txdb.Interface
	blockDB blockdb.Interface
}

func NewService(dbFactory database.Factory) (Service, error) {
	cache, err := newCache()
	if err != nil {
		return nil, err
	}

	s := &service{
		cache:     cache,
		txDB:      dbFactory.TxDB(),
		blockDB:   dbFactory.BlockDB(),
		historyTx: newHistoryTx(cache),
	}

	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *service) init() error {
	h, err := s.blockDB.GetHeaderByLatest()
	if err != nil {
		return err
	}

	if h == nil {
		return nil
	}

	var step uint64
	for step < 100 || (h.Height-step) > 0 {
		height := h.Height - step
		b, err := s.blockDB.GetBodyByHeight(height)
		if err != nil {
			return err
		}

		if err := s.historyTx.AddTxs(height, b.GetTransactions()); err != nil {
			return err
		}
		step++
	}

	return nil
}

func (s *service) Untx(untx *types.UnverifiedTransaction) (common.Hash, error) {
	tx := untx.GetTransaction()
	data, err := proto.Marshal(untx.GetTransaction())
	if err != nil {
		return common.Hash{}, err
	}

	txHash := common.BytesToHash(data)
	pk, err := s.cache.getPublicKey(txHash)
	if err != nil {
		glog.Error(err)
	}
	if pk == nil {
		pk, err = s.verifyTxSig(txHash.Bytes(), untx.GetSignature(), untx.GetCrypto())
		if err != nil {
			glog.Error(err)
			return common.Hash{}, err
		}

		s.cache.setPublicKey(txHash, pk)
	}

	//TODO: black verify

	if err := s.checkTxParams(tx, pk, txHash); err != nil {
		return common.Hash{}, err
	}

	signTx := &types.SignedTransaction{
		TransactionWithSig: untx,
		TxHash:             txHash.Bytes(),
		Signer:             ethcrypto.CompressPubkey(pk),
	}

	if err := s.txDB.AddPool(signTx); err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s *service) AddBlock(height uint64) error {
	body, err := s.blockDB.GetBodyByHeight(height)
	if err != nil {
		return err
	}

	s.historyTx.AddTxs(height, body.GetTransactions())
	return nil
}

func (s *service) checkTxParams(tx *types.Transaction, signer *ecdsa.PublicKey, txHash common.Hash) error {
	if tx.ChainId != s.chainID {
		return errors.New("bad chain id")
	}

	if len(tx.Nonce) > nonceLenLimit {
		return errors.New("invalid nonce")
	}

	if err := s.checkValidUntilBlock(tx.ValidUntilBlock); err != nil {
		return err
	}

	if err := s.checkHistoryTxs(txHash); err != nil {
		return err
	}

	if err := s.checkQuota(tx.Quota, ethcrypto.PubkeyToAddress(*signer)); err != nil {
		return err
	}

	return nil
}

func (s *service) checkValidUntilBlock(validUntilBlock uint64) error {
	if validUntilBlock < s.historyTx.maxHeight+1 || validUntilBlock >= (s.historyTx.maxHeight+1+validUntilBlockLimit) {
		return errors.New("invalid until block")
	}

	return nil
}

func (s *service) checkHistoryTxs(hash common.Hash) error {
	exists, err := s.historyTx.Contains(hash)
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
