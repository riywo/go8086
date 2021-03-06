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

type RegAddress int

const (
	RegAdd_Direct RegAddress = iota
	RegAdd_BX_SI
	RegAdd_BX_DI
	RegAdd_BP_SI
	RegAdd_BP_DI
	RegAdd_SI
	RegAdd_DI
	RegAdd_BP
	RegAdd_BX
)

var RMRegAddressMap = map[RM]RegAddress{
	RM000: RegAdd_BX_SI,
	RM001: RegAdd_BX_DI,
	RM010: RegAdd_BP_SI,
	RM011: RegAdd_BP_DI,
	RM100: RegAdd_SI,
	RM101: RegAdd_DI,
	RM110: RegAdd_BP,
	RM111: RegAdd_BX,
}

var RegAddressMap = map[RegAddress][]*Register{
	RegAdd_Direct: []*Register{},
	RegAdd_BX_SI:  []*Register{BX, SI},
	RegAdd_BX_DI:  []*Register{BX, DI},
	RegAdd_BP_SI:  []*Register{BP, SI},
	RegAdd_BP_DI:  []*Register{BP, DI},
	RegAdd_SI:     []*Register{SI},
	RegAdd_DI:     []*Register{DI},
	RegAdd_BP:     []*Register{BP},
	RegAdd_BX:     []*Register{BX},
}

var RegAddressSegment = map[RegAddress]*SegmentRegister{
	RegAdd_Direct: DS,
	RegAdd_BX_SI:  DS,
	RegAdd_BX_DI:  DS,
	RegAdd_BP_SI:  SS,
	RegAdd_BP_DI:  SS,
	RegAdd_SI:     DS,
	RegAdd_DI:     DS,
	RegAdd_BP:     SS,
	RegAdd_BX:     DS,
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

func (bs Bytes) write32(value uint32) {
	bs[0] = byte(value & 0x00ff)
	bs[1] = byte(value >> 8)
	bs[2] = byte(value >> 16)
	bs[3] = byte(value >> 24)
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
			opr1 = NewMemory(RegAdd_Direct, imm, w, sreg)
			readBytes = bs[0:3]
		} else {
			opr1 = NewMemory(RMRegAddressMap[rm], nil, w, sreg)
		}
	case Mod01:
		imm := NewImmediate(bs[1:].read8(), Sign, Bit8)
		opr1 = NewMemory(RMRegAddressMap[rm], imm, w, sreg)
		readBytes = bs[0:2]
	case Mod10:
		imm := NewImmediate(bs[1:].read16(), Sign, Bit16)
		opr1 = NewMemory(RMRegAddressMap[rm], imm, w, sreg)
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
		opr1, opr2 = AL, NewMemory(RegAdd_Direct, imm, Bit8, sreg)
	case Bit16:
		opr1, opr2 = AX, NewMemory(RegAdd_Direct, imm, Bit16, sreg)
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
