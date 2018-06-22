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
	"testing"
)

func TestEncode(t *testing.T) {
	c := NewCodec()
	data, err := c.Encode("test", []byte("string"))
	if err != nil {
		t.Fatal(err)
	}

	key, data, err := c.Decode(bytes.NewReader(data))
	if key != "test" {
		t.Fatalf("key %s is not equal to test", key)
	}

	if string(data) != "string" {
		t.Fatalf("data %s is not equal to string", string(data))
	}
}
