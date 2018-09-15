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

	pub1 := crypto.FromECDSAPub(&priv1.PublicKey)
	pub2 := crypto.FromECDSAPub(&priv2.PublicKey)
	pub3 := crypto.FromECDSAPub(&priv3.PublicKey)

	pubkList := [][]byte{pub1, pub2, pub3}
	valSet, err := params.NewValidatorSet(pubkList)
	if err != nil {
		t.Fatal(err)
	}
	signer := params.NewSinger(priv1)
	signer2 = params.NewSinger(priv2)
	signer3 = params.NewSinger(priv3)
	mock := &mockExtension{signer: signer}
	rs = NewRoundStep(1, valSet, signer, mock)

	for {
		select {}
	}
}

func (mock *mockExtension) ProposalBlock(height uint64, signer *params.Singer) (*pb.Block, error) {
	log.Infof("ProposalBlock height %d", height)
	return &pb.Block{
		Header: &pb.BlockHeader{
			Height: height,
		},
		Body: &pb.BlockBody{},
	}, nil
}

func (mock *mockExtension) ValidateProposalBlock(proposal *pb.Proposal) error {
	log.Infof("ValidateProposalBlock proposal height %d", proposal.GetHeight())
	return nil
}

func (mock *mockExtension) BroadcastProposal(proposal *pb.Proposal, signature []byte) error {
	log.Infof("BroadcastProposal height=%d round=%d", proposal.GetHeight(), proposal.GetRound())

	vote2, sig2 := buildVote(proposal.GetHeight(), proposal.GetRound(), pb.VoteType_Prevote, proposal.GetBlock(), signer2)
	vote3, sig3 := buildVote(proposal.GetHeight(), proposal.GetRound(), pb.VoteType_Prevote, proposal.GetBlock(), signer3)
	rs.SetVote(vote2, sig2)
	rs.SetVote(vote3, sig3)

	vote2, sig2 = buildVote(proposal.GetHeight(), proposal.GetRound(), pb.VoteType_precommit, proposal.GetBlock(), signer2)
	vote3, sig3 = buildVote(proposal.GetHeight(), proposal.GetRound(), pb.VoteType_precommit, proposal.GetBlock(), signer3)
	rs.SetVote(vote2, sig2)
	rs.SetVote(vote3, sig3)
	return nil
}

func (mock *mockExtension) BroadcastVote(vote *pb.Vote, signature []byte) error {
	return nil
}

func (mock *mockExtension) SetVoteByNode(targetNode []byte, vote *pb.Vote, signature []byte) error {
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

func (mock *mockExtension) GetValidatorSet(height uint64) (*params.ValidatorSet, error) {
	return nil, nil
}

func buildVote(height, round uint64, vt pb.VoteType, block *pb.Block, signer *params.Singer) (*pb.Vote, []byte) {
	blockHash, err := hash.ProtoToSha3(block)
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}

	return vote, signature
}
