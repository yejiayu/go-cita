package protocol

import (
	"net"
	"testing"
)

func TestIpToInt(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	num := ipToInt32(ip)
	t.Log(num)
}
