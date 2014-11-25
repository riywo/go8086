package go8086

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func CallMINIXSyscall(vm *VM) {
	m := MinixMessage(vm.SS(vm.reg["bx"]))
	syscallType := MINIXSyscall(m.Get(m_type))
	DebugLog("syscall: %d, message: %02x\n", syscallType, m[0:24])
	f := minixSyscallFuncMap[syscallType]
	if f == nil {
		fmt.Fprintf(os.Stderr, "Not implemented syscall: %d\n", syscallType)
		os.Exit(1)
	} else {
		result, err := f(vm, m)
		if err != nil {
			result = -1
		}
		m.Set(m_type, int32(result))
		vm.reg["ax"] = uint16(result)
	}
	DebugLog("syscall: %d, message: %02x\n", syscallType, m[0:24])
}

type MINIXSyscall int16

const (
	MINIX_exit      MINIXSyscall = 1
	MINIX_read      MINIXSyscall = 3
	MINIX_write     MINIXSyscall = 4
	MINIX_open      MINIXSyscall = 5
	MINIX_close     MINIXSyscall = 6
	MINIX_creat     MINIXSyscall = 8
	MINIX_unlink    MINIXSyscall = 10
	MINIX_time      MINIXSyscall = 13
	MINIX_chmod     MINIXSyscall = 15
	MINIX_brk       MINIXSyscall = 17
	MINIX_stat      MINIXSyscall = 18
	MINIX_lseek     MINIXSyscall = 19
	MINIX_getpid    MINIXSyscall = 20
	MINIX_getuid    MINIXSyscall = 24
	MINIX_access    MINIXSyscall = 33
	MINIX_getgid    MINIXSyscall = 47
	MINIX_ioctl     MINIXSyscall = 54
	MINIX_sigaction MINIXSyscall = 71
)

type MINIXSyscallFunc func(*VM, MinixMessage) (int, error)

