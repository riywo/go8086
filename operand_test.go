package go8086

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var registerDisasmTests = []struct {
	in  *Register
	out string
}{
	{AX, "ax"}, {CX, "cx"}, {DX, "dx"}, {BX, "bx"}, {SP, "sp"}, {BP, "bp"}, {SI, "si"}, {DI, "di"},
	{AL, "al"}, {CL, "cl"}, {DL, "dl"}, {BL, "bl"}, {AH, "ah"}, {CH, "ch"}, {DH, "dh"}, {BH, "bh"},
}

func TestRegisterDisasm(t *testing.T) {
	for _, test := range registerDisasmTests {
		assert.Equal(t, test.out, test.in.Disasm())
	}
}

func TestRegisterReadWrite(t *testing.T) {
	vm := NewVM()
	AX.Write(vm, 0x1234)
	assert.Equal(t, 0x1234, AX.Read(vm))
	assert.Equal(t, 0x12, AH.Read(vm))
	assert.Equal(t, 0x34, AL.Read(vm))
	AH.Write(vm, 0xff)
	assert.Equal(t, 0xff34, AX.Read(vm))
	assert.Equal(t, 0xff, AH.Read(vm))
	assert.Equal(t, 0x34, AL.Read(vm))
	AL.Write(vm, 0x00)
	assert.Equal(t, 0xff00, AX.Read(vm))
	assert.Equal(t, 0xff, AH.Read(vm))
	assert.Equal(t, 0x00, AL.Read(vm))
}

func TestSegmentRegisterReadWrite(t *testing.T) {
	vm := NewVM()
	CS.Write(vm, 0x1234)
	assert.Equal(t, 0x1234, CS.Read(vm))
}

var immediateDisasmTests = []struct {
	value  uint16
	signed Signed
	w      Bit
	out    string
}{
	{0xffff, Unsign, Bit16, "0xffff"},
	{0xffff, Sign, Bit16, "-0x1"},
	{0xffff, Unsign, Bit8, "0xff"},
	{0xffff, Sign, Bit8, "-0x1"},
	{0x1, Sign, Bit8, "+0x1"},
}

func TestImmediateDisasm(t *testing.T) {
	for _, test := range immediateDisasmTests {
		i := NewImmediate(test.value, test.signed, test.w)
		assert.Equal(t, test.out, i.Disasm())
	}
}

func TestImmediateRead(t *testing.T) {
	vm := NewVM()
	assert.Equal(t, 0x1234, NewImmediate(0x1234, Unsign, Bit16).Read(vm))
	assert.Equal(t, 0xffff, NewImmediate(0xffff, Sign, Bit16).Read(vm))
	assert.Equal(t, 0xffff, NewImmediate(0xff, Sign, Bit8).Read(vm))
}

var memoryDisasmTests = []struct {
	regad RegAddress
	disp  *Immediate
	w     Bit
	sreg  *SegmentRegister
	out   string
}{
	{RegAdd_Direct, NewImmediate(0xffff, Unsign, Bit16), Bit16, nil, "[0xffff]"},
	{RegAdd_BX, NewImmediate(0xffff, Sign, Bit16), Bit16, nil, "[bx-0x1]"},
	{RegAdd_BX_SI, nil, Bit16, nil, "[bx+si]"},
	{RegAdd_BX_SI, NewImmediate(0x1234, Sign, Bit16), Bit16, nil, "[bx+si+0x1234]"},
	{RegAdd_Direct, NewImmediate(0xffff, Unsign, Bit16), Bit16, ES, "[es:0xffff]"},
	{RegAdd_BX, NewImmediate(0xffff, Sign, Bit16), Bit16, CS, "[cs:bx-0x1]"},
	{RegAdd_BX_SI, nil, Bit16, SS, "[ss:bx+si]"},
	{RegAdd_BP_SI, NewImmediate(0x1234, Sign, Bit16), Bit16, DS, "[ds:bp+si+0x1234]"},
}

func TestMemoryDisasm(t *testing.T) {
	for _, test := range memoryDisasmTests {
		m := NewMemory(test.regad, test.disp, test.w, test.sreg)
		assert.Equal(t, test.out, m.Disasm())
	}
}

var memoryReadWriteTests = []struct {
	regad RegAddress
	disp  *Immediate
	w     Bit
	sreg  *SegmentRegister
	ea    uint16
	out   uint16
}{
	{RegAdd_Direct, NewImmediate(0x1234, Unsign, Bit16), Bit16, nil, 0x1234, 0x3534},
	{RegAdd_BX, NewImmediate(0x1234, Sign, Bit16), Bit16, nil, 0x1236, 0x3736},
	{RegAdd_BX_SI, nil, Bit16, nil, 0x0006, 0x0706},
	{RegAdd_BX_SI, NewImmediate(0x1234, Sign, Bit16), Bit16, nil, 0x123a, 0x3b3a},
	{RegAdd_BX_SI, NewImmediate(0xffff, Sign, Bit16), Bit16, nil, 0x0005, 0x0605},
	{RegAdd_BX_SI, NewImmediate(0xff, Sign, Bit8), Bit16, nil, 0x0005, 0x0605},
}

func TestMemoryRead(t *testing.T) {
	for _, test := range memoryReadWriteTests {
		vm := NewVM()
		BX.Write(vm, 0x0002)
		SI.Write(vm, 0x0004)
		for i, _ := range vm.DS(0)[0:0xffff] {
			vm.DS(0)[i] = byte(i)
		}
		m := NewMemory(test.regad, test.disp, test.w, test.sreg)
		assert.Equal(t, test.ea, m.EffectiveAddress(vm))
		assert.Equal(t, test.out, m.Read(vm))
	}
}

func TestMemoryWrite(t *testing.T) {
	for _, test := range memoryReadWriteTests {
		vm := NewVM()
		BX.Write(vm, 0x0002)
		SI.Write(vm, 0x0004)
		for i, _ := range vm.DS(0)[0:0xffff] {
			vm.DS(0)[i] = byte(i)
		}
		m := NewMemory(test.regad, test.disp, test.w, test.sreg)
		m.Write(vm, m.Read(vm)+1)
		assert.Equal(t, test.out+1, m.Read(vm))
	}
}
