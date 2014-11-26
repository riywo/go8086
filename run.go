package go8086

import (
	"fmt"
	"os"
)

type opcodeRunFunc func(*Opcode, *VM)

//http://stackoverflow.com/a/8037485/2052892
func CalcADC(a, b, c uint16, w Bit) (res, cf, of uint16) {
	if w == Bit8 {
		a &= 0xff
		b &= 0xff
	}
	res = a + b + c
	cf = 0
	var max uint16 = 0xffff
	if w == Bit8 {
		res &= 0xff
		max = 0xff
	}
	if c == 1 {
		if a >= max-b {
			cf = 1
		}
	} else {
		if a > max-b {
			cf = 1
		}
	}
	of = SignOf(res^a^b, w) ^ cf
	return
}

var opcodeRunFuncMap = map[Mnemonic]opcodeRunFunc{
	ADD: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), opr2.Read(vm)
		res, cf, of := CalcADC(a, b, 0, w)
		opr1.Write(vm, res)
		vm.SetFlag(CF, cf == 1)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	ADC: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), opr2.Read(vm)
		res, cf, of := CalcADC(a, b, vm.GetFlag(CF), w)
		opr1.Write(vm, res)
		vm.SetFlag(CF, cf == 1)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	SUB: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), opr2.Read(vm)
		res, cf, of := CalcADC(a, ^b, 1, w)
		opr1.Write(vm, res)
		vm.SetFlag(CF, cf == 0)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	SBB: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), opr2.Read(vm)
		res, cf, of := CalcADC(a, ^b, vm.GetFlag(CF)^1, w)
		opr1.Write(vm, res)
		vm.SetFlag(CF, cf == 0)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	CMP: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), opr2.Read(vm)
		res, cf, of := CalcADC(a, ^b, 1, w)
		DebugLog("a: %04x b: %04x res: %04x", a, b, res)
		vm.SetFlag(CF, cf == 0)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	INC: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), uint16(1)
		res, _, of := CalcADC(a, b, 0, w)
		opr1.Write(vm, res)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	DEC: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), uint16(1)
		res, _, of := CalcADC(a, ^b, 1, w)
		opr1.Write(vm, res)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	NOT: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		res := ^opr1.Read(vm)
		opr1.Write(vm, res)
	},
	NEG: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		w := opr1.Bit()
		a, b := uint16(0), opr1.Read(vm)
		res, _, of := CalcADC(a, ^b, 1, w)
		opr1.Write(vm, res)
		vm.SetFlag(CF, b != 0)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	AND: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old & op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	OR: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old | op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	XOR: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old ^ op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	TEST: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		w := opr1.Bit()
		a, b := opr1.Read(vm), opr2.Read(vm)
		res := a & b
		DebugLog("a: %04x b: %04x res: %04x", a, b, res)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	MOV: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(WritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		opr1.Write(vm, opr2.Read(vm))
	},
	XCHG: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		opr2 := op.opr2.(ReadWritableOperand)
		v1 := opr1.Read(vm)
		v2 := opr2.Read(vm)
		opr1.Write(vm, v2)
		opr2.Write(vm, v1)
	},
	LEA: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(*Register)
		opr2 := op.opr2.(*Memory)
		opr1.Write(vm, opr2.EffectiveAddress(vm))
	},
	JNC: func(op *Opcode, vm *VM) {
		if vm.GetFlag(CF) == 0 {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JC: func(op *Opcode, vm *VM) {
		if vm.GetFlag(CF) == 1 {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JNZ: func(op *Opcode, vm *VM) {
		if vm.GetFlag(ZF) == 0 {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JZ: func(op *Opcode, vm *VM) {
		if vm.GetFlag(ZF) == 1 {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JNL: func(op *Opcode, vm *VM) {
		if vm.GetFlag(SF) == vm.GetFlag(OF) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JL: func(op *Opcode, vm *VM) {
		if vm.GetFlag(SF) != vm.GetFlag(OF) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JNG: func(op *Opcode, vm *VM) {
		if (vm.GetFlag(ZF) == 1) || (vm.GetFlag(SF) != vm.GetFlag(OF)) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JG: func(op *Opcode, vm *VM) {
		if (vm.GetFlag(ZF) == 0) && (vm.GetFlag(SF) == vm.GetFlag(OF)) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JA: func(op *Opcode, vm *VM) {
		if (vm.GetFlag(CF) == 0) && (vm.GetFlag(ZF) == 0) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JNA: func(op *Opcode, vm *VM) {
		if (vm.GetFlag(CF) == 1) || (vm.GetFlag(ZF) == 1) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JCXZ: func(op *Opcode, vm *VM) {
		if CX.Read(vm) == 0 {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	MUL: func(op *Opcode, vm *VM) {
		opr := op.opr1.(ReadableOperand)
		switch opr.Bit() {
		case Bit8:
			src1 := AL.Read(vm)
			src2 := opr.Read(vm)
			res := src1 * src2
			AX.Write(vm, res)
			vm.SetFlag(CF, res>>8 != 0)
			vm.SetFlag(OF, res>>8 != 0)
		case Bit16:
			src1 := uint32(AX.Read(vm))
			src2 := uint32(opr.Read(vm))
			res := src1 * src2
			AX.Write(vm, uint16(res))
			DX.Write(vm, uint16(res>>16))
			vm.SetFlag(CF, res>>16 != 0)
			vm.SetFlag(OF, res>>16 != 0)
		}
	},
	DIV: func(op *Opcode, vm *VM) {
		opr := op.opr1.(ReadableOperand)
		switch opr.Bit() {
		case Bit8:
			dividend := AX.Read(vm)
			divisor := opr.Read(vm)
			quotient := dividend / divisor
			remainder := dividend % divisor
			AL.Write(vm, quotient)
			AH.Write(vm, remainder)
		case Bit16:
			dividend := (uint32(DX.Read(vm)) << 16) | uint32(AX.Read(vm))
			divisor := uint32(opr.Read(vm))
			quotient := dividend / divisor
			remainder := dividend % divisor
			AX.Write(vm, uint16(quotient))
			DX.Write(vm, uint16(remainder))
		}
	},
	IDIV: func(op *Opcode, vm *VM) {
		opr := op.opr1.(ReadableOperand)
		switch opr.Bit() {
		case Bit8:
			dividend := int16(AX.Read(vm))
			divisor := int16(opr.Read(vm))
			quotient := dividend / divisor
			remainder := dividend % divisor
			AL.Write(vm, uint16(quotient))
			AH.Write(vm, uint16(remainder))
		case Bit16:
			dividend := (int32(DX.Read(vm)) << 16) | int32(AX.Read(vm))
			divisor := int32(opr.Read(vm))
			quotient := dividend / divisor
			remainder := dividend % divisor
			AX.Write(vm, uint16(quotient))
			DX.Write(vm, uint16(remainder))
		}
	},
	PUSH: func(op *Opcode, vm *VM) {
		vm.Push(op.opr1.(ReadableOperand).Read(vm))
	},
	PUSHF: func(op *Opcode, vm *VM) {
		vm.Push(vm.flag)
	},
	POP: func(op *Opcode, vm *VM) {
		op.opr1.(WritableOperand).Write(vm, vm.Pop())
	},
	POPF: func(op *Opcode, vm *VM) {
		vm.flag = vm.Pop()
	},
	CALL: func(op *Opcode, vm *VM) {
		vm.Push(vm.ip)
		if isMemory(op.opr1) || isRegister(op.opr1) {
			vm.ip = op.opr1.(ReadableOperand).Read(vm)
		} else {
			vm.ip += op.opr1.(ReadableOperand).Read(vm)
		}
	},
	JMP: func(op *Opcode, vm *VM) {
		if isMemory(op.opr1) || isRegister(op.opr1) {
			vm.ip = op.opr1.(ReadableOperand).Read(vm)
		} else {
			vm.ip += op.opr1.(ReadableOperand).Read(vm)
		}
	},
	RET: func(op *Opcode, vm *VM) {
		vm.ip = vm.Pop()
		if isImmediate(op.opr1) {
			vm.reg["sp"] += op.opr1.(*Immediate).Read(vm)
		}
	},
	LOOP: func(op *Opcode, vm *VM) {
		vm.reg["cx"] -= 1
		if vm.reg["cx"] != 0 {
			vm.ip += op.opr1.(ReadableOperand).Read(vm)
		}
	},
	STD: func(op *Opcode, vm *VM) {
		vm.FlagON(DF)
	},
	CLD: func(op *Opcode, vm *VM) {
		vm.FlagOFF(DF)
	},
	SCASB: func(op *Opcode, vm *VM) {
		a, b := AL.Read(vm), vm.ES(vm.reg["di"]).read8()
		w := Bit8
		res, cf, of := CalcADC(a, ^b, 1, w)
		vm.SetFlag(CF, cf == 0)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
		if vm.GetFlag(DF) == 1 {
			vm.reg["di"] -= 1
		} else {
			vm.reg["di"] += 1
		}
	},
	CMPSB: func(op *Opcode, vm *VM) {
		a, b := vm.DS(vm.reg["si"]).read8(), vm.ES(vm.reg["di"]).read8()
		w := Bit8
		res, cf, of := CalcADC(a, ^b, 1, w)
		vm.SetFlag(CF, cf == 0)
		vm.SetFlag(OF, of == 1)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
		if vm.GetFlag(DF) == 1 {
			vm.reg["di"] -= 1
			vm.reg["si"] -= 1
		} else {
			vm.reg["di"] += 1
			vm.reg["si"] += 1
		}
	},
	STOSB: func(op *Opcode, vm *VM) {
		vm.ES(vm.reg["di"]).write8(AL.Read(vm))
		if vm.GetFlag(DF) == 1 {
			vm.reg["di"] -= 1
		} else {
			vm.reg["di"] += 1
		}
	},
	MOVSB: func(op *Opcode, vm *VM) {
		vm.ES(vm.reg["di"]).write8(vm.DS(vm.reg["si"]).read8())
		if vm.GetFlag(DF) == 1 {
			vm.reg["di"] -= 1
			vm.reg["si"] -= 1
		} else {
			vm.reg["di"] += 1
			vm.reg["si"] += 1
		}
	},
	MOVSW: func(op *Opcode, vm *VM) {
		vm.ES(vm.reg["di"]).write16(vm.DS(vm.reg["si"]).read16())
		if vm.GetFlag(DF) == 1 {
			vm.reg["di"] -= 2
			vm.reg["si"] -= 2
		} else {
			vm.reg["di"] += 2
			vm.reg["si"] += 2
		}
	},
	CBW: func(op *Opcode, vm *VM) {
		src := int8(AL.Read(vm))
		dst := int16(src)
		vm.reg["ax"] = uint16(dst)
	},
	CWD: func(op *Opcode, vm *VM) {
		src := int16(AX.Read(vm))
		dst := int32(src)
		vm.reg["ax"] = uint16(dst & 0xffff)
		vm.reg["dx"] = uint16(dst >> 16)
	},
	HLT: func(op *Opcode, vm *VM) {
		panic("HLT")
	},
	INT: func(op *Opcode, vm *VM) {
		n := op.opr1.(ReadableOperand).Read(vm)
		switch n {
		case 32:
			CallMINIXSyscall(vm)
		default:
			fmt.Fprintf(os.Stderr, "Not implemented: %s\n", op.Disasm())
			os.Exit(1)
		}
	},
	SHL: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		w := opr1.Bit()
		count := op.opr2.(*Counter).Count(vm)
		if count == 0 {
			return
		}
		old := opr1.Read(vm) << (count - 1)
		res := old << 1
		opr1.Write(vm, res)
		vm.SetFlag(CF, SignOf(old, w) == 1)
		if count == 1 {
			vm.SetFlag(OF, SignOf(res, w)^SignOf(old, w) == 1)
		}
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	SHR: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		w := opr1.Bit()
		count := op.opr2.(*Counter).Count(vm)
		if count == 0 {
			return
		}
		old := opr1.Read(vm) >> (count - 1)
		res := old >> 1
		opr1.Write(vm, res)
		vm.SetFlag(CF, (old&1) == 1)
		if count == 1 {
			vm.SetFlag(OF, SignOf(old, w) == 1)
		}
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	SAR: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		w := opr1.Bit()
		count := op.opr2.(*Counter).Count(vm)
		if count == 0 {
			return
		}
		shiftR := func(x uint16) (res uint16) {
			res = x >> 1
			if SignOf(x, w) == 1 {
				switch w {
				case Bit8:
					res |= 0x80
				case Bit16:
					res |= 0x8000
				}
			}
			return
		}
		old := opr1.Read(vm)
		for i := 0; i < int(count)-1; i++ {
			old = shiftR(old)
		}
		res := shiftR(old)
		opr1.Write(vm, res)
		vm.SetFlag(CF, (old&1) == 1)
		if count == 1 {
			vm.FlagOFF(OF)
		}
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, w) == 1)
		vm.SetFlag(PF, ParityOf(res) == 1)
	},
	RCL: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		if op.opr2.(*Counter).v == Count1 && isBit16(op.opr1) {
			old := opr1.Read(vm)
			res := old << 1
			if vm.GetFlag(CF) == 1 {
				res = res | 1
			}
			opr1.Write(vm, res)
			vm.SetFlag(CF, (old>>15) == 1)
			vm.SetFlag(OF, (res>>15)^vm.GetFlag(CF) == 1)
			vm.SetFlag(ZF, res == 0)
			vm.SetFlag(SF, SignOf(res, opr1.Bit()) == 1)
			vm.SetFlag(PF, ParityOf(res) == 1)
		} else {
			fmt.Fprintf(os.Stderr, "Not implemented: %s\n", op.Disasm())
			os.Exit(1)
		}
	},
	RCR: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		if op.opr2.(*Counter).v == Count1 && isBit16(op.opr1) {
			old := opr1.Read(vm)
			res := old >> 1
			if vm.GetFlag(CF) == 1 {
				res = res | 0x8000
			}
			opr1.Write(vm, res)
			vm.SetFlag(CF, (old&1) == 1)
			vm.SetFlag(OF, (res&1)^(res&2) == 1)
			vm.SetFlag(ZF, res == 0)
			vm.SetFlag(SF, SignOf(res, opr1.Bit()) == 1)
			vm.SetFlag(PF, ParityOf(res) == 1)
		} else {
			fmt.Fprintf(os.Stderr, "Not implemented: %s\n", op.Disasm())
			os.Exit(1)
		}
	},
}
