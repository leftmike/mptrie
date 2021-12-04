package mptrie

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/sha3"
)

type node interface {
	encode() []byte // XXX: needed?
	hash(rf bool) []byte
	toString(w io.Writer, depth int)
}

func keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
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
	buf := leaf.encode()
	if rf {
		return keccak256(buf)
	}
	if len(buf) < 32 {
		return buf
	}
	return encodeBytes(nil, keccak256(buf))
}

func (leaf *leafNode) toString(w io.Writer, depth int) {
	fmt.Fprint(w, strings.Repeat("  ", depth))
	fmt.Fprintf(w, "%v = %v\n", leaf.suffixKey, leaf.value)
}

func (mpt *MPTrie) newLeafNode(sk nibbleKey, val []byte) node {
	return &leafNode{
		suffixKey:  sk,
		value:      val,
		generation: mpt.generation,
	}
}

type extensionNode struct {
	subKey     nibbleKey
	child      *branchNode // Child will always be a branch node.
	generation int64
}

func (extension *extensionNode) encode() []byte {
	return encodeTuple(nil, encodeBytes(nil, encodeHexPrefix(extension.subKey, false)),
		extension.child.encode())
}

func (extension *extensionNode) hash(rf bool) []byte {
	buf := encodeTuple(nil, encodeBytes(nil, encodeHexPrefix(extension.subKey, false)),
		extension.child.hash(false))
	if rf {
		return keccak256(buf)
	}
	if len(buf) < 32 {
		return buf
	}
	return encodeBytes(nil, keccak256(buf))
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

func (branch *branchNode) noChildren() bool {
	for _, n := range branch.children {
		if n != nil {
			return false
		}
	}
	return true
}

func (branch *branchNode) onlyChild() (nibbleKey, node) {
	var child node
	var ck byte
	for nk, n := range branch.children {
		if n != nil {
			if child != nil {
				return nil, nil
			}
			child = n
			ck = byte(nk)
		}
	}
	return []byte{ck}, child
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
		return keccak256(buf)
	}
	if len(buf) < 32 {
		return buf
	}
	return encodeBytes(nil, keccak256(buf))
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
