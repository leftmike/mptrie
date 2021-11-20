package mptrie

type nibbleKey []byte

type hexPrefixKey []byte

func keyToNibbleKey(key []byte) nibbleKey {
	nk := make(nibbleKey, 0, len(key)*2)
	for _, b := range key {
		nk = append(nk, byte(b>>4))
		nk = append(nk, byte(b%16))
	}
	return nk
}

func commonPrefix(k1, k2 nibbleKey) int {
	l := 0
	for l < len(k1) && l < len(k2) && k1[l] == k2[l] {
		l += 1
	}

	return l
}
