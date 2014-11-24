package go8086

import (
	"fmt"
	"os"
)

type Flag uint

const (
	OF Flag = 10
	DF Flag = 9
	IF Flag = 8
	TF Flag = 7
	SF Flag = 6
	ZF Flag = 5
	AF Flag = 3
	PF Flag = 1
	CF Flag = 0
)

type VM struct {
	reg  map[string]uint16
	sreg map[string]uint16
	ip   uint16
	flag uint16
	mem  Bytes
}

func NewVM() (vm *VM) {
	vm = new(VM)
	vm.reg = make(map[string]uint16)
	vm.sreg = make(map[string]uint16)
	vm.mem = make(Bytes, 0x100000)
	for _, reg := range regs[Bit16] {
		switch reg {
		case SP:
			vm.reg[reg.name] = 0xff0e
		default:
			vm.reg[reg.name] = 0
		}
	}
	for _, sreg := range sregs {
		switch sreg {
		case CS:
			vm.sreg[sreg.name] = 0x1000
		default:
			vm.sreg[sreg.name] = 0
		}
	}
	DebugLog("CS: %04x SS: %04x DS: %04x ES: %04x", vm.sreg["cs"], vm.sreg["ss"], vm.sreg["ds"], vm.sreg["es"])
	return
}

func (vm *VM) Mem(sreg *SegmentRegister, offset uint16) Bytes {
	ea := (uint32(sreg.Read(vm)) << 4) + uint32(offset)
	return vm.mem[ea:]
}

func (vm *VM) CS(offset uint16) Bytes {
	return vm.Mem(CS, offset)
}

func (vm *VM) DS(offset uint16) Bytes {
	return vm.Mem(DS, offset)
}

func (vm *VM) SS(offset uint16) Bytes {
	return vm.Mem(SS, offset)
}

func (vm *VM) ES(offset uint16) Bytes {
	return vm.Mem(ES, offset)
}

func (vm *VM) Push(value uint16) {
	vm.reg["sp"] -= 2
	vm.SS(vm.reg["sp"]).write16(value)
}

func (vm *VM) Pop() (value uint16) {
	value = vm.SS(vm.reg["sp"]).read16()
	vm.reg["sp"] += 2
	return
}

func (vm *VM) GetFlag(f Flag) uint16 {
	return (vm.flag >> f) & 1
}

func (vm *VM) FlagON(f Flag) {
	vm.flag = vm.flag | (1 << f)
}

func (vm *VM) FlagOFF(f Flag) {
	vm.flag = vm.flag & ((1 << f) ^ 0xffff)
}

func (vm *VM) SetFlag(f Flag, condition bool) {
	if condition {
		vm.FlagON(f)
	} else {
		vm.FlagOFF(f)
	}
}

func (vm *VM) getOpcode() (op *Opcode) {
	return getOpcode(nil, vm.ip, vm.CS(vm.ip))
}

func (vm *VM) Run() {
	for {
		op := vm.getOpcode()
		vm.Debug(op)
		vm.ip += uint16(len(op.bytes))
		op.Run(vm)
	}
}

func (vm *VM) Debug(op *Opcode) {
	if !Debug {
		return
	}
	f := vm.GetFlag
	fmt.Fprintf(os.Stderr, "%04x AX:%04x CX:%04x DX:%04x BX:%04x SP:%04x BP:%04x SI:%04x DI:%04x O%dD%dI%dT%dS%dZ%dA%dP%dC%d %-30s %s\n",
		vm.ip,
		vm.reg["ax"], vm.reg["cx"], vm.reg["dx"], vm.reg["bx"],
		vm.reg["sp"], vm.reg["bp"], vm.reg["si"], vm.reg["di"],
		f(OF), f(DF), f(IF), f(TF), f(SF), f(ZF), f(AF), f(PF), f(CF),
		op.Disasm(),
		vm.DebugStack(),
	)
}

func (vm *VM) stackSlice() (s []uint16) {
	top := vm.reg["sp"]
	for {
		if top < vm.reg["sp"] {
			return
		}
		s = append(s, vm.SS(top).read16())
		top += 2
	}
}

func (vm *VM) DebugStack() (s string) {
	for _, v := range vm.stackSlice() {
		s = fmt.Sprintf("%04x ", v) + s
	}
	return
}
