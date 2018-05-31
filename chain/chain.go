package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/yejiayu/go-cita/chain/logs"
	"github.com/yejiayu/go-cita/types"
)

type Interface interface {
}

func New(pType types.ProofType) Interface {
	return &chain{}
}

type chain struct {
	BloomsConfig *logs.Config

	CurrentHeader  *types.BlockHeader
	CurrentHeight  uint64
	MaxStoreHeight uint64
	BlockMap       map[uint64]*blockInQueue

	// block cache
	Headers map[uint64]*types.BlockHeader
	Bodies  map[uint64]*types.BlockBody

	// caches
	BlockHashes         map[common.Hash]uint64
	TransactionAdresses map[common.Hash]common.Address
	// pub blocks_blooms: RwLock<HashMap<LogGroupPosition, LogBloomGroup>>,
	Receipts      map[common.Hash]*types.Receipt
	Nodes         []common.Address
	BlockInterval uint64

	QuatoLimit        uint64
	AccountQuatoLimit *accountGasLimit
	CheckQuota        bool

	// pub cache_man: Mutex<CacheManager<CacheId>>,
	// pub polls_filter: Arc<Mutex<PollManager<PollFilter>>>,
	//
	// /// Proof type
	// pub prooftype: u8,
}

type accountGasLimit struct {
	CommonQuatoLimit   uint64
	SpecificQuatoLimit map[string]uint64
	UnknownFields      *types.UnknownFields
	CachedSize         uint32
}
