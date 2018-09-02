package vm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/database/ethdb"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/vm/evm"
)

type Executor interface {
	Call(ctx context.Context, header *pb.BlockHeader, signedTxs []*pb.SignedTransaction) ([]*pb.Receipt, []byte, error)
	StaticCall(ctx context.Context, height uint64, from, to, data []byte) ([]byte, error)
}

func NewExecutor(factory database.Factory) Executor {
	return &executor{
		ethDB:   factory.EthDB(),
		blockDB: factory.BlockDB(),

		chainConfig: params.MainnetChainConfig,
	}
}

type executor struct {
	ethDB   ethdb.Database
	blockDB block.Interface

	chainConfig *params.ChainConfig
}

func (e *executor) Call(ctx context.Context, header *pb.BlockHeader, signedTxs []*pb.SignedTransaction) ([]*pb.Receipt, []byte, error) {
	rootHash := common.BytesToHash(header.GetStateRoot())
	stateDB, err := state.New(rootHash, state.NewDatabase(e.ethDB))
	if err != nil {
		return nil, nil, err
	}
	blockHash := common.BytesToHash(header.Prevhash)

	receipts := make([]*pb.Receipt, len(signedTxs))
	for i, signedTx := range signedTxs {
		snap := stateDB.Snapshot()
		tx := signedTx.GetTransactionWithSig().GetTransaction()

		var toAddr *common.Address
		if tx.GetTo() != "" {
			temp := common.HexToAddress(tx.GetTo())
			if (common.Address{}) != temp {
				toAddr = &temp
			}
		}
		value := big.NewInt(0)
		if len(tx.GetValue()) > 0 {
			value = value.SetBytes(tx.GetValue())
		}
		msg := NewMessage(
			common.BytesToAddress(signedTx.GetSigner()), toAddr,
			tx.GetData(), header.GetGasLimit(), value,
			signedTx.GetTxHash(),
		)
		vmCtx := NewVMContext(header, msg, e.blockDB)

		stateDB.Prepare(common.BytesToHash(signedTx.GetTxHash()), blockHash, i)
		evm := evm.New(vmCtx, stateDB, e.chainConfig, vm.Config{})

		receipt, err := e.applyMessage(stateDB, evm, msg)
		if err != nil {
			stateDB.RevertToSnapshot(snap)
		}
		receipts[i] = receipt
	}

	root := stateDB.IntermediateRoot(true)
	stateDB.Commit(true)
	stateDB.Database().TrieDB().Commit(root, true)
	return receipts, root.Bytes(), nil
}

func (e *executor) StaticCall(ctx context.Context, height uint64, from, to, data []byte) ([]byte, error) {
	header, err := e.blockDB.GetHeaderByHeight(ctx, height)
	if err != nil {
		return nil, err
	}
	if header == nil {
		return nil, fmt.Errorf("not found header at height %d", height)
	}

	root := common.BytesToHash(header.GetStateRoot())
	stateDB, err := state.New(root, state.NewDatabase(e.ethDB))
	if err != nil {
		return nil, err
	}

	fromAddr := common.BytesToAddress(from)
	var toAddr *common.Address
	if len(to) > 0 {
		temp := common.BytesToAddress(to)
		toAddr = &temp
	}
	msg := NewMessage(fromAddr, toAddr, data, header.GetGasLimit(), big.NewInt(0), []byte{})
	vmCtx := NewVMContext(header, msg, e.blockDB)

	evm := evm.New(vmCtx, stateDB, e.chainConfig, vm.Config{})
	ret, _, err := evm.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())
	return ret, err
}

func (e *executor) applyMessage(stateDB *state.StateDB, evm evm.EVM, msg *Message) (*pb.Receipt, error) {
	sender := vm.AccountRef(msg.From())
	receipt := &pb.Receipt{StateRoot: []byte{}}
	var quotaUsed uint64
	var err error

	if msg.To() == nil {
		var addr common.Address
		_, addr, quotaUsed, err = evm.Create(sender, msg.Data(), msg.Gas(), msg.Value())
		receipt.ContractAddress = addr.Bytes()
	} else {
		_, quotaUsed, err = evm.Call(sender, *msg.To(), msg.Data(), msg.Gas(), msg.Value())
	}
	stateDB.Finalise(true)
	if err != nil {
		receipt.Error = err.Error()
	}
	receipt.GasUsed = quotaUsed
	ethLogs := stateDB.GetLogs(common.BytesToHash(msg.TxHash()))
	logs := make([]*pb.LogEntry, len(ethLogs))
	for i, l := range ethLogs {
		topics := make([][]byte, len(l.Topics))
		for j, topic := range l.Topics {
			topics[j] = topic.Bytes()
		}

		logs[i] = &pb.LogEntry{
			Address: l.Address.Bytes(),
			Topics:  topics,
			Data:    l.Data,
		}
	}
	receipt.Logs = logs

	bin := ethTypes.LogsBloom(ethLogs)
	receipt.LogBloom = bin.Bytes()

	return receipt, err
}

func NewVMContext(header *pb.BlockHeader, msg *Message, blockDB block.Interface) vm.Context {
	blockNumber := &big.Int{}
	blockNumber.SetUint64(header.GetHeight())

	t := &big.Int{}
	t.SetInt64(time.Now().Unix())

	return vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFunc(blockDB),
		Origin:      msg.From(),
		GasLimit:    header.GetGasLimit(),
		BlockNumber: blockNumber,
		Time:        t,
		Difficulty:  nil,
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

// GetHashFunc returns the hash corresponding to n
func GetHashFunc(blockDB block.Interface) func(height uint64) common.Hash {
	cache := make(map[uint64]common.Hash)
	ctx := context.Background()

	return func(height uint64) common.Hash {
		h, ok := cache[height]
		if ok {
			return h
		}

		header, err := blockDB.GetHeaderByHeight(ctx, height)
		if err != nil {
			return common.Hash{}
		}
		data, err := hash.ProtoToSha3(header)
		if err != nil {
			return common.Hash{}
		}
		return common.BytesToHash(data.Bytes())
	}
}
