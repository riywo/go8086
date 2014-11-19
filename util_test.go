package go8086

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var parityTests = []struct {
	in  uint16
	out bool
}{
	{0x0000, false},
	{0x0001, true},
	{0x0002, true},
	{0x0003, false},
	{0x0004, true},
	{0x0005, false},
	{0x0006, false},
	{0x0007, true},
	{0x0008, true},
	{0x0009, false},
	{0x000a, false},
	{0x000b, true},
	{0x000c, false},
	{0x000d, true},
	{0x000e, true},
	{0x000f, false},
	{0x0000, false},
	{0x0011, false},
	{0x0022, false},
	{0x0033, false},
	{0x0044, false},
	{0x0055, false},
	{0x0066, false},
	{0x0077, false},
	{0x0088, false},
	{0x0099, false},
	{0x00aa, false},
	{0x00bb, false},
	{0x00cc, false},
	{0x00dd, false},
	{0x00ee, false},
	{0x00ff, false},
	{0x0000, false},
	{0x0111, true},
	{0x0222, true},
	{0x0333, false},
	{0x0444, true},
	{0x0555, false},
	{0x0666, false},
	{0x0777, true},
	{0x0888, true},
	{0x0999, false},
	{0x0aaa, false},
	{0x0bbb, true},
	{0x0ccc, false},
	{0x0ddd, true},
	{0x0eee, true},
	{0x0fff, false},
	{0x1000, true},
	{0x2111, false},
	{0x3222, true},
	{0x4333, true},
	{0x5444, true},
	{0x6555, false},
	{0x7666, true},
	{0x8777, false},
	{0x9888, true},
	{0xa999, false},
	{0xbaaa, true},
	{0xcbbb, true},
	{0xdccc, true},
	{0xeddd, false},
	{0xfeee, true},
	{0x0fff, false},
}

func TestParity(t *testing.T) {
	for _, test := range parityTests {
		assert.Equal(t, test.out, Parity(test.in))
	}
}
