package runtime

// GetStandardFuncs provides the standard library functions for the Cutter VM.
// Each function is implemented as a sequence of VM instructions.
// Not all functions from standardfunctions.md are implemented due to limitations in the current VM instruction set.
func GetStandardFuncs() map[string][]VMInstr {
	StandardFuncs := make(map[string][]VMInstr)

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

	// Branching and Control Flow
	StandardFuncs["ifel"] = []VMInstr{
		{Op: OpBrch, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1), Oprand3: makeIntValueObj(2)},
		{Op: OpRslSet, Oprand1: makeIntValueObj(0)},
	}

	// Memory and Variable Manipulation
	StandardFuncs["set"] = []VMInstr{
		{Op: OpStrInd, Oprand1: makeIntValueObj(0), Oprand2: makeIntValueObj(1)},
	}

	// System Functions
	StandardFuncs["exit"] = []VMInstr{
		{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)},
		{Op: OpHlt}, // Stop execution
	}

	return StandardFuncs
}