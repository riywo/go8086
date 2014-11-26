package go8086

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func CallMINIXSyscall(vm *VM) {
	m := MinixMessage(vm.SS(vm.reg["bx"]))
	syscallType := MINIXSyscall(m.Get(m_type))
	f := minixSyscallFuncMap[syscallType]
	if f == nil {
		TraceLog(syscallType)("called   message: %02x", m[0:24])
		ErrorLog("Not implemented syscall: %d", syscallType)
		os.Exit(1)
	} else {
		TraceLog(syscallType)("called   message: %02x", m[0:24])
		result, err := f(vm, m, TraceLog(syscallType))
		if err != nil {
			result = -1
			TraceLog(syscallType)("error: %v", err)
		}
		m.Set(m_type, int32(result))
		vm.reg["ax"] = uint16(result)
		TraceLog(syscallType)("finished message: %02x result: %d", m[0:24], result)
	}
}

type MINIXSyscall int16

const (
	MINIX_exit      MINIXSyscall = 1
	MINIX_fork      MINIXSyscall = 2
	MINIX_read      MINIXSyscall = 3
	MINIX_write     MINIXSyscall = 4
	MINIX_open      MINIXSyscall = 5
	MINIX_close     MINIXSyscall = 6
	MINIX_wait      MINIXSyscall = 7
	MINIX_creat     MINIXSyscall = 8
	MINIX_unlink    MINIXSyscall = 10
	MINIX_time      MINIXSyscall = 13
	MINIX_chmod     MINIXSyscall = 15
	MINIX_brk       MINIXSyscall = 17
	MINIX_stat      MINIXSyscall = 18
	MINIX_lseek     MINIXSyscall = 19
	MINIX_getpid    MINIXSyscall = 20
	MINIX_getuid    MINIXSyscall = 24
	MINIX_fstat     MINIXSyscall = 28
	MINIX_access    MINIXSyscall = 33
	MINIX_pipe      MINIXSyscall = 42
	MINIX_getgid    MINIXSyscall = 47
	MINIX_signal    MINIXSyscall = 48
	MINIX_ioctl     MINIXSyscall = 54
	MINIX_fcntl     MINIXSyscall = 55
	MINIX_exec      MINIXSyscall = 59
	MINIX_sigaction MINIXSyscall = 71
)

var minixSyscallString = map[MINIXSyscall]string{
	MINIX_exit:      "exit",
	MINIX_fork:      "fork",
	MINIX_read:      "read",
	MINIX_write:     "write",
	MINIX_open:      "open",
	MINIX_close:     "close",
	MINIX_wait:      "wait",
	MINIX_creat:     "creat",
	MINIX_unlink:    "unlink",
	MINIX_time:      "time",
	MINIX_chmod:     "chmod",
	MINIX_brk:       "brk",
	MINIX_stat:      "stat",
	MINIX_lseek:     "lseek",
	MINIX_getpid:    "getpid",
	MINIX_getuid:    "getuid",
	MINIX_fstat:     "fstat",
	MINIX_access:    "access",
	MINIX_pipe:      "pipe",
	MINIX_getgid:    "getgid",
	MINIX_signal:    "signal",
	MINIX_ioctl:     "ioctl",
	MINIX_fcntl:     "fcntl",
	MINIX_exec:      "exec",
	MINIX_sigaction: "sigaction",
}

func (s MINIXSyscall) String() string {
	return minixSyscallString[s]
}

type MINIXSyscallFunc func(*VM, MinixMessage, TraceLogger) (int, error)

