package go8086

type setOpcodeFunc func(xs Bytes, op *Opcode) (bs Bytes)

func setOpcodeImm(mn Mnemonic, w Bit, signed Signed) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.mn = mn
		switch w {
		case Bit8:
			op.opr1, bs = NewImmediate(xs.read8(), signed, Bit8), xs[0:1]
		case Bit16:
			op.opr1, bs = NewImmediate(xs.read16(), signed, Bit16), xs[0:2]
		}
		return
	}
}

func setOpcodeDirectFarAddress(mn Mnemonic) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.mn = mn
		offset := xs.read16()
		segment := xs[2:].read16()
		op.opr1 = NewDirectFarAddress(segment, offset)
		bs = xs[0:4]
		return
	}
}

func setOpcodeByModRM(mn Mnemonic, w Bit, d Direction) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.mn = mn
		_, opr, reg, bs := xs.GetOperandByModRM(w, op.sreg)
		switch d {
		case FromReg:
			op.opr1, op.opr2 = opr, reg
		case ToReg:
			op.opr1, op.opr2 = reg, opr
		}
		return
	}
}

func setOpcodeByModRMSReg(mn Mnemonic, d Direction) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		reg, opr, _, bs := xs.GetOperandByModRM(Bit16, op.sreg)
		if reg < Reg100 {
			sreg := getSegmentRegister(SReg(reg & 3))
			switch d {
			case FromReg:
				op.opr1, op.opr2 = opr, sreg
			case ToReg:
				op.opr1, op.opr2 = sreg, opr
			}
			op.mn = mn
		}
		return
	}
}

func setOpcodeByModRMLoad(mn Mnemonic) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		_, opr, reg, bs := xs.GetOperandByModRM(Bit16, op.sreg)
		if !isRegister(opr) {
			op.mn, op.opr1, op.opr2 = mn, reg, opr
		}
		return
	}
}

func setOpcodeImmAcc(mn Mnemonic, w Bit, immw Bit, d Direction) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		acc, imm, bs := xs.GetOperandOfAccImm(w, immw)
		switch d {
		case FromReg:
			op.opr1, op.opr2 = imm, acc
		case ToReg:
			op.opr1, op.opr2 = acc, imm
		}
		op.mn = mn
		return
	}
}

func setOpcodeMemAcc(mn Mnemonic, w Bit, d Direction) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		acc, mem, bs := xs.GetOperandOfAccMem(w, op.sreg)
		switch d {
		case FromReg:
			op.opr1, op.opr2 = mem, acc
		case ToReg:
			op.opr1, op.opr2 = acc, mem
		}
		op.mn = mn
		return
	}
}

func setOpcodeRegImm(mn Mnemonic, w Bit, reg Reg) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.opr1, op.opr2, bs = xs.GetOperandOfRegImm(w, reg)
		op.mn = mn
		return
	}
}

func setOpcodeOneOperandMultiMnemonics(w Bit, mns ...Mnemonic) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		reg, opr1, opr2, bs := xs.GetOperandByModRM(w, op.sreg)
		mn := mns[reg]
		switch mn {
		case TEST:
			_, op.opr1, op.opr2, bs = xs.GetOperandByModRMData(w, Unsign, op.sreg)
			op.mn = mn
		case CALL:
			if reg == Reg010 {
				op.mn, op.opr1 = mn, opr1
			} else if reg == Reg011 && !isMemory(opr2) {
				op.mn = mn
				op.opr1 = NewIndirectFarAddress(opr1.(*Memory))
			}
		case JMP:
			if reg == Reg100 {
				op.mn, op.opr1 = mn, opr1
			} else if reg == Reg101 && !isMemory(opr2) {
				op.mn = mn
				op.opr1 = NewIndirectFarAddress(opr1.(*Memory))
			}
		case NIL:
		default:
			op.mn, op.opr1 = mn, opr1
		}
		return
	}
}

func setOpcodeCountMultiMnemonics(v Count, w Bit, mns ...Mnemonic) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		reg, opr, _, bs := xs.GetOperandByModRM(w, op.sreg)
		mn := mns[reg]
		if mn != NIL {
			op.opr1, op.opr2 = opr, NewCounter(v, w)
			op.mn = mn
		}
		return
	}
}

func setOpcodeMultiMnemonics(w Bit, s Signed, mns ...Mnemonic) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		reg, opr1, opr2, bs := xs.GetOperandByModRMData(w, s, op.sreg)
		mn := mns[reg]
		if mn != NIL {
			op.mn, op.opr1, op.opr2 = mn, opr1, opr2
		}
		return
	}
}

