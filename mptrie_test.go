package mptrie_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/leftmike/mptrie"
)

func TestClone(t *testing.T) {
	mpt := mptrie.New()
	if mpt.Clone() == mpt {
		t.Error("mpt.Clone() == mpt")
	}
	if mpt.Clone().Clone() == mpt {
		t.Error("mpt.Clone().Clone() == mpt")
	}
	if mpt.Clone().Clone() == mpt.Clone() {
		t.Error("mpt.Clone().Clone() == mpt.Clone()")
	}
}

type testOp int

const (
	testDelete testOp = iota
	testGet
	testHash
	testPut
)

type testCase struct {
	op       testOp
	k, v, h  []byte
	notFound bool
}

func testMPTrie(t *testing.T, mpt *mptrie.MPTrie, cases []testCase) {
	t.Helper()

	for _, c := range cases {
		switch c.op {
		case testDelete:
			err := mpt.Delete(c.k)
			if c.notFound {
				if err != mptrie.ErrNotFound {
					t.Errorf("mpt.Delete(%v) returned %v, expected not found", c.k, err)
				}
			} else if err != nil {
				t.Errorf("mpt.Delete(%v) failed with %s", c.k, err)
			}

			if testing.Verbose() {
				fmt.Println(mpt.String())
			}

		case testGet:
			v, err := mpt.Get(c.k)
			if c.notFound {
				if err != mptrie.ErrNotFound {
					t.Errorf("mpt.Get(%v) returned %v, expected not found", c.k, err)
				}
			} else if err != nil {
				t.Errorf("mpt.Get(%v) failed with %s", c.k, err)
			} else if !bytes.Equal(c.v, v) {
				t.Errorf("mpt.Get(%v): got %v, want %v", c.k, v, c.v)
			}

		case testHash:
			h := mpt.Hash()
			if !bytes.Equal(c.h, h) {
				t.Errorf("mpt.Hash(): got %#v, want %#v", h, c.h)
			}

		case testPut:
			err := mpt.Put(c.k, c.v)
			if err != nil {
				t.Errorf("mpt.Put(%v, %v) failed with %s", c.k, c.v, err)
			}

			if testing.Verbose() {
				fmt.Println(mpt.String())
			}

		default:
			panic(fmt.Sprintf("unexpected test op: %v", c.op))
		}
	}

	if testing.Verbose() {
		fmt.Println(strings.Repeat("-", 16))
	}
}

