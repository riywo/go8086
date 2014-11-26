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

func Run(file string, args, env []string) {
	aout := NewMinixAout(file)
	vm := aout.NewVM(args, env)
	vm.Run()
}

func ErrorLog(format string, a ...interface{}) {
	log := fmt.Sprintf("%d [Error] ", Pid())
	log += fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, log)
}
