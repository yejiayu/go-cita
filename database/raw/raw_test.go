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
}
