package config

type Config struct {
	ID    uint32
	Port  uint32
	Peers []Peer
}

type Peer struct {
	ID   uint32
	IP   string
	Port uint32
}
