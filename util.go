package go8086

func ParityOf(v uint16) uint16 {
	v ^= v >> 8
	v ^= v >> 4
	v &= 0xf
	return (0x6996 >> v) & 1
}

func SignOf(v uint16, w Bit) (s uint16) {
	switch w {
	case Bit8:
		s = (v >> 7) & 1
	case Bit16:
		s = v >> 15
	}
	return
}
