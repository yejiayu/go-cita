package logic

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"sync"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/yejiayu/go-cita/auth/pool"
	"github.com/yejiayu/go-cita/types"
)

type Auth interface {
	Untx(untx *types.UnverifiedTransaction) error
}

const (
	nonceLenLimit        = 128
	validUntilBlockLimit = 100
)

type auth struct {
	mu sync.Mutex

	chainID         uint32
	blockQuotaLimit uint64

	cache  *cache
	txPool pool.TxPool
}

func NewAuth() (Auth, error) {
	cache, err := newCache()
	if err != nil {
		return nil, err
	}

	return &auth{
		txPool: pool.NewTxPool(),
		cache:  cache,
	}, nil
}

func (a *auth) Untx(untx *types.UnverifiedTransaction) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	tx := untx.GetTransaction()
	data, err := proto.Marshal(untx.GetTransaction())
	if err != nil {
		return err
	}

	txHash := common.BytesToHash(data)
	pk, err := a.cache.getPublicKey(txHash)
	if err != nil {
		glog.Error(err)
	}
	if pk == nil {
		pk, err = a.verifyTxSig(txHash.Bytes(), untx.GetSignature(), untx.GetCrypto())
		if err != nil {
			return err
		}

		a.cache.setPublicKey(txHash, pk)
	}

	//TODO: black verify

	if err := a.checkTxParams(tx, pk); err != nil {
		return err
	}

	signTx := &types.SignedTransaction{
		TransactionWithSig: untx,
		TxHash:             txHash.Bytes(),
		Signer:             ethcrypto.CompressPubkey(pk),
	}

	return a.txPool.Add(signTx)
}

func (a *auth) checkTxParams(tx *types.Transaction, signer *ecdsa.PublicKey) error {
	if tx.ChainId != a.chainID {
		return errors.New("bad chain id")
	}

	if len(tx.Nonce) > nonceLenLimit {
		return errors.New("invalid nonce")
	}

	if err := a.checkValidUntilBlockFromCache(tx.ValidUntilBlock); err != nil {
		return err
	}

	if err := a.checkQuota(tx.Quota, ethcrypto.PubkeyToAddress(*signer)); err != nil {
		return err
	}

	return nil
}

// TODO: check valid until block
func (a *auth) checkValidUntilBlockFromCache(validUntilBlock uint64) error {
	return nil
}

// TODO: check history txs
func (a *auth) checkHistoryTxs(hash common.Hash) error {
	return nil
}

// TODO: check quota
func (a *auth) checkQuota(quota uint64, address common.Address) error {
	if quota > a.blockQuotaLimit {
		return errors.New("quota  not enough")
	}
	return nil
}

func (a *auth) verifyTxSig(hash, signature []byte, crypto types.Crypto) (*ecdsa.PublicKey, error) {
	switch crypto {
	case types.Crypto_SECP:
		return ethcrypto.SigToPub(hash, signature)
	}

	return nil, fmt.Errorf("%s is Unexpected crypto", crypto.String())
}