func setOpcodeReg(mn Mnemonic, reg Reg) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.mn = mn
		op.opr1 = getRegister(Bit16, reg)
		return
	}
}

func setOpcodeSReg(mn Mnemonic, sreg SReg) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.mn = mn
		op.opr1 = getSegmentRegister(sreg)
		return
	}
}

func setSRegOverridePrefix(sreg SReg) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		*op = *getOpcode(getSegmentRegister(sreg), op.address+1, xs)
		return
	}
}

func setOpcodeRegReg(mn Mnemonic, register1, register2 *Register) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.mn, op.opr1, op.opr2 = mn, register1, register2
		return
	}
}

func setOpcodeNoOperand(mn Mnemonic, followed ...byte) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		if len(followed) == 0 {
			op.mn = mn
		} else if followed[0] == xs[0] {
			op.mn = mn
			bs = xs[0:1]
		}
		return
	}
}

func setOpcodePrefix(mn Mnemonic) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.following = getOpcode(op.sreg, op.address+1, xs)
		op.mn = mn
		op.bytes = op.following.bytes
		return
	}
}

func setOpcodeDb(x byte) setOpcodeFunc {
	return func(xs Bytes, op *Opcode) (bs Bytes) {
		op.mn, bs = DB, Bytes{x}
		return
	}
}