var minixSyscallFuncMap = map[MINIXSyscall]MINIXSyscallFunc{
	MINIX_exit: func(vm *VM, m MinixMessage) (result int, err error) {
		status := m.Get(m1_i1)
		syscall.Exit(int(status))
		return
	},
	MINIX_read: func(vm *VM, m MinixMessage) (result int, err error) {
		fd := m.Get(m1_i1)
		nbytes := m.Get(m1_i2)
		buffer := uint16(m.Get(m1_p1))
		data := vm.DS(buffer)[0:nbytes]
		result, err = syscall.Read(int(fd), data)
		DebugLog("read : %04x %04x %04x data: %02x result: %04x", fd, nbytes, buffer, data[0:10], result)
		return
	},
	MINIX_write: func(vm *VM, m MinixMessage) (result int, err error) {
		fd := m.Get(m1_i1)
		nbytes := m.Get(m1_i2)
		buffer := uint16(m.Get(m1_p1))
		data := vm.DS(buffer)[0:nbytes]
		result, err = syscall.Write(int(fd), data)
		DebugLog("write: %04x %04x %04x data: %02x result: %04x", fd, nbytes, buffer, data[0:10], result)
		return
	},
	MINIX_open: func(vm *VM, m MinixMessage) (result int, err error) {
		names := ""
		flags := m.Get(m1_i2)
		if flags&syscall.O_CREAT != 0 {
			os.Exit(1)
		} else {
			names = m.Get_m3_name(vm)
		}

		result, err = syscall.Open(names, int(flags), 0)
		DebugLog("open : flags: %d names: %s result: %04x", flags, names, result)
		return
	},
	MINIX_close: func(vm *VM, m MinixMessage) (result int, err error) {
		fd := m.Get(m1_i1)
		result = 0
		err = syscall.Close(int(fd))
		DebugLog("close: fd: %04x result: %04x", fd, result)
		return
	},
	MINIX_creat: func(vm *VM, m MinixMessage) (result int, err error) {
		mode := m.Get(m3_i2)
		names := m.Get_m3_name(vm)
		result, err = syscall.Open(names, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, uint32(mode))
		DebugLog("creat : mode: %d names: %s result: %04x err: %v", mode, names, result, err)
		return
	},
	MINIX_unlink: func(vm *VM, m MinixMessage) (result int, err error) {
		names := m.Get_m3_name(vm)
		err = syscall.Unlink(names)
		DebugLog("unlink : %s, %v", names, err)
		return
	},
	MINIX_time: func(vm *VM, m MinixMessage) (result int, err error) {
		time_t := time.Now().Unix()
		DebugLog("time : %x", time_t)
		m.Set(m2_l1, int32(time_t))
		return
	},
	MINIX_chmod: func(vm *VM, m MinixMessage) (result int, err error) {
		mode := m.Get(m3_i2)
		names := m.Get_m3_name(vm)
		err = syscall.Chmod(names, uint32(mode))
		DebugLog("chmod : mode: %d names: %s result: %04x", mode, names, result)
		return
	},
	MINIX_brk: func(vm *VM, m MinixMessage) (result int, err error) {
		nd := m.Get(m1_p1)
		if nd > 0x10000 || uint16(nd) >= vm.reg["sp"] {
			result = -1
		} else {
			m.Set(m2_p1, nd)
		}
		DebugLog("brk  : nd: %04x result: %04x", nd, result)
		return
	},
	MINIX_stat: func(vm *VM, m MinixMessage) (result int, err error) {
		bytes := m.Get(m1_i1)
		name := m.Get(m1_p1)
		buf := m.Get(m1_p2)
		names := string(vm.SS(uint16(name))[0 : bytes-1])
		stat := syscall.Stat_t{}
		err = syscall.Stat(names, &stat)
		vm.SS(uint16(buf))[0:].write16(uint16(stat.Dev))
		vm.SS(uint16(buf))[2:].write16(uint16(stat.Ino))
		vm.SS(uint16(buf))[4:].write16(uint16(stat.Mode))
		vm.SS(uint16(buf))[6:].write16(uint16(stat.Nlink))
		vm.SS(uint16(buf))[8:].write16(uint16(stat.Uid))
		vm.SS(uint16(buf))[10:].write16(uint16(stat.Gid))
		vm.SS(uint16(buf))[12:].write16(uint16(stat.Rdev))
		vm.SS(uint16(buf))[14:].write32(uint32(stat.Size))
		vm.SS(uint16(buf))[18:].write32(uint32(stat.Atimespec.Sec))
		vm.SS(uint16(buf))[22:].write32(uint32(stat.Mtimespec.Sec))
		vm.SS(uint16(buf))[26:].write32(uint32(stat.Ctimespec.Sec))
		DebugLog("stat : names: %s stat: %+v err: %v", names, stat, err)
		return
	},
	MINIX_lseek: func(vm *VM, m MinixMessage) (result int, err error) {
		fd := m.Get(m2_i1)
		offset := m.Get(m2_l1)
		whence := m.Get(m2_i2)
		new_offset, err := syscall.Seek(int(fd), int64(offset), int(whence))
		m.Set(m2_l1, int32(new_offset))
		DebugLog("lseek: fd: %d, offset: %d, whence: %d, new_offset: %08x", fd, offset, whence, new_offset)
		return
	},
	MINIX_getpid: func(vm *VM, m MinixMessage) (result int, err error) {
		pid := syscall.Getpid()
		result = (pid << 4) % 30000
		DebugLog("getpid : result: %04x", result)
		return
	},
	MINIX_getuid: func(vm *VM, m MinixMessage) (result int, err error) {
		result = syscall.Getuid()
		DebugLog("getuid : result: %04x", result)
		return
	},
	MINIX_access: func(vm *VM, m MinixMessage) (result int, err error) {
		mode := m.Get(m3_i2)
		names := m.Get_m3_name(vm)
		err = syscall.Access(names, uint32(mode))
		DebugLog("access : mode: %d names: %s result: %04x", mode, names, result)
		return
	},
	MINIX_getgid: func(vm *VM, m MinixMessage) (result int, err error) {
		result = syscall.Getgid()
		DebugLog("getgid : result: %04x", result)
		return
	},
	MINIX_ioctl: func(vm *VM, m MinixMessage) (result int, err error) {
		result = -1
		return
	},
	MINIX_sigaction: func(vm *VM, m MinixMessage) (result int, err error) {
		return
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
		p += uint16(len(arg) + 1)
	}
	if len(chars)%2 != 0 {
		chars = append(chars, 0x0)
	}

	stack_len := 2 + 2*len(args) + 2 + len(chars)
	stack := make(Bytes, stack_len)
	top := sp - uint16(stack_len)
	stack.write16(uint16(len(args)))
	for i, _ := range args {
		ptr := ptrs[i] + 2 + 2*uint16(len(args)) + 2 + top
		stack[2+2*i:].write16(ptr)
	}
	stack = append(stack, 0x0)
	stack[2+2*len(args)+2:].write(chars)

	vm.SS(top).write(stack)
	vm.reg["sp"] -= uint16(stack_len)
	DebugLog("SP: %04x", vm.reg["sp"])
	DebugLog("Stack: %d", vm.SS(vm.reg["sp"])[0:stack_len])
}

type MinixMessage Bytes

type MinixMessageAccessor int

const (
	m_source MinixMessageAccessor = iota
	m_type
	m1_i1
	m1_i2
	m1_i3
	m1_p1
	m1_p2
	m1_p3
	m2_i1
	m2_i2
	m2_i3
	m2_l1
	m2_l2
	m2_p1
	m3_i1
	m3_i2
	m3_p1
	m3_ca1
	m4_l1
	m4_l2
	m4_l3
	m4_l4
	m4_l5
	m5_c1
	m5_c2
	m5_i1
	m5_i2
	m5_l1
	m5_l2
	m5_l3
	m6_i1
	m6_i2
	m6_i3
	m6_l1
	m6_f1
)

