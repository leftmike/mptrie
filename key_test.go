package mptrie

import (
	"bytes"
	"testing"
)

func TestKeyToNibbleKey(t *testing.T) {
	nk := keyToNibbleKey([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF})
	want := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF}

	if !bytes.Equal(nk, want) {
		t.Errorf("keyToNibbleKey: got %v, want %v", nk, want)
	}
}

func TestCommonPrefix(t *testing.T) {
	cases := []struct {
		k1, k2 nibbleKey
		l      int
	}{
		{nibbleKey("ABCDE"), nibbleKey("ABCXYZ"), 3},
		{nibbleKey("ABCDE"), nibbleKey("ABC"), 3},
		{nibbleKey("ABC"), nibbleKey("ABCXYZ"), 3},
		{nibbleKey(""), nibbleKey("ABCXYZ"), 0},
		{nibbleKey("ABC"), nibbleKey(""), 0},
		{nibbleKey("ABC"), nibbleKey("XYZ"), 0},
	}

	for _, c := range cases {
		l := commonPrefix(c.k1, c.k2)
		if l != c.l {
			t.Errorf("commonPrefix(%v, %v): got %d, want %d", c.k1, c.k2, l, c.l)
		}
	}
}
