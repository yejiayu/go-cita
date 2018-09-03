package evm

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/log"
)

var code1 = "0x6080604052600436106100da5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633408e47081146100df57806363ffec6e1461010d5780636fbf656a146101465780638ec1aaed14610168578063984dc34b1461018f578063abb1dc44146101b8578063b4a0e24c14610311578063c0c41f2214610326578063d5cd402c14610346578063d722b0bc14610378578063dd12b51f14610402578063df51aa4914610417578063e19709e01461042c578063e7f43c681461046a578063f87f44b91461047f575b600080fd5b3480156100eb57600080fd5b506100f461049f565b6040805163ffffffff9092168252519081900360200190f35b34801561011957600080fd5b506101226104ac565b6040518082600181111561013257fe5b60ff16815260200191505060405180910390f35b34801561015257600080fd5b5061016660048035602481019101356104c1565b005b34801561017457600080fd5b5061017d6104d2565b60408051918252519081900360200190f35b34801561019b57600080fd5b506101a46104d8565b604080519115158252519081900360200190f35b3480156101c457600080fd5b506101cd6104e7565b60405180806020018060200180602001848103845287818151815260200191508051906020019080838360005b838110156102125781810151838201526020016101fa565b50505050905090810190601f16801561023f5780820380516001836020036101000a031916815260200191505b50848103835286518152865160209182019188019080838360005b8381101561027257818101518382015260200161025a565b50505050905090810190601f16801561029f5780820380516001836020036101000a031916815260200191505b50848103825285518152855160209182019187019080838360005b838110156102d25781810151838201526020016102ba565b50505050905090810190601f1680156102ff5780820380516001836020036101000a031916815260200191505b50965050505050505060405180910390f35b34801561031d57600080fd5b506101a46106a1565b34801561033257600080fd5b5061016660048035602481019101356106dc565b34801561035257600080fd5b5061035b6106e8565b6040805167ffffffffffffffff9092168252519081900360200190f35b34801561038457600080fd5b5061038d6106f8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156103c75781810151838201526020016103af565b50505050905090810190601f1680156103f45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561040e57600080fd5b506101a461078b565b34801561042357600080fd5b5061038d610794565b34801561043857600080fd5b506104416107f5565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b34801561047657600080fd5b5061038d610818565b34801561048b57600080fd5b506101666004803560248101910135610879565b60035463ffffffff165b90565b60065468010000000000000000900460ff1690565b6104cd60048383610881565b505050565b60005490565b60015462010000900460ff1690565b60078054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093849384938301828280156105735780601f1061054857610100808354040283529160200191610573565b820191906000526020600020905b81548152906001019060200180831161055657829003601f168201915b505060088054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152969950919450925084019050828280156106045780601f106105d957610100808354040283529160200191610604565b820191906000526020600020905b8154815290600101906020018083116105e757829003601f168201915b505060098054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152969850919450925084019050828280156106955780601f1061066a57610100808354040283529160200191610695565b820191906000526020600020905b81548152906001019060200180831161067857829003601f168201915b50505050509050909192565b600154600090610100900460ff1680156106d75750600060065468010000000000000000900460ff1660018111156106d557fe5b145b905090565b6104cd60028383610881565b60065467ffffffffffffffff1690565b60028054604080516020601f60001961010060018716150201909416859004938401819004810282018101909252828152606093909290918301828280156107815780601f1061075657610100808354040283529160200191610781565b820191906000526020600020905b81548152906001019060200180831161076457829003601f168201915b5050505050905090565b60015460ff1690565b60058054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156107815780601f1061075657610100808354040283529160200191610781565b6001546301000000900473ffffffffffffffffffffffffffffffffffffffff1690565b60048054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156107815780601f1061075657610100808354040283529160200191610781565b6104cd600583835b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106108c25782800160ff198235161785556108ef565b828001600101855582156108ef579182015b828111156108ef5782358255916020019190600101906108d4565b506108fb9291506108ff565b5090565b6104a991905b808211156108fb57600081556001016109055600a165627a7a72305820d91cb76167d5a6a60fea8f161e99d41ad6ea384c2f9a501163b79c48441ab0620029"

