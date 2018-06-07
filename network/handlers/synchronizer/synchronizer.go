package synchronizer

import (
	"sync"

	"github.com/yejiayu/go-cita/types"
)

// Synchronizer is Get messages and determine if need to synchronize or broadcast the current node status
type Synchronizer interface {
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
	UpdateCurrentStatus(latestStatus *types.Status) error
	UpdateGlobalStatus(latestStatus *types.Status, origin uint32) error
	ProcessSync(syncResponse *types.SyncResponse) error
}

func NewSynchronizer() Synchronizer {
	return &synchronizer{}
}

type synchronizer struct {
	mu sync.Mutex

	CurrentStatus   *types.Status
	GlobalStatus    *types.Status
	SyncEndHeight   uint64
	IsSynchronizing bool

	LatestStatusList map[uint64][]uint32
	BlockList        map[uint64]*types.Block

	LocalSyncCount uint8
}

func (s *synchronizer) UpdateCurrentStatus(latestStatus *types.Status) error {
	return nil
}

func (s *synchronizer) UpdateGlobalStatus(latestStatus *types.Status, origin uint32) error {
	return nil
}

func (s *synchronizer) ProcessSync(syncResponse *types.SyncResponse) error {
	return nil
}
