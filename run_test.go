package go8086

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var opcodeADDwithOperandTests = []struct {
	opr1 Operand
	opr2 Operand
	out  uint16
}{
	{AX, CX, 0x2244},
	{AL, CL, 0x44},
	{AX, NewMemory(RegAdd_Direct, NewImmediate(0x1234, Unsign, Bit16), Bit16, nil), 0x0011 + 0x3534},
	{AX, NewMemory(RegAdd_BX, NewImmediate(0x1234, Sign, Bit16), Bit16, nil), 0x0011 + 0x3736},
	{AX, NewMemory(RegAdd_BX_SI, NewImmediate(0x1234, Sign, Bit16), Bit16, nil), 0x0011 + 0x3b3a},
	{AL, NewMemory(RegAdd_BX_SI, NewImmediate(0x1234, Sign, Bit16), Bit8, nil), (0x0011 + 0x3b3a) & 0x00ff},
	{NewMemory(RegAdd_Direct, NewImmediate(0x0000, Unsign, Bit16), Bit8, nil), NewImmediate(0xff, Sign, Bit8), 0xff},
}

func TestOpcodeADDwithOperand(t *testing.T) {
	for _, test := range opcodeADDwithOperandTests {
		vm := NewVM()
		AX.Write(vm, 0x0011)
		CX.Write(vm, 0x2233)
		BX.Write(vm, 0x0002)
		SI.Write(vm, 0x0004)
		for i, _ := range vm.DS(0)[0:0xffff] {
			vm.DS(0)[i] = byte(i)
		}
		op := Opcode{mn: ADD, opr1: test.opr1, opr2: test.opr2}
		op.Run(vm)
		assert.Equal(t, test.out, test.opr1.(ReadableOperand).Read(vm))
	}
}

