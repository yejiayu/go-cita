package vm

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/database"
)

const TokenABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"supplyAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

const sysConfigABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getChainId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getEconomicalModel\n\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_operator\",\"type\":\"string\"}],\"name\":\"setOperator\",\"outputs\":[],\"pay\nable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getDelayBlockNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"ty\npe\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTokenInfo\",\"outputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\"},{\"name\":\"avatar\",\"type\":\"string\"}],\"payable\":false,\"stateMuta\nbility\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getQuotaCheck\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"in\nputs\":[{\"name\":\"_chainName\",\"type\":\"string\"}],\"name\":\"setChainName\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getBlockInterval\",\"ou\ntputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getChainName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"\nstateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPermissionCheck\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"const\nant\":true,\"inputs\":[],\"name\":\"getWebsite\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"\nname\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_website\",\"type\":\"string\"}],\"name\":\"setWebsite\",\"outputs\":[],\"payable\":false,\"sta\nteMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_delayBlockNumber\",\"type\":\"uint256\"},{\"name\":\"_checkPermission\",\"type\":\"bool\"},{\"name\":\"_checkQuota\",\"type\":\"bool\"},{\"name\":\"_chainName\",\"t\nype\":\"string\"},{\"name\":\"_chainId\",\"type\":\"uint32\"},{\"name\":\"_operator\",\"type\":\"string\"},{\"name\":\"_website\",\"type\":\"string\"},{\"name\":\"_blockInterval\",\"type\":\"uint64\"},{\"name\":\"_economical\",\"type\":\"uint8\"},{\n\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_avatar\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

func TestExecutor(t *testing.T) {
	factory, err := database.NewFactory("redis", []string{"127.0.0.1:6379"})
	if err != nil {
		t.Fatal(err)
	}

	count := 5
	mockData := make([]*pb.SignedTransaction, count)
	abi, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[],"name":"add","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"addEvent","type":"event"}]`))
	if err != nil {
		t.Fatal(err)
	}

	data, err := abi.Pack("add")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		tx := &pb.UnverifiedTransaction{
			Transaction: &pb.Transaction{
				To:              "0x0f572e5295c57f15886f9b263e2f6d2d6c7b5ec6",
				Nonce:           fmt.Sprintf("%d", i),
				Quota:           9999,
				ValidUntilBlock: 100,
				Data:            data,
				Value:           big.NewInt(0).Bytes(),
				ChainId:         10,
				Version:         2,
			},
			Signature: []byte("0xf39b9183786e3347b48c1647fd6058c1832ddc9e3c909743fead6f4092dc47d825e6b20be4eb7e2a272b73c2e8a670719113a94f4e4989beb4e74492aa5b332000"),
			Crypto:    pb.Crypto_SECP,
		}
		txHash, err := hash.ProtoToSha3(tx.GetTransaction())
		if err != nil {
			t.Fatal(err)
		}
		signedTx := &pb.SignedTransaction{
			TransactionWithSig: tx,
			TxHash:             txHash.Bytes(),
			Signer:             []byte("0xf39b9183786e3347b48c1647fd6058c1832ddc9e3c909743fead6f4092dc47d825e6b20be4eb7e2a272b73c2e8a670719113a94f4e4989beb4e74492aa5b332000"),
		}
		mockData[i] = signedTx
	}

	header, err := factory.BlockDB().GetHeaderByLatest(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	executor := NewExecutor(factory)
	receipts, root, err := executor.Call(context.Background(), header, mockData)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(root)
	for _, r := range receipts {
		if r.Error != "" {
			t.Fatal(err)
		}
	}
}
