package mptrie

import (
	"bytes"
	"testing"
)

func TestEncodeBytes(t *testing.T) {
	cases := []struct {
		b, buf []byte
	}{
		{b: []byte{}, buf: []byte{0x80}},
		{b: nil, buf: []byte{0x80}},
		{b: []byte{0x00}, buf: []byte{0x00}},
		{b: []byte{0x01}, buf: []byte{0x01}},
		{b: []byte{0x02}, buf: []byte{0x02}},
		{b: []byte{0x7F}, buf: []byte{0x7F}},
		{b: []byte{0x80}, buf: []byte{0x81, 0x80}},
		{b: []byte{0x8A}, buf: []byte{0x81, 0x8A}},
		{b: []byte{0x00, 0x11}, buf: []byte{0x82, 0x00, 0x11}},
		{b: []byte{0, 2, 4, 6}, buf: []byte{0x84, 0, 2, 4, 6}},
		{b: makeByteSlice(54), buf: append([]byte{0xB6}, makeByteSlice(54)...)},
		{b: makeByteSlice(55), buf: append([]byte{0xB7}, makeByteSlice(55)...)},
		{b: makeByteSlice(56), buf: append([]byte{0xB8, 0x38}, makeByteSlice(56)...)},
		{b: makeByteSlice(57), buf: append([]byte{0xB8, 0x39}, makeByteSlice(57)...)},
		{b: makeByteSlice(255), buf: append([]byte{0xB8, 0xFF}, makeByteSlice(255)...)},
		{
			b:   makeByteSlice(0xFEED),
			buf: append([]byte{0xB9, 0xFE, 0xED}, makeByteSlice(0xFEED)...),
		},
		{
			b:   makeByteSlice(0x123456),
			buf: append([]byte{0xBA, 0x12, 0x34, 0x56}, makeByteSlice(0x123456)...),
		},
	}

	for _, c := range cases {
		buf := encodeBytes(nil, c.b)
		if !bytes.Equal(buf, c.buf) {
			t.Errorf("encodeBytes(%v): got %v, want %v", c.b, buf, c.buf)
		}
	}
}

func makeByteSlice(l int) []byte {
	bs := make([]byte, l)
	for bi := range bs {
		bs[bi] = byte(bi % 256)
	}
	return bs
}

func TestEncodeUint(t *testing.T) {
	cases := []struct {
		u uint64
		b []byte
	}{
		{u: 0x7A, b: []byte{0x7A}},
		{u: 0xABCD, b: []byte{0xAB, 0xCD}},
		{u: 0xFFFFFF, b: []byte{0xFF, 0xFF, 0xFF}},
		{u: 0, b: []byte{0}},
		{u: 0xFFFEFDFC, b: []byte{0xFF, 0xFE, 0xFD, 0xFC}},
		{u: 0x1020304050, b: []byte{0x10, 0x20, 0x30, 0x40, 0x50}},
		{u: 0x112233445566, b: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}},
		{u: 0xA1B2C3D4E5F607, b: []byte{0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07}},
		{u: 0x0102030405060708, b: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}},
	}

	for _, c := range cases {
		l, b := encodeUint(nil, c.u)
		if l != len(c.b) {
			t.Errorf("encodeUint(%x): got %d, want %d", c.u, l, len(c.b))
		}
		if !bytes.Equal(c.b, b) {
			t.Errorf("encodeUint(%x): got %v, want %v", c.u, b, c.b)
		}
	}
}