var runArithmeticTests = []struct {
	mn  Mnemonic
	in1 uint16
	in2 uint16
	out uint16
	CF  uint16
	OF  uint16
	AF  uint16
	ZF  uint16
	SF  uint16
}{
	{ADD, 0x0000, 0x0000, 0x0000, 0, 0, 0, 1, 0},
	{ADD, 0x0000, 0x0001, 0x0001, 0, 0, 0, 0, 0},
	{ADD, 0x0000, 0x7fff, 0x7fff, 0, 0, 0, 0, 0},
	{ADD, 0x0000, 0x8000, 0x8000, 0, 0, 0, 0, 1},
	{ADD, 0x0000, 0x8001, 0x8001, 0, 0, 0, 0, 1},
	{ADD, 0x0000, 0xffff, 0xffff, 0, 0, 0, 0, 1},
	{ADD, 0x0001, 0x0000, 0x0001, 0, 0, 0, 0, 0},
	{ADD, 0x0001, 0x0001, 0x0002, 0, 0, 0, 0, 0},
	{ADD, 0x0001, 0x7fff, 0x8000, 0, 1, 0, 0, 1},
	{ADD, 0x0001, 0x8000, 0x8001, 0, 0, 0, 0, 1},
	{ADD, 0x0001, 0x8001, 0x8002, 0, 0, 0, 0, 1},
	{ADD, 0x0001, 0xffff, 0x0000, 1, 0, 0, 1, 0},
	{ADD, 0x7fff, 0x0000, 0x7fff, 0, 0, 0, 0, 0},
	{ADD, 0x7fff, 0x0001, 0x8000, 0, 1, 0, 0, 1},
	{ADD, 0x7fff, 0x7fff, 0xfffe, 0, 1, 0, 0, 1},
	{ADD, 0x7fff, 0x8000, 0xffff, 0, 0, 0, 0, 1},
	{ADD, 0x7fff, 0x8001, 0x0000, 1, 0, 0, 1, 0},
	{ADD, 0x7fff, 0xffff, 0x7ffe, 1, 0, 0, 0, 0},
	{ADD, 0x8000, 0x0000, 0x8000, 0, 0, 0, 0, 1},
	{ADD, 0x8000, 0x0001, 0x8001, 0, 0, 0, 0, 1},
	{ADD, 0x8000, 0x7fff, 0xffff, 0, 0, 0, 0, 1},
	{ADD, 0x8000, 0x8000, 0x0000, 1, 1, 0, 1, 0},
	{ADD, 0x8000, 0x8001, 0x0001, 1, 1, 0, 0, 0},
	{ADD, 0x8000, 0xffff, 0x7fff, 1, 1, 0, 0, 0},
	{ADD, 0x8001, 0x0000, 0x8001, 0, 0, 0, 0, 1},
	{ADD, 0x8001, 0x0001, 0x8002, 0, 0, 0, 0, 1},
	{ADD, 0x8001, 0x7fff, 0x0000, 1, 0, 0, 1, 0},
	{ADD, 0x8001, 0x8000, 0x0001, 1, 1, 0, 0, 0},
	{ADD, 0x8001, 0x8001, 0x0002, 1, 1, 0, 0, 0},
	{ADD, 0x8001, 0xffff, 0x8000, 1, 0, 0, 0, 1},
	{ADD, 0xffff, 0x0000, 0xffff, 0, 0, 0, 0, 1},
	{ADD, 0xffff, 0x0001, 0x0000, 1, 0, 0, 1, 0},
	{ADD, 0xffff, 0x7fff, 0x7ffe, 1, 0, 0, 0, 0},
	{ADD, 0xffff, 0x8000, 0x7fff, 1, 1, 0, 0, 0},
	{ADD, 0xffff, 0x8001, 0x8000, 1, 0, 0, 0, 1},
	{ADD, 0xffff, 0xffff, 0xfffe, 1, 0, 0, 0, 1},

	{SUB, 0x0000, 0x0000, 0x0000, 0, 0, 0, 1, 0},
	{SUB, 0x0000, 0x0001, 0xffff, 1, 0, 0, 0, 1},
	{SUB, 0x0000, 0x7fff, 0x8001, 1, 0, 0, 0, 1},
	{SUB, 0x0000, 0x8000, 0x8000, 1, 1, 0, 0, 1},
	{SUB, 0x0000, 0x8001, 0x7fff, 1, 0, 0, 0, 0},
	{SUB, 0x0000, 0xffff, 0x0001, 1, 0, 0, 0, 0},
	{SUB, 0x0001, 0x0000, 0x0001, 0, 0, 0, 0, 0},
	{SUB, 0x0001, 0x0001, 0x0000, 0, 0, 0, 1, 0},
	{SUB, 0x0001, 0x7fff, 0x8002, 1, 0, 0, 0, 1},
	{SUB, 0x0001, 0x8000, 0x8001, 1, 1, 0, 0, 1},
	{SUB, 0x0001, 0x8001, 0x8000, 1, 1, 0, 0, 1},
	{SUB, 0x0001, 0xffff, 0x0002, 1, 0, 0, 0, 0},
	{SUB, 0x7fff, 0x0000, 0x7fff, 0, 0, 0, 0, 0},
	{SUB, 0x7fff, 0x0001, 0x7ffe, 0, 0, 0, 0, 0},
	{SUB, 0x7fff, 0x7fff, 0x0000, 0, 0, 0, 1, 0},
	{SUB, 0x7fff, 0x8000, 0xffff, 1, 1, 0, 0, 1},
	{SUB, 0x7fff, 0x8001, 0xfffe, 1, 1, 0, 0, 1},
	{SUB, 0x7fff, 0xffff, 0x8000, 1, 1, 0, 0, 1},
	{SUB, 0x8000, 0x0000, 0x8000, 0, 0, 0, 0, 1},
	{SUB, 0x8000, 0x0001, 0x7fff, 0, 1, 0, 0, 0},
	{SUB, 0x8000, 0x7fff, 0x0001, 0, 1, 0, 0, 0},
	{SUB, 0x8000, 0x8000, 0x0000, 0, 0, 0, 1, 0},
	{SUB, 0x8000, 0x8001, 0xffff, 1, 0, 0, 0, 1},
	{SUB, 0x8000, 0xffff, 0x8001, 1, 0, 0, 0, 1},
	{SUB, 0x8001, 0x0000, 0x8001, 0, 0, 0, 0, 1},
	{SUB, 0x8001, 0x0001, 0x8000, 0, 0, 0, 0, 1},
	{SUB, 0x8001, 0x7fff, 0x0002, 0, 1, 0, 0, 0},
	{SUB, 0x8001, 0x8000, 0x0001, 0, 0, 0, 0, 0},
	{SUB, 0x8001, 0x8001, 0x0000, 0, 0, 0, 1, 0},
	{SUB, 0x8001, 0xffff, 0x8002, 1, 0, 0, 0, 1},
	{SUB, 0xffff, 0x0000, 0xffff, 0, 0, 0, 0, 1},
	{SUB, 0xffff, 0x0001, 0xfffe, 0, 0, 0, 0, 1},
	{SUB, 0xffff, 0x7fff, 0x8000, 0, 0, 0, 0, 1},
	{SUB, 0xffff, 0x8000, 0x7fff, 0, 0, 0, 0, 0},
	{SUB, 0xffff, 0x8001, 0x7ffe, 0, 0, 0, 0, 0},
	{SUB, 0xffff, 0xffff, 0x0000, 0, 0, 0, 1, 0},

	{CMP, 0x0000, 0x0000, 0x0000, 0, 0, 0, 1, 0},
	{CMP, 0x0000, 0x0001, 0x0000, 1, 0, 0, 0, 1},
	{CMP, 0x0000, 0x7fff, 0x0000, 1, 0, 0, 0, 1},
	{CMP, 0x0000, 0x8000, 0x0000, 1, 1, 0, 0, 1},
	{CMP, 0x0000, 0x8001, 0x0000, 1, 0, 0, 0, 0},
	{CMP, 0x0000, 0xffff, 0x0000, 1, 0, 0, 0, 0},
	{CMP, 0x0001, 0x0000, 0x0001, 0, 0, 0, 0, 0},
	{CMP, 0x0001, 0x0001, 0x0001, 0, 0, 0, 1, 0},
	{CMP, 0x0001, 0x7fff, 0x0001, 1, 0, 0, 0, 1},
	{CMP, 0x0001, 0x8000, 0x0001, 1, 1, 0, 0, 1},
	{CMP, 0x0001, 0x8001, 0x0001, 1, 1, 0, 0, 1},
	{CMP, 0x0001, 0xffff, 0x0001, 1, 0, 0, 0, 0},
	{CMP, 0x7fff, 0x0000, 0x7fff, 0, 0, 0, 0, 0},
	{CMP, 0x7fff, 0x0001, 0x7fff, 0, 0, 0, 0, 0},
	{CMP, 0x7fff, 0x7fff, 0x7fff, 0, 0, 0, 1, 0},
	{CMP, 0x7fff, 0x8000, 0x7fff, 1, 1, 0, 0, 1},
	{CMP, 0x7fff, 0x8001, 0x7fff, 1, 1, 0, 0, 1},
	{CMP, 0x7fff, 0xffff, 0x7fff, 1, 1, 0, 0, 1},
	{CMP, 0x8000, 0x0000, 0x8000, 0, 0, 0, 0, 1},
	{CMP, 0x8000, 0x0001, 0x8000, 0, 1, 0, 0, 0},
	{CMP, 0x8000, 0x7fff, 0x8000, 0, 1, 0, 0, 0},
	{CMP, 0x8000, 0x8000, 0x8000, 0, 0, 0, 1, 0},
	{CMP, 0x8000, 0x8001, 0x8000, 1, 0, 0, 0, 1},
	{CMP, 0x8000, 0xffff, 0x8000, 1, 0, 0, 0, 1},
	{CMP, 0x8001, 0x0000, 0x8001, 0, 0, 0, 0, 1},
	{CMP, 0x8001, 0x0001, 0x8001, 0, 0, 0, 0, 1},
	{CMP, 0x8001, 0x7fff, 0x8001, 0, 1, 0, 0, 0},
	{CMP, 0x8001, 0x8000, 0x8001, 0, 0, 0, 0, 0},
	{CMP, 0x8001, 0x8001, 0x8001, 0, 0, 0, 1, 0},
	{CMP, 0x8001, 0xffff, 0x8001, 1, 0, 0, 0, 1},
	{CMP, 0xffff, 0x0000, 0xffff, 0, 0, 0, 0, 1},
	{CMP, 0xffff, 0x0001, 0xffff, 0, 0, 0, 0, 1},
	{CMP, 0xffff, 0x7fff, 0xffff, 0, 0, 0, 0, 1},
	{CMP, 0xffff, 0x8000, 0xffff, 0, 0, 0, 0, 0},
	{CMP, 0xffff, 0x8001, 0xffff, 0, 0, 0, 0, 0},
	{CMP, 0xffff, 0xffff, 0xffff, 0, 0, 0, 1, 0},
}

