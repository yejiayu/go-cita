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
	"errors"
	"sync"
	"time"

	"github.com/yejiayu/go-cita/common/hash"
	"github.com/yejiayu/go-cita/log"
	"github.com/yejiayu/go-cita/pb"

	"github.com/yejiayu/go-cita/consensus/tendermint/params"
)

var (
	ErrInvalidProposalSignature = errors.New("Error invalid proposal signature")
	ErrInvalidProposalPOLRound  = errors.New("Error invalid proposal POL round")
	ErrAddingVote               = errors.New("Error adding vote")
	ErrVoteHeightMismatch       = errors.New("Error vote height mismatch")
)

type Extension interface {
	ProposalBlock(height uint64, signer *params.Singer) (*pb.Block, error)
	ValidateProposalBlock(proposal *pb.Proposal) error

	BroadcastProposal(proposal *pb.Proposal, signature []byte) error
	BroadcastVote(vote *pb.Vote, signature []byte) error

	Commit(block *pb.Block) error

	GetValidatorSet(height uint64) (*params.ValidatorSet, error)

	WAL(data []byte) error
}

// RoundStepType enumerates the state of the consensus state machine
type RoundStepType uint8 // These must be numeric, ordered.

// RoundStepType
const (
	RoundStepNewHeight     = RoundStepType(0x01) // Wait til CommitTime + timeoutCommit
	RoundStepNewRound      = RoundStepType(0x02) // Setup new round and go to RoundStepPropose
	RoundStepPropose       = RoundStepType(0x03) // Did propose, gossip proposal
	RoundStepPrevote       = RoundStepType(0x04) // Did prevote, gossip prevotes
	RoundStepPrevoteWait   = RoundStepType(0x05) // Did receive any +2/3 prevotes, start timeout
	RoundStepPrecommit     = RoundStepType(0x06) // Did precommit, gossip precommits
	RoundStepPrecommitWait = RoundStepType(0x07) // Did receive any +2/3 precommits, start timeout
	RoundStepCommit        = RoundStepType(0x08) // Entered commit state machine
	// NOTE: RoundStepNewHeight acts as RoundStepCommitWait.
)

// String returns a string
func (rs RoundStepType) String() string {
	switch rs {
	case RoundStepNewHeight:
		return "RoundStepNewHeight"
	case RoundStepNewRound:
		return "RoundStepNewRound"
	case RoundStepPropose:
		return "RoundStepPropose"
	case RoundStepPrevote:
		return "RoundStepPrevote"
	case RoundStepPrevoteWait:
		return "RoundStepPrevoteWait"
	case RoundStepPrecommit:
		return "RoundStepPrecommit"
	case RoundStepPrecommitWait:
		return "RoundStepPrecommitWait"
	case RoundStepCommit:
		return "RoundStepCommit"
	default:
		return "RoundStepUnknown" // Cannot panic.
	}
}

// RoundStep defines the internal consensus state.
// NOTE: Not thread safe. Should only be manipulated by functions downstream
// of the cs.receiveRoutine
type RoundStep interface {
	SetProposal(proposal *pb.Proposal, signature []byte) error
	SetVote(vote *pb.Vote, signature []byte) error
}

// RoundStep defines the internal consensus state.
// NOTE: Not thread safe. Should only be manipulated by functions downstream
// of the cs.receiveRoutine
type roundStep struct {
	mu        sync.RWMutex
	stepCh    chan RoundStepType
	extension Extension
	ticker    TimeoutTicker

	Height     uint64        `json:"height"`
	Round      uint64        `json:"round"`
	RoundStep  RoundStepType `json:"step"`
	StartTime  time.Time     `json:"start_time"`
	CommitTime time.Time     `json:"commit_time"`

	Validators *params.ValidatorSet `json:"validators"`
	Singer     *params.Singer       `json:"singer"`

	Proposal  *pb.Proposal `json:"proposal"`
	Signature []byte       `json:"signature"`

	LockedRound   uint64         `json:"locked_round"`
	LockedBlock   *pb.Block      `json:"locked_block"`
	HeightVoteSet *HeightVoteSet `json:"height_vote_set"`
}

