package mptrie

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/sha3"
)

type node interface {
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