func TestRunArithmetic(t *testing.T) {
	for _, test := range runArithmeticTests {
		vm := NewVM()
		AX.Write(vm, test.in1)
		DX.Write(vm, test.in2)
		op := Opcode{mn: test.mn, opr1: AX, opr2: DX}
		op.Run(vm)
		msg := fmt.Sprintf(" - %s %#04x,%#04x", test.mn, test.in1, test.in2)
		assert.Equal(t, test.out, AX.Read(vm), "out"+msg)
		assert.Equal(t, test.CF, vm.GetFlag(CF), "CF"+msg)
		assert.Equal(t, test.OF, vm.GetFlag(OF), "OF"+msg)
		assert.Equal(t, test.AF, vm.GetFlag(AF), "AF"+msg)
		assert.Equal(t, test.ZF, vm.GetFlag(ZF), "ZF"+msg)
		assert.Equal(t, test.SF, vm.GetFlag(SF), "SF"+msg)
	}
}

var runArithmeticWithCarryTests = []struct {
	mn   Mnemonic
	in1  uint16
	in2  uint16
	inCF uint16
	out  uint16
	CF   uint16
	OF   uint16
	AF   uint16
	ZF   uint16
	SF   uint16
}{
	{ADC, 0x0000, 0x0000, 0, 0x0000, 0, 0, 0, 1, 0},
	{ADC, 0x0000, 0x0001, 0, 0x0001, 0, 0, 0, 0, 0},
	{ADC, 0x0000, 0x7fff, 0, 0x7fff, 0, 0, 0, 0, 0},
	{ADC, 0x0000, 0x8000, 0, 0x8000, 0, 0, 0, 0, 1},
	{ADC, 0x0000, 0x8001, 0, 0x8001, 0, 0, 0, 0, 1},
	{ADC, 0x0000, 0xffff, 0, 0xffff, 0, 0, 0, 0, 1},
	{ADC, 0x0001, 0x0000, 0, 0x0001, 0, 0, 0, 0, 0},
	{ADC, 0x0001, 0x0001, 0, 0x0002, 0, 0, 0, 0, 0},
	{ADC, 0x0001, 0x7fff, 0, 0x8000, 0, 1, 0, 0, 1},
	{ADC, 0x0001, 0x8000, 0, 0x8001, 0, 0, 0, 0, 1},
	{ADC, 0x0001, 0x8001, 0, 0x8002, 0, 0, 0, 0, 1},
	{ADC, 0x0001, 0xffff, 0, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0x7fff, 0x0000, 0, 0x7fff, 0, 0, 0, 0, 0},
	{ADC, 0x7fff, 0x0001, 0, 0x8000, 0, 1, 0, 0, 1},
	{ADC, 0x7fff, 0x7fff, 0, 0xfffe, 0, 1, 0, 0, 1},
	{ADC, 0x7fff, 0x8000, 0, 0xffff, 0, 0, 0, 0, 1},
	{ADC, 0x7fff, 0x8001, 0, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0x7fff, 0xffff, 0, 0x7ffe, 1, 0, 0, 0, 0},
	{ADC, 0x8000, 0x0000, 0, 0x8000, 0, 0, 0, 0, 1},
	{ADC, 0x8000, 0x0001, 0, 0x8001, 0, 0, 0, 0, 1},
	{ADC, 0x8000, 0x7fff, 0, 0xffff, 0, 0, 0, 0, 1},
	{ADC, 0x8000, 0x8000, 0, 0x0000, 1, 1, 0, 1, 0},
	{ADC, 0x8000, 0x8001, 0, 0x0001, 1, 1, 0, 0, 0},
	{ADC, 0x8000, 0xffff, 0, 0x7fff, 1, 1, 0, 0, 0},
	{ADC, 0x8001, 0x0000, 0, 0x8001, 0, 0, 0, 0, 1},
	{ADC, 0x8001, 0x0001, 0, 0x8002, 0, 0, 0, 0, 1},
	{ADC, 0x8001, 0x7fff, 0, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0x8001, 0x8000, 0, 0x0001, 1, 1, 0, 0, 0},
	{ADC, 0x8001, 0x8001, 0, 0x0002, 1, 1, 0, 0, 0},
	{ADC, 0x8001, 0xffff, 0, 0x8000, 1, 0, 0, 0, 1},
	{ADC, 0xffff, 0x0000, 0, 0xffff, 0, 0, 0, 0, 1},
	{ADC, 0xffff, 0x0001, 0, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0xffff, 0x7fff, 0, 0x7ffe, 1, 0, 0, 0, 0},
	{ADC, 0xffff, 0x8000, 0, 0x7fff, 1, 1, 0, 0, 0},
	{ADC, 0xffff, 0x8001, 0, 0x8000, 1, 0, 0, 0, 1},
	{ADC, 0xffff, 0xffff, 0, 0xfffe, 1, 0, 0, 0, 1},

	{ADC, 0x0000, 0x0000, 1, 0x0001, 0, 0, 0, 0, 0},
	{ADC, 0x0000, 0x0001, 1, 0x0002, 0, 0, 0, 0, 0},
	{ADC, 0x0000, 0x7fff, 1, 0x8000, 0, 1, 0, 0, 1},
	{ADC, 0x0000, 0x8000, 1, 0x8001, 0, 0, 0, 0, 1},
	{ADC, 0x0000, 0x8001, 1, 0x8002, 0, 0, 0, 0, 1},
	{ADC, 0x0000, 0xffff, 1, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0x0001, 0x0000, 1, 0x0002, 0, 0, 0, 0, 0},
	{ADC, 0x0001, 0x0001, 1, 0x0003, 0, 0, 0, 0, 0},
	{ADC, 0x0001, 0x7fff, 1, 0x8001, 0, 1, 0, 0, 1},
	{ADC, 0x0001, 0x8000, 1, 0x8002, 0, 0, 0, 0, 1},
	{ADC, 0x0001, 0x8001, 1, 0x8003, 0, 0, 0, 0, 1},
	{ADC, 0x0001, 0xffff, 1, 0x0001, 1, 0, 0, 0, 0},
	{ADC, 0x7fff, 0x0000, 1, 0x8000, 0, 1, 0, 0, 1},
	{ADC, 0x7fff, 0x0001, 1, 0x8001, 0, 1, 0, 0, 1},
	{ADC, 0x7fff, 0x7fff, 1, 0xffff, 0, 1, 0, 0, 1},
	{ADC, 0x7fff, 0x8000, 1, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0x7fff, 0x8001, 1, 0x0001, 1, 0, 0, 0, 0},
	{ADC, 0x7fff, 0xffff, 1, 0x7fff, 1, 0, 0, 0, 0},
	{ADC, 0x8000, 0x0000, 1, 0x8001, 0, 0, 0, 0, 1},
	{ADC, 0x8000, 0x0001, 1, 0x8002, 0, 0, 0, 0, 1},
	{ADC, 0x8000, 0x7fff, 1, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0x8000, 0x8000, 1, 0x0001, 1, 1, 0, 0, 0},
	{ADC, 0x8000, 0x8001, 1, 0x0002, 1, 1, 0, 0, 0},
	{ADC, 0x8000, 0xffff, 1, 0x8000, 1, 0, 0, 0, 1},
	{ADC, 0x8001, 0x0000, 1, 0x8002, 0, 0, 0, 0, 1},
	{ADC, 0x8001, 0x0001, 1, 0x8003, 0, 0, 0, 0, 1},
	{ADC, 0x8001, 0x7fff, 1, 0x0001, 1, 0, 0, 0, 0},
	{ADC, 0x8001, 0x8000, 1, 0x0002, 1, 1, 0, 0, 0},
	{ADC, 0x8001, 0x8001, 1, 0x0003, 1, 1, 0, 0, 0},
	{ADC, 0x8001, 0xffff, 1, 0x8001, 1, 0, 0, 0, 1},
	{ADC, 0xffff, 0x0000, 1, 0x0000, 1, 0, 0, 1, 0},
	{ADC, 0xffff, 0x0001, 1, 0x0001, 1, 0, 0, 0, 0},
	{ADC, 0xffff, 0x7fff, 1, 0x7fff, 1, 0, 0, 0, 0},
	{ADC, 0xffff, 0x8000, 1, 0x8000, 1, 0, 0, 0, 1},
	{ADC, 0xffff, 0x8001, 1, 0x8001, 1, 0, 0, 0, 1},
	{ADC, 0xffff, 0xffff, 1, 0xffff, 1, 0, 0, 0, 1},

	{SBB, 0x0000, 0x0000, 0, 0x0000, 0, 0, 0, 1, 0},
	{SBB, 0x0000, 0x0001, 0, 0xffff, 1, 0, 0, 0, 1},
	{SBB, 0x0000, 0x7fff, 0, 0x8001, 1, 0, 0, 0, 1},
	{SBB, 0x0000, 0x8000, 0, 0x8000, 1, 1, 0, 0, 1},
	{SBB, 0x0000, 0x8001, 0, 0x7fff, 1, 0, 0, 0, 0},
	{SBB, 0x0000, 0xffff, 0, 0x0001, 1, 0, 0, 0, 0},
	{SBB, 0x0001, 0x0000, 0, 0x0001, 0, 0, 0, 0, 0},
	{SBB, 0x0001, 0x0001, 0, 0x0000, 0, 0, 0, 1, 0},
	{SBB, 0x0001, 0x7fff, 0, 0x8002, 1, 0, 0, 0, 1},
	{SBB, 0x0001, 0x8000, 0, 0x8001, 1, 1, 0, 0, 1},
	{SBB, 0x0001, 0x8001, 0, 0x8000, 1, 1, 0, 0, 1},
	{SBB, 0x0001, 0xffff, 0, 0x0002, 1, 0, 0, 0, 0},
	{SBB, 0x7fff, 0x0000, 0, 0x7fff, 0, 0, 0, 0, 0},
	{SBB, 0x7fff, 0x0001, 0, 0x7ffe, 0, 0, 0, 0, 0},
	{SBB, 0x7fff, 0x7fff, 0, 0x0000, 0, 0, 0, 1, 0},
	{SBB, 0x7fff, 0x8000, 0, 0xffff, 1, 1, 0, 0, 1},
	{SBB, 0x7fff, 0x8001, 0, 0xfffe, 1, 1, 0, 0, 1},
	{SBB, 0x7fff, 0xffff, 0, 0x8000, 1, 1, 0, 0, 1},
	{SBB, 0x8000, 0x0000, 0, 0x8000, 0, 0, 0, 0, 1},
	{SBB, 0x8000, 0x0001, 0, 0x7fff, 0, 1, 0, 0, 0},
	{SBB, 0x8000, 0x7fff, 0, 0x0001, 0, 1, 0, 0, 0},
	{SBB, 0x8000, 0x8000, 0, 0x0000, 0, 0, 0, 1, 0},
	{SBB, 0x8000, 0x8001, 0, 0xffff, 1, 0, 0, 0, 1},
	{SBB, 0x8000, 0xffff, 0, 0x8001, 1, 0, 0, 0, 1},
	{SBB, 0x8001, 0x0000, 0, 0x8001, 0, 0, 0, 0, 1},
	{SBB, 0x8001, 0x0001, 0, 0x8000, 0, 0, 0, 0, 1},
	{SBB, 0x8001, 0x7fff, 0, 0x0002, 0, 1, 0, 0, 0},
	{SBB, 0x8001, 0x8000, 0, 0x0001, 0, 0, 0, 0, 0},
	{SBB, 0x8001, 0x8001, 0, 0x0000, 0, 0, 0, 1, 0},
	{SBB, 0x8001, 0xffff, 0, 0x8002, 1, 0, 0, 0, 1},
	{SBB, 0xffff, 0x0000, 0, 0xffff, 0, 0, 0, 0, 1},
	{SBB, 0xffff, 0x0001, 0, 0xfffe, 0, 0, 0, 0, 1},
	{SBB, 0xffff, 0x7fff, 0, 0x8000, 0, 0, 0, 0, 1},
	{SBB, 0xffff, 0x8000, 0, 0x7fff, 0, 0, 0, 0, 0},
	{SBB, 0xffff, 0x8001, 0, 0x7ffe, 0, 0, 0, 0, 0},
	{SBB, 0xffff, 0xffff, 0, 0x0000, 0, 0, 0, 1, 0},

	{SBB, 0x0000, 0x0000, 1, 0xffff, 1, 0, 0, 0, 1},
	{SBB, 0x0000, 0x0001, 1, 0xfffe, 1, 0, 0, 0, 1},
	{SBB, 0x0000, 0x7fff, 1, 0x8000, 1, 0, 0, 0, 1},
	{SBB, 0x0000, 0x8000, 1, 0x7fff, 1, 0, 0, 0, 0},
	{SBB, 0x0000, 0x8001, 1, 0x7ffe, 1, 0, 0, 0, 0},
	{SBB, 0x0000, 0xffff, 1, 0x0000, 1, 0, 0, 1, 0},
	{SBB, 0x0001, 0x0000, 1, 0x0000, 0, 0, 0, 1, 0},
	{SBB, 0x0001, 0x0001, 1, 0xffff, 1, 0, 0, 0, 1},
	{SBB, 0x0001, 0x7fff, 1, 0x8001, 1, 0, 0, 0, 1},
	{SBB, 0x0001, 0x8000, 1, 0x8000, 1, 1, 0, 0, 1},
	{SBB, 0x0001, 0x8001, 1, 0x7fff, 1, 0, 0, 0, 0},
	{SBB, 0x0001, 0xffff, 1, 0x0001, 1, 0, 0, 0, 0},
	{SBB, 0x7fff, 0x0000, 1, 0x7ffe, 0, 0, 0, 0, 0},
	{SBB, 0x7fff, 0x0001, 1, 0x7ffd, 0, 0, 0, 0, 0},
	{SBB, 0x7fff, 0x7fff, 1, 0xffff, 1, 0, 0, 0, 1},
	{SBB, 0x7fff, 0x8000, 1, 0xfffe, 1, 1, 0, 0, 1},
	{SBB, 0x7fff, 0x8001, 1, 0xfffd, 1, 1, 0, 0, 1},
	{SBB, 0x7fff, 0xffff, 1, 0x7fff, 1, 0, 0, 0, 0},
	{SBB, 0x8000, 0x0000, 1, 0x7fff, 0, 1, 0, 0, 0},
	{SBB, 0x8000, 0x0001, 1, 0x7ffe, 0, 1, 0, 0, 0},
	{SBB, 0x8000, 0x7fff, 1, 0x0000, 0, 1, 0, 1, 0},
	{SBB, 0x8000, 0x8000, 1, 0xffff, 1, 0, 0, 0, 1},
	{SBB, 0x8000, 0x8001, 1, 0xfffe, 1, 0, 0, 0, 1},
	{SBB, 0x8000, 0xffff, 1, 0x8000, 1, 0, 0, 0, 1},
	{SBB, 0x8001, 0x0000, 1, 0x8000, 0, 0, 0, 0, 1},
	{SBB, 0x8001, 0x0001, 1, 0x7fff, 0, 1, 0, 0, 0},
	{SBB, 0x8001, 0x7fff, 1, 0x0001, 0, 1, 0, 0, 0},
	{SBB, 0x8001, 0x8000, 1, 0x0000, 0, 0, 0, 1, 0},
	{SBB, 0x8001, 0x8001, 1, 0xffff, 1, 0, 0, 0, 1},
	{SBB, 0x8001, 0xffff, 1, 0x8001, 1, 0, 0, 0, 1},
	{SBB, 0xffff, 0x0000, 1, 0xfffe, 0, 0, 0, 0, 1},
	{SBB, 0xffff, 0x0001, 1, 0xfffd, 0, 0, 0, 0, 1},
	{SBB, 0xffff, 0x7fff, 1, 0x7fff, 0, 1, 0, 0, 0},
	{SBB, 0xffff, 0x8000, 1, 0x7ffe, 0, 0, 0, 0, 0},
	{SBB, 0xffff, 0x8001, 1, 0x7ffd, 0, 0, 0, 0, 0},
	{SBB, 0xffff, 0xffff, 1, 0xffff, 1, 0, 0, 0, 1},
}