func NewRoundStep(height uint64, valSet *params.ValidatorSet, singer *params.Singer, extension Extension) RoundStep {
	rs := &roundStep{
		stepCh:    make(chan RoundStepType),
		extension: extension,
		ticker:    NewTimeoutTicker(),

		Height: height,

		Validators: valSet,
		Singer:     singer,

		HeightVoteSet: NewHeightVoteSet(height, 0, valSet),
	}

	rs.ticker.Start()
	go rs.timeout()

	rs.updateStep(RoundStepNewRound)
	return rs
}

func (rs *roundStep) SetProposal(proposal *pb.Proposal, signature []byte) error {
	if rs.Proposal != nil {
		return nil
	}

	log.Infof("SetProposal height=%d round=%d self height=%d, round=%d", proposal.GetHeight(), proposal.GetRound(), rs.Height, rs.Round)
	height := proposal.GetHeight()
	round := proposal.GetRound()

	if height != rs.Height || round != rs.Round {
		return nil
	}

	blockHash, err := hash.ProtoToSha3(proposal.GetBlock())
	if err != nil {
		return err
	}

	log.Infof("SetProposal %d, proposal round %d", height, round)
	if !rs.Validators.GetProposer(height, round).VerifySignature(blockHash, signature) {
		return ErrInvalidProposalSignature
	}

	rs.Proposal = proposal
	rs.Height = proposal.GetHeight()
	rs.Round = proposal.GetRound()

	go rs.enterPrevote()
	return nil
}

func (rs *roundStep) SetVote(vote *pb.Vote, signature []byte) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if rs.Height != vote.GetHeight() {
		return nil
	}

	if !rs.HeightVoteSet.AddVote(vote, signature) {
		return nil
	}

	if vote.GetRound() < rs.Round ||
		(vote.GetRound() == rs.Round && vote.GetVoteType() == pb.VoteType_Prevote && rs.RoundStep > RoundStepPrevoteWait) {
		newVote, newSignature := rs.genVote(vote.GetHeight(), vote.GetRound(), vote.GetVoteType(), nil)
		if newVote == nil || newSignature == nil {
			return nil
		}

		go rs.extension.BroadcastVote(newVote, newSignature)
		return nil
	}

	log.Infof("vote added, address=%s, height=%d, round=%d, vote_type=%s", hash.ToHex(vote.GetAddress()), vote.GetHeight(), vote.GetRound(), vote.GetVoteType())

	switch vote.GetVoteType() {
	case pb.VoteType_Prevote:
		if rs.RoundStep == RoundStepPrevote || rs.RoundStep == RoundStepPrevoteWait {
			if _, ok := rs.HeightVoteSet.Prevotes(rs.Round).TwoThirdsMajority(); ok {
				if rs.RoundStep == RoundStepPrevoteWait || rs.RoundStep == RoundStepPrevote {
					go rs.enterPrecommit()
				}
			} else if rs.HeightVoteSet.Prevotes(rs.Round).HasTwoThirdsAny() {
				if rs.RoundStep == RoundStepPrevote {
					rs.updateStep(RoundStepPrevoteWait)
				}
			}
		}
	case pb.VoteType_precommit:
		if rs.RoundStep == RoundStepPrecommit || rs.RoundStep == RoundStepPrecommitWait {
			if _, ok := rs.HeightVoteSet.Precommits(rs.Round).TwoThirdsMajority(); ok {
				if rs.RoundStep == RoundStepPrecommitWait || rs.RoundStep == RoundStepPrecommit {
					go rs.enterCommit()
				}
			} else if rs.HeightVoteSet.Precommits(rs.Round).HasTwoThirdsAny() {
				if rs.RoundStep == RoundStepPrecommit {
					rs.updateStep(RoundStepPrecommitWait)
				}
			}
		}
	}

	return nil
}

