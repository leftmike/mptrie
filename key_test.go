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

func TestEncodeHexPrefix(t *testing.T) {
	cases := []struct {
		nk  nibbleKey
		tf  bool
		buf []byte
	}{
		{nk: []byte{0x01}, tf: true, buf: []byte{0x31}},
		{nk: []byte{0x01}, tf: false, buf: []byte{0x11}},
		{nk: []byte{0x01, 0x02}, tf: true, buf: []byte{0x20, 0x12}},
		{nk: []byte{0x01, 0x02}, tf: false, buf: []byte{0x00, 0x12}},
		{nk: []byte{0x0A, 0x0B, 0x0C}, tf: true, buf: []byte{0x3A, 0xBC}},
		{nk: []byte{0x0A, 0x0B, 0x0C}, tf: false, buf: []byte{0x1A, 0xBC}},
		{nk: []byte{0x0A, 0x0B, 0x0C, 0x0D}, tf: true, buf: []byte{0x20, 0xAB, 0xCD}},
		{nk: []byte{0x0A, 0x0B, 0x0C, 0x0D}, tf: false, buf: []byte{0x00, 0xAB, 0xCD}},
		{nk: []byte{0x01, 0x02, 0x03, 0x04, 0x05}, tf: true, buf: []byte{0x31, 0x23, 0x45}},
		{nk: []byte{0x01, 0x02, 0x03, 0x04, 0x05}, tf: false, buf: []byte{0x11, 0x23, 0x45}},
	}

	for _, c := range cases {
		buf := encodeHexPrefix(c.nk, c.tf)
		if !bytes.Equal(buf, c.buf) {
			t.Errorf("encodeHexPrefix(%v %v): got %v, want %v", c.nk, c.tf, buf, c.buf)
		}
	}
}
