// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package tendermint

import (
	"context"
	"errors"
	"sync"

	"github.com/yejiayu/go-cita/common/crypto"
	cfg "github.com/yejiayu/go-cita/config/consensus"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/database"
	"github.com/yejiayu/go-cita/database/block"

	"github.com/yejiayu/go-cita/consensus/tendermint/params"
	"github.com/yejiayu/go-cita/consensus/wal"
)

var (
	ErrInvalidProposalSignature = errors.New("Error invalid proposal signature")
	ErrInvalidProposalPOLRound  = errors.New("Error invalid proposal POL round")
	ErrAddingVote               = errors.New("Error adding vote")
	ErrVoteHeightMismatch       = errors.New("Error vote height mismatch")
)

type Interface interface {
	Run() error

	SetProposal(ctx context.Context, proposal *pb.Proposal, signature []byte) error
	SetVote(ctx context.Context, vote *pb.Vote, signature []byte) error
}

func New(dbFactory database.Factory) (Interface, error) {
	privKey, err := cfg.GetPrivKey()
	if err != nil {
		return nil, err
	}

	wal, err := wal.NewFileWAL("file_wal")
	if err != nil {
		return nil, err
	}

	t := &tendermint{
		blockDB: dbFactory.BlockDB(),

		// chainClient: chainClient,
		// authClient:  authClient,
		wal: wal,

		singer: params.NewSinger(privKey),
	}
	return t, nil
}

type tendermint struct {
	mu sync.RWMutex

	blockDB block.Interface
	wal     wal.Interface

	chainClient pb.ChainClient
	authClient  pb.AuthClient

	vs     *params.ValidatorSet
	singer *params.Singer

	rs RoundStep
}

func (t *tendermint) Run() error {
	return nil
}

// SetProposal. RoundStepPropose -> RoundStepPrevote
func (t *tendermint) SetProposal(ctx context.Context, proposal *pb.Proposal, signature []byte) error {
	return t.rs.SetProposal(proposal, signature)
}

func (t *tendermint) SetVote(ctx context.Context, vote *pb.Vote, signature []byte) error {
	return t.rs.SetVote(vote, signature)
}

// func (t *tendermint) doVote(ctx context.Context, block *pb.Block) error {
// 	blockHash := hash.Hash{}
// 	if block != nil {
// 		var err error
// 		blockHash, err = hash.ProtoToSha3(block)
// 		if err != nil {
// 			log.Error(err)
// 		}
// 	}
//
// 	vote := &pb.Vote{
// 		Height: t.height,
// 		Round:  t.round,
// 		Step:   uint32(t.step),
//
// 		Address: t.singer.Address().Bytes(),
// 		Hash:    blockHash.Bytes(),
// 	}
//
// 	sig, err := t.singer.SignVote(vote)
// 	if err != nil {
// 		return err
// 	}
//
// 	// t.networkClient.broadcastVote(vote, sig)
// 	return t.SetVote(ctx, vote, sig)
// }

func (t *tendermint) updateNodes(ctx context.Context) error {
	head, err := t.blockDB.GetHeaderByLatest(ctx)
	if err != nil {
		return err
	}

	res, err := t.chainClient.NodeList(ctx, &pb.NodeListReq{Height: head.GetHeight()})
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	vals := make([]*params.Validator, len(res.Nodes))
	for i, node := range res.Nodes {
		pubkey, err := crypto.DecompressPubkey(node)
		if err != nil {
			return err
		}
		vals[i] = params.NewValidator(uint32(i), pubkey)
	}

	t.vs = params.NewValidatorSet(vals)

	return nil
}

// implement Extension
func (t *tendermint) ValidateProposalBlock(proposal *pb.Proposal) error {
	return nil
}

func (t *tendermint) BroadcastVote(vote *pb.Vote) error {
	return nil
}

func (t *tendermint) Commit(block *pb.Block) error {
	return nil
}
