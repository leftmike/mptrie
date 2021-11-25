package mptrie

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("mptrie: key not found")
	emptyHash   = keccak256([]byte{0x80})
)

type MPTrie struct {
	root       node
	generation int64
	hash       []byte
}

func New() *MPTrie {
	return &MPTrie{}
}

func (mpt *MPTrie) String() string {
	if mpt.root == nil {
		return "<nil>\n"
	}

	var sb strings.Builder
	mpt.root.toString(&sb, 0)
	return sb.String()
}

func (mpt *MPTrie) Clone() *MPTrie {
	mpt.generation += 1
	clone := *mpt
	return &clone
}

func (mpt *MPTrie) Delete(key []byte) error {
	mpt.hash = nil

	nk := keyToNibbleKey(key)
	_ = nk
	return nil
}

func (mpt *MPTrie) Get(key []byte) ([]byte, error) {
	nk := keyToNibbleKey(key)
	n := mpt.root

	for n != nil {
		if branch, ok := n.(*branchNode); ok {
			if len(nk) == 0 {
				if branch.value == nil {
					return nil, ErrNotFound
				}
				return branch.value, nil
			}

			n = branch.children[nk[0]]
			nk = nk[1:]
		} else if extension, ok := n.(*extensionNode); ok {
			l := len(extension.subKey)
			if len(nk) < l || !bytes.Equal(nk[:l], extension.subKey) {
				return nil, ErrNotFound
			}

			nk = nk[l:]
			n = extension.child
		} else if leaf, ok := n.(*leafNode); ok {
			if bytes.Equal(nk, leaf.suffixKey) {
				return leaf.value, nil
			}

			return nil, ErrNotFound
		} else {
			panic(fmt.Sprintf("unexpected mptrie node: %#v", n))
		}
	}

	return nil, ErrNotFound
}

func (mpt *MPTrie) Hash() []byte {
	if mpt.hash != nil {
		return mpt.hash
	}

	if mpt.root == nil {
		return emptyHash
	}

	// mpt.hash = mpt.root.hash()
	return mpt.hash
}

func (mpt *MPTrie) Put(key, val []byte) error {
	mpt.hash = nil

	nk := keyToNibbleKey(key)
	pn := &mpt.root

	for (*pn) != nil {
		if branch, ok := (*pn).(*branchNode); ok {
			if len(nk) == 0 {
				branch.value = val
				return nil
			}

			pn = &branch.children[nk[0]]
			nk = nk[1:]
		} else if extension, ok := (*pn).(*extensionNode); ok {
			cpl := commonPrefix(nk, extension.subKey)
			if cpl == len(extension.subKey) {
				pn = &extension.child
				nk = nk[cpl:]
			} else {
				if cpl > 0 {
					newExtension := mpt.newExtensionNode(extension.subKey[:cpl])
					*pn = newExtension
					pn = &newExtension.child
				}

				newBranch := mpt.newBranchNode()
				*pn = newBranch
				if len(extension.subKey) == cpl+1 {
					newBranch.children[extension.subKey[cpl]] = extension.child
				} else {
					newBranch.children[extension.subKey[cpl]] = extension
					extension.subKey = extension.subKey[cpl+1:]
				}

				if len(nk) == cpl {
					newBranch.value = val
					return nil
				}

				pn = &newBranch.children[nk[cpl]]
				nk = nk[cpl+1:]
				break
			}
		} else if leaf, ok := (*pn).(*leafNode); ok {
			if bytes.Equal(nk, leaf.suffixKey) {
				leaf.value = val
				return nil
			}

			cpl := commonPrefix(nk, leaf.suffixKey)
			if cpl > 0 {
				newExtension := mpt.newExtensionNode(nk[:cpl])
				*pn = newExtension
				pn = &newExtension.child
				nk = nk[cpl:]
				leaf.suffixKey = leaf.suffixKey[cpl:]
			}

			newBranch := mpt.newBranchNode()
			*pn = newBranch
			if len(leaf.suffixKey) == 0 {
				newBranch.value = leaf.value
			} else {
				newBranch.children[leaf.suffixKey[0]] = leaf
				leaf.suffixKey = leaf.suffixKey[1:]

				if len(nk) == 0 {
					newBranch.value = val
					return nil
				}
			}

			pn = &newBranch.children[nk[0]]
			nk = nk[1:]
			break
		} else {
			panic(fmt.Sprintf("unexpected mptrie node: %#v", *pn))
		}
	}

	*pn = mpt.newLeafNode(nk, val)
	return nil
}
