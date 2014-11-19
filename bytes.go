package go8086

type Mod int

const (
	Mod00 Mod = iota
	Mod01
	Mod10
	Mod11
)

type Reg int

const (
	Reg000 Reg = iota
	Reg001
	Reg010
	Reg011
	Reg100
	Reg101
	Reg110
	Reg111
)

type RM int

const (
	RM000 RM = iota
	RM001
	RM010
	RM011
	RM100
	RM101
	RM110
	RM111
)

var regAddress = [8][]*Register{
	[]*Register{BX, SI},
	[]*Register{BX, DI},
	[]*Register{BP, SI},
	[]*Register{BP, DI},
	[]*Register{SI},
	[]*Register{DI},
	[]*Register{BP},
	[]*Register{BX},
}

type SReg int

const (
	SReg00 SReg = iota
	SReg01
	SReg10
	SReg11
)

type Bytes []byte

func (bs Bytes) read32() uint32 {
	return uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)
}

func (bs Bytes) read16() uint16 {
	return uint16(bs[0]) | (uint16(bs[1]) << 8)
}

func (bs Bytes) read8() uint16 {
	return uint16(bs[0])
}

func (bs Bytes) write16(value uint16) {
	bs[0] = byte(value & 0x00ff)
	bs[1] = byte(value >> 8)
}

func (bs Bytes) write8(value uint16) {
	bs[0] = byte(value & 0x00ff)
}

func (bs Bytes) write(vs Bytes) {
	for i, v := range vs {
		bs[i] = v
	}
}

func (bs Bytes) modRM(mod Mod, reg Reg, rm RM, w Bit, sreg *SegmentRegister) (opr1, opr2 Operand, readBytes Bytes) {
	opr2 = getRegister(w, reg)
	readBytes = bs[0:1]
	switch mod {
	case Mod00:
		if rm == RM110 {
			imm := NewImmediate(bs[1:].read16(), Unsign, Bit16)
			opr1 = NewMemory(nil, imm, w, sreg)
			readBytes = bs[0:3]
		} else {
			opr1 = NewMemory(regAddress[rm], nil, w, sreg)
		}
	case Mod01:
		imm := NewImmediate(bs[1:].read8(), Sign, Bit8)
		opr1 = NewMemory(regAddress[rm], imm, w, sreg)
		readBytes = bs[0:2]
	case Mod10:
		imm := NewImmediate(bs[1:].read16(), Sign, Bit16)
		opr1 = NewMemory(regAddress[rm], imm, w, sreg)
		readBytes = bs[0:3]
	case Mod11:
		opr1 = getRegister(w, Reg(rm))
	}
	return
}

func ModRegRM(b byte) (Mod, Reg, RM) {
	x := uint8(b)
	return Mod(x >> 6), Reg((x >> 3) & 7), RM(x & 7)
}

func (bs Bytes) GetOperandByModRM(w Bit, sreg *SegmentRegister) (reg Reg, opr1, opr2 Operand, readBytes Bytes) {
	mod, reg, rm := ModRegRM(bs[0])
	opr1, opr2, readBytes = bs.modRM(mod, reg, rm, w, sreg)
	return
}

func (bs Bytes) GetOperandByModRMData(w Bit, signed Signed, sreg *SegmentRegister) (reg Reg, opr1, opr2 Operand, readBytes Bytes) {
	mod, reg, rm := ModRegRM(bs[0])
	opr1, _, readBytes = bs.modRM(mod, reg, rm, w, sreg)
	data := bs[len(readBytes):]
	switch w {
	case Bit8:
		opr2 = NewImmediate(data.read8(), Unsign, Bit8)
		readBytes = append(readBytes, data[0:1]...)
	case Bit16:
		if signed {
			opr2 = NewImmediate(data.read8(), Sign, Bit8)
			readBytes = append(readBytes, data[0:1]...)
		} else {
			opr2 = NewImmediate(data.read16(), Unsign, Bit16)
			readBytes = append(readBytes, data[0:2]...)
		}
	}
	return
}

func (bs Bytes) GetOperandOfAccImm(w Bit, wImm Bit) (opr1, opr2 Operand, readBytes Bytes) {
	var imm *Immediate
	switch wImm {
	case Bit8:
		imm, readBytes = NewImmediate(bs.read8(), Unsign, Bit8), bs[0:1]
	case Bit16:
		imm, readBytes = NewImmediate(bs.read16(), Unsign, Bit16), bs[0:2]
	}
	switch w {
	case Bit8:
		opr1, opr2 = AL, imm
	case Bit16:
		opr1, opr2 = AX, imm
	}
	return
}

func (bs Bytes) GetOperandOfAccMem(w Bit, sreg *SegmentRegister) (opr1, opr2 Operand, readBytes Bytes) {
	imm := NewImmediate(bs.read16(), Unsign, Bit16)
	readBytes = bs[0:2]
	switch w {
	case Bit8:
		opr1, opr2 = AL, NewMemory(nil, imm, Bit8, sreg)
	case Bit16:
		opr1, opr2 = AX, NewMemory(nil, imm, Bit16, sreg)
	}
	return
}

func (bs Bytes) GetOperandOfRegImm(w Bit, reg Reg) (opr1, opr2 Operand, readBytes Bytes) {
	switch w {
	case Bit8:
		opr1, opr2, readBytes = getRegister(Bit8, reg), NewImmediate(bs.read8(), Unsign, Bit8), bs[0:1]
	case Bit16:
		opr1, opr2, readBytes = getRegister(Bit16, reg), NewImmediate(bs.read16(), Unsign, Bit16), bs[0:2]
	}
	return
}
