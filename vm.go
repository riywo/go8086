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
	reg    map[string]uint16
	sreg   map[string]uint16
	ip     uint16
	flag   uint16
	memCS  Bytes  //temporary
	memDS  Bytes  //temporary
	initSP uint16 //temporary
}

func NewVM() (vm *VM) {
	vm = new(VM)
	vm.reg = make(map[string]uint16)
	vm.sreg = make(map[string]uint16)
	vm.memCS = make(Bytes, 0x10000) //temporary
	vm.memDS = make(Bytes, 0x10000) //temporary
	for _, reg := range regs[Bit16] {
		switch reg {
		case SP:
			vm.reg[reg.name] = 0xff0e
		default:
			vm.reg[reg.name] = 0
		}
	}
	for _, sreg := range sregs {
		vm.sreg[sreg.name] = 0
	}
	return
}

func (vm *VM) CS(offset uint16) Bytes {
	return vm.memCS[offset:] //temporary
}

func (vm *VM) DS(offset uint16) Bytes {
	return vm.memDS[offset:] //temporary
}

func (vm *VM) Push(value uint16) {
	vm.reg["sp"] -= 2
	vm.DS(vm.reg["sp"]).write16(value)
}

func (vm *VM) Pop() (value uint16) {
	value = vm.DS(vm.reg["sp"]).read16()
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
	return getOpcode(nil, vm.ip, vm.memCS[vm.ip:])
}

func (vm *VM) Run() {
	vm.initSP = vm.reg["sp"]
	for {
		op := vm.getOpcode()
		vm.Debug()
		vm.ip += uint16(len(op.bytes))
		op.Run(vm)
	}
}

func (vm *VM) Debug() {
	if !Debug {
		return
	}
	f := vm.GetFlag
	fmt.Fprintf(os.Stderr, "%04x AX:%04x CX:%04x DX:%04x BX:%04x SP:%04x BP:%04x SI:%04x DI:%04x O%dD%dI%dT%dS%dZ%dA%dP%dC%d %-30s %s\n",
		vm.ip,
		vm.reg["ax"], vm.reg["cx"], vm.reg["dx"], vm.reg["bx"],
		vm.reg["sp"], vm.reg["bp"], vm.reg["si"], vm.reg["di"],
		f(OF), f(DF), f(IF), f(TF), f(SF), f(ZF), f(AF), f(PF), f(CF),
		vm.getOpcode().Disasm(),
		vm.DebugStack(),
	)
}

func (vm *VM) stackSlice() (s []uint16) {
	top := vm.DS(vm.reg["sp"])
	for {
		if len(top) < 0xffff-int(vm.initSP)+2 {
			return
		}
		s = append(s, top.read16())
		top = top[2:]
	}
}

func (vm *VM) DebugStack() (s string) {
	for _, v := range vm.stackSlice() {
		s = fmt.Sprintf("%04x ", v) + s
	}
	return
}
