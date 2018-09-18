package params

import (
	"bytes"
	"sync"
	"time"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"
)

type SignatureVote struct {
	Vote      *pb.Vote  `json:"vote"`
	Signature []byte    `json:"signature"`
	Time      time.Time `json:"time"`
}

type VoteSet struct {
	mu     sync.RWMutex
	height uint64
	round  uint64

	votesByBlock map[hash.Hash][]*SignatureVote
	valSet       *ValidatorSet
}

func NewVoteSet(height, round uint64, valSet *ValidatorSet) *VoteSet {
	return &VoteSet{
		height: height,
		round:  round,

		votesByBlock: make(map[hash.Hash][]*SignatureVote),
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

	voteHash, err := hash.ProtoToSha3(vote)
	if err != nil {
		log.Error(err)
		return false
	}

	if validator := vs.valSet.GetByAddress(address); validator != nil {
		if !validator.VerifySignature(voteHash, signature) {
			return false
		}
	} else {
		return false
	}

	for _, v := range vs.votesByBlock[blockHash] {
		if bytes.Equal(v.Vote.GetAddress(), address.Bytes()) {
			if time.Now().Sub(v.Time) < time.Second*3 {
				return false
			}
			v.Time = time.Now()
		}
	}

	vs.votesByBlock[blockHash] = append(vs.votesByBlock[blockHash], &SignatureVote{
		Vote:      vote,
		Signature: signature,
		Time:      time.Now(),
	})
	log.Infof("vote set, round %d, %+v", vote.GetRound(), vs.votesByBlock[blockHash])

	return true
}

func (vs *VoteSet) HasTwoThirdsAny() bool {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	valLen := len(vs.valSet.Validators)
	totalVotes := 0
	log.Infof("validators length %d", len(vs.valSet.Validators))
	for blockHash, votes := range vs.votesByBlock {
		log.Infof("HasTwoThirdsAny blockHash %s, votes %d", blockHash.String(), len(votes))
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
	log.Infof("validators length %d", len(vs.valSet.Validators))
	for blockHash, votes := range vs.votesByBlock {
		log.Infof("TwoThirdsMajority blockHash %s, votes %d", blockHash.String(), len(votes))
		if len(votes) > (valLen * 2 / 3) {
			return blockHash, true
		}
	}

	return hash.Hash{}, false
}

func (vs *VoteSet) GetVotes(blockHash hash.Hash) []*SignatureVote {
	return vs.votesByBlock[blockHash]
}
