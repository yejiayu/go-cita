package types

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TxAction string

const (
	/// Just store the data.
	Store TxAction = "Store"
	/// Create creates new contract.
	Create TxAction = "Create"
	/// Calls contract at given address.
	/// In the case of a transfer, this is the receiver's address.'
	TxCall TxAction = "Call"
	/// Store the contract ABI
	TxAbiStore TxAction = "AbiStore"
	/// Create creates new contract for grpc.
	TxGoCreate TxAction = "GoCreate"
)

type Transaction struct {
	/// Nonce.
	Nonce string `json:"nonce"`
	/// Gas paid up front for transaction execution.
	Quota *big.Int `json:"quota"`
	/// Action, can be either call or contract create.
	Action TxAction `json:"-"`
	/// Transfered value.
	Value *big.Int `json:"value"`
	/// Transaction data.
	Data []byte `json:"data"`
	// /// valid before this block number
	// pub block_limit: BlockNumber,
	// /// Unique chain_id
	// pub chain_id: u32,
	// /// transaction version
	// pub version: u32,
}

type CryptoType string

const (
	CryptoTypeSECP = "SECP"
	CryptoTypeSM2  = "SM2"
)

type UnverifiedTransaction struct {
	Unsigned   *Transaction
	Signature  types.Signer
	CryptoType CryptoType
	Hash       common.Hash
}

type SignedTransaction struct {
	Transaction *UnverifiedTransaction
	Sender      common.Hash
	Public      *ecdsa.PublicKey
}
