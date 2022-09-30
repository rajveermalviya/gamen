package win32

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func loword(x uint32) uint16 {
	return uint16(x & 0xFFFF)
}

func hiword(x uint32) uint16 {
	return uint16((x >> 16) & 0xFFFF)
}

func decodeUtf16(s uint16) rune {
	const (
		// 0xd800-0xdc00 encodes the high 10 bits of a pair.
		// 0xdc00-0xe000 encodes the low 10 bits of a pair.
		// the value is those 20 bits plus 0x10000.
		surr1 = 0xd800
		surr3 = 0xe000

		// Unicode replacement character
		replacementChar = '\uFFFD'
	)

	var a rune
	switch r := s; {
	case r < surr1, surr3 <= r:
		// normal rune
		a = rune(r)
	default:
		// invalid surrogate sequence
		a = replacementChar
	}
	return a
}
