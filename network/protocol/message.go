package protocol

import (
	"net"
	"strconv"
	"strings"
)

// Message for communication between microservices or nodes.
//
// # Message in Bytes:
//
// +---------------------------------------------------------------+
// | Bytes | Type | Function                                       |
// |-------+------+------------------------------------------------|
// |   0   |  u8  | Reserved (For Version?)                        |
// |   1   |  u8  | Reserved                                       |
// |   2   |  u8  | Reserved                                       |
// |-------+------+------------------------------------------------|
// |   3   |  u4  | Reserved                                       |
// |       |  u1  | Reserved                                       |
// |       |  u1  | Compress: true 1, false 0                      |
// |       |  u2  | OperateType                                    |
// |-------+------+------------------------------------------------|
// |  4~7  |  u32 | Origin                                         |
// |-------+------+------------------------------------------------|
// |  8~   |      | Payload (Serialized Data with Compress)        |
// +-------+------+------------------------------------------------+
//
// We DO NOT have to known the contents of payloads (uncompress and deserialize them) if we just
// want to distribute them.
// So we use first 8 bytes to store `OperateType` and `Origin`.
// And we uncompress and deserialize the payloads only before when we use the contents of them.

type OpType uint8

const (
	OpTypeBroadcast = 0
	OpTypeSingle    = 1
	OpTypeSubtract  = 2
)

type Message struct {
	Reserved1 uint8
	Reserved2 uint8
	Reserved3 uint8
	Reserved4 uint8
	Origin    uint32
	Payload   []byte
}

func NewMessage(opType OpType, origin net.IP, data []byte) *Message {
	reserved4 := (0 & 0xFC) + (uint8(opType) & 0x3)

	return &Message{
		Reserved1: 0,
		Reserved2: 0,
		Reserved3: 0,
		Reserved4: reserved4,
		Origin:    ipToInt32(origin),
		Payload:   data,
	}
}

func (m *Message) Optype() OpType {
	return OpType(m.Reserved4 & 0x3)
}

// Deserialize with protobuf
// func (m *Message) Payload()  {
//
// }

func ipToInt32(ip net.IP) uint32 {
	ips := strings.Split(ip.String(), ".")

	var intIP uint32
	for k, v := range ips {
		i, _ := strconv.Atoi(v)

		intIP = intIP | uint32(i)<<uint32(8*(3-k))
	}

	return intIP
}