func TestRunArithmeticWithCarry(t *testing.T) {
	for _, test := range runArithmeticWithCarryTests {
		vm := NewVM()
		AX.Write(vm, test.in1)
		DX.Write(vm, test.in2)
		vm.SetFlag(CF, test.inCF == 1)
		op := Opcode{mn: test.mn, opr1: AX, opr2: DX}
		op.Run(vm)
		msg := fmt.Sprintf(" - %s %#04x,%#04x CF:%d", test.mn, test.in1, test.in2, test.inCF)
		assert.Equal(t, test.out, AX.Read(vm), "out"+msg)
		assert.Equal(t, test.CF, vm.GetFlag(CF), "CF"+msg)
		assert.Equal(t, test.OF, vm.GetFlag(OF), "OF"+msg)
		assert.Equal(t, test.AF, vm.GetFlag(AF), "AF"+msg)
		assert.Equal(t, test.ZF, vm.GetFlag(ZF), "ZF"+msg)
		assert.Equal(t, test.SF, vm.GetFlag(SF), "SF"+msg)
	}
}

var runArithmetic8BitTests = []struct {
	mn  Mnemonic
	in1 uint16
	in2 uint16
	out uint16
	CF  uint16
	OF  uint16
	AF  uint16
	ZF  uint16
	SF  uint16
}{
	{ADD, 0x00, 0x00, 0x00, 0, 0, 0, 1, 0},
	{ADD, 0x00, 0x01, 0x01, 0, 0, 0, 0, 0},
	{ADD, 0x00, 0x7f, 0x7f, 0, 0, 0, 0, 0},
	{ADD, 0x00, 0x80, 0x80, 0, 0, 0, 0, 1},
	{ADD, 0x00, 0x81, 0x81, 0, 0, 0, 0, 1},
	{ADD, 0x00, 0xff, 0xff, 0, 0, 0, 0, 1},
	{ADD, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0},
	{ADD, 0x01, 0x01, 0x02, 0, 0, 0, 0, 0},
	{ADD, 0x01, 0x7f, 0x80, 0, 1, 0, 0, 1},
	{ADD, 0x01, 0x80, 0x81, 0, 0, 0, 0, 1},
	{ADD, 0x01, 0x81, 0x82, 0, 0, 0, 0, 1},
	{ADD, 0x01, 0xff, 0x00, 1, 0, 0, 1, 0},
	{ADD, 0x7f, 0x00, 0x7f, 0, 0, 0, 0, 0},
	{ADD, 0x7f, 0x01, 0x80, 0, 1, 0, 0, 1},
	{ADD, 0x7f, 0x7f, 0xfe, 0, 1, 0, 0, 1},
	{ADD, 0x7f, 0x80, 0xff, 0, 0, 0, 0, 1},
	{ADD, 0x7f, 0x81, 0x00, 1, 0, 0, 1, 0},
	{ADD, 0x7f, 0xff, 0x7e, 1, 0, 0, 0, 0},
	{ADD, 0x80, 0x00, 0x80, 0, 0, 0, 0, 1},
	{ADD, 0x80, 0x01, 0x81, 0, 0, 0, 0, 1},
	{ADD, 0x80, 0x7f, 0xff, 0, 0, 0, 0, 1},
	{ADD, 0x80, 0x80, 0x00, 1, 1, 0, 1, 0},
	{ADD, 0x80, 0x81, 0x01, 1, 1, 0, 0, 0},
	{ADD, 0x80, 0xff, 0x7f, 1, 1, 0, 0, 0},
	{ADD, 0x81, 0x00, 0x81, 0, 0, 0, 0, 1},
	{ADD, 0x81, 0x01, 0x82, 0, 0, 0, 0, 1},
	{ADD, 0x81, 0x7f, 0x00, 1, 0, 0, 1, 0},
	{ADD, 0x81, 0x80, 0x01, 1, 1, 0, 0, 0},
	{ADD, 0x81, 0x81, 0x02, 1, 1, 0, 0, 0},
	{ADD, 0x81, 0xff, 0x80, 1, 0, 0, 0, 1},
	{ADD, 0xff, 0x00, 0xff, 0, 0, 0, 0, 1},
	{ADD, 0xff, 0x01, 0x00, 1, 0, 0, 1, 0},
	{ADD, 0xff, 0x7f, 0x7e, 1, 0, 0, 0, 0},
	{ADD, 0xff, 0x80, 0x7f, 1, 1, 0, 0, 0},
	{ADD, 0xff, 0x81, 0x80, 1, 0, 0, 0, 1},
	{ADD, 0xff, 0xff, 0xfe, 1, 0, 0, 0, 1},

	{SUB, 0x00, 0x00, 0x00, 0, 0, 0, 1, 0},
	{SUB, 0x00, 0x01, 0xff, 1, 0, 0, 0, 1},
	{SUB, 0x00, 0x7f, 0x81, 1, 0, 0, 0, 1},
	{SUB, 0x00, 0x80, 0x80, 1, 1, 0, 0, 1},
	{SUB, 0x00, 0x81, 0x7f, 1, 0, 0, 0, 0},
	{SUB, 0x00, 0xff, 0x01, 1, 0, 0, 0, 0},
	{SUB, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0},
	{SUB, 0x01, 0x01, 0x00, 0, 0, 0, 1, 0},
	{SUB, 0x01, 0x7f, 0x82, 1, 0, 0, 0, 1},
	{SUB, 0x01, 0x80, 0x81, 1, 1, 0, 0, 1},
	{SUB, 0x01, 0x81, 0x80, 1, 1, 0, 0, 1},
	{SUB, 0x01, 0xff, 0x02, 1, 0, 0, 0, 0},
	{SUB, 0x7f, 0x00, 0x7f, 0, 0, 0, 0, 0},
	{SUB, 0x7f, 0x01, 0x7e, 0, 0, 0, 0, 0},
	{SUB, 0x7f, 0x7f, 0x00, 0, 0, 0, 1, 0},
	{SUB, 0x7f, 0x80, 0xff, 1, 1, 0, 0, 1},
	{SUB, 0x7f, 0x81, 0xfe, 1, 1, 0, 0, 1},
	{SUB, 0x7f, 0xff, 0x80, 1, 1, 0, 0, 1},
	{SUB, 0x80, 0x00, 0x80, 0, 0, 0, 0, 1},
	{SUB, 0x80, 0x01, 0x7f, 0, 1, 0, 0, 0},
	{SUB, 0x80, 0x7f, 0x01, 0, 1, 0, 0, 0},
	{SUB, 0x80, 0x80, 0x00, 0, 0, 0, 1, 0},
	{SUB, 0x80, 0x81, 0xff, 1, 0, 0, 0, 1},
	{SUB, 0x80, 0xff, 0x81, 1, 0, 0, 0, 1},
	{SUB, 0x81, 0x00, 0x81, 0, 0, 0, 0, 1},
	{SUB, 0x81, 0x01, 0x80, 0, 0, 0, 0, 1},
	{SUB, 0x81, 0x7f, 0x02, 0, 1, 0, 0, 0},
	{SUB, 0x81, 0x80, 0x01, 0, 0, 0, 0, 0},
	{SUB, 0x81, 0x81, 0x00, 0, 0, 0, 1, 0},
	{SUB, 0x81, 0xff, 0x82, 1, 0, 0, 0, 1},
	{SUB, 0xff, 0x00, 0xff, 0, 0, 0, 0, 1},
	{SUB, 0xff, 0x01, 0xfe, 0, 0, 0, 0, 1},
	{SUB, 0xff, 0x7f, 0x80, 0, 0, 0, 0, 1},
	{SUB, 0xff, 0x80, 0x7f, 0, 0, 0, 0, 0},
	{SUB, 0xff, 0x81, 0x7e, 0, 0, 0, 0, 0},
	{SUB, 0xff, 0xff, 0x00, 0, 0, 0, 1, 0},
}

