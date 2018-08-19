package params

import (
	"bytes"
	"sync"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/pb"
)

type VoteSet struct {
	mu     sync.RWMutex
	height uint64
	round  uint64

	votesByBlock map[hash.Hash][]*pb.Vote
	valSet       *ValidatorSet
}

func NewVoteSet(height, round uint64, valSet *ValidatorSet) *VoteSet {
	return &VoteSet{
		height: height,
		round:  round,

		votesByBlock: make(map[hash.Hash][]*pb.Vote),
		valSet:       valSet,
	}
}

func (vs *VoteSet) AddVote(vote *pb.Vote, signature []byte) bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vs.height != vote.GetHeight() || vs.round != vote.GetRound() {
		return false
	}

	address := hash.BytesToAddress(vote.GetAddress())
	blockHash := hash.BytesToHash(vote.GetHash())
	if !vs.valSet.GetByAddress(address).VerifySignature(blockHash, signature) {
		return false
	}

	for _, votes := range vs.votesByBlock {
		for _, v := range votes {
			if bytes.Equal(v.GetAddress(), address.Bytes()) {
				return false
			}
		}
	}

	vs.votesByBlock[blockHash] = append(vs.votesByBlock[blockHash], vote)
	return true
}

func (vs *VoteSet) HasTwoThirdsAny() bool {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	valLen := len(vs.valSet.Validators)
	totalVotes := 0
	for _, votes := range vs.votesByBlock {
		totalVotes = totalVotes + len(votes)
		if totalVotes > (valLen * 2 / 3) {
			return true
		}
	}
	return false
}

// If there was a +2/3 majority for blockID, return blockID and true.
// Else, return the empty BlockID{} and false.
func (vs *VoteSet) TwoThirdsMajority() (hash.Hash, bool) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	valLen := len(vs.valSet.Validators)
	for blockHash, votes := range vs.votesByBlock {
		if len(votes) > (valLen * 2 / 3) {
			return blockHash, true
		}
	}

	return hash.Hash{}, false
}
