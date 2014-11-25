package go8086

import (
	"fmt"
	"os"
)

type Mnemonic int

const (
	NIL Mnemonic = iota
	ADD
	ADC
	SUB
	SBB
	CMP
	AND
	OR
	XOR
	INC
	DEC
	PUSH
	POP
	MOV
	XCHG
	IN
	OUT
	LEA
	LDS
	LES
	TEST
	NOT
	NEG
	MUL
	IMUL
	DIV
	IDIV
	SHL
	SHR
	SAR
	ROL
	ROR
	RCL
	RCR
	CALL
	JMP
	RET
	RETF
	JZ
	JL
	JNG
	JC
	JNA
	JPE
	JO
	JS
	JNZ
	JNL
	JG
	JNC
	JA
	JPO
	JNO
	JNS
	LOOP
	LOOPE
	LOOPNE
	JCXZ
	INT
	REP
	REPNE
	LOCK
	WAIT
	XLAT
	LAHF
	SAHF
	PUSHF
	POPF
	AAM
	AAD
	AAA
	DAA
	AAS
	DAS
	CBW
	CWD
	MOVSB
	MOVSW
	CMPSB
	CMPSW
	SCASB
	SCASW
	LODSB
	LODSW
	STOSB
	STOSW
	INT3
	INTO
	IRET
	CLC
	CMC
	STC
	CLD
	STD
	CLI
	STI
	HLT
	NOP
	DB
)

var mnemonicString = map[Mnemonic]string{
	ADD:    "add",
	ADC:    "adc",
	SUB:    "sub",
	SBB:    "sbb",
	CMP:    "cmp",
	AND:    "and",
	OR:     "or",
	XOR:    "xor",
	INC:    "inc",
	DEC:    "dec",
	PUSH:   "push",
	POP:    "pop",
	MOV:    "mov",
	XCHG:   "xchg",
	IN:     "in",
	OUT:    "out",
	LEA:    "lea",
	LDS:    "lds",
	LES:    "les",
	TEST:   "test",
	NOT:    "not",
	NEG:    "neg",
	MUL:    "mul",
	IMUL:   "imul",
	DIV:    "div",
	IDIV:   "idiv",
	SHL:    "shl",
	SHR:    "shr",
	SAR:    "sar",
	ROL:    "rol",
	ROR:    "ror",
	RCL:    "rcl",
	RCR:    "rcr",
	CALL:   "call",
	JMP:    "jmp",
	RET:    "ret",
	RETF:   "retf",
	JZ:     "jz",
	JL:     "jl",
	JNG:    "jng",
	JC:     "jc",
	JNA:    "jna",
	JPE:    "jpe",
	JO:     "jo",
	JS:     "js",
	JNZ:    "jnz",
	JNL:    "jnl",
	JG:     "jg",
	JNC:    "jnc",
	JA:     "ja",
	JPO:    "jpo",
	JNO:    "jno",
	JNS:    "jns",
	LOOP:   "loop",
	LOOPE:  "loope",
	LOOPNE: "loopne",
	JCXZ:   "jcxz",
	INT:    "int",
	REP:    "rep",
	REPNE:  "repne",
	LOCK:   "lock",
	WAIT:   "wait",
	XLAT:   "xlatb",
	LAHF:   "lahf",
	SAHF:   "sahf",
	PUSHF:  "pushfw",
	POPF:   "popfw",
	AAM:    "aam",
	AAD:    "aad",
	AAA:    "aaa",
	DAA:    "daa",
	AAS:    "aas",
	DAS:    "das",
	CBW:    "cbw",
	CWD:    "cwd",
	MOVSB:  "movsb",
	MOVSW:  "movsw",
	CMPSB:  "cmpsb",
	CMPSW:  "cmpsw",
	SCASB:  "scasb",
	SCASW:  "scasw",
	LODSB:  "lodsb",
	LODSW:  "lodsw",
	STOSB:  "stosb",
	STOSW:  "stosw",
	INT3:   "int3",
	INTO:   "into",
	IRET:   "iretw",
	CLC:    "clc",
	CMC:    "cmc",
	STC:    "stc",
	CLD:    "cld",
	STD:    "std",
	CLI:    "cli",
	STI:    "sti",
	HLT:    "hlt",
	NOP:    "nop",
}

func (mn Mnemonic) String() string {
	return mnemonicString[mn]
}

type Opcode struct {
	mn        Mnemonic
	opr1      Operand
	opr2      Operand
	bytes     Bytes
	address   uint16
	sreg      *SegmentRegister
	following *Opcode
}

func (op *Opcode) Disasm() (asm string) {
	if op.following != nil {
		asm = op.mn.String() + " " + op.following.Disasm()
	} else {
		f := disasmFuncMap[op.mn]
		if f == nil {
			f = disasmDefault
		}
		asm = f(op)
	}
	return
}

func (op *Opcode) Run(vm *VM) {
	switch op.mn {
	case REPNE:
		s := ""
		data := vm.ES(vm.reg["di"])
		if op.following.mn == MOVSB {
			data = vm.DS(vm.reg["si"])
		}
		for _, c := range data[0:50] {
			s += fmt.Sprintf("%c", c)
		}
		DebugLog("rep %s: %s", op.following.Disasm(), s)
		for {
			op.following.Run(vm)
			CX.Write(vm, CX.Read(vm)-1)
			if CX.Read(vm) == 0 || vm.GetFlag(ZF) == 1 {
				return
			}
		}
	default:
		f := opcodeRunFuncMap[op.mn]
		if f != nil {
			f(op, vm)
		} else {
			fmt.Fprintf(os.Stderr, "Not implemented: %s\n", op.Disasm())
			os.Exit(1)
		}
	}
	return
}