func TestRunArithmetic8Bit(t *testing.T) {
	for _, test := range runArithmetic8BitTests {
		vm := NewVM()
		AL.Write(vm, test.in1)
		DH.Write(vm, test.in2)
		op := Opcode{mn: test.mn, opr1: AL, opr2: DH}
		op.Run(vm)
		msg := fmt.Sprintf(" - %s %#04x,%#04x", test.mn, test.in1, test.in2)
		assert.Equal(t, test.out, AL.Read(vm), "out"+msg)
		assert.Equal(t, test.CF, vm.GetFlag(CF), "CF"+msg)
		assert.Equal(t, test.OF, vm.GetFlag(OF), "OF"+msg)
		assert.Equal(t, test.AF, vm.GetFlag(AF), "AF"+msg)
		assert.Equal(t, test.ZF, vm.GetFlag(ZF), "ZF"+msg)
		assert.Equal(t, test.SF, vm.GetFlag(SF), "SF"+msg)
	}
}

var runIncDecNegTests = []struct {
	mn  Mnemonic
	in  uint16
	out uint16
	CF  uint16
	OF  uint16
	AF  uint16
	ZF  uint16
	SF  uint16
}{
	{INC, 0x0000, 0x0001, 0, 0, 0, 0, 0},
	{INC, 0x7fff, 0x8000, 0, 1, 0, 0, 1},
	{INC, 0xffff, 0x0000, 0, 0, 0, 1, 0},
	{DEC, 0x0001, 0x0000, 0, 0, 0, 1, 0},
	{DEC, 0x8000, 0x7fff, 0, 1, 0, 0, 0},
	{DEC, 0x0000, 0xffff, 0, 0, 0, 0, 1},
	{NEG, 0x0000, 0x0000, 0, 0, 0, 1, 0},
	{NEG, 0x0001, 0xffff, 1, 0, 0, 0, 1},
	{NEG, 0x7fff, 0x8001, 1, 0, 0, 0, 1},
	{NEG, 0x8000, 0x8000, 1, 1, 0, 0, 1},
}

