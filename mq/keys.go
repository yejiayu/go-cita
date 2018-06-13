package mq

// RoutingKey is rabbit routingkey
type RoutingKey string

// network
const (
	NetworkSyncBlock = "sync.blocks"

	NetworkVerifyTxRequest = "net.NetworkVerifyTxRequest"
)

// chain
const (
	ChainSyncStatus = "chain.status"
)

// auth
const (
	AuthVerifyTxResponse = "auth.verify_tx_response"
)
