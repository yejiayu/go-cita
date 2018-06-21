package raw

import "testing"

func TestPUT(t *testing.T) {
	db, err := New([]string{"47.75.129.215:2379", "47.75.129.215:2380", "47.75.129.215:2381"})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Put([]byte("test"), []byte("test")); err != nil {
		t.Fatal(err)
	}

	val, err := db.Get([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	if string(val) != "test" {
		t.Fatalf("expect val is test, but got %s", string(val))
	}
	t.Log("done")
}

func TestScan(t *testing.T) {
	db, err := New([]string{"47.75.129.215:2379", "47.75.129.215:2380", "47.75.129.215:2381"})
	if err != nil {
		t.Fatal(err)
	}

	db.Put([]byte{1}, []byte("1"))
	db.Put([]byte{2}, []byte("2"))
	db.Put([]byte{3}, []byte("3"))
	db.Put([]byte{4}, []byte("4"))

	db.Scan([]byte{1}, limit int)
}