func (rs *roundStep) updateStep(step RoundStepType) {
	log.Infof("update step to %s, height=%d, round=%d", step.String(), rs.Height, rs.Round)
	rs.RoundStep = step

	switch step {
	case RoundStepNewHeight:
		go rs.enterNewHeight()
	case RoundStepNewRound:
		go rs.enterNewRound()
	case RoundStepPropose:
		go rs.enterPropose()
	case RoundStepPrevoteWait:
		go rs.enterPrecommit()
	case RoundStepPrecommitWait:
		go rs.enterCommit()
	}
}

func (rs *roundStep) timeout() {
	for ti := range rs.ticker.Chan() {
		rs.mu.Lock()
		height := rs.Height
		round := rs.Round
		step := rs.RoundStep
		rs.mu.Unlock()

		if ti.Height != height || ti.Round != round {
			continue
		}

		switch ti.Step {
		case RoundStepPropose:
			go rs.enterPrevote()
		case RoundStepPrevote:
			if step == RoundStepPrevote {
				go rs.enterPrevote()
			}
		case RoundStepPrevoteWait:
			go rs.enterPrecommit()
		case RoundStepPrecommit:
			if step == RoundStepPrecommit {
				go rs.enterPrecommit()
			}
		case RoundStepPrecommitWait:
			go rs.enterCommit()
		}
	}
}

func (rs *roundStep) enterNewHeight() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	t := (time.Second * 3) - rs.CommitTime.Sub(rs.StartTime)
	if t > 0 {
		log.Infof("enterNewHeight sleep %s", t.String())
		time.Sleep(t)
	}

	rs.StartTime = time.Now()
	rs.Height++
	rs.Round = 0
	rs.Proposal = nil
	rs.LockedBlock = nil
	rs.LockedRound = 0

	log.Infof("enterNewHeight, height=%d, round=%d, step=%s", rs.Height, rs.Round, rs.RoundStep.String())
	// TODO: update validators
	rs.HeightVoteSet.Reset(rs.Height, rs.Validators)

	rs.updateStep(RoundStepNewRound)
}

func (rs *roundStep) enterNewRound() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	if rs.LockedBlock == nil {
		rs.Proposal = nil
	}
	rs.HeightVoteSet.SetRound(rs.Round)

	log.Infof("enterNewRound, height=%d, round=%d, step=%s", rs.Height, rs.Round, rs.RoundStep.String())
	rs.updateStep(RoundStepPropose)
}

func (rs *roundStep) enterPropose() {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	log.Infof("enterPropose, height=%d, round=%d, step=%s", rs.Height, rs.Round, rs.RoundStep.String())

	defer rs.scheduleTimeout(rs.Height, rs.Round, RoundStepPropose, time.Second)
	proposer := rs.Validators.GetProposer(rs.Height, rs.Round)
	if rs.Singer.Address() == proposer.Address {
		log.Infof("is proposer, height %d、round %d", rs.Height, rs.Round)

		var block *pb.Block
		if rs.LockedBlock != nil {
			block = rs.LockedBlock
		} else {
			var err error
			block, err = rs.extension.ProposalBlock(rs.Height, rs.Singer)
			if err != nil {
				log.Panic(err)
			}
		}

		proposal := &pb.Proposal{
			Block:     block,
			Islock:    rs.LockedBlock != nil,
			LockRound: rs.LockedRound,
			Round:     rs.Round,
			Height:    rs.Height,
		}

		signature, err := rs.Singer.SignBlock(block)
		if err != nil {
			log.Panic(err)
		}

		go rs.extension.BroadcastProposal(proposal, signature)
		if err := rs.SetProposal(proposal, signature); err != nil {
			log.Panic(err)
		}
	}
}

