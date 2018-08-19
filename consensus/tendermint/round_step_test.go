package tendermint

import (
	"testing"

	"github.com/yejiayu/go-cita/common/crypto"
	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/consensus/tendermint/params"
)

type mockExtension struct {
	signer *params.Singer
}

var rs RoundStep

var signer2 *params.Singer
var signer3 *params.Singer

func TestRoundStep(t *testing.T) {
	priv1, err := crypto.HexToECDSA("add757cf60afa08fc54376db9cd1f313f2d20d907f3ac984f227ea0835fc0111")
	if err != nil {
		t.Fatal(err)
	}
	priv2, err := crypto.HexToECDSA("4cba2ced3a58574ce410fc703aaddda03574206b0f7bb1689e1ad257f5516e41")
	if err != nil {
		t.Fatal(err)
	}
	priv3, err := crypto.HexToECDSA("3dc8448b091626437abeded1551ad5667359ca094ad85a7595e521ece2be0e30")
	if err != nil {
		t.Fatal(err)
	}

	val1 := params.NewValidator(1, &priv1.PublicKey)
	val2 := params.NewValidator(2, &priv2.PublicKey)
	val3 := params.NewValidator(3, &priv3.PublicKey)

	valSet := params.NewValidatorSet([]*params.Validator{val1, val2, val3})
	signer := params.NewSinger(priv1)
	signer2 = params.NewSinger(priv2)
	signer3 = params.NewSinger(priv3)
	mock := &mockExtension{signer: signer}
	rs = NewRoundStep(1, valSet, signer, mock)

	for {
		select {}
	}
}

func (mock *mockExtension) ProposalBlock(height uint64) (*pb.Block, error) {
	log.Infof("proposal block height %d", height)
	return &pb.Block{
		Header: &pb.BlockHeader{
			Height: height,
		},
		Body: &pb.BlockBody{},
	}, nil
}

func (mock *mockExtension) ValidateProposalBlock(proposal *pb.Proposal) error {
	log.Infof("proposal height %d", proposal.GetHeight())
	return nil
}

func (mock *mockExtension) BroadcastProposal(proposal *pb.Proposal) error {
	log.Infof("BroadcastProposal height=%d round=%d", proposal.GetHeight(), proposal.GetRound())

	vote2, sig2 := buildVote(proposal.GetHeight(), proposal.GetRound(), pb.VoteType_Prevote, proposal.GetBlock(), signer2)
	vote3, sig3 := buildVote(proposal.GetHeight(), proposal.GetRound(), pb.VoteType_Prevote, proposal.GetBlock(), signer3)
	rs.SetVote(vote2, sig2)
	rs.SetVote(vote3, sig3)
	return nil
}

func (mock *mockExtension) BroadcastVote(vote *pb.Vote) error {
	log.Infof("BroadcastVote, height=%d round=%d", vote.GetHeight(), vote.GetRound())
	return nil
}

func (mock *mockExtension) Commit(block *pb.Block) error {
	log.Infof("commit block, height %d", block.GetHeader().GetHeight())
	return nil
}

func (mock *mockExtension) WAL(data []byte) error {
	log.Infof("wal")
	return nil
}

func buildVote(height, round uint64, vt pb.VoteType, block *pb.Block, signer *params.Singer) (*pb.Vote, []byte) {
	blockHash, err := hash.ProtoToSha3(block)
	if err != nil {
		log.Fatal(err)
	}

	vote := &pb.Vote{
		Height:   height,
		Round:    round,
		VoteType: vt,
		Address:  signer.Address().Bytes(),
		Hash:     blockHash.Bytes(),
	}

	signature, err := signer.SignVote(vote)
	if err != nil {
		log.Fatal(err)
	}

	return vote, signature
}
