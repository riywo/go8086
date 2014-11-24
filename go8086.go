package go8086

import (
	"fmt"
	"os"
)

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
	vm.CS(0).write(aout.text)
	vm.DS(0).write(aout.data)
	MinixStackArgs(vm, args)
	vm.Run()
}

var Debug bool = false

func DebugLog(format string, a ...interface{}) {
	if Debug {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "[Debug] "+format, a...)
		fmt.Fprintln(os.Stderr, "")
	}
}
