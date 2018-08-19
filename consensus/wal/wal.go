package wal

import (
	"context"
	"errors"
)

var (
	ErrNotExists = errors.New("WAL not exists")
)

type LogType uint8

const (
	LogTypeProposal         = 1
	LogTypeVote             = 2
	LogTypeState            = 3
	LogTypePrevHash         = 4
	LogTypeCommits          = 5
	LogTypeVerifiedProposal = 6
	LogTypeAuthTxs          = 7
)

func (l LogType) String() string {
	switch l {
	case LogTypeProposal:
		return "LogTypeProposal"
	case LogTypeVote:
		return "LogTypeVote"
	case LogTypeState:
		return "LogTypeState"
	case LogTypePrevHash:
		return "LogTypePrevHash"
	case LogTypeCommits:
		return "LogTypeCommits"
	case LogTypeVerifiedProposal:
		return "LogTypeVerifiedProposal"
	case LogTypeAuthTxs:
		return "LogTypeAuthTxs"
	default:
		return "LogTypeUnknown" // Cannot panic.
	}
}

type Interface interface {
	SetHeight(ctx context.Context, height uint64) error
	Save(ctx context.Context, logType LogType, data []byte) error
	Load(ctx context.Context) (LogType, []byte, error)
}

func IsNotExists(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "WAL not exists"
}
