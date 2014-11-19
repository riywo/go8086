package go8086

import (
	"fmt"
)

type disasmFunc func(*Opcode) string

var disasmFuncMap = map[Mnemonic]disasmFunc{
	ADD:    disasmMemRegImmWithPrefix,
	ADC:    disasmMemRegImmWithPrefix,
	SUB:    disasmMemRegImmWithPrefix,
	SBB:    disasmMemRegImmWithPrefix,
	CMP:    disasmMemRegImmWithPrefix,
	AND:    disasmMemRegImmWithPrefix,
	OR:     disasmMemRegImmWithPrefix,
	XOR:    disasmMemRegImmWithPrefix,
	INC:    disasmMemRegImmWithPrefix,
	DEC:    disasmMemRegImmWithPrefix,
	PUSH:   disasmMemRegImmWithPrefix,
	POP:    disasmMemRegImmWithPrefix,
	MOV:    disasmMemRegImmWithPrefix,
	XCHG:   disasmMemRegImmWithPrefix,
	TEST:   disasmMemRegImmWithPrefix,
	NOT:    disasmMemRegImmWithPrefix,
	NEG:    disasmMemRegImmWithPrefix,
	MUL:    disasmMemRegImmWithPrefix,
	IMUL:   disasmMemRegImmWithPrefix,
	DIV:    disasmMemRegImmWithPrefix,
	IDIV:   disasmMemRegImmWithPrefix,
	SHL:    disasmMemRegImmWithPrefix,
	SHR:    disasmMemRegImmWithPrefix,
	SAR:    disasmMemRegImmWithPrefix,
	ROL:    disasmMemRegImmWithPrefix,
	ROR:    disasmMemRegImmWithPrefix,
	RCL:    disasmMemRegImmWithPrefix,
	RCR:    disasmMemRegImmWithPrefix,
	CALL:   disasmAddressWithPrefix,
	JMP:    disasmAddressWithPrefix,
	JZ:     disasmAddress,
	JL:     disasmAddress,
	JNG:    disasmAddress,
	JC:     disasmAddress,
	JNA:    disasmAddress,
	JPE:    disasmAddress,
	JO:     disasmAddress,
	JS:     disasmAddress,
	JNZ:    disasmAddress,
	JNL:    disasmAddress,
	JG:     disasmAddress,
	JNC:    disasmAddress,
	JA:     disasmAddress,
	JPO:    disasmAddress,
	JNO:    disasmAddress,
	JNS:    disasmAddress,
	LOOP:   disasmAddress,
	LOOPE:  disasmAddress,
	LOOPNE: disasmAddress,
	JCXZ:   disasmAddress,
	DB:     disasmDb,
}

var disasmMemRegImmWithPrefix = func(op *Opcode) (asm string) {
	asm = op.mn.String()
	pfx1, pfx2 := "", ""
	if isMemory(op.opr1) && (isImmediate(op.opr2) || isCounter(op.opr2) || op.opr2 == nil) {
		switch {
		case isBit8(op.opr1):
			pfx1 = "byte "
		case isBit16(op.opr1):
			pfx1 = "word "
			if op.opr2 != nil && isBit8(op.opr2) {
				pfx2 = "byte "
			}
		}
	}
	if op.opr1 != nil {
		asm += " " + pfx1 + op.opr1.Disasm()
	}
	if op.opr2 != nil {
		asm += "," + pfx2 + op.opr2.Disasm()
	}
	if op.sreg != nil && !isMemory(op.opr1) && (op.opr2 == nil || !isMemory(op.opr2)) {
		asm = op.sreg.Disasm() + " " + asm
	}
	return
}

var disasmAddressWithPrefix = func(op *Opcode) (asm string) {
	asm = op.mn.String()
	pfx := ""
	if isMemory(op.opr1) || isImmediate(op.opr1) || isDirectFarAddress(op.opr1) || isIndirectFarAddress(op.opr1) {
		switch op.opr1.Bit() {
		case Bit8:
			pfx = "short "
		case Bit16:
			pfx = "word "
		}
	}
	if isImmediate(op.opr1) {
		disp := op.opr1.(*Immediate).value
		if isBit8(op.opr1) {
			disp = uint16(int8(op.opr1.(*Immediate).value))
		}
		realAddress := uint16(op.address) + uint16(len(op.bytes)) + disp
		asm += " " + pfx + fmt.Sprintf("%#x", realAddress)
	} else {
		asm += " " + pfx + op.opr1.Disasm()
	}
	return
}

var disasmAddress = func(op *Opcode) (asm string) {
	realAddress := uint16(op.address) + uint16(len(op.bytes)) + op.opr1.(*Immediate).value
	asm = op.mn.String() + " " + fmt.Sprintf("%#x", realAddress)
	return
}

var disasmDb = func(op *Opcode) string {
	return fmt.Sprintf("db %#02x", op.bytes[0])
}

var disasmDefault = func(op *Opcode) (asm string) {
	asm = op.mn.String()
	if op.sreg != nil {
		asm = op.sreg.Disasm() + " " + asm
	}
	if op.opr1 != nil {
		asm += " " + op.opr1.Disasm()
	}
	if op.opr2 != nil {
		asm += "," + op.opr2.Disasm()
	}
	return
}
