package mptrie

func encodeBytes(buf []byte, b []byte) []byte {
	if len(b) == 1 && b[0] < 128 {
		return append(buf, b[0])
	} else if len(b) < 56 {
		buf = append(buf, byte(128+len(b)))
		return append(buf, b...)
	}

	li := len(buf)
	buf = append(buf, 0)
	var l int
	l, buf = encodeUint(buf, uint64(len(b)))
	buf[li] = byte(183 + l)
	return append(buf, b...)
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

func encodeSequence(buf []byte, s ...[]byte) []byte {
	return nil
}
