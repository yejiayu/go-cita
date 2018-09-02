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
	"math"
	"sync"
	"time"

	"github.com/yejiayu/go-cita/common/hash"
	cfg "github.com/yejiayu/go-cita/config/consensus"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/consensus/tendermint/params"
	"github.com/yejiayu/go-cita/consensus/wal"
)

type Interface interface {
	SetProposal(ctx context.Context, proposal *pb.Proposal, signature []byte) error
	SetVote(ctx context.Context, vote *pb.Vote, signature []byte) error
}

func New(authClient pb.AuthClient, chainClient pb.ChainClient, networkClient pb.NetworkClient) Interface {
	t := &tendermint{
		authClient:    authClient,
		chainClient:   chainClient,
		networkClient: networkClient,

		quotaLimit: cfg.GetQuotaLimit(),
		txCount:    cfg.GetTxCount(),
	}

	privKey, err := cfg.GetPrivKey()
	if err != nil {
		log.Panic(err)
	}
	t.singer = params.NewSinger(privKey)

	wal, err := wal.NewFileWAL("file_wal")
	if err != nil {
		log.Panic(err)
	}
	t.wal = wal

	if err := t.init(); err != nil {
		log.Panic(err)
	}
	return t
}

type tendermint struct {
	authClient    pb.AuthClient
	chainClient   pb.ChainClient
	networkClient pb.NetworkClient

	wal wal.Interface

	quotaLimit uint64
	txCount    uint32
	singer     *params.Singer

	rs RoundStep

	mu             sync.RWMutex
	lastHeader     *pb.BlockHeader
	lastHeaderHash hash.Hash
	valSet         *params.ValidatorSet
}

func (t *tendermint) init() error {
	ctx := context.Background()

	res, err := t.chainClient.GetBlockHeader(ctx, &pb.GetBlockHeaderReq{
		Height: math.MaxUint64,
	})
	if err != nil {
		return err
	}
	if res.GetHeader() == nil {
		log.Panic("Please confirm the genesis block exists")
	}

	lastHeader := res.GetHeader()
	lastHeaderHash, err := hash.ProtoToSha3(lastHeader)
	if err != nil {
		return err
	}
	t.lastHeader = lastHeader
	t.lastHeaderHash = lastHeaderHash

	valSet, err := t.GetValidatorSet(t.lastHeader.GetHeight())
	if err != nil {
		return err
	}
	t.valSet = valSet
	t.rs = NewRoundStep(t.lastHeader.GetHeight()+1, t.valSet, t.singer, t)
	return nil
}

// func (t *tendermint) createGenesis(ctx context.Context) (*pb.BlockHeader, error) {
// 	_, err := t.chainClient.NewBlock(ctx, &pb.NewBlockReq{
// 		Block: &pb.Block{
// 			Header: &pb.BlockHeader{
// 				Height: 0,
// 			},
// 			Body: &pb.BlockBody{},
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	res, err := t.chainClient.GetBlockHeader(ctx, &pb.GetBlockHeaderReq{
// 		Height: math.MaxUint64,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res.GetHeader(), nil
// }

// SetProposal. RoundStepPropose -> RoundStepPrevote
func (t *tendermint) SetProposal(ctx context.Context, proposal *pb.Proposal, signature []byte) error {
	return t.rs.SetProposal(proposal, signature)
}

func (t *tendermint) SetVote(ctx context.Context, vote *pb.Vote, signature []byte) error {
	return t.rs.SetVote(vote, signature)
}

// ---------impl extension---------

func (t *tendermint) ProposalBlock(height uint64, signer *params.Singer) (*pb.Block, error) {
	if height-1 != t.lastHeader.GetHeight() {
		log.Panicf("ProposalBlock, proposal height is %d, but last height is %d", height, t.lastHeader.GetHeight())
	}

	ctx := context.Background()
	res, err := t.authClient.GetTxFromPool(ctx, &pb.GetTxFromPoolReq{
		QuotaLimit: t.quotaLimit,
		TxCount:    t.txCount,
	})
	if err != nil {
		return nil, err
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	var txRoot []byte
	for _, txHash := range res.GetTxHashes() {
		txRoot = append(txRoot, txHash...)
	}
	block := &pb.Block{
		Version: 2,
		Header: &pb.BlockHeader{
			Prevhash:         t.lastHeaderHash.Bytes(),
			Timestamp:        uint64(time.Now().Unix()),
			Height:           height,
			TransactionsRoot: hash.BytesToSha3(txRoot).Bytes(),
			GasUsed:          res.GetQuotaUsed(),
			GasLimit:         t.quotaLimit,
			Proposer:         signer.Address().Bytes(),
		},
		Body: &pb.BlockBody{
			TxHashes: res.GetTxHashes(),
		},
	}

	return block, nil
}

func (t *tendermint) ValidateProposalBlock(proposal *pb.Proposal) error {
	log.Info("defaultExtension.ValidateProposalBlock")
	return nil
}

func (t *tendermint) BroadcastProposal(proposal *pb.Proposal, signature []byte) error {
	_, err := t.networkClient.BroadcastProposal(context.Background(), &pb.BroadcastProposalReq{
		Proposal:  proposal,
		Signature: signature,
	})
	return err
}

func (t *tendermint) BroadcastVote(vote *pb.Vote, signature []byte) error {
	_, err := t.networkClient.BroadcastVote(context.Background(), &pb.BroadcastVotelReq{
		Vote:      vote,
		Signature: signature,
	})
	return err
}

func (t *tendermint) Commit(block *pb.Block) error {
	ctx := context.Background()
	_, err := t.chainClient.NewBlock(ctx, &pb.NewBlockReq{
		Block: block,
	})
	if err != nil {
		return err
	}

	_, err = t.authClient.ClearPool(ctx, &pb.ClearPoolReq{
		Height: block.GetHeader().GetHeight(),
	})
	if err != nil {
		return err
	}

	res, err := t.chainClient.GetBlockHeader(ctx, &pb.GetBlockHeaderReq{
		Height: math.MaxUint64,
	})
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.lastHeader = res.GetHeader()
	lastHeaderHash, err := hash.ProtoToSha3(t.lastHeader)
	if err != nil {
		return err
	}
	t.lastHeaderHash = lastHeaderHash
	return err
}

func (e *tendermint) WAL(data []byte) error {
	return nil
}

func (e *tendermint) GetValidatorSet(height uint64) (*params.ValidatorSet, error) {
	res, err := e.chainClient.GetValidators(context.Background(), &pb.GetValidatorsReq{
		Height: height,
	})
	if err != nil {
		return nil, err
	}

	return params.NewValidatorSet(res.GetVals())
}
