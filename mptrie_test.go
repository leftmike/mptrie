package mptrie_test

import (
	"bytes"
	"fmt"
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
	testPut
)

type testCase struct {
	op       testOp
	k, v     []byte
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

			//fmt.Println(mpt.String())

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

		case testPut:
			err := mpt.Put(c.k, c.v)
			if err != nil {
				t.Errorf("mpt.Put(%v, %v) failed with %s", c.k, c.v, err)
			}

			//fmt.Println(mpt.String())

		default:
			panic(fmt.Sprintf("unexpected test op: %v", c.op))
		}
	}

	//fmt.Println(strings.Repeat("-", 16))
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
		})

	key10 := []byte{0x01, 0x23, 0x45, 0x67, 0x89}
	val10 := []byte{0xAA, 0xAA}
	key11 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x01}
	val11 := []byte{0xBB, 0xBB}
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
		})

	testMPTrie(t, mptrie.New(),
		[]testCase{
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testDelete, k: key2, notFound: true},
			{op: testDelete, k: key1},
			{op: testGet, k: key1, notFound: true},
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
		})
}
