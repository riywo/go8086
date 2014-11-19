package go8086

type Bit int

const (
	Bit8 Bit = iota
	Bit16
)

type Count int

const (
	Count1 Count = iota
	CountCL
)

type Signed bool

const (
	Sign   Signed = true
	Unsign Signed = false
)

type Direction int

const (
	FromReg Direction = iota
	ToReg
)

func Run(bs []byte, args []string) {
	vm := NewVM()
	aout := NewMinixAout(Bytes(bs))
	vm.memCS.write(aout.text)
	vm.memDS.write(aout.data)
	MinixStackArgs(vm, args)
	vm.Run()
}
