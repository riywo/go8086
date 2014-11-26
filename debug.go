package go8086

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

var Debug bool = false
var Trace bool = false

func Pid() int {
	return (os.Getpid() << 4) % 30000
}

func DebugLog(format string, a ...interface{}) {
	if Debug {
		log := fmt.Sprintf("%d [Debug] ", Pid())
		log += fmt.Sprintf(format, a...)
		fmt.Fprintln(os.Stderr, log)
	}
}

type TraceLogger func(format string, a ...interface{})

func TraceLog(syscallType MINIXSyscall) TraceLogger {
	return func(format string, a ...interface{}) {
		if Trace {
			log := fmt.Sprintf("%d [Trace] %-10s ", Pid(), syscallType)
			log += fmt.Sprintf(format, a...)
			fmt.Fprintln(os.Stderr, log)
		}
	}
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
	fmt.Fprintf(os.Stderr, "%d %04x AX:%s CX:%s DX:%s BX:%s SP:%s BP:%s SI:%s DI:%s %s%s%s%s%s%s%s%s%s %-30s %s\n",
		Pid(),
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
		if top < vm.reg["sp"] {
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
		if p == vm.reg["bx"] {
			str = bxString(v)
		}
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
