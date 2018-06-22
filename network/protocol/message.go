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

package protocol

import (
	"bytes"
	"encoding/binary"

	"github.com/golang/protobuf/proto"

	"github.com/yejiayu/go-cita/types"
)

// Message for communication between microservices or nodes.
//
// # Message in Bytes:
//
// +-------------`11`--------------------------------------------------+
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

type Message interface {
	Optype() OpType
	Origin() uint32
	Payload() []byte
	Raw() []byte

	UnmarshalStatus() (*types.Status, error)
	UnmarshalSyncResponse() (*types.SyncResponse, error)
}

func NewMessage(opType OpType, origin uint32, data []byte) Message {
	reserved4 := (0 & 0xFC) + (uint8(opType) & 0x3)

	return &message{
		reserved1: 0,
		reserved2: 0,
		reserved3: 0,
		reserved4: reserved4,
		origin:    origin,
		payload:   data,
	}
}

func NewMessageWithRaw(raw []byte) Message {
	return &message{
		reserved1: raw[0],
		reserved2: raw[1],
		reserved3: raw[2],
		reserved4: raw[3],
		origin:    binary.BigEndian.Uint32(raw[3:7]),
		payload:   raw[7:],
	}
}

type message struct {
	reserved1 uint8
	reserved2 uint8
	reserved3 uint8
	reserved4 uint8
	origin    uint32
	payload   []byte
}

func (m *message) Optype() OpType {
	return OpType(m.reserved4 & 0x3)
}

func (m *message) Origin() uint32 {
	return m.origin
}

func (m *message) Raw() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(byte(m.reserved1))
	buf.WriteByte(byte(m.reserved2))
	buf.WriteByte(byte(m.reserved3))
	buf.WriteByte(byte(m.reserved4))

	originBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(originBytes, m.origin)
	buf.Write(originBytes)
	buf.Write(m.payload)
	return buf.Bytes()
}

func (m *message) Payload() []byte {
	return m.payload
}

func (m *message) UnmarshalStatus() (*types.Status, error) {
	var status types.Status

	if err := proto.Unmarshal(m.Payload(), &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (m *message) UnmarshalSyncResponse() (*types.SyncResponse, error) {
	var res types.SyncResponse

	if err := proto.Unmarshal(m.Payload(), &res); err != nil {
		return nil, err
	}

	return &res, nil
}