func TestRunIncDecNeg(t *testing.T) {
	for _, test := range runIncDecNegTests {
		vm := NewVM()
		AX.Write(vm, test.in)
		op := Opcode{mn: test.mn, opr1: AX}
		op.Run(vm)
		msg := fmt.Sprintf(" - %s %#04x", test.mn, test.in)
		assert.Equal(t, test.out, AX.Read(vm), "out"+msg)
		assert.Equal(t, test.CF, vm.GetFlag(CF), "CF"+msg)
		assert.Equal(t, test.OF, vm.GetFlag(OF), "OF"+msg)
		assert.Equal(t, test.AF, vm.GetFlag(AF), "AF"+msg)
		assert.Equal(t, test.ZF, vm.GetFlag(ZF), "ZF"+msg)
		assert.Equal(t, test.SF, vm.GetFlag(SF), "SF"+msg)
	}
}

var runDivIdiv16Tests = []struct {
	mn        Mnemonic
	dx        uint16
	ax        uint16
	divisor   uint16
	quotient  uint16
	remainder uint16
}{
	{DIV, 0x0000, 0x0002, 0x0002, 0x0001, 0x0000},
	{DIV, 0x0000, 0x0003, 0x0002, 0x0001, 0x0001},
	{DIV, 0x0001, 0x1172, 0x000a, 0x1b58, 0x0002},
	{IDIV, 0x0000, 0x0002, 0x0002, 0x0001, 0x0000},
	{IDIV, 0x0000, 0x0003, 0x0002, 0x0001, 0x0001},
	{IDIV, 0xffff, 0xfff7, 0x0004, 0xfffe, 0xffff},
}

func TestRunDivIdiv16(t *testing.T) {
	for _, test := range runDivIdiv16Tests {
		vm := NewVM()
		DX.Write(vm, test.dx)
		AX.Write(vm, test.ax)
		BX.Write(vm, test.divisor)
		op := Opcode{mn: test.mn, opr1: BX}
		op.Run(vm)
		msg := fmt.Sprintf(" - %s %#04x DX:AX=%04x:%04x", test.mn, test.divisor, test.dx, test.ax)
		assert.Equal(t, test.quotient, AX.Read(vm), "quotient"+msg)
		assert.Equal(t, test.remainder, DX.Read(vm), "remainder"+msg)
	}
}

var runCbwTests = []struct {
	al uint16
	ax uint16
}{
	{0x01, 0x0001},
	{0x8c, 0xff8c},
}

func TestRunCbw(t *testing.T) {
	for _, test := range runCbwTests {
		vm := NewVM()
		AL.Write(vm, test.al)
		op := Opcode{mn: CBW}
		op.Run(vm)
		msg := fmt.Sprintf(" - %s AL: %02x", CBW, test.al)
		assert.Equal(t, test.ax, AX.Read(vm), "AX"+msg)
	}
}
