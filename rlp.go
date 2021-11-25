package mptrie

func encodeBytes(buf []byte, bs []byte) []byte {
	if len(bs) == 1 && bs[0] < 128 {
		return append(buf, bs[0])
	} else if len(bs) < 56 {
		buf = append(buf, byte(128+len(bs)))
		return append(buf, bs...)
	}

	li := len(buf)
	buf = append(buf, 0)
	var l int
	l, buf = encodeUint(buf, uint64(len(bs)))
	buf[li] = byte(183 + l)
	return append(buf, bs...)
}

func encodeUint(buf []byte, u uint64) (int, []byte) {
	if u <= 0xFF {
		return 1, append(buf, byte(u))
	} else if u <= 0xFFFF {
		buf = append(buf, byte(u>>8))
		return 2, append(buf, byte(u))
	} else if u <= 0xFFFFFF {
		buf = append(buf, byte(u>>16))
		buf = append(buf, byte(u>>8))
		return 3, append(buf, byte(u))
	} else if u <= 0xFFFFFFFF {
		buf = append(buf, byte(u>>24))
		buf = append(buf, byte(u>>16))
		buf = append(buf, byte(u>>8))
		return 4, append(buf, byte(u))
	} else if u <= 0xFFFFFFFFFF {
		buf = append(buf, byte(u>>32))
		buf = append(buf, byte(u>>24))
		buf = append(buf, byte(u>>16))
		buf = append(buf, byte(u>>8))
		return 5, append(buf, byte(u))
	} else if u <= 0xFFFFFFFFFFFF {
		buf = append(buf, byte(u>>40))
		buf = append(buf, byte(u>>32))
		buf = append(buf, byte(u>>24))
		buf = append(buf, byte(u>>16))
		buf = append(buf, byte(u>>8))
		return 6, append(buf, byte(u))
	} else if u <= 0xFFFFFFFFFFFFFF {
		buf = append(buf, byte(u>>48))
		buf = append(buf, byte(u>>40))
		buf = append(buf, byte(u>>32))
		buf = append(buf, byte(u>>24))
		buf = append(buf, byte(u>>16))
		buf = append(buf, byte(u>>8))
		return 7, append(buf, byte(u))
	}

	buf = append(buf, byte(u>>56))
	buf = append(buf, byte(u>>48))
	buf = append(buf, byte(u>>40))
	buf = append(buf, byte(u>>32))
	buf = append(buf, byte(u>>24))
	buf = append(buf, byte(u>>16))
	buf = append(buf, byte(u>>8))
	return 8, append(buf, byte(u))
}

func encodeTuple(buf []byte, ts ...[]byte) []byte {
	ets := make([][]byte, len(ts))
	for ti, bs := range ts {
		ets[ti] = encodeBytes(nil, bs)
	}

	var tl uint64
	for _, bs := range ets {
		tl += uint64(len(bs))
	}

	if tl < 56 {
		buf = append(buf, byte(192+tl))
	} else {
		li := len(buf)
		buf = append(buf, 0)
		var l int
		l, buf = encodeUint(buf, tl)
		buf[li] = byte(247 + l)
	}

	for _, bs := range ets {
		buf = append(buf, bs...)
	}
	return buf
}