// enterPrevote. RoundStepPropose -> RoundStepPrevote
// Prevote Step (height:H,round:R)
// Upon entering Prevote, each validator broadcasts its prevote vote.
//
// First, if the validator is locked on a block since LastLockRound but now has a PoLC for something else at round PoLC-Round where LastLockRound < PoLC-Round < R, then it unlocks.
// If the validator is still locked on a block, it prevotes that.
// Else, if the proposed block from Propose(H,R) is good, it prevotes that.
// Else, if the proposal is invalid or wasn't received on time, it prevotes <nil>.
// The Prevote step ends:
//
// After +2/3 prevotes for a particular block or <nil>. --> goto Precommit(H,R)
// After timeoutPrevote after receiving any +2/3 prevotes. --> goto Precommit(H,R)
// After common exit conditions
func (rs *roundStep) enterPrevote() {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	log.Infof("enterPrevote, height=%d, round=%d, step=%s", rs.Height, rs.Round, rs.RoundStep.String())

	defer func() {
		rs.updateStep(RoundStepPrevote)
		rs.scheduleTimeout(rs.Height, rs.Round, RoundStepPrevote, time.Second)
	}()
	// First, if the validator is locked on a block since LastLockRound
	// but now has a PoLC for something
	// else at round PoLC-Round where LastLockRound < PoLC-Round < R, then it unlocks
	if rs.LockedBlock != nil {
		polcRound, _, ok := rs.HeightVoteSet.POLInfo()
		if ok {
			if rs.LockedRound < polcRound && polcRound <= rs.Round {
				rs.LockedBlock = nil
				rs.LockedRound = 0
			}
		}

		// If the validator is still locked on a block, it prevotes that.
		if rs.LockedBlock != nil {
			rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_Prevote, rs.LockedBlock)
			return
		}
	}

	// Else, if the proposal is invalid or wasn't received on time, it prevotes <nil>.
	if rs.Proposal == nil {
		rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_Prevote, nil)
		return
	}

	if err := rs.extension.ValidateProposalBlock(rs.Proposal); err != nil {
		log.Errorf("ValidateProposalBlock err %s", err.Error())
		rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_Prevote, nil)
		return
	}

	// Else, if the proposed block from Propose(H,R) is good, it prevotes that.
	rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_Prevote, rs.Proposal.GetBlock())
}

func (rs *roundStep) enterPrecommit() {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	log.Infof("enterPrecommit, height=%d, round=%d, step=%s", rs.Height, rs.Round, rs.RoundStep.String())

	defer func() {
		rs.updateStep(RoundStepPrecommit)
		rs.scheduleTimeout(rs.Height, rs.Round, RoundStepPrecommit, time.Second)
	}()

	votes := rs.HeightVoteSet.Prevotes(rs.Round)
	if votes == nil {
		rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_precommit, nil)
		return
	}

	// If we don't have a polka, we must precommit nil.
	blockHash, ok := votes.TwoThirdsMajority()
	if !ok {
		if rs.LockedBlock != nil {
			log.Info("enterPrecommit: No +2/3 prevotes during enterPrecommit while we're locked. Precommitting nil")
		} else {
			log.Info("enterPrecommit: No +2/3 prevotes during enterPrecommit. Precommitting nil.")
		}
		rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_precommit, nil)
		return
	}

	// +2/3 prevoted nil. Unlock and precommit nil.
	if hash.IsZeroHash(blockHash) {
		if rs.LockedBlock == nil {
			log.Info("enterPrecommit: +2/3 prevoted for nil.")
		} else {
			log.Info("enterPrecommit: +2/3 prevoted for nil. Unlocking")
			rs.LockedRound = 0
			rs.LockedBlock = nil
		}
		rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_precommit, nil)
		return
	}

	// At this point, +2/3 prevoted for a particular block.

	// If we're already locked on that block, precommit it, and update the LockedRound
	lockedHash, err := hash.ProtoToSha3(rs.LockedBlock)
	if err == nil && lockedHash == blockHash {
		rs.LockedRound = rs.Round
		rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_precommit, rs.LockedBlock)
		return
	}

	// If +2/3 prevoted for proposal block, stage and precommit it
	proposalHash, err := hash.ProtoToSha3(rs.Proposal.GetBlock())
	if err == nil && proposalHash == blockHash {
		rs.LockedBlock = rs.Proposal.GetBlock()
		rs.LockedRound = rs.Round
		rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_precommit, rs.Proposal.GetBlock())
		return
	}

	// There was a polka in this round for a block we don't have.
	// Fetch that block, unlock, and precommit nil.
	// The +2/3 prevotes for this round is the POL for our unlock.
	// TODO: In the future save the POL prevotes for justification.
	rs.LockedBlock = nil
	rs.LockedRound = 0
	rs.addVoteByBlock(rs.Height, rs.Round, pb.VoteType_precommit, nil)
}

