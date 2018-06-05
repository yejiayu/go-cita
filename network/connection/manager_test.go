package connection

import (
	"testing"

	networkConfig "github.com/yejiayu/go-cita/config/network"
)

func TestManager(t *testing.T) {
	config := networkConfig.Config{
		ID: 1,
		Peers: []networkConfig.Peer{
			{
				ID:   2,
				IP:   "115.239.210.27",
				Port: 80,
			},
		},
	}
	m := NewManager(config)
	m.Run()

	peers := m.Available()
	if len(peers) != 1 {
		t.Logf("The expectation of the peer length is one, but got %d", len(peers))
	}

	config2 := networkConfig.Config{
		ID: 1,
		Peers: []networkConfig.Peer{
			{
				ID:   3,
				IP:   "140.205.220.96",
				Port: 80,
			},
		},
	}

	m.UpdateConfig(config2)

	peers = m.Available()

	if len(peers) != 1 {
		t.Logf("The expectation of the peer length is one, but got %d", len(peers))
	}

	peer := peers[0]
	if peer.ID != 3 || peer.IP != "115.239.210.27" || peer.Port != 80 {
		t.Logf("Not the expected peer, %-v", peer)
	}
}
