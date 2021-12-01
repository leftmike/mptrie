package mptrie

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("mptrie: key not found")
	emptyBytes  = encodeBytes(nil, nil)
	emptyHash   = keccak256(emptyBytes)
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

func (mpt *MPTrie) deleteBranch(branch *branchNode, nk nibbleKey) (node, error) {
	if len(nk) == 0 {
		if branch.value == nil {
			return nil, ErrNotFound
		}
		branch.value = nil
	} else if branch.children[nk[0]] != nil {
		n, err := mpt.deleteNode(branch.children[nk[0]], nk[1:])
		if err != nil {
			return nil, err
		}
		branch.children[nk[0]] = n
		if n != nil {
			return branch, nil
		}
	} else {
		return nil, ErrNotFound
	}

	// A child or the value was deleted; maybe this branch can be deleted as well.

	if branch.value == nil {
		if ck, onlyChild := branch.onlyChild(); onlyChild != nil {
			if child, ok := onlyChild.(*branchNode); ok {
				extension := mpt.newExtensionNode(ck)
				extension.child = child
				return extension, nil
			} else if child, ok := onlyChild.(*extensionNode); ok {
				child.subKey = append(ck, child.subKey...)
				return child, nil
			} else if child, ok := onlyChild.(*leafNode); ok {
				child.suffixKey = append(ck, child.suffixKey...)
				return child, nil
			}

			panic(fmt.Sprintf("unexpected mptrie node: %#v", onlyChild))
		}
	} else {
		if branch.noChildren() {
			return mpt.newLeafNode([]byte{}, branch.value), nil
		}
	}

	return branch, nil

}

func (mpt *MPTrie) deleteExtension(extension *extensionNode, nk nibbleKey) (node, error) {
	l := len(extension.subKey)
	if len(nk) < l || !bytes.Equal(nk[:l], extension.subKey) {
		return nil, ErrNotFound
	}

	n, err := mpt.deleteNode(extension.child, nk[l:])
	if err != nil {
		return nil, err
	}
	if n == nil { // XXX: not possible
		return nil, nil
	} else if child, ok := n.(*extensionNode); ok {
		extension.subKey = append(extension.subKey, child.subKey...)
		extension.child = child.child
		return extension, nil
	} else if child, ok := n.(*leafNode); ok {
		child.suffixKey = append(extension.subKey, child.suffixKey...)
		return child, nil
	} else if _, ok := n.(*branchNode); !ok {
		panic(fmt.Sprintf("extension.child must be a branch node: %#v", n))
	}

	extension.child = n
	return extension, nil
}

func (mpt *MPTrie) deleteNode(n node, nk nibbleKey) (node, error) {
	if branch, ok := n.(*branchNode); ok {
		return mpt.deleteBranch(branch, nk)
	} else if extension, ok := n.(*extensionNode); ok {
		return mpt.deleteExtension(extension, nk)
	} else if leaf, ok := n.(*leafNode); ok {
		if bytes.Equal(nk, leaf.suffixKey) {
			return nil, nil
		}

		return nil, ErrNotFound
	}

	panic(fmt.Sprintf("unexpected mptrie node: %#v", n))
}

func (mpt *MPTrie) Delete(key []byte) error {
	mpt.hash = nil

	if mpt.root == nil {
		return ErrNotFound
	}

	n, err := mpt.deleteNode(mpt.root, keyToNibbleKey(key))
	if err != nil {
		return err
	}
	mpt.root = n
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

	mpt.hash = mpt.root.hash(true)
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

func (mpt *MPTrie) Encode() []byte {
	if mpt.root == nil {
		return nil
	}
	return mpt.root.encode()
}