func dispatchByFirstByte(x byte) (f setOpcodeFunc) {
	switch x {
	case 0x00:
		f = setOpcodeByModRM(ADD, Bit8, FromReg)
	case 0x01:
		f = setOpcodeByModRM(ADD, Bit16, FromReg)
	case 0x02:
		f = setOpcodeByModRM(ADD, Bit8, ToReg)
	case 0x03:
		f = setOpcodeByModRM(ADD, Bit16, ToReg)
	case 0x04:
		f = setOpcodeImmAcc(ADD, Bit8, Bit8, ToReg)
	case 0x05:
		f = setOpcodeImmAcc(ADD, Bit16, Bit16, ToReg)
	case 0x06:
		f = setOpcodeSReg(PUSH, SReg00)
	case 0x07:
		f = setOpcodeSReg(POP, SReg00)
	case 0x08:
		f = setOpcodeByModRM(OR, Bit8, FromReg)
	case 0x09:
		f = setOpcodeByModRM(OR, Bit16, FromReg)
	case 0x0A:
		f = setOpcodeByModRM(OR, Bit8, ToReg)
	case 0x0B:
		f = setOpcodeByModRM(OR, Bit16, ToReg)
	case 0x0C:
		f = setOpcodeImmAcc(OR, Bit8, Bit8, ToReg)
	case 0x0D:
		f = setOpcodeImmAcc(OR, Bit16, Bit16, ToReg)
	case 0x0E:
		f = setOpcodeSReg(PUSH, SReg01)

	case 0x10:
		f = setOpcodeByModRM(ADC, Bit8, FromReg)
	case 0x11:
		f = setOpcodeByModRM(ADC, Bit16, FromReg)
	case 0x12:
		f = setOpcodeByModRM(ADC, Bit8, ToReg)
	case 0x13:
		f = setOpcodeByModRM(ADC, Bit16, ToReg)
	case 0x14:
		f = setOpcodeImmAcc(ADC, Bit8, Bit8, ToReg)
	case 0x15:
		f = setOpcodeImmAcc(ADC, Bit16, Bit16, ToReg)
	case 0x16:
		f = setOpcodeSReg(PUSH, SReg10)
	case 0x17:
		f = setOpcodeSReg(POP, SReg10)
	case 0x18:
		f = setOpcodeByModRM(SBB, Bit8, FromReg)
	case 0x19:
		f = setOpcodeByModRM(SBB, Bit16, FromReg)
	case 0x1A:
		f = setOpcodeByModRM(SBB, Bit8, ToReg)
	case 0x1B:
		f = setOpcodeByModRM(SBB, Bit16, ToReg)
	case 0x1C:
		f = setOpcodeImmAcc(SBB, Bit8, Bit8, ToReg)
	case 0x1D:
		f = setOpcodeImmAcc(SBB, Bit16, Bit16, ToReg)
	case 0x1E:
		f = setOpcodeSReg(PUSH, SReg11)
	case 0x1F:
		f = setOpcodeSReg(POP, SReg11)

	case 0x20:
		f = setOpcodeByModRM(AND, Bit8, FromReg)
	case 0x21:
		f = setOpcodeByModRM(AND, Bit16, FromReg)
	case 0x22:
		f = setOpcodeByModRM(AND, Bit8, ToReg)
	case 0x23:
		f = setOpcodeByModRM(AND, Bit16, ToReg)
	case 0x24:
		f = setOpcodeImmAcc(AND, Bit8, Bit8, ToReg)
	case 0x25:
		f = setOpcodeImmAcc(AND, Bit16, Bit16, ToReg)
	case 0x26:
		f = setSRegOverridePrefix(SReg00)
	case 0x27:
		f = setOpcodeNoOperand(DAA)
	case 0x28:
		f = setOpcodeByModRM(SUB, Bit8, FromReg)
	case 0x29:
		f = setOpcodeByModRM(SUB, Bit16, FromReg)
	case 0x2A:
		f = setOpcodeByModRM(SUB, Bit8, ToReg)
	case 0x2B:
		f = setOpcodeByModRM(SUB, Bit16, ToReg)
	case 0x2C:
		f = setOpcodeImmAcc(SUB, Bit8, Bit8, ToReg)
	case 0x2D:
		f = setOpcodeImmAcc(SUB, Bit16, Bit16, ToReg)
	case 0x2E:
		f = setSRegOverridePrefix(SReg01)
	case 0x2F:
		f = setOpcodeNoOperand(DAS)

	case 0x30:
		f = setOpcodeByModRM(XOR, Bit8, FromReg)
	case 0x31:
		f = setOpcodeByModRM(XOR, Bit16, FromReg)
	case 0x32:
		f = setOpcodeByModRM(XOR, Bit8, ToReg)
	case 0x33:
		f = setOpcodeByModRM(XOR, Bit16, ToReg)
	case 0x34:
		f = setOpcodeImmAcc(XOR, Bit8, Bit8, ToReg)
	case 0x35:
		f = setOpcodeImmAcc(XOR, Bit16, Bit16, ToReg)
	case 0x36:
		f = setSRegOverridePrefix(SReg10)
	case 0x37:
		f = setOpcodeNoOperand(AAA)
	case 0x38:
		f = setOpcodeByModRM(CMP, Bit8, FromReg)
	case 0x39:
		f = setOpcodeByModRM(CMP, Bit16, FromReg)
	case 0x3A:
		f = setOpcodeByModRM(CMP, Bit8, ToReg)
	case 0x3B:
		f = setOpcodeByModRM(CMP, Bit16, ToReg)
	case 0x3C:
		f = setOpcodeImmAcc(CMP, Bit8, Bit8, ToReg)
	case 0x3D:
		f = setOpcodeImmAcc(CMP, Bit16, Bit16, ToReg)
	case 0x3E:
		f = setSRegOverridePrefix(SReg11)
	case 0x3F:
		f = setOpcodeNoOperand(AAS)

	case 0x40:
		f = setOpcodeReg(INC, Reg000)
	case 0x41:
		f = setOpcodeReg(INC, Reg001)
	case 0x42:
		f = setOpcodeReg(INC, Reg010)
	case 0x43:
		f = setOpcodeReg(INC, Reg011)
	case 0x44:
		f = setOpcodeReg(INC, Reg100)
	case 0x45:
		f = setOpcodeReg(INC, Reg101)
	case 0x46:
		f = setOpcodeReg(INC, Reg110)
	case 0x47:
		f = setOpcodeReg(INC, Reg111)
	case 0x48:
		f = setOpcodeReg(DEC, Reg000)
	case 0x49:
		f = setOpcodeReg(DEC, Reg001)
	case 0x4A:
		f = setOpcodeReg(DEC, Reg010)
	case 0x4B:
		f = setOpcodeReg(DEC, Reg011)
	case 0x4C:
		f = setOpcodeReg(DEC, Reg100)
	case 0x4D:
		f = setOpcodeReg(DEC, Reg101)
	case 0x4E:
		f = setOpcodeReg(DEC, Reg110)
	case 0x4F:
		f = setOpcodeReg(DEC, Reg111)

	case 0x50:
		f = setOpcodeReg(PUSH, Reg000)
	case 0x51:
		f = setOpcodeReg(PUSH, Reg001)
	case 0x52:
		f = setOpcodeReg(PUSH, Reg010)
	case 0x53:
		f = setOpcodeReg(PUSH, Reg011)
	case 0x54:
		f = setOpcodeReg(PUSH, Reg100)
	case 0x55:
		f = setOpcodeReg(PUSH, Reg101)
	case 0x56:
		f = setOpcodeReg(PUSH, Reg110)
	case 0x57:
		f = setOpcodeReg(PUSH, Reg111)
	case 0x58:
		f = setOpcodeReg(POP, Reg000)
	case 0x59:
		f = setOpcodeReg(POP, Reg001)
	case 0x5A:
		f = setOpcodeReg(POP, Reg010)
	case 0x5B:
		f = setOpcodeReg(POP, Reg011)
	case 0x5C:
		f = setOpcodeReg(POP, Reg100)
	case 0x5D:
		f = setOpcodeReg(POP, Reg101)
	case 0x5E:
		f = setOpcodeReg(POP, Reg110)
	case 0x5F:
		f = setOpcodeReg(POP, Reg111)

	case 0x70:
		f = setOpcodeImm(JO, Bit8, Sign)
	case 0x71:
		f = setOpcodeImm(JNO, Bit8, Sign)
	case 0x72:
		f = setOpcodeImm(JC, Bit8, Sign)
	case 0x73:
		f = setOpcodeImm(JNC, Bit8, Sign)
	case 0x74:
		f = setOpcodeImm(JZ, Bit8, Sign)
	case 0x75:
		f = setOpcodeImm(JNZ, Bit8, Sign)
	case 0x76:
		f = setOpcodeImm(JNA, Bit8, Sign)
	case 0x77:
		f = setOpcodeImm(JA, Bit8, Sign)
	case 0x78:
		f = setOpcodeImm(JS, Bit8, Sign)
	case 0x79:
		f = setOpcodeImm(JNS, Bit8, Sign)
	case 0x7A:
		f = setOpcodeImm(JPE, Bit8, Sign)
	case 0x7B:
		f = setOpcodeImm(JPO, Bit8, Sign)
	case 0x7C:
		f = setOpcodeImm(JL, Bit8, Sign)
	case 0x7D:
		f = setOpcodeImm(JNL, Bit8, Sign)
	case 0x7E:
		f = setOpcodeImm(JNG, Bit8, Sign)
	case 0x7F:
		f = setOpcodeImm(JG, Bit8, Sign)

	case 0x80:
		f = setOpcodeMultiMnemonics(Bit8, Unsign, ADD, OR, ADC, SBB, AND, SUB, XOR, CMP)
	case 0x81:
		f = setOpcodeMultiMnemonics(Bit16, Unsign, ADD, OR, ADC, SBB, AND, SUB, XOR, CMP)
	case 0x83:
		f = setOpcodeMultiMnemonics(Bit16, Sign, ADD, OR, ADC, SBB, AND, SUB, XOR, CMP)
	case 0x84:
		f = setOpcodeByModRM(TEST, Bit8, FromReg)
	case 0x85:
		f = setOpcodeByModRM(TEST, Bit16, FromReg)
	case 0x86:
		f = setOpcodeByModRM(XCHG, Bit8, ToReg)
	case 0x87:
		f = setOpcodeByModRM(XCHG, Bit16, ToReg)
	case 0x88:
		f = setOpcodeByModRM(MOV, Bit8, FromReg)
	case 0x89:
		f = setOpcodeByModRM(MOV, Bit16, FromReg)
	case 0x8A:
		f = setOpcodeByModRM(MOV, Bit8, ToReg)
	case 0x8B:
		f = setOpcodeByModRM(MOV, Bit16, ToReg)
	case 0x8C:
		f = setOpcodeByModRMSReg(MOV, FromReg)
	case 0x8D:
		f = setOpcodeByModRMLoad(LEA)
	case 0x8E:
		f = setOpcodeByModRMSReg(MOV, ToReg)
	case 0x8F:
		f = setOpcodeOneOperandMultiMnemonics(Bit16, POP)

	case 0x90:
		f = setOpcodeNoOperand(NOP)
	case 0x91:
		f = setOpcodeRegReg(XCHG, AX, CX)
	case 0x92:
		f = setOpcodeRegReg(XCHG, AX, DX)
	case 0x93:
		f = setOpcodeRegReg(XCHG, AX, BX)
	case 0x94:
		f = setOpcodeRegReg(XCHG, AX, SP)
	case 0x95:
		f = setOpcodeRegReg(XCHG, AX, BP)
	case 0x96:
		f = setOpcodeRegReg(XCHG, AX, SI)
	case 0x97:
		f = setOpcodeRegReg(XCHG, AX, DI)
	case 0x98:
		f = setOpcodeNoOperand(CBW)
	case 0x99:
		f = setOpcodeNoOperand(CWD)
	case 0x9A:
		f = setOpcodeDirectFarAddress(CALL)
	case 0x9B:
		f = setOpcodePrefix(WAIT)
	case 0x9C:
		f = setOpcodeNoOperand(PUSHF)
	case 0x9D:
		f = setOpcodeNoOperand(POPF)
	case 0x9E:
		f = setOpcodeNoOperand(SAHF)
	case 0x9F:
		f = setOpcodeNoOperand(LAHF)

	case 0xA0:
		f = setOpcodeMemAcc(MOV, Bit8, ToReg)
	case 0xA1:
		f = setOpcodeMemAcc(MOV, Bit16, ToReg)
	case 0xA2:
		f = setOpcodeMemAcc(MOV, Bit8, FromReg)
	case 0xA3:
		f = setOpcodeMemAcc(MOV, Bit16, FromReg)
	case 0xA4:
		f = setOpcodeNoOperand(MOVSB)
	case 0xA5:
		f = setOpcodeNoOperand(MOVSW)
	case 0xA6:
		f = setOpcodeNoOperand(CMPSB)
	case 0xA7:
		f = setOpcodeNoOperand(CMPSW)
	case 0xA8:
		f = setOpcodeImmAcc(TEST, Bit8, Bit8, ToReg)
	case 0xA9:
		f = setOpcodeImmAcc(TEST, Bit16, Bit16, ToReg)
	case 0xAA:
		f = setOpcodeNoOperand(STOSB)
	case 0xAB:
		f = setOpcodeNoOperand(STOSW)
	case 0xAC:
		f = setOpcodeNoOperand(LODSB)
	case 0xAD:
		f = setOpcodeNoOperand(LODSW)
	case 0xAE:
		f = setOpcodeNoOperand(SCASB)
	case 0xAF:
		f = setOpcodeNoOperand(SCASW)

	case 0xB0:
		f = setOpcodeRegImm(MOV, Bit8, Reg000)
	case 0xB1:
		f = setOpcodeRegImm(MOV, Bit8, Reg001)
	case 0xB2:
		f = setOpcodeRegImm(MOV, Bit8, Reg010)
	case 0xB3:
		f = setOpcodeRegImm(MOV, Bit8, Reg011)
	case 0xB4:
		f = setOpcodeRegImm(MOV, Bit8, Reg100)
	case 0xB5:
		f = setOpcodeRegImm(MOV, Bit8, Reg101)
	case 0xB6:
		f = setOpcodeRegImm(MOV, Bit8, Reg110)
	case 0xB7:
		f = setOpcodeRegImm(MOV, Bit8, Reg111)
	case 0xB8:
		f = setOpcodeRegImm(MOV, Bit16, Reg000)
	case 0xB9:
		f = setOpcodeRegImm(MOV, Bit16, Reg001)
	case 0xBA:
		f = setOpcodeRegImm(MOV, Bit16, Reg010)
	case 0xBB:
		f = setOpcodeRegImm(MOV, Bit16, Reg011)
	case 0xBC:
		f = setOpcodeRegImm(MOV, Bit16, Reg100)
	case 0xBD:
		f = setOpcodeRegImm(MOV, Bit16, Reg101)
	case 0xBE:
		f = setOpcodeRegImm(MOV, Bit16, Reg110)
	case 0xBF:
		f = setOpcodeRegImm(MOV, Bit16, Reg111)

	case 0xC2:
		f = setOpcodeImm(RET, Bit16, Unsign)
	case 0xC3:
		f = setOpcodeNoOperand(RET)
	case 0xC4:
		f = setOpcodeByModRMLoad(LES)
	case 0xC5:
		f = setOpcodeByModRMLoad(LDS)
	case 0xC6:
		f = setOpcodeMultiMnemonics(Bit8, Unsign, MOV)
	case 0xC7:
		f = setOpcodeMultiMnemonics(Bit16, Unsign, MOV)
	case 0xCA:
		f = setOpcodeImm(RETF, Bit16, Unsign)
	case 0xCB:
		f = setOpcodeNoOperand(RETF)
	case 0xCC:
		f = setOpcodeNoOperand(INT3)
	case 0xCD:
		f = setOpcodeImm(INT, Bit8, Unsign)
	case 0xCE:
		f = setOpcodeNoOperand(INTO)
	case 0xCF:
		f = setOpcodeNoOperand(IRET)

	case 0xD0:
		f = setOpcodeCountMultiMnemonics(Count1, Bit8, ROL, ROR, RCL, RCR, SHL, SHR, NIL, SAR)
	case 0xD1:
		f = setOpcodeCountMultiMnemonics(Count1, Bit16, ROL, ROR, RCL, RCR, SHL, SHR, NIL, SAR)
	case 0xD2:
		f = setOpcodeCountMultiMnemonics(CountCL, Bit8, ROL, ROR, RCL, RCR, SHL, SHR, NIL, SAR)
	case 0xD3:
		f = setOpcodeCountMultiMnemonics(CountCL, Bit16, ROL, ROR, RCL, RCR, SHL, SHR, NIL, SAR)
	case 0xD4:
		f = setOpcodeNoOperand(AAM, 0x0A)
	case 0xD5:
		f = setOpcodeNoOperand(AAD, 0x0A)
	case 0xD7:
		f = setOpcodeNoOperand(XLAT)

	case 0xE0:
		f = setOpcodeImm(LOOPNE, Bit8, Sign)
	case 0xE1:
		f = setOpcodeImm(LOOPE, Bit8, Sign)
	case 0xE2:
		f = setOpcodeImm(LOOP, Bit8, Sign)
	case 0xE3:
		f = setOpcodeImm(JCXZ, Bit8, Sign)
	case 0xE4:
		f = setOpcodeImmAcc(IN, Bit8, Bit8, ToReg)
	case 0xE5:
		f = setOpcodeImmAcc(IN, Bit16, Bit8, ToReg)
	case 0xE6:
		f = setOpcodeImmAcc(OUT, Bit8, Bit8, FromReg)
	case 0xE7:
		f = setOpcodeImmAcc(OUT, Bit16, Bit8, FromReg)
	case 0xE8:
		f = setOpcodeImm(CALL, Bit16, Sign)
	case 0xE9:
		f = setOpcodeImm(JMP, Bit16, Sign)
	case 0xEA:
		f = setOpcodeDirectFarAddress(JMP)
	case 0xEB:
		f = setOpcodeImm(JMP, Bit8, Sign)
	case 0xEC:
		f = setOpcodeRegReg(IN, AL, DX)
	case 0xED:
		f = setOpcodeRegReg(IN, AX, DX)
	case 0xEE:
		f = setOpcodeRegReg(OUT, DX, AL)
	case 0xEF:
		f = setOpcodeRegReg(OUT, DX, AX)

	case 0xF0:
		f = setOpcodePrefix(LOCK)
	case 0xF2:
		f = setOpcodePrefix(REPNE)
	case 0xF3:
		f = setOpcodePrefix(REP)
	case 0xF4:
		f = setOpcodeNoOperand(HLT)
	case 0xF5:
		f = setOpcodeNoOperand(CMC)
	case 0xF6:
		f = setOpcodeOneOperandMultiMnemonics(Bit8, TEST, NIL, NOT, NEG, MUL, IMUL, DIV, IDIV)
	case 0xF7:
		f = setOpcodeOneOperandMultiMnemonics(Bit16, TEST, NIL, NOT, NEG, MUL, IMUL, DIV, IDIV)
	case 0xF8:
		f = setOpcodeNoOperand(CLC)
	case 0xF9:
		f = setOpcodeNoOperand(STC)
	case 0xFA:
		f = setOpcodeNoOperand(CLI)
	case 0xFB:
		f = setOpcodeNoOperand(STI)
	case 0xFC:
		f = setOpcodeNoOperand(CLD)
	case 0xFD:
		f = setOpcodeNoOperand(STD)
	case 0xFE:
		f = setOpcodeOneOperandMultiMnemonics(Bit8, INC, DEC)
	case 0xFF:
		f = setOpcodeOneOperandMultiMnemonics(Bit16, INC, DEC, CALL, CALL, JMP, JMP, PUSH)
	default:
		f = setOpcodeDb(x)
	}
	return
}

func getOpcode(sreg *SegmentRegister, address uint16, bs Bytes) (op *Opcode) {
	defer func() {
		if err := recover(); err != nil {
			op = new(Opcode)
			op.mn = DB
			op.bytes = bs[0:1]
			op.address = address
		}
	}()

	x, xs := bs[0], bs[1:]
	op = new(Opcode)
	op.sreg = sreg

	readBytes := dispatchByFirstByte(x)(xs, op)
	op.bytes = append(Bytes{x}, op.bytes...)
	op.bytes = append(op.bytes, readBytes...)
	op.address = address

	if op.mn == NIL {
		panic("")
	}
	return
}
