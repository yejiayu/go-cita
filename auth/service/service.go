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
	Untx(untx *types.UnverifiedTransaction) error
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

	return s, nil
}

func (s *service) Untx(untx *types.UnverifiedTransaction) error {
	tx := untx.GetTransaction()
	data, err := proto.Marshal(untx.GetTransaction())
	if err != nil {
		return err
	}

	txHash := common.BytesToHash(data)
	pk, err := s.cache.getPublicKey(txHash)
	if err != nil {
		glog.Error(err)
	}
	if pk == nil {
		pk, err = s.verifyTxSig(txHash.Bytes(), untx.GetSignature(), untx.GetCrypto())
		if err != nil {
			return err
		}

		s.cache.setPublicKey(txHash, pk)
	}

	//TODO: black verify

	if err := s.checkTxParams(tx, pk, txHash); err != nil {
		return err
	}

	signTx := &types.SignedTransaction{
		TransactionWithSig: untx,
		TxHash:             txHash.Bytes(),
		Signer:             ethcrypto.CompressPubkey(pk),
	}

	return s.txDB.AddPool(signTx)
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
	if quota > s.blockQuotaLimit {
		return errors.New("quota not enough")
	}
	return nil
}

func (s *service) verifyTxSig(hash, signature []byte, crypto types.Crypto) (*ecdsa.PublicKey, error) {
	switch crypto {
	case types.Crypto_SECP:
		return ethcrypto.SigToPub(hash, signature)
	}

	return nil, fmt.Errorf("%s is Unexpected crypto", crypto.String())
}
