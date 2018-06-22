// Copyright (C) 2018 yejiayu

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
	m.Run(make(chan error))

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
