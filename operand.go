package go8086

import (
	"fmt"
)

type Operand interface {
	Disasm() string
	Bit() Bit
}

type ReadableOperand interface {
	Read(*VM) uint16
	Bit() Bit
}

type WritableOperand interface {
	Write(*VM, uint16)
	Bit() Bit
}

type ReadWritableOperand interface {
	Read(*VM) uint16
	Write(*VM, uint16)
	Bit() Bit
}

type Register struct {
	name    string
	w       Bit
	reg16   *Register
	bytePos BytePosition
}

func NewRegister16(name string) *Register {
	return &Register{name: name, w: Bit16}
}

func NewRegister8(name string, reg16 *Register, bytePos BytePosition) *Register {
	return &Register{name: name, w: Bit8, reg16: reg16, bytePos: bytePos}
}

func (r *Register) Bit() Bit {
	return r.w
}

func (r *Register) Disasm() string {
	return r.name
}

func (r *Register) Read(vm *VM) (value uint16) {
	switch r.w {
	case Bit8:
		switch r.bytePos {
		case LSB:
			value = r.reg16.Read(vm) & 0x00ff
		case MSB:
			value = r.reg16.Read(vm) >> 8
		}
	case Bit16:
		value = vm.reg[r.name]
	}
	return
}

func (r *Register) Write(vm *VM, value uint16) {
	switch r.w {
	case Bit8:
		switch r.bytePos {
		case LSB:
			v16msb := r.reg16.Read(vm) & 0xff00
			r.reg16.Write(vm, v16msb|value)
		case MSB:
			v16lsb := r.reg16.Read(vm) & 0x00ff
			r.reg16.Write(vm, (value<<8)|v16lsb)
		}
	case Bit16:
		vm.reg[r.name] = value
	}
	return
}

type BytePosition int

const (
	LSB BytePosition = iota
	MSB
)

var AX *Register = NewRegister16("ax")
var CX *Register = NewRegister16("cx")
var DX *Register = NewRegister16("dx")
var BX *Register = NewRegister16("bx")
var SP *Register = NewRegister16("sp")
var BP *Register = NewRegister16("bp")
var SI *Register = NewRegister16("si")
var DI *Register = NewRegister16("di")
var AL *Register = NewRegister8("al", AX, LSB)
var CL *Register = NewRegister8("cl", CX, LSB)
var DL *Register = NewRegister8("dl", DX, LSB)
var BL *Register = NewRegister8("bl", BX, LSB)
var AH *Register = NewRegister8("ah", AX, MSB)
var CH *Register = NewRegister8("ch", CX, MSB)
var DH *Register = NewRegister8("dh", DX, MSB)
var BH *Register = NewRegister8("bh", BX, MSB)

var regs = [2][8]*Register{
	[8]*Register{AL, CL, DL, BL, AH, CH, DH, BH},
	[8]*Register{AX, CX, DX, BX, SP, BP, SI, DI},
}

func getRegister(w Bit, r Reg) *Register {
	return regs[w][r]
}

type SegmentRegister struct {
	name string
}

func NewSegmentRegister(name string) *SegmentRegister {
	r := &SegmentRegister{name: name}
	return r
}

func (r *SegmentRegister) Bit() Bit {
	return Bit16
}

func (r *SegmentRegister) Disasm() string {
	return r.name
}

func (r *SegmentRegister) Read(vm *VM) (value uint16) {
	return vm.sreg[r.name]
}

func (r *SegmentRegister) Write(vm *VM, value uint16) {
	vm.sreg[r.name] = value
	return
}

var ES *SegmentRegister = NewSegmentRegister("es")
var CS *SegmentRegister = NewSegmentRegister("cs")
var SS *SegmentRegister = NewSegmentRegister("ss")
var DS *SegmentRegister = NewSegmentRegister("ds")

var sregs = [4]*SegmentRegister{ES, CS, SS, DS}

func getSegmentRegister(r SReg) *SegmentRegister {
	return sregs[r]
}

type Immediate struct {
	value  uint16
	signed Signed
	w      Bit
}

func NewImmediate(value uint16, signed Signed, w Bit) (imm *Immediate) {
	imm = &Immediate{value: value, signed: signed, w: w}
	if w == Bit8 {
		imm.value &= 0xff
		if signed && imm.value >= 0x80 {
			imm.value |= 0xff00
		}
	}
	return
}

func (i *Immediate) Bit() Bit {
	return i.w
}

func (i *Immediate) Disasm() (asm string) {
	if i.w == Bit8 && i.signed {
		asm = fmt.Sprintf("%#+x", int8(i.value))
	} else if i.w == Bit8 && !i.signed {
		asm = fmt.Sprintf("%#x", uint8(i.value))
	} else if i.w == Bit16 && i.signed {
		asm = fmt.Sprintf("%#+x", int16(i.value))
	} else if i.w == Bit16 && !i.signed {
		asm = fmt.Sprintf("%#x", i.value)
	}
	return
}

