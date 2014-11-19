package go8086

import (
	"fmt"
	"os"
)

type opcodeRunFunc func(*Opcode, *VM)

var opcodeRunFuncMap = map[Mnemonic]opcodeRunFunc{
	ADD: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old + op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.SetFlag(CF, old > res)
		vm.SetFlag(OF, old > res)
		vm.SetFlag(AF, LSBOf(old) > LSBOf(res))
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	SUB: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old - op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.SetFlag(CF, old < res)
		vm.SetFlag(OF, old < res)
		vm.SetFlag(AF, LSBOf(old) < LSBOf(res))
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	CMP: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadableOperand)
		old := opr1.Read(vm)
		res := old - op.opr2.(ReadableOperand).Read(vm)
		vm.SetFlag(CF, old < res)
		vm.SetFlag(OF, old < res)
		vm.SetFlag(AF, LSBOf(old) < LSBOf(res))
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	AND: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old & op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	OR: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old | op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	XOR: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old ^ op.opr2.(ReadableOperand).Read(vm)
		opr1.Write(vm, res)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	TEST: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadableOperand)
		old := opr1.Read(vm)
		res := old & op.opr2.(ReadableOperand).Read(vm)
		vm.FlagOFF(CF)
		vm.FlagOFF(OF)
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	MOV: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(WritableOperand)
		opr2 := op.opr2.(ReadableOperand)
		opr1.Write(vm, opr2.Read(vm))
	},
	LEA: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(*Register)
		opr2 := op.opr2.(*Memory)
		opr1.Write(vm, opr2.EffectiveAddress(vm))
	},
	INC: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old + 1
		opr1.Write(vm, res)
		vm.SetFlag(OF, old > res)
		vm.SetFlag(AF, LSBOf(old) > LSBOf(res))
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	DEC: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		old := opr1.Read(vm)
		res := old - 1
		opr1.Write(vm, res)
		vm.SetFlag(OF, old > res)
		vm.SetFlag(AF, LSBOf(old) > LSBOf(res))
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
	},
	JNC: func(op *Opcode, vm *VM) {
		if !vm.GetFlag(CF) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JNZ: func(op *Opcode, vm *VM) {
		if !vm.GetFlag(ZF) {
			vm.ip += op.opr1.(*Immediate).Read(vm)
		}
	},
	JZ: func(op *Opcode, vm *VM) {
		if vm.GetFlag(ZF) {
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
	PUSH: func(op *Opcode, vm *VM) {
		vm.Push(op.opr1.(ReadableOperand).Read(vm))
	},
	POP: func(op *Opcode, vm *VM) {
		op.opr1.(WritableOperand).Write(vm, vm.Pop())
	},
	CALL: func(op *Opcode, vm *VM) {
		vm.Push(vm.ip)
		vm.ip += op.opr1.(ReadableOperand).Read(vm)
	},
	JMP: func(op *Opcode, vm *VM) {
		vm.ip += op.opr1.(ReadableOperand).Read(vm)
	},
	RET: func(op *Opcode, vm *VM) {
		vm.ip = vm.Pop()
	},
	CLD: func(op *Opcode, vm *VM) {
		vm.FlagOFF(DF)
	},
	SCASB: func(op *Opcode, vm *VM) {
		opr1 := AL
		old := opr1.Read(vm)
		res := old - vm.DS(vm.reg["di"]).read8()
		vm.SetFlag(CF, old < res)
		vm.SetFlag(OF, old < res)
		vm.SetFlag(AF, LSBOf(old) < LSBOf(res))
		vm.SetFlag(ZF, res == 0)
		vm.SetFlag(SF, SignOf(res, opr1.Bit()))
		vm.SetFlag(PF, Parity(res))
		if vm.GetFlag(DF) {
			DI.Write(vm, DI.Read(vm)-1)
		} else {
			DI.Write(vm, DI.Read(vm)+1)
		}
	},
	HLT: func(op *Opcode, vm *VM) {
		panic("HLT")
	},
	INT: func(op *Opcode, vm *VM) {
		CallMINIXSyscall(vm)
	},
	SHL: func(op *Opcode, vm *VM) {
		opr1 := op.opr1.(ReadWritableOperand)
		if op.opr2.(*Counter).v == Count1 && isBit16(op.opr1) {
			old := opr1.Read(vm)
			res := old << 1
			opr1.Write(vm, res)
			vm.SetFlag(CF, (old>>15) == 1)
			vm.SetFlag(OF, (old>>15) != (res>>15))
			vm.SetFlag(ZF, res == 0)
			vm.SetFlag(SF, SignOf(res, opr1.Bit()))
			vm.SetFlag(PF, Parity(res))
		} else {
			fmt.Printf("Not implemented: %s\n", op.Disasm())
			os.Exit(1)
		}
	},
}
