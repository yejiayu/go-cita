package params

import (
	"testing"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

func TestSigner(t *testing.T) {
	priv, err := crypto.HexToECDSA("add757cf60afa08fc54376db9cd1f313f2d20d907f3ac984f227ea0835fc0111")
	if err != nil {
		t.Fatal(err)
	}

	signer := NewSinger(priv)
	block := &pb.Block{
		Header: &pb.BlockHeader{
			Height: 1,
		},
	}
	signature, err := signer.SignBlock(block)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(signature))

	val := NewValidator(0, &priv.PublicKey)

	h, err := hash.ProtoToSha3(block)
	if err != nil {
		t.Fatal(err)
	}
	if !val.VerifySignature(h, signature) {
		t.Fatal("invalid signature")
	}

	t.Log("sign valid")
}
