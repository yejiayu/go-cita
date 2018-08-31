package vm

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/vm/evm"
)

type Executor interface {
	Call(ctx context.Context, header *pb.BlockHeader, signedTxs []*pb.SignedTransaction) ([]*pb.Receipt, []byte, error)
}

func NewExecutor(factory database.Factory) Executor {
	return &executor{factory: factory}
}

type executor struct {
	factory database.Factory
}

func (e *executor) Call(ctx context.Context, header *pb.BlockHeader, signedTxs []*pb.SignedTransaction) ([]*pb.Receipt, []byte, error) {
	rootHash := common.BytesToHash(header.GetStateRoot())
	stateDB, err := state.New(rootHash, state.NewDatabase(e.factory.EthDB()))
	if err != nil {
		return nil, nil, err
	}
	blockDB := e.factory.BlockDB()
	blockHash := common.BytesToHash(header.Prevhash)

	receipts := make([]*pb.Receipt, len(signedTxs))
	for i, signedTx := range signedTxs {
		snap := stateDB.Snapshot()

		msg := NewMessage(header.GetGasLimit(), signedTx)
		vmCtx := NewVMContext(header, msg, blockDB)

		stateDB.Prepare(common.BytesToHash(signedTx.GetTxHash()), blockHash, i)
		evm := evm.New(vmCtx, stateDB, params.MainnetChainConfig, vm.Config{})

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

func NewVMContext(header *pb.BlockHeader, msg core.Message, blockDB block.Interface) vm.Context {
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