var minixSyscallFuncMap = map[MINIXSyscall]MINIXSyscallFunc{
	MINIX_exit: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		status := m.Get(m1_i1)
		logger("status: %d", status)
		syscall.Exit(int(status))
		return
	},
	MINIX_fork: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		ret, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
		pid := syscall.Getpid()
		if pid != int(ret) {
			result = (int(ret) << 4) % 30000
		} else {
			result = 0
		}
		return
	},
	MINIX_read: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		fd := m.Get(m1_i1)
		nbytes := m.Get(m1_i2)
		buffer := uint16(m.Get(m1_p1))
		data := vm.DS(buffer)[0:nbytes]
		result, err = syscall.Read(int(fd), data)
		logger("fd: %d nbytes: %d buffer: %04x data: %s", fd, nbytes, buffer, strconv.Quote(string(data[0:100])))
		return
	},
	MINIX_write: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		fd := m.Get(m1_i1)
		nbytes := m.Get(m1_i2)
		buffer := uint16(m.Get(m1_p1))
		data := vm.DS(buffer)[0:nbytes]
		result, err = syscall.Write(int(fd), data)
		logger("fd: %d nbytes: %d buffer: %04x data: %s", fd, nbytes, buffer, strconv.Quote(string(data)))
		return
	},
	MINIX_open: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		names := ""
		flags := m.Get(m1_i2)
		if flags&syscall.O_CREAT != 0 {
			os.Exit(1)
		} else {
			names = WithMinixPathPrefix(m.Get_m3_name(vm))
		}

		result, err = syscall.Open(names, int(flags), 0)
		logger("flags: %d names: %s", flags, names)
		return
	},
	MINIX_close: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		fd := m.Get(m1_i1)
		result = 0
		logger("fd: %d", fd)
		err = syscall.Close(int(fd))
		return
	},
	MINIX_wait: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		status := syscall.WaitStatus(0)
		result, err = syscall.Wait4(-1, &status, 0, nil)
		result = (result << 4) % 30000
		m.Set(m2_i1, int32(status))
		logger("status: %d", status)
		return
	},
	MINIX_creat: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		mode := m.Get(m3_i2)
		names := WithMinixPathPrefix(m.Get_m3_name(vm))
		result, err = syscall.Open(names, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, uint32(mode))
		logger("mode: %d names: %s", mode, names)
		return
	},
	MINIX_unlink: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		names := WithMinixPathPrefix(m.Get_m3_name(vm))
		err = syscall.Unlink(names)
		logger("names: %s", names)
		return
	},
	MINIX_time: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		time_t := time.Now().Unix()
		m.Set(m2_l1, int32(time_t))
		logger("time_t : %d", time_t)
		return
	},
	MINIX_chmod: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		mode := m.Get(m3_i2)
		names := WithMinixPathPrefix(m.Get_m3_name(vm))
		err = syscall.Chmod(names, uint32(mode))
		logger("mode: %d names: %s", mode, names)
		return
	},
	MINIX_brk: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		nd := m.Get(m1_p1)
		if nd > 0x10000 || uint16(nd) >= vm.reg["sp"] {
			result = -1
		} else {
			m.Set(m2_p1, nd)
		}
		logger("nd: %04x", nd)
		return
	},
	MINIX_stat: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		bytes := m.Get(m1_i1)
		name := m.Get(m1_p1)
		buf := m.Get(m1_p2)
		names := WithMinixPathPrefix(string(vm.SS(uint16(name))[0 : bytes-1]))
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
		logger("names: %s stat: %+v", names, stat)
		return
	},
	MINIX_lseek: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		fd := m.Get(m2_i1)
		offset := m.Get(m2_l1)
		whence := m.Get(m2_i2)
		new_offset, err := syscall.Seek(int(fd), int64(offset), int(whence))
		m.Set(m2_l1, int32(new_offset))
		logger("fd: %d, offset: %d, whence: %d, new_offset: %d", fd, offset, whence, new_offset)
		return
	},
	MINIX_getpid: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		pid := syscall.Getpid()
		result = (pid << 4) % 30000
		return
	},
	MINIX_getuid: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		result = syscall.Getuid()
		return
	},
	MINIX_fstat: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		fd := m.Get(m1_i1)
		buf := m.Get(m1_p1)
		stat := syscall.Stat_t{}
		err = syscall.Fstat(int(fd), &stat)
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
		logger("fd: %d stat: %+v", fd, stat)
		return
	},
	MINIX_access: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		mode := m.Get(m3_i2)
		names := WithMinixPathPrefix(m.Get_m3_name(vm))
		err = syscall.Access(names, uint32(mode))
		logger("mode: %d names: %s", mode, names)
		return
	},
	MINIX_pipe: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		fields := []int{0, 0}
		err = syscall.Pipe(fields)
		m.Set(m1_i1, int32(fields[0]))
		m.Set(m1_i2, int32(fields[1]))
		logger("fields: %d", fields)
		return
	},
	MINIX_getgid: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		result = syscall.Getgid()
		return
	},
	MINIX_signal: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		return
	},
	MINIX_ioctl: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		result = -1
		return
	},
	MINIX_fcntl: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		result = -1
		return
	},
	MINIX_exec: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		bytes := m.Get(m1_i1)
		name := m.Get(m1_p1)
		names := WithMinixPathPrefix(string(vm.SS(uint16(name))[0 : bytes-1]))
		frame_size := m.Get(m1_i2)
		frame := vm.SS(uint16(m.Get(m1_p2)))[0:frame_size]

		argc := frame[0:].read16()
		args := []string{}
		for i := 0; i < int(argc); i++ {
			ptr := frame[i*2+2:].read16()
			for j, c := range frame[ptr:] {
				if c == 0x0 {
					args = append(args, string(frame[ptr:int(ptr)+j]))
					break
				}
			}
		}

		envc := (frame[2:].read16() - (argc+3)*2) / 2
		envs := []string{}
		for i := 0; i < int(envc); i++ {
			ptr := frame[i*2+2+int(argc)*2+2:].read16()
			for j, c := range frame[ptr:] {
				if c == 0x0 {
					envs = append(envs, string(frame[ptr:int(ptr)+j]))
					break
				}
			}
		}

		logger("names: %s args: %v envs: %v", names, args, envs)
		aout := NewMinixAout(names)
		aout.InitVM(vm, args, envs)
		return
	},
	MINIX_sigaction: func(vm *VM, m MinixMessage, logger TraceLogger) (result int, err error) {
		return
	},
}

