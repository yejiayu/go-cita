package handlers

import (
	"fmt"

	"github.com/yejiayu/go-cita/network/handlers/synchronizer"
	"github.com/yejiayu/go-cita/network/protocol"
)

const (
	// sync
	// Synchronizer >> Status update_global_status
	synchronizerStatus = "synchronizer.status"
	// Synchronizer >> SyncResponse process_sync
	synchronizerSyncResponse = "synchronizer.sync_response"
	// Synchronizer >> SyncRequest
	synchronizerSyncRequest = "synchronize.sync_request"
	// Chain >> Status
	chainStatus = "chainStatus"

	// auth
	// Auth >> Request
	authRequest = "auth.request"

	// consensus
	// Consensus >> SignedProposal
	consensusSignedProposal = "consensus.signed_proposal"
	// Consensus >> RawBytes
	consensusRawBytes = "consensus.raw_bytes"
)

type HandleFunc func(key string, msg *protocol.Message) error

type Interface interface {
	Call(key string, msg *protocol.Message) error
}

type handler struct {
	syncHandler synchronizer.Synchronizer
}

func New() Interface {
	return &handler{
		syncHandler: synchronizer.NewSynchronizer(),
	}
}

func (h *handler) Call(key string, msg *protocol.Message) error {
	switch key {
	case synchronizerStatus:
		status, err := msg.UnmarshalStatus()
		if err != nil {
			return err
		}

		return h.syncHandler.UpdateGlobalStatus(status, msg.Origin())

	case chainStatus:
		status, err := msg.UnmarshalStatus()
		if err != nil {
			return err
		}

		return h.syncHandler.UpdateCurrentStatus(status)
	case synchronizerSyncResponse:
		syncRes, err := msg.UnmarshalSyncResponse()
		if err != nil {
			return err
		}

		return h.syncHandler.ProcessSync(syncRes)
	case synchronizerSyncRequest:
	}

	return fmt.Errorf("method %s not found", key)
}
