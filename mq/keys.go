package mq

// RoutingKey is rabbit routingkey
type RoutingKey string

// network
const (
	NetworkSyncBlock = "sync.blocks"

	NetworkUntx = "sync.untx"
)

// chain
const (
	ChainSyncStatus = "chain.status"
)

// auth
const (
	AuthUntx = "auth.untx"
)
