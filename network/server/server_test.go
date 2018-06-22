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

package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	netMsgStart := 100 & 0xffffffff00000000
	fmt.Println(netMsgStart)
	conn, err := net.Dial("TCP", "47.75.129.215:4001")
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))
}
