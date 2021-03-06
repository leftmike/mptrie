package mptrie

import (
	"bytes"
	"testing"
)

func TestKeccak256(t *testing.T) {
	cases := []struct {
		b, h []byte
	}{
		{
			b: []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
			h: []byte{0x7f, 0x4b, 0x26, 0x3f, 0xb8, 0xb6, 0x2, 0x16, 0x72, 0xf, 0xdb, 0x72, 0xca,
				0xcf, 0xd3, 0x1a, 0x67, 0xe8, 0x9e, 0x70, 0x9e, 0x47, 0xc7, 0x72, 0x3d, 0x9, 0xd4,
				0x3b, 0x32, 0x2c, 0x3c, 0xbc},
		},
		{
			b: []byte("0123456789012345678901234567890123456789"),
			h: []byte{0x4f, 0xab, 0xaf, 0x72, 0x4b, 0x17, 0x3, 0xd9, 0x5d, 0x7c, 0x17, 0xa7, 0xf,
				0x71, 0xbe, 0x58, 0xca, 0xb7, 0xbc, 0x8f, 0x32, 0xfc, 0x45, 0x1, 0x40, 0x3d, 0x9c,
				0x50, 0x8f, 0x13, 0xdf, 0xf9},
		},
	}

	for _, c := range cases {
		h := keccak256(c.b)
		if !bytes.Equal(h, c.h) {
			t.Errorf("Keccack256(%#v): got %#v, want %#v", c.b, h, c.h)
		}
	}
}
