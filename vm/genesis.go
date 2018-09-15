package vm

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"
)

type Genesis struct {
	Timestamp uint64 `json:"timestamp`
	Prevhash  string `json:"prevhash"`

	Alloc GenesisAlloc `json:"alloc"`
}

type GenesisAlloc map[string]GenesisAccount

type GenesisAccount struct {
	Code    string            `json:"code,omitempty"`
	Storage map[string]string `json:"storage,omitempty"`
	Nonce   string            `json:"nonce,omitempty"`

	// Balance *big.Int   `json:"balance"`
	// Coinbase   common.Address      `json:"coinbase"`
}

func SetupGenesis(factory database.Factory, genesis *Genesis, execContract bool) error {
	ctx := context.Background()

	blockDB := factory.BlockDB()
	header, err := blockDB.GetHeaderByLatest(ctx)
	if err != nil {
		return err
	}

	if header != nil {
		log.Info("genesis block exists")
		return nil
	}

	log.Info("setup genesis block")
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(factory.EthDB()))
	if err != nil {
		return err
	}

	if genesis == nil {
		genesis = &Genesis{
			Timestamp: uint64(time.Now().Unix()),
			Prevhash:  common.Hash{}.String(),
		}
	}
	for addrStr, account := range genesis.Alloc {
		addr := common.HexToAddress(addrStr)
		if execContract {
			code, err := genesisCreateContract(stateDB, addr, common.FromHex(account.Code))
			if err != nil {
				log.Panic(err)
			}
			stateDB.SetCode(addr, code)
		} else {
			stateDB.SetCode(addr, common.FromHex(account.Code))
		}
		for key, value := range account.Storage {
			stateDB.SetState(addr, common.HexToHash(key), common.HexToHash(value))
		}
	}

	root := stateDB.IntermediateRoot(false)
	block := &pb.Block{
		Header: &pb.BlockHeader{
			Prevhash:   common.HexToHash(genesis.Prevhash).Bytes(),
			Timestamp:  genesis.Timestamp,
			QuotaLimit: 99999999,
			Height:     0,
			StateRoot:  root.Bytes(),
		},
		Body: &pb.BlockBody{},
	}

	stateDB.Commit(false)
	stateDB.Database().TrieDB().Commit(root, true)
	return blockDB.AddBlock(context.Background(), block, nil)
}

func genesisCreateContract(stateDB *state.StateDB, addr common.Address, code []byte) ([]byte, error) {
	ctx := vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash: func(height uint64) common.Hash {
			return common.Hash{}
		},
		Origin:      common.Address{},
		GasLimit:    99999999,
		GasPrice:    big.NewInt(1),
		BlockNumber: big.NewInt(0),
		Time:        big.NewInt(time.Now().Unix()),
		Difficulty:  big.NewInt(0),
	}
	evm := vm.NewEVM(ctx, stateDB, params.MainnetChainConfig, vm.Config{})

	contract := vm.NewContract(vm.AccountRef(common.Address{}), vm.AccountRef(addr), big.NewInt(0), 999999)
	contract.SetCallCode(&addr, crypto.Keccak256Hash(code), code)

	return vm.NewEVMInterpreter(evm, vm.Config{}).Run(contract, nil)
}
