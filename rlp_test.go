package mptrie

import (
	"bytes"
	"testing"
)

func TestEncodeBytes(t *testing.T) {
	cases := []struct {
		bs, buf []byte
	}{
		{bs: []byte{}, buf: []byte{0x80}},
		{bs: nil, buf: []byte{0x80}},
		{bs: []byte{0x00}, buf: []byte{0x00}},
		{bs: []byte{0x01}, buf: []byte{0x01}},
		{bs: []byte{0x02}, buf: []byte{0x02}},
		{bs: []byte{0x7F}, buf: []byte{0x7F}},
		{bs: []byte{0x80}, buf: []byte{0x81, 0x80}},
		{bs: []byte{0x8A}, buf: []byte{0x81, 0x8A}},
		{bs: []byte{0x00, 0x11}, buf: []byte{0x82, 0x00, 0x11}},
		{bs: []byte{0, 2, 4, 6}, buf: []byte{0x84, 0, 2, 4, 6}},
		{bs: makeByteSlice(54), buf: append([]byte{0xB6}, makeByteSlice(54)...)},
		{bs: makeByteSlice(55), buf: append([]byte{0xB7}, makeByteSlice(55)...)},
		{bs: makeByteSlice(56), buf: append([]byte{0xB8, 0x38}, makeByteSlice(56)...)},
		{bs: makeByteSlice(57), buf: append([]byte{0xB8, 0x39}, makeByteSlice(57)...)},
		{bs: makeByteSlice(255), buf: append([]byte{0xB8, 0xFF}, makeByteSlice(255)...)},
		{
			bs:  makeByteSlice(0xFEED),
			buf: append([]byte{0xB9, 0xFE, 0xED}, makeByteSlice(0xFEED)...),
		},
		{
			bs:  makeByteSlice(0x123456),
			buf: append([]byte{0xBA, 0x12, 0x34, 0x56}, makeByteSlice(0x123456)...),
		},
	}

	for _, c := range cases {
		buf := encodeBytes(nil, c.bs)
		if !bytes.Equal(buf, c.buf) {
			t.Errorf("encodeBytes(%v): got %v, want %v", c.bs, buf, c.buf)
		}

		buf = encodeBytes([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, c.bs)
		want := append([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, c.buf...)
		if !bytes.Equal(want, buf) {
			t.Errorf("encodeBytes(%v): got %v, want %v", c.bs, buf, want)
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
		u   uint64
		buf []byte
	}{
		{u: 0x7A, buf: []byte{0x7A}},
		{u: 0xABCD, buf: []byte{0xAB, 0xCD}},
		{u: 0xFFFFFF, buf: []byte{0xFF, 0xFF, 0xFF}},
		{u: 0, buf: []byte{0}},
		{u: 0xFFFEFDFC, buf: []byte{0xFF, 0xFE, 0xFD, 0xFC}},
		{u: 0x1020304050, buf: []byte{0x10, 0x20, 0x30, 0x40, 0x50}},
		{u: 0x112233445566, buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}},
		{u: 0xA1B2C3D4E5F607, buf: []byte{0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07}},
		{u: 0x0102030405060708, buf: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}},
	}

	for _, c := range cases {
		l, b := encodeUint(nil, c.u)
		if l != len(c.buf) {
			t.Errorf("encodeUint(%x): got %d, want %d", c.u, l, len(c.buf))
		}
		if !bytes.Equal(c.buf, b) {
			t.Errorf("encodeUint(%x): got %v, want %v", c.u, b, c.buf)
		}
	}
}

func TestEncodeTuple(t *testing.T) {
	cases := []struct {
		ts  [][]byte
		buf []byte
	}{
		{ts: [][]byte{}, buf: []byte{0xC0}},
		{ts: nil, buf: []byte{0xC0}},
		{
			ts:  [][]byte{{0x01, 0x02, 0x03}, nil, {0x04, 0x05, 0x06, 0x07}},
			buf: []byte{0xCA, 0x83, 0x01, 0x02, 0x03, 0x80, 0x84, 0x04, 0x05, 0x06, 0x07},
		},
		{ts: [][]byte{nil}, buf: []byte{0xC1, 0x80}},
		{ts: [][]byte{nil, nil}, buf: []byte{0xC2, 0x80, 0x80}},
		{
			ts:  [][]byte{nil, {0xFF}, nil, {0xAA}},
			buf: []byte{0xC6, 0x80, 0x81, 0xFF, 0x80, 0x81, 0xAA},
		},
		{
			ts: [][]byte{makeByteSlice(11), makeByteSlice(20), makeByteSlice(20)},
			buf: []byte{0xF6, 0x8B, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
				0x0A, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12,
				0x13},
		},
		{
			ts: [][]byte{makeByteSlice(12), makeByteSlice(20), makeByteSlice(20)},
			buf: []byte{0xF7, 0x8C, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x94, 0x00,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D,
				0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13},
		},
		{
			ts: [][]byte{makeByteSlice(13), makeByteSlice(20), makeByteSlice(20)},
			buf: []byte{0xF8, 0x38, 0x8D, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13,
				0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
				0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13},
		},
		{
			ts: [][]byte{makeByteSlice(14), makeByteSlice(20), makeByteSlice(20)},
			buf: []byte{0xF8, 0x39, 0x8E, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
				0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12,
				0x13, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
				0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13},
		},
		{
			ts: [][]byte{makeByteSlice(15), makeByteSlice(20), makeByteSlice(20)},
			buf: []byte{0xF8, 0x3A, 0x8F, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11,
				0x12, 0x13, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
				0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13},
		},
		{
			ts: [][]byte{makeByteSlice(20), makeByteSlice(20), makeByteSlice(20)},
			buf: []byte{0xF8, 0x3F, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x94,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
				0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x94, 0x00, 0x01, 0x02, 0x03, 0x04,
				0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11,
				0x12, 0x13},
		},
	}

	for _, c := range cases {
		buf := encodeTuple(nil, c.ts...)
		if !bytes.Equal(buf, c.buf) {
			t.Errorf("encodeTuple(%v): got %v, want %v", c.ts, buf, c.buf)
		}

		buf = encodeTuple([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, c.ts...)
		want := append([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}, c.buf...)
		if !bytes.Equal(want, buf) {
			t.Errorf("encodeTuple(%v): got %v, want %v", c.ts, buf, want)
		}
	}
}
