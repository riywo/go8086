package go8086

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
	vm.Init()
	return
}

func (vm *VM) Init() {
	vm.reg = make(map[string]uint16)
	vm.sreg = make(map[string]uint16)
	vm.ip = 0
	vm.flag = 0
	vm.mem = make(Bytes, 0x100000)
	for _, reg := range regs[Bit16] {
		switch reg {
		case SP:
			vm.reg[reg.name] = 0xfffe
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
