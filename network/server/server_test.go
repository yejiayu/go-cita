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