func (m MinixMessage) Get(accessor MinixMessageAccessor) (v int32) {
	switch accessor {
	case m_source:
		v = int32(Bytes(m[0:]).read16())
	case m_type:
		v = int32(Bytes(m[2:]).read16())
	case m1_i1, m2_i1, m3_i1, m6_i1:
		v = int32(Bytes(m[4:]).read16())
	case m1_i2, m2_i2, m3_i2, m6_i2:
		v = int32(Bytes(m[6:]).read16())
	case m1_i3, m2_i3, m6_i3:
		v = int32(Bytes(m[8:]).read16())
	case m1_p1:
		v = int32(Bytes(m[10:]).read16())
	case m1_p2:
		v = int32(Bytes(m[12:]).read16())
	case m1_p3:
		v = int32(Bytes(m[14:]).read16())
	case m2_l1:
		v = int32(Bytes(m[10:]).read32())
	case m2_l2:
		v = int32(Bytes(m[14:]).read32())
	case m2_p1:
		v = int32(Bytes(m[18:]).read16())
	case m3_p1:
		v = int32(Bytes(m[8:]).read16())
	case m4_l1:
		v = int32(Bytes(m[4:]).read32())
	case m4_l2:
		v = int32(Bytes(m[8:]).read32())
	case m4_l3:
		v = int32(Bytes(m[12:]).read32())
	case m4_l4:
		v = int32(Bytes(m[16:]).read32())
	case m4_l5:
		v = int32(Bytes(m[20:]).read32())
	case m5_c1:
		v = int32(Bytes(m[4:]).read8())
	case m5_c2:
		v = int32(Bytes(m[5:]).read8())
	case m5_i1:
		v = int32(Bytes(m[6:]).read16())
	case m5_i2:
		v = int32(Bytes(m[8:]).read16())
	case m5_l1:
		v = int32(Bytes(m[10:]).read32())
	case m5_l2:
		v = int32(Bytes(m[14:]).read32())
	case m5_l3:
		v = int32(Bytes(m[18:]).read32())
	case m6_l1:
		v = int32(Bytes(m[10:]).read32())
	case m6_f1:
		v = int32(Bytes(m[14:]).read16())
	}
	return
}

func (m MinixMessage) Set(accessor MinixMessageAccessor, v int32) {
	switch accessor {
	case m_source:
		Bytes(m[0:]).write16(uint16(v))
	case m_type:
		Bytes(m[2:]).write16(uint16(v))
	case m1_i1, m2_i1, m3_i1, m6_i1:
		Bytes(m[4:]).write16(uint16(v))
	case m1_i2, m2_i2, m3_i2, m6_i2:
		Bytes(m[6:]).write16(uint16(v))
	case m1_i3, m2_i3, m6_i3:
		Bytes(m[8:]).write16(uint16(v))
	case m1_p1:
		Bytes(m[10:]).write16(uint16(v))
	case m1_p2:
		Bytes(m[12:]).write16(uint16(v))
	case m1_p3:
		Bytes(m[14:]).write16(uint16(v))
	case m2_l1:
		Bytes(m[10:]).write32(uint32(v))
	case m2_l2:
		Bytes(m[14:]).write32(uint32(v))
	case m2_p1:
		Bytes(m[18:]).write16(uint16(v))
	case m3_p1:
		Bytes(m[8:]).write16(uint16(v))
	case m4_l1:
		Bytes(m[4:]).write32(uint32(v))
	case m4_l2:
		Bytes(m[8:]).write32(uint32(v))
	case m4_l3:
		Bytes(m[12:]).write32(uint32(v))
	case m4_l4:
		Bytes(m[16:]).write32(uint32(v))
	case m4_l5:
		Bytes(m[20:]).write32(uint32(v))
	case m5_c1:
		Bytes(m[4:]).write8(uint16(v))
	case m5_c2:
		Bytes(m[5:]).write8(uint16(v))
	case m5_i1:
		Bytes(m[6:]).write16(uint16(v))
	case m5_i2:
		Bytes(m[8:]).write16(uint16(v))
	case m5_l1:
		Bytes(m[10:]).write32(uint32(v))
	case m5_l2:
		Bytes(m[14:]).write32(uint32(v))
	case m5_l3:
		Bytes(m[18:]).write32(uint32(v))
	case m6_l1:
		Bytes(m[10:]).write32(uint32(v))
	case m6_f1:
		Bytes(m[14:]).write16(uint16(v))
	}
}

func (m MinixMessage) Get_m3_ca1() Bytes {
	return Bytes(m)[10:24]
}

func (m MinixMessage) Get_m3_name(vm *VM) string {
	k := m.Get(m3_i1)
	if k <= 14 {
		return string(m.Get_m3_ca1()[0 : k-1])
	} else {
		return string(vm.DS(uint16(m.Get(m3_p1)))[0 : k-1])
	}
}
