package auth

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/auth/pool"
	"github.com/yejiayu/go-cita/mq"
	"github.com/yejiayu/go-cita/types"
)

const (
	nonceLenLimit        = 128
	validUntilBlockLimit = 100
)

type Interface interface {
	Run(quit chan<- error)
}

func New(queue mq.Queue) Interface {
	txPool := pool.NewTxPool()

	return &auth{
		queue:  queue,
		txPool: txPool,
	}
}

type auth struct {
	mu     sync.Mutex
	queue  mq.Queue
	txPool pool.TxPool

	chainID         uint32
	blockQuotaLimit uint64

	pkCache map[common.Hash]*ecdsa.PublicKey
}

func (a *auth) Run(quit chan<- error) {
	go a.handleMQ(quit)
}

func (a *auth) handleMQ(quit chan<- error) {
	delivery, err := a.queue.Sub()
	if err != nil {
		quit <- err
		return
	}

	for msg := range delivery {
		key := mq.RoutingKey(msg.RoutingKey)
		switch key {
		case mq.NetworkVerifyTxRequest:
			var untx types.UnverifiedTransaction
			if err := proto.Unmarshal(msg.Body, &untx); err != nil {
				glog.Error(err)
				continue
			}

			if err := a.authUntx(&untx); err != nil {
				glog.Error(err)
				continue
			}
		}
	}
}

func (a *auth) authUntx(untx *types.UnverifiedTransaction) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	tx := untx.GetTransaction()
	data, err := proto.Marshal(untx.GetTransaction())
	if err != nil {
		return err
	}

	txHash := common.BytesToHash(data)
	pk, ok := a.pkCache[txHash]
	if !ok {
		var err error
		pk, err = a.verifyTxSig(txHash.Bytes(), untx.GetSignature(), untx.GetCrypto())
		if err != nil {
			return err
		}

		a.pkCache[txHash] = pk
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
