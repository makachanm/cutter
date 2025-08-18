package runtime

const (
	SYS_IO_FLUSH        = 1
	SYS_SET_FUNC_RETURN = 2 // New syscall for setting function return value
)

func doSyscall(vm *VM, instr VMInstr) {
	switch instr.Oprand1.IntData {
	case SYS_IO_FLUSH:
		stdout := vm.Mem.GetObj("stdout")
		vm.IO.WriteObjectToStream(*stdout)
		vm.Mem.SetObj("stdout", VMDataObject{})
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
	}
}