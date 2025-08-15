package runtime

const (
	SYS_IO_FLUSH = 1
)

func doSyscall(vm *VM, instr VMInstr) {
	switch instr.Oprand1.IntData {
	case SYS_IO_FLUSH:
		stdout := vm.Mem.GetObj("stdout")
		vm.IO.WriteObjectToStream(*stdout)
		vm.Mem.SetObj("stdout", VMDataObject{})
	}
}