func (i *Immediate) Read(vm *VM) uint16 {
	return i.value
}

type Counter struct {
	v Count
	w Bit
}

func NewCounter(v Count, w Bit) *Counter {
	return &Counter{v: v, w: w}
}

func (c *Counter) Bit() Bit {
	return c.w
}

func (c *Counter) Disasm() (asm string) {
	switch c.v {
	case Count1:
		asm = "1"
	case CountCL:
		asm = CL.Disasm()
	}
	return
}

func (c *Counter) Count(vm *VM) (n uint16) {
	switch c.v {
	case Count1:
		n = 1
	case CountCL:
		n = CL.Read(vm)
	}
	return
}

type Memory struct {
	regad RegAddress
	disp  *Immediate
	w     Bit
	sreg  *SegmentRegister
}

func NewMemory(regad RegAddress, disp *Immediate, w Bit, sreg *SegmentRegister) *Memory {
	m := &Memory{regad: regad, disp: disp, w: w, sreg: sreg}
	if sreg == nil {
		m.sreg = RegAddressSegment[regad]
	}
	return m
}

func (m *Memory) Bit() Bit {
	return m.w
}

func (m *Memory) Disasm() string {
	ea := ""
	regad := RegAddressMap[m.regad]
	if m.regad == RegAdd_Direct {
		ea = m.disp.Disasm()
	} else {
		ea = regad[0].Disasm()
		for _, reg := range regad[1:] {
			ea += "+" + reg.Disasm()
		}
		if m.disp != nil {
			ea += m.disp.Disasm()
		}
	}
	if m.sreg != RegAddressSegment[m.regad] {
		ea = m.sreg.Disasm() + ":" + ea
	}
	return "[" + ea + "]"
}

func (m *Memory) Read(vm *VM) (value uint16) {
	switch m.w {
	case Bit8:
		value = m.Mem(vm).read8()
	case Bit16:
		value = m.Mem(vm).read16()
	}
	return
}

func (m *Memory) Write(vm *VM, value uint16) {
	switch m.w {
	case Bit8:
		m.Mem(vm).write8(value)
	case Bit16:
		m.Mem(vm).write16(value)
	}
	return
}

func (m *Memory) Mem(vm *VM) (mem Bytes) {
	ea := m.EffectiveAddress(vm)
	sreg := m.sreg
	return vm.Mem(sreg, ea)
}

func (m *Memory) EffectiveAddress(vm *VM) (ea uint16) {
	regad := RegAddressMap[m.regad]
	if len(regad) == 0 {
		ea = m.disp.Read(vm)
	} else {
		ea = regad[0].Read(vm)
		for _, reg := range regad[1:] {
			ea += reg.Read(vm)
		}
		if m.disp != nil {
			ea += m.disp.Read(vm)
		}
	}
	return
}

type DirectFarAddress struct {
	segment *Immediate
	offset  *Immediate
}

func NewDirectFarAddress(seg, off uint16) (dfa *DirectFarAddress) {
	dfa = new(DirectFarAddress)
	dfa.segment = NewImmediate(seg, Unsign, Bit16)
	dfa.offset = NewImmediate(off, Unsign, Bit16)
	return
}

func (dfa *DirectFarAddress) Disasm() string {
	return dfa.segment.Disasm() + ":" + dfa.offset.Disasm()
}

func (dfa *DirectFarAddress) Bit() Bit {
	return Bit16
}

type IndirectFarAddress struct {
	memory *Memory
}

func NewIndirectFarAddress(memory *Memory) (idfa *IndirectFarAddress) {
	idfa = new(IndirectFarAddress)
	idfa.memory = memory
	return
}

func (idfa *IndirectFarAddress) Disasm() string {
	return "far " + idfa.memory.Disasm()
}

func (idfa *IndirectFarAddress) Bit() Bit {
	return Bit16
}

func isRegister(opr Operand) (ok bool) {
	_, ok = opr.(*Register)
	return
}

func isSegmentRegister(opr Operand) (ok bool) {
	_, ok = opr.(*SegmentRegister)
	return
}

func isCounter(opr Operand) (ok bool) {
	_, ok = opr.(*Counter)
	return
}

func isImmediate(opr Operand) (ok bool) {
	_, ok = opr.(*Immediate)
	return
}

func isMemory(opr Operand) (ok bool) {
	_, ok = opr.(*Memory)
	return
}

func isDirectFarAddress(opr Operand) (ok bool) {
	_, ok = opr.(*DirectFarAddress)
	return
}

func isIndirectFarAddress(opr Operand) (ok bool) {
	_, ok = opr.(*IndirectFarAddress)
	return
}

func isBit8(opr Operand) bool {
	return opr.Bit() == Bit8
}

func isBit16(opr Operand) bool {
	return opr.Bit() == Bit16
}
