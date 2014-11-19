package go8086

func Parity(v uint16) bool {
	v ^= v >> 8
	v ^= v >> 4
	v &= 0xf
	return ((0x6996 >> v) & 1) == 1
}

func SignOf(v uint16, w Bit) (s bool) {
	switch w {
	case Bit8:
		s = int8(w) < 0
	case Bit16:
		s = int16(w) < 0
	}
	return
}

func LSBOf(v uint16) uint16 {
	return v & 0xff
}

func MSBOf(v uint16) uint16 {
	return (v >> 8) & 0xff
}
