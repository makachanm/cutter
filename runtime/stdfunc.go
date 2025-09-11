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
	StandardFuncs["bigger"] = []VMInstr{
		{Op: OpCmpGt, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["smaller"] = []VMInstr{
		{Op: OpCmpLt, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["bigsame"] = []VMInstr{
		{Op: OpCmpGte, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}
	StandardFuncs["smallsame"] = []VMInstr{
		{Op: OpCmpLte, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: tempReg},
		{Op: OpRslSet, Oprand1: tempReg},
	}

	// String Functions
	StandardFuncs["strcontact"] = StandardFuncs["add"]
	StandardFuncs["strlen"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_STR_LEN)},
	}
	StandardFuncs["stridx"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_STR_MATCH)},
	}
	StandardFuncs["strsub"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_STR_SUB)},
	}
	StandardFuncs["strrep"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_STR_REPLACE)},
	}
	StandardFuncs["strexp"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_STR_REGEXP)},
	}

	// Branching and Control Flow
	StandardFuncs["ifel"] = []VMInstr{
		{Op: OpBrch, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(2)},
	}

	// Memory and Variable Manipulation
	StandardFuncs["set"] = []VMInstr{
		// Arguments are expected to be in registers 0 and 1
		// Reg 0: object name (string)
		// Reg 1: value to set (VMDataObject)

		// Call the syscall to set the object's value
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_MEM_SET)},
	}

	StandardFuncs["echo"] = []VMInstr{
		{Op: OpStr, Oprand1: makeStrValueObj("stdout"), Oprand2: makeIntValueObj(int64(0))},
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)},
		{Op: OpRslSet, Oprand1: VMDataObject{}},
		{Op: OpClearReg},
	}

	// System Functions
	StandardFuncs["exit"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)},
		{Op: OpHlt}, // Stop execution
	}

	// Array Functions
	StandardFuncs["arrmake"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_ARR_MAKE)}, // SYS_ARRAY_MAKE
	}
	StandardFuncs["arrpush"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_ARR_PUSH)}, // SYS_ARRAY_PUSH
	}
	StandardFuncs["arrset"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_ARR_SET)}, // SYS_ARRAY_SET
	}
	StandardFuncs["arrget"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_ARR_GET)}, // SYS_ARRAY_GET
	}
	StandardFuncs["arrlen"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_ARR_LEN)}, // SYS_ARRAY_LEN
	}

	StandardFuncs["getenv"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_GET_ENV)}, // SYS_GET_ENV
	}
	StandardFuncs["exec"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_EXEC_CMD)}, // SYS_EXEC_CMD
	}
	StandardFuncs["getos"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_GET_OS_TYPE)}, // SYS_GET
	}

	// Conversion Functions
	StandardFuncs["convint"] = []VMInstr{
		{Op: OpCstInt, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["convreal"] = []VMInstr{
		{Op: OpCstReal, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["convstr"] = []VMInstr{
		{Op: OpCstStr, Oprand1: makeIntValueObj(0)},
	}

	return StandardFuncs
}