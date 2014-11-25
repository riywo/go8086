package go8086

import (
	"fmt"
	"github.com/fatih/color"
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

var flagMap = map[Flag]string{
	OF: "O",
	DF: "D",
	IF: "I",
	TF: "T",
	SF: "S",
	ZF: "Z",
	AF: "A",
	PF: "P",
	CF: "C",
}

func (f Flag) String() string {
	return flagMap[f]
}

type VM struct {
	reg    map[string]uint16
	sreg   map[string]uint16
	ip     uint16
	flag   uint16
	mem    Bytes
	initSP uint16 //temporary
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
	vm.initSP = vm.reg["sp"]
	for {
		op := vm.getOpcode()
		vm.Debug(op)
		vm.ip += uint16(len(op.bytes))
		op.Run(vm)
	}
}

func axString(x uint16) string {
	return color.New(color.FgRed).SprintfFunc()("%04x", x)
}
func cxString(x uint16) string {
	return color.New(color.FgGreen).SprintfFunc()("%04x", x)
}
func dxString(x uint16) string {
	return color.New(color.FgYellow).SprintfFunc()("%04x", x)
}
func bxString(x uint16) string {
	return color.New(color.FgCyan).SprintfFunc()("%04x", x)
}
func spString(x uint16) string {
	return color.New(color.FgRed).Add(color.Underline).SprintfFunc()("%04x", x)
}
func bpString(x uint16) string {
	return color.New(color.FgGreen).Add(color.Underline).SprintfFunc()("%04x", x)
}
func siString(x uint16) string {
	return color.New(color.FgYellow).Add(color.Underline).SprintfFunc()("%04x", x)
}
func diString(x uint16) string {
	return color.New(color.FgCyan).Add(color.Underline).SprintfFunc()("%04x", x)
}

func (vm *VM) Debug(op *Opcode) {
	if !Debug {
		return
	}
	f := func(fl Flag) string {
		v := vm.GetFlag(fl)
		if v == 1 {
			return color.MagentaString("%s", fl.String())
		} else {
			return fl.String()
		}
	}
	fmt.Fprintf(os.Stderr, "%04x AX:%s CX:%s DX:%s BX:%s SP:%s BP:%s SI:%s DI:%s %s%s%s%s%s%s%s%s%s %-30s %s\n",
		vm.ip,
		axString(vm.reg["ax"]),
		cxString(vm.reg["cx"]),
		dxString(vm.reg["dx"]),
		bxString(vm.reg["bx"]),
		spString(vm.reg["sp"]),
		bpString(vm.reg["bp"]),
		siString(vm.reg["si"]),
		diString(vm.reg["di"]),
		f(OF), f(DF), f(IF), f(TF), f(SF), f(ZF), f(AF), f(PF), f(CF),
		op.Disasm(),
		vm.DebugStack(),
	)
}

func (vm *VM) stackSlice() (s []uint16) {
	top := vm.reg["sp"]
	for {
		if top < vm.reg["sp"] || top >= vm.initSP {
			return
		}
		s = append(s, vm.SS(top).read16())
		top += 2
	}
}

func (vm *VM) DebugStack() (s string) {
	for i, v := range vm.stackSlice() {
		str := fmt.Sprintf("%04x", v)
		p := uint16(2*i) + vm.reg["sp"]
		if p == vm.reg["sp"] {
			str = spString(v)
		}
		if p == vm.reg["bp"] {
			str = bpString(v)
		}
		if p == vm.reg["si"] {
			str = siString(v)
		}
		if p == vm.reg["di"] {
			str = diString(v)
		}
		s = str + " " + s
	}
	return
}
