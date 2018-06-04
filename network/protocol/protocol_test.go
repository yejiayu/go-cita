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
