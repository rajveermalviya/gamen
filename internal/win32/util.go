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