func (rs *roundStep) enterCommit() {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if rs.RoundStep == RoundStepCommit || rs.RoundStep < RoundStepPrecommit {
		return
	}
	rs.updateStep(RoundStepCommit)

	log.Infof("enterCommit, height=%d, round=%d, step=%s", rs.Height, rs.Round, rs.RoundStep.String())
	votes := rs.HeightVoteSet.Precommits(rs.Round)
	if votes == nil {
		rs.Round++
		rs.updateStep(RoundStepNewRound)
		return
	}

	blockHash, ok := votes.TwoThirdsMajority()
	if !ok {
		log.Info("votes not +2/3")
		rs.Round++
		rs.updateStep(RoundStepNewRound)
		return
	}

	if hash.IsZeroHash(blockHash) {
		log.Info("votes +2/3 is zero")
		rs.Round++
		rs.updateStep(RoundStepNewRound)
		return
	}

	lockedHash, err := hash.ProtoToSha3(rs.LockedBlock)
	if err == nil && lockedHash == blockHash {
		log.Info("votes commit locked block")
		if err = rs.extension.Commit(rs.LockedBlock); err != nil {
			log.Panic(err)
		}

		rs.CommitTime = time.Now()
		rs.updateStep(RoundStepNewHeight)
		return
	}

	proposalHash, err := hash.ProtoToSha3(rs.Proposal.GetBlock())
	if err == nil && proposalHash == blockHash {
		log.Info("votes commit block")
		if err = rs.extension.Commit(rs.Proposal.GetBlock()); err != nil {
			log.Panic(err)
		}

		rs.CommitTime = time.Now()
		rs.updateStep(RoundStepNewHeight)
		return
	}

	rs.Round++
	rs.updateStep(RoundStepNewRound)
}

func (rs *roundStep) addVoteByBlock(height, round uint64, vt pb.VoteType, block *pb.Block) {
	vote, signature := rs.genVote(height, round, vt, block)
	go rs.extension.BroadcastVote(vote, signature)
	log.Info("addVoteByBlock  ", hash.BytesToHash(vote.GetHash()).String())
	go rs.SetVote(vote, signature)
}

func (rs *roundStep) genVote(height, round uint64, vt pb.VoteType, block *pb.Block) (*pb.Vote, []byte) {
	var err error
	var blockHash hash.Hash

	if block != nil {
		blockHash, err = hash.ProtoToSha3(block)
		if err != nil {
			log.Error(err)
			blockHash = hash.Hash{}
		}
	}

	vote := &pb.Vote{
		Height:   height,
		Round:    round,
		VoteType: vt,
		Address:  rs.Singer.Address().Bytes(),
		Hash:     blockHash.Bytes(),
	}
	signature, err := rs.Singer.SignVote(vote)
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	return vote, signature
}

func (rs *roundStep) scheduleTimeout(height, round uint64, step RoundStepType, duration time.Duration) {
	var d time.Duration
	if step == RoundStepPropose {
		d = time.Duration(3000+500*round) * time.Millisecond
	} else {
		d = time.Duration(1000+500*round) * time.Millisecond
	}

	ti := timeoutInfo{
		Duration: d,
		Height:   height,
		Round:    round,
		Step:     step,
	}

	rs.ticker.ScheduleTimeout(ti)
}
