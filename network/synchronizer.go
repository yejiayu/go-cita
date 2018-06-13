package network

import (
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/mq"
	"github.com/yejiayu/go-cita/network/connection"
	"github.com/yejiayu/go-cita/types"
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

func newSynchronizer(cm connection.Manager) *synchronizer {
	return &synchronizer{cm: cm}
}

// Synchronizer is Get messages and determine if need to synchronize or broadcast the current node status
type synchronizer struct {
	mu sync.Mutex

	queue mq.Queue
	cm    connection.Manager

	CurrentStatus   *types.Status
	GlobalStatus    *types.Status
	SyncEndHeight   uint64
	IsSynchronizing bool

	LatestStatusList map[uint64][]uint32
	BlockList        map[uint64]*types.Block

	LocalSyncCount uint8
}

// After receiving the `Chain >> Status`, it is processed as follows:
// 1. The chain height suddenly becomes lower than the original,
//    which means that the library is deleted, that is,
//    the synchronization is restarted from the received height.
// 2. The height of the chain is greater than or equal to the original height,
//    less than or equal to the height that the network has synchronized
//        - It is equal to the original height more than 2 times, indicating that
//          the data in the chain or executor is lost, and the block information is
//          sent again from the buffer. If no buffer exists, the data is requested again from other nodes.
//        - When the maximum height of data has been synchronized, a synchronization
//          request is sent to other nodes to confirm that the current node has successfully
//          completed the data synchronization work.
// 3. The height is greater than the maximum height that has been synchronized, ie it is the latest data itself.
func (s *synchronizer) UpdateCurrentStatus(latestStatus *types.Status) error {
	currentHeight := s.CurrentStatus.Height
	latestHeight := latestStatus.Height

	if currentHeight == latestHeight {
		return nil
	}

	s.CurrentStatus = latestStatus
	if err := s.broadcastStatus(s.CurrentStatus); err != nil {
		return err
	}

	return nil
}

func (s *synchronizer) UpdateGlobalStatus(latestStatus *types.Status, origin uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// currentHeight := s.CurrentStatus.Height
	// globalHeight := s.GlobalStatus.Height
	// latestHeight = latestStatus.Height
	//
	// if globalHeight < latestHeight {
	// 	s.GlobalStatus = latestStatus
	// }
	//
	// if latestHeight < currentHeight {
	// 	return nil
	// }

	return nil
}

func (s *synchronizer) ProcessSync(syncResponse *types.SyncResponse) error {
	data, err := proto.Marshal(syncResponse)
	if err != nil {
		return err
	}

	return s.queue.Pub(mq.NetworkSyncBlock, data)
}

func (s *synchronizer) broadcastStatus(status *types.Status) error {
	data, err := proto.Marshal(status)
	if err != nil {
		return err
	}

	s.cm.Broadcast(synchronizerStatus, data)
	return nil
}
