package runtime

import "fmt"

const (
	SYS_IO_FLUSH = 1
)

func doSyscall(vm *VM, instr VMInstr) {
	switch instr.Oprand1.IntData {
	case SYS_IO_FLUSH:
		stdout := vm.Mem.GetObj("stdout")
		fmt.Print(stdout.StringData)
		vm.Mem.SetObj("stdout", VMDataObject{Type: STRING, StringData: ""})
	}
}
