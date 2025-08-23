package runtime

import (
	"regexp"
	"strings"
)

const (
	SYS_SET_FUNC_RETURN = 1 // New syscall for setting function return value
	SYS_IO_FLUSH        = 2
	SYS_STR_LEN         = 3
	SYS_STR_SUB         = 4
	SYS_STR_MATCH       = 5
	SYS_STR_REPLACE     = 6
	SYS_STR_REGEXP      = 7
)

func doSyscall(vm *VM, instr VMInstr) {
	switch instr.Oprand1.IntData {
	case SYS_SET_FUNC_RETURN:
		// Expect function name in register 0 and return value in register 1
		funcNameObj := vm.Reg.GetRegister(0)
		returnValue := vm.Reg.GetRegister(1)

		if funcNameObj.Type != STRING {
			panic("SYS_SET_FUNC_RETURN: First argument must be a string (function name)")
		}

		funcName := funcNameObj.StringData

		// Get the target function object
		targetFunc := vm.Mem.GetFunc(funcName)

		// Create new instructions for the target function
		// These instructions will simply set the result register to the desired return value
		newInstructions := []VMInstr{
			{Op: OpDefFunc, Oprand1: makeStrValueObj(funcName)},               // Define the function entry point
			{Op: OpRegSet, Oprand1: makeIntValueObj(0), Oprand2: returnValue}, // Set Reg 0 to the return value
			{Op: OpRslSet, Oprand1: makeIntValueObj(0)},                       // Set result to Reg 0
			{Op: OpReturn}, // Return from the function
		}

		targetPos := targetFunc.JumpPc - 1
		for i, instr := range newInstructions {
			vm.Program[targetPos+i] = instr
		}

		// Set the result of the syscall itself (e.g., true for success)
		vm.Reg.InsertResult(VMDataObject{Type: BOOLEAN, BoolData: true})

	case SYS_IO_FLUSH:
		stdout := vm.Mem.GetObj("stdout")
		vm.IO.WriteObjectToStream(*stdout)
		vm.Mem.SetObj("stdout", VMDataObject{})
	case SYS_STR_LEN:
		str := vm.Reg.GetRegister(0)
		if str.Type != STRING {
			panic("SYS_STR_LEN: First argument must be a string")
		}
		vm.Reg.InsertResult(VMDataObject{Type: INTGER, IntData: int64(len(str.StringData))})
	case SYS_STR_SUB:
		str := vm.Reg.GetRegister(0)
		start := vm.Reg.GetRegister(1)
		end := vm.Reg.GetRegister(2)
		if str.Type != STRING || start.Type != INTGER || end.Type != INTGER {
			panic("SYS_STR_SUB: Invalid arguments")
		}
		if start.IntData < 0 || end.IntData > int64(len(str.StringData)) || start.IntData > end.IntData {
			panic("SYS_STR_SUB: Index out of bounds")
		}
		vm.Reg.InsertResult(VMDataObject{Type: STRING, StringData: str.StringData[start.IntData:end.IntData]})
	case SYS_STR_MATCH:
		str := vm.Reg.GetRegister(0)
		substr := vm.Reg.GetRegister(1)
		if str.Type != STRING || substr.Type != STRING {
			panic("SYS_STR_MATCH: Invalid arguments")
		}
		vm.Reg.InsertResult(VMDataObject{Type: INTGER, IntData: int64(strings.Index(str.StringData, substr.StringData))})
	case SYS_STR_REPLACE:
		str := vm.Reg.GetRegister(0)
		old := vm.Reg.GetRegister(1)
		new := vm.Reg.GetRegister(2)
		if str.Type != STRING || old.Type != STRING || new.Type != STRING {
			panic("SYS_STR_REPLACE: Invalid arguments")
		}
		vm.Reg.InsertResult(VMDataObject{Type: STRING, StringData: strings.ReplaceAll(str.StringData, old.StringData, new.StringData)})
	case SYS_STR_REGEXP:
		str := vm.Reg.GetRegister(0)
		pattern := vm.Reg.GetRegister(1)
		if str.Type != STRING || pattern.Type != STRING {
			panic("SYS_STR_REGEXP: Invalid arguments")
		}
		re := regexp.MustCompile(pattern.StringData)
		matches := re.FindAllString(str.StringData, -1)
		vm.Reg.InsertResult(VMDataObject{Type: STRING, StringData: strings.Join(matches, " ")})
	}
}