func TestBasic(t *testing.T) {
	key1 := []byte{0x00, 0x12, 0x34}
	val1 := []byte{0x01, 0x23, 0x45}
	key2 := []byte{0xA0, 0x12, 0x34}
	val2 := []byte{0xA1, 0x23, 0x45}

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, notFound: true},
			{op: testPut, k: key2, v: val2},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{
				op: testHash,
				h: []byte{0xe0, 0x17, 0x26, 0x94, 0x41, 0xb7, 0x76, 0x97, 0xa3, 0x2f, 0x0, 0x62,
					0xdc, 0x3f, 0x7a, 0xaa, 0xa9, 0x7, 0x1, 0x6c, 0x64, 0x2a, 0x8e, 0x82, 0xae,
					0x28, 0xc1, 0x72, 0x59, 0xb2, 0x3e, 0xc7},
			},
		})

	key3 := []byte{0x00, 0x23, 0x45}
	val3 := []byte{0x01, 0x01, 0x01}

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, notFound: true},
			{op: testPut, k: key3, v: val3},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{
				op: testHash,
				h: []byte{0x48, 0x30, 0xba, 0x5a, 0x8c, 0x0, 0xae, 0x34, 0xeb, 0x70, 0xe8, 0x90,
					0x30, 0xa4, 0x91, 0x84, 0x28, 0xea, 0x11, 0xba, 0x49, 0x88, 0xc3, 0xa8, 0x4a,
					0xf4, 0x30, 0xee, 0xaa, 0xb1, 0xca, 0xfc},
			},
		})

	key4 := []byte{0x00, 0x34, 0x56}
	val4 := []byte{0x02, 0x03, 0x04}

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, notFound: true},
			{op: testPut, k: key2, v: val2},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{op: testGet, k: key4, notFound: true},
			{op: testPut, k: key4, v: val4},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{op: testGet, k: key4, v: val4},
			{
				op: testHash,
				h: []byte{0x30, 0x7c, 0x7f, 0x43, 0xab, 0x18, 0x12, 0xa6, 0x78, 0x57, 0x45, 0x26,
					0xf2, 0x5f, 0xb1, 0xa6, 0x11, 0xf0, 0xb3, 0x8e, 0xa5, 0xc, 0x59, 0x7d, 0x3a,
					0x44, 0x77, 0x6b, 0x8e, 0x7e, 0x60, 0x13},
			},
		})

	key5 := []byte{0x00}
	val5 := []byte{0x11, 0x22, 0x33, 0x44}

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, notFound: true},
			{op: testPut, k: key3, v: val3},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, notFound: true},
			{op: testPut, k: key4, v: val4},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, notFound: true},
			{op: testPut, k: key5, v: val5},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, v: val5},
			{op: testGet, k: []byte{0x01}, notFound: true},
			{
				op: testHash,
				h: []byte{0xdc, 0xda, 0x5e, 0x3, 0x35, 0xa0, 0xcd, 0x7f, 0xf, 0x78, 0xe0, 0x8f,
					0x66, 0x53, 0xe7, 0x65, 0x27, 0xd, 0xf1, 0xc9, 0x81, 0xaf, 0x89, 0xf3, 0xc2,
					0x0, 0xb5, 0x43, 0x79, 0xfe, 0x2e, 0x1a},
			},
		})

	key6 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x00}
	val6 := []byte{0x66, 0x66}
	key7 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x01}
	val7 := []byte{0x77, 0x77}
	key8 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x10}
	val8 := []byte{0x88, 0x88}
	key9 := []byte{0x01, 0x23, 0x45, 0x00}
	val9 := []byte{0x99, 0x99}
	val9a := []byte{0x9a, 0x9A}

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testGet, k: key6, notFound: true},
			{op: testPut, k: key6, v: val6},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, notFound: true},
			{op: testPut, k: key7, v: val7},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, notFound: true},
			{op: testPut, k: key8, v: val8},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, v: val8},
			{op: testGet, k: key9, notFound: true},
			{op: testPut, k: key9, v: val9},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, v: val8},
			{op: testGet, k: key9, v: val9},
			{op: testPut, k: key9, v: val9a},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, v: val8},
			{op: testGet, k: key9, v: val9a},
			{
				op: testHash,
				h: []byte{0x9d, 0xf, 0xef, 0xbb, 0xce, 0x5a, 0xe7, 0x63, 0xbe, 0x29, 0x77, 0xf0,
					0x81, 0x80, 0xf, 0x22, 0x9e, 0xa8, 0x4e, 0x4c, 0xcf, 0xe0, 0x6d, 0xb7, 0x81,
					0xfb, 0x62, 0xf0, 0x28, 0x19, 0xa4, 0x96},
			},
		})

	key10 := []byte{0x01, 0x23, 0x45, 0x67, 0x89}
	val10 := []byte{0xAA, 0xAA}
	key11 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x01}
	val11 := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC,
		0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB,
		0xCC, 0xDD, 0xEE, 0xFF}
	key12 := []byte{0x01, 0x23, 0x45, 0x67}
	val12 := []byte{0xCC, 0xCC}

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testGet, k: key10, notFound: true},
			{op: testPut, k: key10, v: val10},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, notFound: true},
			{op: testPut, k: key11, v: val11},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, v: val11},
			{op: testGet, k: key12, notFound: true},
			{op: testPut, k: key12, v: val12},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, v: val11},
			{op: testGet, k: key12, v: val12},
			{
				op: testHash,
				h: []byte{0x9a, 0x4b, 0xc7, 0x37, 0x1d, 0x9f, 0xa1, 0xd7, 0xfb, 0x73, 0x0, 0xf2,
					0x8a, 0x1d, 0x47, 0x1d, 0xb9, 0x2e, 0x18, 0x8a, 0xf4, 0xbb, 0xcd, 0x7c, 0x39,
					0xd2, 0x6e, 0x64, 0x4, 0x29, 0x2c, 0x5a},
			},
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testGet, k: key10, notFound: true},
			{op: testPut, k: key10, v: val10},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key12, notFound: true},
			{op: testPut, k: key12, v: val12},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key12, v: val12},
			{op: testGet, k: key11, notFound: true},
			{op: testPut, k: key11, v: val11},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, v: val11},
			{op: testGet, k: key12, v: val12},
			{
				op: testHash,
				h: []byte{0x9a, 0x4b, 0xc7, 0x37, 0x1d, 0x9f, 0xa1, 0xd7, 0xfb, 0x73, 0x0, 0xf2,
					0x8a, 0x1d, 0x47, 0x1d, 0xb9, 0x2e, 0x18, 0x8a, 0xf4, 0xbb, 0xcd, 0x7c, 0x39,
					0xd2, 0x6e, 0x64, 0x4, 0x29, 0x2c, 0x5a},
			},
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testDelete, k: key2, notFound: true},
			{op: testDelete, k: key1},
			{op: testGet, k: key1, notFound: true},
			{
				op: testHash,
				h: []byte{0x56, 0xe8, 0x1f, 0x17, 0x1b, 0xcc, 0x55, 0xa6, 0xff, 0x83, 0x45, 0xe6,
					0x92, 0xc0, 0xf8, 0x6e, 0x5b, 0x48, 0xe0, 0x1b, 0x99, 0x6c, 0xad, 0xc0, 0x1,
					0x62, 0x2f, 0xb5, 0xe3, 0x63, 0xb4, 0x21},
			},
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testPut, k: key1, v: val1},
			{op: testPut, k: key2, v: val2},

			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},

			{op: testDelete, k: key3, notFound: true},
			{op: testDelete, k: key1},

			{op: testGet, k: key1, notFound: true},
			{op: testGet, k: key2, v: val2},

			{
				op: testHash,
				h: []byte{0xb, 0x55, 0xd3, 0x5, 0x29, 0x2c, 0xff, 0x54, 0x1b, 0x71, 0x3d, 0x8c,
					0xf4, 0xdc, 0x36, 0xbb, 0xa4, 0x92, 0x24, 0xde, 0xdb, 0x4f, 0x88, 0x9d, 0xf3,
					0x14, 0x27, 0xe7, 0x4, 0xab, 0xaf, 0x17},
			},
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testPut, k: key1, v: val1},
			{op: testPut, k: key2, v: val2},
			{op: testPut, k: key3, v: val3},
			{op: testPut, k: key4, v: val4},
			{op: testPut, k: key5, v: val5},
			{op: testPut, k: key6, v: val6},

			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, v: val5},
			{op: testGet, k: key6, v: val6},

			{op: testDelete, k: key6},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, v: val5},
			{op: testGet, k: key6, notFound: true},

			{op: testDelete, k: key5},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, notFound: true},
			{op: testGet, k: key6, notFound: true},

			{op: testDelete, k: key2},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, notFound: true},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, notFound: true},
			{op: testGet, k: key6, notFound: true},

			{
				op: testHash,
				h: []byte{0xb8, 0xea, 0x27, 0x74, 0x14, 0xe, 0x9f, 0x82, 0xd0, 0xf3, 0x6c, 0x2d,
					0x43, 0x98, 0x5c, 0x14, 0xbb, 0xb0, 0x54, 0xaa, 0x22, 0x8c, 0x3b, 0x62, 0x54,
					0x21, 0x97, 0x4b, 0xe1, 0xb9, 0x76, 0xe},
			},
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testPut, k: key1, v: val1},
			{op: testPut, k: key3, v: val3},
			{op: testPut, k: key4, v: val4},

			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},

			{op: testDelete, k: key5, notFound: true},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},

			{op: testDelete, k: key3},

			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, notFound: true},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, notFound: true},

			{op: testDelete, k: key3, notFound: true},

			{
				op: testHash,
				h: []byte{0xa8, 0x8d, 0x8b, 0x11, 0x6d, 0x81, 0x35, 0x4e, 0xa8, 0x67, 0x44, 0x8c,
					0xa2, 0xcd, 0x57, 0x52, 0xe1, 0xd1, 0xf8, 0xba, 0x90, 0x3f, 0xb1, 0xb0, 0x1c,
					0x34, 0x72, 0xe4, 0xf9, 0x73, 0x8c, 0x57},
			},
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testDelete, k: key1, notFound: true},

			{op: testPut, k: key1, v: val1},
			{op: testPut, k: key5, v: val5},

			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key5, v: val5},

			{op: testDelete, k: key1},
			{op: testGet, k: key1, notFound: true},
			{op: testGet, k: key5, v: val5},

			{
				op: testHash,
				h: []byte{0x8b, 0xd3, 0x7b, 0xe0, 0x43, 0x3d, 0x3e, 0xb4, 0xb0, 0x18, 0xf6, 0x16,
					0x4b, 0x4c, 0x96, 0x26, 0x8c, 0x21, 0xa6, 0x2f, 0x29, 0x6e, 0x99, 0x8e, 0x24,
					0x8e, 0xf, 0x53, 0x89, 0x67, 0x87, 0x38},
			},
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testPut, k: key6, v: val6},
			{op: testPut, k: key7, v: val7},
			{op: testPut, k: key8, v: val8},

			{op: testDelete, k: key9, notFound: true},
			{op: testPut, k: key9, v: val9},

			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, v: val8},
			{op: testGet, k: key9, v: val9},

			{
				op: testHash,
				h: []byte{0xb2, 0xcc, 0x73, 0xce, 0x1a, 0x58, 0x7e, 0x93, 0x30, 0xb6, 0x15, 0x46,
					0x92, 0x27, 0x5e, 0x75, 0xad, 0xa4, 0xbc, 0xac, 0x86, 0xe8, 0x3a, 0xbd, 0x33,
					0xc7, 0xef, 0xfa, 0xd8, 0x69, 0x75, 0xd9},
			},
		})

	key13 := []byte{0x01, 0x23, 0x45, 0x67}
	val13 := []byte{0x13, 0x13}
	key14 := []byte{0x01, 0x23, 0x01, 0x23, 0x45}
	val14 := []byte{0x14, 0x14}
	key15 := []byte{0x01, 0x23, 0x01, 0x23, 0x01}
	val15 := []byte{0x15, 0x15}

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testPut, k: key13, v: val13},
			{op: testPut, k: key14, v: val14},
			{op: testPut, k: key15, v: val15},
			{op: testGet, k: key13, v: val13},
			{op: testGet, k: key14, v: val14},
			{op: testGet, k: key15, v: val15},

			{op: testDelete, k: key13},

			{op: testGet, k: key13, notFound: true},
			{op: testGet, k: key14, v: val14},
			{op: testGet, k: key15, v: val15},

			{
				op: testHash,
				h: []byte{0x6e, 0xa8, 0xe9, 0x73, 0xf2, 0xf0, 0xe8, 0xf, 0xe9, 0x49, 0x7a, 0xdd,
					0x64, 0xa3, 0x2d, 0x11, 0xa5, 0x7e, 0xc, 0x14, 0x8, 0x94, 0xaa, 0xdf, 0x95,
					0x2c, 0x60, 0x38, 0x19, 0xfc, 0xce, 0x5d},
			},
		})
}
