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

func encodeHexPrefix(nk nibbleKey, tf bool) []byte {
	var ni int
	var b byte
	if tf {
		b |= 0x20
	}

	l := len(nk)
	if l%2 == 1 {
		b |= 0x10
		b |= nk[0]
		ni = 1
	}

	buf := make([]byte, 0, (l/2)+1)
	buf = append(buf, b)

	for ni < l {
		buf = append(buf, (nk[ni]<<4)|nk[ni+1])
		ni += 2
	}

	return buf
}
