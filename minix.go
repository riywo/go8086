package go8086

import (
	"fmt"
	"os"
	"syscall"
)

func CallMINIXSyscall(vm *VM) {
	syscallType := MINIXSyscall(vm.DS(BX.Read(vm) + 2).read16())
	f := minixSyscallFuncMap[syscallType]
	if f == nil {
		fmt.Printf("Not implemented syscall: %d", syscallType)
		os.Exit(1)
	} else {
		f(vm)
	}
}

type MINIXSyscall uint16

const (
	MINIX_write MINIXSyscall = 4
)

type MINIXSyscallFunc func(*VM)

var minixSyscallFuncMap = map[MINIXSyscall]MINIXSyscallFunc{
	MINIX_write: func(vm *VM) {
		fd := int(vm.DS(BX.Read(vm) + 4).read16())
		offset := vm.DS(BX.Read(vm) + 10).read16()
		count := vm.DS(BX.Read(vm) + 6).read16()
		data := vm.DS(offset)[0:count]
		//fmt.Printf("%d %04x %04x %v\n", fd, offset, count, data)
		syscall.Write(fd, data)
		//os.Exit(0)
	},
}

type MinixAout struct {
	a_hdrlen uint8
	a_text   int32
	a_data   int32
	text     Bytes
	data     Bytes
}

func NewMinixAout(bs Bytes) (aout *MinixAout) {
	aout = new(MinixAout)
	aout.a_hdrlen = uint8(bs[4])
	aout.a_text = int32(bs[8:].read32())
	aout.a_data = int32(bs[12:].read32())
	aout.text = bs[int32(aout.a_hdrlen) : int32(aout.a_hdrlen)+aout.a_text]
	aout.data = bs[int32(aout.a_hdrlen)+aout.a_text : int32(aout.a_hdrlen)+aout.a_text+aout.a_data]
	return
}

func MinixStackArgs(vm *VM, args []string) {
	sp := vm.reg["sp"]

	chars := Bytes{}
	ptrs := []uint16{}
	p := uint16(0)
	for _, arg := range args {
		chars = append(chars, Bytes(arg)...)
		chars = append(chars, 0x0)
		ptrs = append(ptrs, p)
		p += uint16(len(chars))
	}
	if len(chars)%2 != 0 {
		chars = append(chars, 0x0)
	}

	stack_len := 2 + 2*len(args) + 2 + len(chars)
	stack := make(Bytes, stack_len)
	top := sp - uint16(stack_len) + 2
	stack.write16(uint16(len(args)))
	for i, _ := range args {
		ptr := ptrs[i] + 2 + uint16(len(args)) + 2 + top
		stack[2+2*i:].write16(ptr)
	}
	stack = append(stack, 0x0)
	stack[2+2*len(args)+2:].write(chars)

	vm.DS(top - 2).write(stack)
	vm.reg["sp"] -= uint16(stack_len)
}
