package mq

// RoutingKey is rabbit routingkey
type RoutingKey string

// network
const (
	SyncUnverifiedTx = "sync.unverified_tx"
	SyncBlock        = "sync.blocks"

	NetworkUnverifiedTx = "net.unverified_tx"
)

// chain
const (
	ChainSyncStatus = "chain.status"
)

// auth
const (
	AuthUnverifiedTx = "auth.unverified_tx"
)
