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

func isSurrogatedCharacter(x rune) bool {
	return x > 0xd800 // Surrogate characters are mounted after 0xd800
}

// surrogatedUtf16toRune recovers code points from high and low surrogates
func surrogatedUtf16toRune(high rune, low rune) rune {
	high -= 0xd800
	low -= 0xdc00
	return (high << 10) + low + 0x10000
}
