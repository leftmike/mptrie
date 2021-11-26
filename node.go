package mptrie

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/sha3"
)

type node interface {
	encode() []byte
	hash(rf bool) []byte
	toString(w io.Writer, depth int)
}

// XXX: change to keccak256; add internal test for it
func Keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, d := range data {
		h.Write(d)
	}
	// XXX
	//return h.Sum(nil)

	r := h.(io.Reader)
	buf := make([]byte, 32)
	r.Read(buf)
	return buf
}

type leafNode struct {
	suffixKey  nibbleKey
	value      []byte
	generation int64
}

func (leaf *leafNode) encode() []byte {
	return encodeTuple(nil, encodeBytes(nil, encodeHexPrefix(leaf.suffixKey, true)),
		encodeBytes(nil, leaf.value))
}

func (leaf *leafNode) hash(rf bool) []byte {
	// XXX: test for len(encodeHexPrefix) > 32 and/or len(leaf.value) > 32
	buf := leaf.encode()
	if rf {
		return Keccak256(buf)
	}
	if len(buf) < 32 {
		return buf
	}
	return encodeBytes(nil, Keccak256(buf))
}

func (leaf *leafNode) toString(w io.Writer, depth int) {
	fmt.Fprint(w, strings.Repeat("  ", depth))
	fmt.Fprintf(w, "%v = %v\n", leaf.suffixKey, leaf.value)
}

func (mpt *MPTrie) newLeafNode(sk nibbleKey, val []byte) *leafNode { // XXX: maybe return node?
	return &leafNode{
		suffixKey:  sk,
		value:      val,
		generation: mpt.generation,
	}
}

type extensionNode struct {
	subKey     nibbleKey
	child      node
	generation int64
}

func (extension *extensionNode) encode() []byte {
	return encodeTuple(nil, encodeBytes(nil, encodeHexPrefix(extension.subKey, false)),
		extension.child.encode())
}

func (extension *extensionNode) hash(rf bool) []byte {
	// XXX: test for len(encodeHexPrefix) > 32 and/or len(child) > 32
	buf := encodeTuple(nil, encodeBytes(nil, encodeHexPrefix(extension.subKey, false)),
		extension.child.hash(false))
	if rf {
		return Keccak256(buf)
	}
	if len(buf) < 32 {
		return buf
	}
	return encodeBytes(nil, Keccak256(buf))
}

func (extension *extensionNode) toString(w io.Writer, depth int) {
	fmt.Fprint(w, strings.Repeat("  ", depth))
	fmt.Fprintf(w, "%v:\n", extension.subKey)
	extension.child.toString(w, depth+1)
}

func (mpt *MPTrie) newExtensionNode(sk nibbleKey) *extensionNode {
	return &extensionNode{
		subKey:     sk,
		child:      nil, // Must be set by the caller.
		generation: mpt.generation,
	}
}

type branchNode struct {
	children   [16]node
	value      []byte
	generation int64
}

func (branch *branchNode) encode() []byte {
	tuple := make([][]byte, 17)
	for ci := range branch.children {
		if branch.children[ci] == nil {
			tuple[ci] = emptyBytes
		} else {
			tuple[ci] = branch.children[ci].encode()
		}
	}
	if branch.value == nil {
		tuple[16] = emptyBytes
	} else {
		tuple[16] = branch.value
	}

	return encodeTuple(nil, tuple...)
}

func (branch *branchNode) hash(rf bool) []byte {
	tuple := make([][]byte, 17)
	for ci := range branch.children {
		if branch.children[ci] == nil {
			tuple[ci] = emptyBytes
		} else {
			tuple[ci] = branch.children[ci].hash(false)
		}
	}
	if branch.value == nil {
		tuple[16] = emptyBytes
	} else {
		tuple[16] = encodeBytes(nil, branch.value)
	}

	buf := encodeTuple(nil, tuple...)
	if rf {
		return Keccak256(buf)
	}
	if len(buf) < 32 {
		return buf
	}
	return encodeBytes(nil, Keccak256(buf))
}

func (branch *branchNode) toString(w io.Writer, depth int) {
	if branch.value != nil {
		fmt.Fprint(w, strings.Repeat("  ", depth))
		fmt.Fprintf(w, "%v\n", branch.value)
	}

	for idx := range branch.children {
		n := branch.children[idx]
		if n != nil {
			fmt.Fprint(w, strings.Repeat("  ", depth))
			if leaf, ok := n.(*leafNode); ok {
				fmt.Fprintf(w, "[%x] %v = %v\n", idx, leaf.suffixKey, leaf.value)
			} else if extension, ok := n.(*extensionNode); ok {
				fmt.Fprintf(w, "[%x] %v:\n", idx, extension.subKey)
				extension.child.toString(w, depth+1)
			} else if _, ok := n.(*branchNode); ok {
				fmt.Fprintf(w, "[%x]\n", idx)
				n.toString(w, depth+1)
			} else {
				panic(fmt.Sprintf("unexpected mptrie node: %#v", n))
			}
		}
	}
}

func (mpt *MPTrie) newBranchNode() *branchNode {
	return &branchNode{
		generation: mpt.generation,
	}
}
