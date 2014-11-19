package go8086

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var opcodeADDTests = []struct {
	opr1 Operand
	opr2 Operand
	out  uint16
}{
	{AX, CX, 0x2244},
	{AL, CL, 0x44},
	{AX, NewMemory(nil, NewImmediate(0x1234, Unsign, Bit16), Bit16, nil), 0x0011 + 0x3534},
	{AX, NewMemory([]*Register{BX}, NewImmediate(0x1234, Sign, Bit16), Bit16, nil), 0x0011 + 0x3736},
	{AX, NewMemory([]*Register{BX, SI}, NewImmediate(0x1234, Sign, Bit16), Bit16, nil), 0x0011 + 0x3b3a},
	{AL, NewMemory([]*Register{BX, SI}, NewImmediate(0x1234, Sign, Bit16), Bit8, nil), (0x0011 + 0x3b3a) & 0x00ff},
	{NewMemory(nil, NewImmediate(0x0000, Unsign, Bit16), Bit8, nil), NewImmediate(0xff, Sign, Bit8), 0xff},
}

func TestOpcodeADD(t *testing.T) {
	for _, test := range opcodeADDTests {
		vm := NewVM()
		AX.Write(vm, 0x0011)
		CX.Write(vm, 0x2233)
		BX.Write(vm, 0x0002)
		SI.Write(vm, 0x0004)
		for i, _ := range vm.memDS {
			vm.memDS[i] = byte(i)
		}
		op := Opcode{mn: ADD, opr1: test.opr1, opr2: test.opr2}
		op.Run(vm)
		assert.Equal(t, test.out, test.opr1.(ReadableOperand).Read(vm))
	}
}
