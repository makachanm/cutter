package runtime

// GetStandardFuncs provides the standard library functions for the Cutter VM.
// Each function is implemented as a sequence of VM instructions.
// Not all functions from standardfunctions.md are implemented due to limitations in the current VM instruction set.
func GetStandardFuncs() map[string][]VMInstr {
	StandardFuncs := make(map[string][]VMInstr)

	// Use register 1023 as a dedicated temporary register for standard functions
	// to avoid clobbering argument registers (0 and 1).
	tempReg := makeIntValueObj(1023)

	// Arithmetic Functions
	StandardFuncs["add"] = []VMInstr{
		{Op: OpAdd, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["sub"] = []VMInstr{
		{Op: OpSub, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["mul"] = []VMInstr{
		{Op: OpMul, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["div"] = []VMInstr{
		{Op: OpDiv, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["mod"] = []VMInstr{
		{Op: OpMod, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}

	// Comparison Functions
	StandardFuncs["same"] = []VMInstr{
		{Op: OpCmpEq, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["notsame"] = []VMInstr{
		{Op: OpCmpNeq, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}

	// String Functions
	StandardFuncs["strcontact"] = StandardFuncs["add"]

	// Branching and Control Flow
	StandardFuncs["ifel"] = []VMInstr{
		{Op: OpBrch, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(2)},
	}

	// Memory and Variable Manipulation
	StandardFuncs["set"] = []VMInstr{
		// Arguments are expected to be in registers 0 and 1
		// Reg 0: function name (string)
		// Reg 1: value to return (VMDataObject)

		// Call the syscall to set the function's return value
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_SET_FUNC_RETURN)},
		// The syscall itself sets the result of the 'set' operation (true for success)
		// So, no need for OpRslSet here, as the syscall already handles it.
	}

	// System Functions
	StandardFuncs["exit"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)},
		{Op: OpHlt}, // Stop execution
	}

	return StandardFuncs
}
