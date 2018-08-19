package tendermint

import (
	"sync"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/consensus/tendermint/params"
)

type RoundVoteSet struct {
	Prevotes   *params.VoteSet
	Precommits *params.VoteSet
}

type HeightVoteSet struct {
	mu            sync.RWMutex
	height        uint64
	round         uint64
	valSet        *params.ValidatorSet
	roundVoteSets map[uint64]*RoundVoteSet
}

func NewHeightVoteSet(height, round uint64, valSet *params.ValidatorSet) *HeightVoteSet {
	return &HeightVoteSet{
		height:        height,
		round:         round,
		valSet:        valSet,
		roundVoteSets: make(map[uint64]*RoundVoteSet),
	}
}

func (hvs *HeightVoteSet) AddVote(vote *pb.Vote, signature []byte) bool {
	hvs.mu.Lock()
	defer hvs.mu.Unlock()

	if vote.GetHeight() != hvs.height {
		return false
	}

	height := vote.GetHeight()
	round := vote.GetRound()
	rvs, ok := hvs.roundVoteSets[round]
	if !ok {
		rvs = &RoundVoteSet{
			Prevotes:   params.NewVoteSet(height, round, hvs.valSet),
			Precommits: params.NewVoteSet(height, round, hvs.valSet),
		}
		hvs.roundVoteSets[round] = rvs
	}

	switch vote.GetVoteType() {
	case pb.VoteType_Prevote:
		return rvs.Prevotes.AddVote(vote, signature)
	case pb.VoteType_precommit:
		return rvs.Precommits.AddVote(vote, signature)
	}

	return false
}

func (hvs *HeightVoteSet) Prevotes(round uint64) *params.VoteSet {
	rvs, ok := hvs.roundVoteSets[round]
	if !ok {
		return nil
	}

	return rvs.Prevotes
}

func (hvs *HeightVoteSet) Precommits(round uint64) *params.VoteSet {
	rvs, ok := hvs.roundVoteSets[round]
	if !ok {
		return nil
	}

	return rvs.Precommits
}

// Last round and blockID that has +2/3 prevotes for a particular block or nil.
// Returns -1 if no such round exists.
func (hvs *HeightVoteSet) POLInfo() (uint64, hash.Hash, bool) {
	hvs.mu.Lock()
	defer hvs.mu.Unlock()

	for r := hvs.round; r >= 0; r-- {
		rvs, ok := hvs.roundVoteSets[r]
		if ok {
			blockHash, ok := rvs.Prevotes.TwoThirdsMajority()
			if ok {
				return r, blockHash, true
			}
		}
	}
	return 0, hash.Hash{}, false
}