type account struct {
	Addr common.Address
}

func (a *account) Address() common.Address {
	return a.Addr
}

func TestABI(t *testing.T) {
	abi, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[],"name":"add","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"CountEvent","type":"event"}]`))
	if err != nil {
		t.Fatal(err)
	}
	addData, err := abi.Pack("add")
	if err != nil {
		t.Fatal(err)
	}
	getData, err := abi.Pack("get")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(common.ToHex(addData))
	t.Log(common.ToHex(getData))
}

func TestCall(t *testing.T) {
	evm := newEvm(t)

	abi, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"getChainId","outputs":[{"name":"","type":"uint32"}],"payable":false,"stateMutability":"pure","type":"function"}]`))
	if err != nil {
		t.Fatal(err)
	}
	data, err := abi.Pack("getChainId")
	if err != nil {
		t.Fatal(err)
	}

	caller := &account{Addr: common.HexToAddress("0xffffffffffffffffffffffffffffffffff020000")}
	addr := common.HexToAddress("0xffffffffffffffffffffffffffffffffff020000")
	ret, gas, err := evm.Call(caller, addr, data, 99999, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("call result", hexutil.Bytes(ret))
	t.Log(gas)
}

func TestCreate(t *testing.T) {
	evm := newEvm(t)
	// 0x608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632f2f485914610046575b600080fd5b34801561005257600080fd5b5061005b6100d6565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561009b578082015181840152602081019050610080565b50505050905090810190601f1680156100c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040805190810160405280600a81526020017f68656c6c6f576f726c64000000000000000000000000000000000000000000008152509050905600a165627a7a723058203f90f5c641bddbfe4dd2d41d8a0c83bb619993f7570d045e5e4543e34c1be46e0029
	data := common.FromHex(`0x608060405234801561001057600080fd5b5061013f806100206000396000f300608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632f2f485914610046575b600080fd5b34801561005257600080fd5b5061005b6100d6565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561009b578082015181840152602081019050610080565b50505050905090810190601f1680156100c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040805190810160405280600a81526020017f68656c6c6f576f726c64000000000000000000000000000000000000000000008152509050905600a165627a7a723058203f90f5c641bddbfe4dd2d41d8a0c83bb619993f7570d045e5e4543e34c1be46e0029`)
	caller := &account{Addr: common.HexToAddress("0xff572e5295c57f15886f9b263e2f6d2d6c7b5ec6")}
	ret, addr, quota, err := evm.Create(caller, data, 99999, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(ret))
	t.Log(addr.String())
	t.Log(quota)

	abi, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"helloworld","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"pure","type":"function"}]`))
	if err != nil {
		t.Fatal(err)
	}

	data, err = abi.Pack("helloworld")
	if err != nil {
		t.Fatal(err)
	}

	ret, quota, err = evm.Call(caller, addr, data, 9999, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(ret), "call create")
	t.Log(quota)
}

func newEvm(t *testing.T) EVM {
	factory, err := database.NewFactory("redis", []string{"127.0.0.1:6379"})
	if err != nil {
		t.Fatal(err)
	}
	header, err := factory.BlockDB().GetHeaderByLatest(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	stateDB, err := state.New(common.BytesToHash(header.GetStateRoot()), state.NewDatabase(factory.EthDB()))
	if err != nil {
		panic(err)
	}

	ctx := vm.Context{
		CanTransfer: func(db vm.StateDB, addr common.Address, amount *big.Int) bool {
			return true
		},
		Transfer: func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
			db.SubBalance(sender, amount)
			db.AddBalance(recipient, amount)
		},
		GetHash: func(height uint64) common.Hash {
			log.Info(height, "get hash")
			return common.Hash{}
		},
		Origin:      common.HexToAddress("0xff572e5295c57f15886f9b263e2f6d2d6c7b5ec6"),
		GasLimit:    99999999,
		GasPrice:    big.NewInt(2),
		BlockNumber: big.NewInt(0),
		Time:        big.NewInt(time.Now().Unix()),
		Difficulty:  big.NewInt(1000),
	}
	chainParams := params.TestChainConfig
	return New(ctx, stateDB, chainParams, vm.Config{})
}