var MinixPathPrefix = ""

func WithMinixPathPrefix(path string) string {
	return filepath.Join(MinixPathPrefix, path)
}

type MinixAout struct {
	a_hdrlen uint8
	a_text   int32
	a_data   int32
	a_bss    int32
	a_entry  int32
	text     Bytes
	data     Bytes
}

func NewMinixAout(file string) (aout *MinixAout) {
	aout = new(MinixAout)
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	aout.a_hdrlen = uint8(Bytes(bs)[4])
	aout.a_text = int32(Bytes(bs)[8:].read32())
	aout.a_data = int32(Bytes(bs)[12:].read32())
	aout.a_bss = int32(Bytes(bs)[16:].read32())
	aout.a_entry = int32(Bytes(bs)[20:].read32())
	aout.text = Bytes(bs)[int32(aout.a_hdrlen) : int32(aout.a_hdrlen)+aout.a_text]
	aout.data = Bytes(bs)[int32(aout.a_hdrlen)+aout.a_text : int32(aout.a_hdrlen)+aout.a_text+aout.a_data]
	return
}

func (aout *MinixAout) NewVM(args, env []string) (vm *VM) {
	vm = NewVM()
	aout.InitVM(vm, args, env)
	return
}

func (aout *MinixAout) InitVM(vm *VM, args, envs []string) {
	vm.Init()
	vm.ip = uint16(aout.a_entry)
	vm.CS(0x0).write(aout.text)
	vm.DS(0x0).write(aout.data)
	aout.StackArgsEnv(vm, args, envs)
	vm.initSP = vm.reg["sp"]
	DebugLog("%02x", aout.data[0:100])
}

func (aout *MinixAout) StackArgsEnv(vm *VM, args, envs []string) {
	sp := vm.reg["sp"]

	chars := Bytes{}
	arg_ptrs := []uint16{}
	env_ptrs := []uint16{}
	p := uint16(0)
	for _, arg := range args {
		chars = append(chars, Bytes(arg)...)
		chars = append(chars, 0x0)
		arg_ptrs = append(arg_ptrs, p)
		p += uint16(len(arg) + 1)
	}
	for _, env := range envs {
		chars = append(chars, Bytes(env)...)
		chars = append(chars, 0x0)
		env_ptrs = append(env_ptrs, p)
		p += uint16(len(env) + 1)
	}
	if len(chars)%2 != 0 {
		chars = append(chars, 0x0)
	}

	stack_len := 2 + 2*len(args) + 2 + 2*len(envs) + 2 + len(chars)
	stack := make(Bytes, stack_len)
	top := sp - uint16(stack_len)

	stack.write16(uint16(len(args)))

	for i, _ := range args {
		ptr := arg_ptrs[i] + 2 + 2*uint16(len(args)+len(envs)+2) + top
		stack[2+2*i:].write16(ptr)
	}
	stack = append(stack, 0x0)

	for i, _ := range envs {
		ptr := env_ptrs[i] + 2 + 2*uint16(len(args)+len(envs)+2) + top
		stack[2+2*len(args)+2+2*i:].write16(ptr)
	}
	stack = append(stack, 0x0)

	stack[2+2*len(args)+2+2*len(envs)+2:].write(chars)

	vm.SS(top).write(stack)
	vm.reg["sp"] -= uint16(stack_len)
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

func (m MinixMessage) Get_m3_name(vm *VM) (names string) {
	k := m.Get(m3_i1)
	if k <= 14 {
		names = string(m.Get_m3_ca1()[0 : k-1])
	} else {
		names = string(vm.DS(uint16(m.Get(m3_p1)))[0 : k-1])
	}
	return
}
