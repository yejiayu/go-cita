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
	"io"

	"github.com/yejiayu/go-cita/errors"
)

// Implementation of the multiplexed line-based protocol.
//
// Frames begin with a 4 byte header, consisting of the numeric request ID
// encoded in network order, followed by the frame payload encoded as a UTF-8
// string and terminated with a '\n' character:
//
// # An example frame:
//
// +------------------------+--------------------------+
// | Type                   | Content                  |
// +------------------------+--------------------------+
// | Symbol for Start       | \xDEADBEEF               |
// | Length of Full Payload | u32                      |
// +------------------------+--------------------------+
// | Length of Key          | u8                       |
// | Key                    | bytes of a str           |
// +------------------------+--------------------------+
// | Message                | a serialize data         |
// +------------------------+--------------------------+

// Start of network messages.
const startMsg = 0xDEADBEEF00000000

// Codec definition of cita protocol parsing.
type Codec interface {
	Decode(read io.Reader) (string, []byte, error)
	Encode(key string, data []byte) ([]byte, error)
}

func NewCodec() Codec {
	return codec{}
}

type codec struct{}

func (c codec) Decode(read io.Reader) (string, []byte, error) {
	for {
		data := make([]byte, 8)
		_, err := io.ReadFull(read, data)
		if err != nil {
			return "", nil, err
		}

		requestID := binary.BigEndian.Uint64(data)
		netMsgStart := requestID & 0xffffffff00000000
		lengthFull := requestID & 0x00000000ffffffff

		if netMsgStart != startMsg {
			return "", nil, errors.IllegalMessage.Build("the start message is not 0xDEADBEEF00000000")
		}

		data = make([]byte, lengthFull)
		if _, err = io.ReadFull(read, data); err != nil {
			return "", nil, err
		}

		buf := bytes.NewBuffer(data)
		b, err := buf.ReadByte()
		if err != nil {
			return "", nil, err
		}
		keyLen := uint8(b)
		keyBytes := make([]byte, keyLen)
		if _, err := io.ReadFull(buf, keyBytes); err != nil {
			return "", nil, err
		}

		return string(keyBytes), buf.Bytes(), nil
	}
}

func (c codec) Encode(key string, data []byte) ([]byte, error) {
	keyLen := uint8(len(key))

	// Use 1 bytes to store the length for key, then store key, the last part is body.
	fullLen := 1 + int(keyLen) + len(data)
	requestID := startMsg + uint64(fullLen)

	buf := bytes.NewBuffer(nil)

	requestIDBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(requestIDBytes, requestID)
	if _, err := buf.Write(requestIDBytes); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(byte(keyLen)); err != nil {
		return nil, err
	}
	if _, err := buf.Write([]byte(key)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
