package go8086

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var parityOfTests = []struct {
	in  uint16
	out uint16
}{
	{0x0000, 0},
	{0x0001, 1},
	{0x0002, 1},
	{0x0003, 0},
	{0x0004, 1},
	{0x0005, 0},
	{0x0006, 0},
	{0x0007, 1},
	{0x0008, 1},
	{0x0009, 0},
	{0x000a, 0},
	{0x000b, 1},
	{0x000c, 0},
	{0x000d, 1},
	{0x000e, 1},
	{0x000f, 0},
	{0x0000, 0},
	{0x0011, 0},
	{0x0022, 0},
	{0x0033, 0},
	{0x0044, 0},
	{0x0055, 0},
	{0x0066, 0},
	{0x0077, 0},
	{0x0088, 0},
	{0x0099, 0},
	{0x00aa, 0},
	{0x00bb, 0},
	{0x00cc, 0},
	{0x00dd, 0},
	{0x00ee, 0},
	{0x00ff, 0},
	{0x0000, 0},
	{0x0111, 1},
	{0x0222, 1},
	{0x0333, 0},
	{0x0444, 1},
	{0x0555, 0},
	{0x0666, 0},
	{0x0777, 1},
	{0x0888, 1},
	{0x0999, 0},
	{0x0aaa, 0},
	{0x0bbb, 1},
	{0x0ccc, 0},
	{0x0ddd, 1},
	{0x0eee, 1},
	{0x0fff, 0},
	{0x1000, 1},
	{0x2111, 0},
	{0x3222, 1},
	{0x4333, 1},
	{0x5444, 1},
	{0x6555, 0},
	{0x7666, 1},
	{0x8777, 0},
	{0x9888, 1},
	{0xa999, 0},
	{0xbaaa, 1},
	{0xcbbb, 1},
	{0xdccc, 1},
	{0xeddd, 0},
	{0xfeee, 1},
	{0x0fff, 0},
}

func TestParityOf(t *testing.T) {
	for _, test := range parityOfTests {
		assert.Equal(t, test.out, ParityOf(test.in))
	}
}

var signOfTests = []struct {
	in  uint16
	bit Bit
	out uint16
}{
	{0x1, Bit16, 0},
	{0xf004, Bit16, 1},
	{0xf004, Bit8, 0},
}

func TestSignOf(t *testing.T) {
	for _, test := range signOfTests {
		assert.Equal(t, test.out, SignOf(test.in, test.bit))
	}
}
