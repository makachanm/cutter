package runtime

// StandardFuncs provides the standard library functions for the Cutter VM.
// Each function is implemented as a sequence of VM instructions.
// Not all functions from standardfunctions.md are implemented due to limitations in the current VM instruction set.
var StandardFuncs map[string][]VMInstr

func init() {
	StandardFuncs = make(map[string][]VMInstr)

	// Arithmetic Functions
	StandardFuncs["add"] = []VMInstr{
		{Op: OpAdd, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(0)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["sub"] = []VMInstr{
		{Op: OpSub, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(0)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["mul"] = []VMInstr{
		{Op: OpMul, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(0)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["div"] = []VMInstr{
		{Op: OpDiv, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(0)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["mod"] = []VMInstr{
		{Op: OpMod, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(0)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}

	// Comparison Functions
	StandardFuncs["same"] = []VMInstr{
		{Op: OpCmpEq, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(0)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["notsame"] = []VMInstr{
		{Op: OpCmpNeq, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(0)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}

	// String Functions
	StandardFuncs["strcontact"] = StandardFuncs["add"]

	StandardFuncs["ifel"] = []VMInstr{
		{Op: OpBrch, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(2)}, // Branch if condition is true
	}

	// VM Internal Functions
	StandardFuncs["setreg"] = []VMInstr{
		{Op: OpRegSet, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1)},
	}
	StandardFuncs["getreg"] = []VMInstr{
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}
	StandardFuncs["setmem"] = []VMInstr{
		{Op: OpStr, Oprand1: makeStrValueObj(""), Oprand2: makeIntValueObj(0)}, // Note: This is not fully functional as the memory location is static.
	}
	StandardFuncs["getmem"] = []VMInstr{
		{Op: OpLdr, Oprand1: makeIntValueObj(0), Oprand2: makeStrValueObj("")}, // Note: This is not fully functional as the memory location is static.
	}

	StandardFuncs["exit"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)},
		{Op: OpHlt}, // Stop execution
	}

	// Functions that cannot be implemented with the current instruction set:
	// set: Requires indirect memory access (setting memory location from a register).
	// ifel: Requires conditional branching instructions.
	// bigger, smaller, bigsame, smallsame: Requires comparison opcodes for greater/less than.
	// Array functions (makeintarr, etc.): Requires array data structure support in the VM.
	// String functions (getstrlen, strslice): Requires dedicated opcodes.
	// I/O functions (getln, etc.): Requires new syscalls.
	// Type conversion functions (tostr, etc.): Requires dedicated opcodes.
	// sleep, exit, panic: Requires new syscalls.
	// Stack and Queue functions: Requires stack and queue manipulation opcodes.
}
