package vm

import (
	"encoding/json"
	"testing"

	"github.com/yejiayu/go-cita/database"
)

var genesisConfig = []byte(`{
    "timestamp": 1535527595251,
    "prevhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "alloc": {
        "0x0f572e5295c57f15886f9b263e2f6d2d6c7b5ec6": {
            "nonce": "1",
            "code": "608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634f2be91f146044575b600080fd5b348015604f57600080fd5b5060566058565b005b600160008082825401925050819055507fdb737147a8412358aa59d4fb68dde7c7aa6b2b63281b8b3d05a960725a9d1f4733600054604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15600a165627a7a7230582040dbd1c0ad63ae20242005d9a2b75fa414e216f00850ac0b7c330ad2c5e0d5460029"
        }
    }
}`)

func TestSetupGensis(t *testing.T) {
	factory, err := database.NewFactory("redis", []string{"127.0.0.1:6379"})
	if err != nil {
		t.Fatal(err)
	}

	var genesis Genesis
	if err := json.Unmarshal(genesisConfig, &genesis); err != nil {
		t.Fatal(err)
	}

	if err := SetupGenesis(factory, &genesis, false); err != nil {
		t.Fatal(err)
	}
}
