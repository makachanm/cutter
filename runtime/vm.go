package runtime

import (
	"strconv"
)

type VM struct {
	Stack   *CallStack
	Program []VMInstr
	Reg     VMArgumentRegisters
	Mem     VMMEMObjectTable
	IO      RuntimeIO

	PC int

	isFuncDefineState bool
}

func NewVM(input []VMInstr) *VM {
	vm := &VM{
		Stack:   NewCallStack(),
		Program: input,
		Reg:     NewRegister(), Mem: NewVMMEMObjTable(),
		IO: NewIO(),
		PC: 0,

		isFuncDefineState: false,
	}

	// Register standard functions
	for name, instrs := range GetStandardFuncs() {
		vm.Mem.MakeFunc(name)
		vm.Mem.SetFunc(name, VMFunctionObject{
			JumpPc:       -1, // Special value to indicate a standard function
			IsStandard:   true,
			Instructions: instrs,
		})
	}
	return vm
}

func (vm *VM) Run() {
	vm.Mem.MakeObj("stdout")

	vm.PC = 0
	for vm.PC < len(vm.Program) {
		instr := vm.Program[vm.PC]

		switch instr.Op {
		case OpDefFunc:
			funcName := instr.Oprand1.StringData
			funcObj := VMFunctionObject{
				JumpPc: vm.PC + 1,
			}
			vm.Mem.MakeFunc(funcName)
			vm.Mem.SetFunc(funcName, funcObj)

			// Skip to the end of the function definition
			for vm.PC < len(vm.Program) && vm.Program[vm.PC].Op != OpReturn {
				vm.PC++
			}

		case OpCall:
			funcName := instr.Oprand1.StringData
			funcObj := vm.Mem.GetFunc(funcName)

			if funcObj.IsStandard {
				// Execute standard function instructions
				for _, stdInstr := range funcObj.Instructions {
					vm.executeInstruction(stdInstr)
				}
				// After executing standard function, continue with the next instruction
				// in the main program.
				vm.PC++ // Move to the next instruction after the OpCall
				continue
			} else {
				// Existing logic for user-defined functions
				vm.Stack.Push(vm.PC)
				vm.PC = funcObj.JumpPc
				continue
			}

		case OpReturn:
			if len(vm.Stack.stack) == 0 {
				panic("Call stack is empty, cannot return.")
			}
			vm.PC = vm.Stack.Pop()

		case OpRegSet:
			vm.Reg.InsertRegister(int(instr.Oprand1.IntData), instr.Oprand2)

		case OpRegMov:
			value := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			vm.Reg.InsertRegister(int(instr.Oprand2.IntData), value)

		case OpMemSet:
			if !vm.Mem.HasObj(instr.Oprand1.StringData) {
				vm.Mem.MakeObj(instr.Oprand1.StringData)
			}
			vm.Mem.SetObj(instr.Oprand1.StringData, instr.Oprand2)

		case OpMemMov:
			value := vm.Mem.GetObj(instr.Oprand1.StringData)
			vm.Mem.SetObj(instr.Oprand2.StringData, *value)

		case OpRslSet:
			vm.Reg.InsertResult(vm.Reg.GetRegister(int(instr.Oprand1.IntData)))

		case OpRslMov:
			value := vm.Reg.GetResult()
			vm.Reg.InsertRegister(int(instr.Oprand1.IntData), value)

		case OpLdr:
			vm.Reg.InsertRegister(int(instr.Oprand1.IntData), *vm.Mem.GetObj(instr.Oprand2.StringData))

		case OpStr:
			if !vm.Mem.HasObj(instr.Oprand1.StringData) {
				vm.Mem.MakeObj(instr.Oprand1.StringData)
			}
			vm.Mem.SetObj(instr.Oprand1.StringData, vm.Reg.GetRegister(int(instr.Oprand2.IntData)))

		case OpRslStr:
			vm.Mem.SetObj(instr.Oprand1.StringData, vm.Reg.GetResult())

		case OpStrReg:
			targetData := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			fromData := vm.Mem.GetObj(vm.Reg.GetRegister(int(instr.Oprand2.IntData)).StringData)

			if !vm.Mem.HasObj(targetData.StringData) {
				vm.Mem.MakeObj(targetData.StringData)
			}
			vm.Mem.SetObj(targetData.StringData, *fromData)

		case OpSyscall:
			doSyscall(vm, instr)

		case OpAdd:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performOperation(r1, r2, func(a, b float64) float64 { return a + b }, func(a, b int64) int64 { return a + b }, func(a, b string) string { return a + b })
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpSub:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performOperation(r1, r2, func(a, b float64) float64 { return a - b }, func(a, b int64) int64 { return a - b }, nil)
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpMul:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performOperation(r1, r2, func(a, b float64) float64 { return a * b }, func(a, b int64) int64 { return a * b }, nil)
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpDiv:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performOperation(r1, r2, func(a, b float64) float64 { return a / b }, func(a, b int64) int64 { return a / b }, nil)
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpMod:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performOperation(r1, r2, nil, func(a, b int64) int64 { return a % b }, nil)
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpAnd:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			if r1.Type == BOOLEAN && r2.Type == BOOLEAN {
				vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1.BoolData && r2.BoolData})
			}

		case OpOr:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			if r1.Type == BOOLEAN && r2.Type == BOOLEAN {
				vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1.BoolData || r2.BoolData})
			}

		case OpNot:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			if r1.Type == BOOLEAN {
				vm.Reg.InsertRegister(int(instr.Oprand2.IntData), VMDataObject{Type: BOOLEAN, BoolData: !r1.BoolData})
			}

		case OpCmpEq:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1 == r2})

		case OpCmpNeq:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1 != r2})

		case OpCmpGt:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performComparison(r1, r2, func(a, b float64) bool { return a > b }, func(a, b int64) bool { return a > b })
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpCmpLt:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performComparison(r1, r2, func(a, b float64) bool { return a < b }, func(a, b int64) bool { return a < b })
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpCmpGte:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performComparison(r1, r2, func(a, b float64) bool { return a >= b }, func(a, b int64) bool { return a >= b })
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpCmpLte:
			r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
			result := performComparison(r1, r2, func(a, b float64) bool { return a <= b }, func(a, b int64) bool { return a <= b })
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)

		case OpBrch:
			condition := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			if condition.Type != BOOLEAN {
				panic("Branch condition must be BOOLEAN type")
			}

			if condition.BoolData {
				vm.Reg.InsertResult(vm.Reg.GetRegister(int(instr.Oprand2.IntData)))
			} else {
				vm.Reg.InsertResult(vm.Reg.GetRegister(int(instr.Oprand3.IntData)))
			}

		case OpJmp:
			vm.PC = int(instr.Oprand1.IntData)
			continue

		case OpJmpIfFalse:
			condition := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
			if condition.Type == BOOLEAN && !condition.BoolData {
				vm.PC = int(instr.Oprand2.IntData)
				continue
			}

		case OpClearReg:
			vm.Reg.ClearRegisters()

		case OpHlt:
			// Stop execution
			return
		}
		vm.PC++
	}
}

func performComparison(r1, r2 VMDataObject, floatOp func(float64, float64) bool, intOp func(int64, int64) bool) VMDataObject {
	// Default to false if types are incompatible
	result := false

	switch r1.Type {
	case INTGER:
		switch r2.Type {
		case INTGER:
			if intOp != nil {
				result = intOp(r1.IntData, r2.IntData)
			}
		case REAL:
			if floatOp != nil {
				result = floatOp(float64(r1.IntData), r2.FloatData)
			}
		}
	case REAL:
		switch r2.Type {
		case INTGER:
			if floatOp != nil {
				result = floatOp(r1.FloatData, float64(r2.IntData))
			}
		case REAL:
			if floatOp != nil {
				result = floatOp(r1.FloatData, r2.FloatData)
			}
		}
	}
	return VMDataObject{Type: BOOLEAN, BoolData: result}
}

func performOperation(r1, r2 VMDataObject, floatOp func(float64, float64) float64, intOp func(int64, int64) int64, strOp func(string, string) string) VMDataObject {
	switch r1.Type {
	case INTGER:
		switch r2.Type {
		case INTGER:
			return VMDataObject{Type: INTGER, IntData: intOp(r1.IntData, r2.IntData)}
		case REAL:
			return VMDataObject{Type: REAL, FloatData: floatOp(float64(r1.IntData), r2.FloatData)}
		case STRING:
			return VMDataObject{Type: STRING, StringData: strOp(strconv.FormatInt(r1.IntData, 10), r2.StringData)}
		}
	case REAL:
		switch r2.Type {
		case INTGER:
			return VMDataObject{Type: REAL, FloatData: floatOp(r1.FloatData, float64(r2.IntData))}
		case REAL:
			return VMDataObject{Type: REAL, FloatData: floatOp(r1.FloatData, r2.FloatData)}
		case STRING:
			return VMDataObject{Type: STRING, StringData: strOp(strconv.FormatFloat(r1.FloatData, 'f', -1, 64), r2.StringData)}
		}
	case STRING:
		switch r2.Type {
		case INTGER:
			return VMDataObject{Type: STRING, StringData: strOp(r1.StringData, strconv.FormatInt(r2.IntData, 10))}
		case REAL:
			return VMDataObject{Type: STRING, StringData: strOp(r1.StringData, strconv.FormatFloat(r2.FloatData, 'f', -1, 64))}
		case STRING:
			if strOp != nil {
				return VMDataObject{Type: STRING, StringData: strOp(r1.StringData, r2.StringData)}
			}
		}
	}
	return VMDataObject{}
}

func (vm *VM) executeInstruction(instr VMInstr) {
	switch instr.Op {
	case OpRegSet:
		vm.Reg.InsertRegister(int(instr.Oprand1.IntData), instr.Oprand2)
	case OpMemSet:
		if !vm.Mem.HasObj(instr.Oprand1.StringData) {
			vm.Mem.MakeObj(instr.Oprand1.StringData)
		}
		vm.Mem.SetObj(instr.Oprand1.StringData, instr.Oprand2)
	case OpRslSet:
		vm.Reg.InsertResult(vm.Reg.GetRegister(int(instr.Oprand1.IntData)))
	case OpRegMov:
		value := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		vm.Reg.InsertRegister(int(instr.Oprand2.IntData), value)
	case OpMemMov:
		value := vm.Mem.GetObj(instr.Oprand1.StringData)
		vm.Mem.SetObj(instr.Oprand2.StringData, *value)
	case OpRslMov:
		value := vm.Reg.GetResult()
		vm.Reg.InsertRegister(int(instr.Oprand1.IntData), value)
	case OpLdr:
		vm.Reg.InsertRegister(int(instr.Oprand1.IntData), *vm.Mem.GetObj(instr.Oprand2.StringData))
	case OpStr:
		if !vm.Mem.HasObj(instr.Oprand1.StringData) {
			vm.Mem.MakeObj(instr.Oprand1.StringData)
		}
		vm.Mem.SetObj(instr.Oprand1.StringData, vm.Reg.GetRegister(int(instr.Oprand2.IntData)))
	case OpSyscall:
		doSyscall(vm, instr)
	case OpAdd:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performOperation(r1, r2, func(a, b float64) float64 { return a + b }, func(a, b int64) int64 { return a + b }, func(a, b string) string { return a + b })
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpSub:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performOperation(r1, r2, func(a, b float64) float64 { return a - b }, func(a, b int64) int64 { return a - b }, nil)
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpMul:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performOperation(r1, r2, func(a, b float64) float64 { return a * b }, func(a, b int64) int64 { return a * b }, nil)
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpDiv:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performOperation(r1, r2, func(a, b float64) float64 { return a / b }, func(a, b int64) int64 { return a / b }, nil)
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpMod:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performOperation(r1, r2, nil, func(a, b int64) int64 { return a % b }, nil)
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpAnd:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		if r1.Type == BOOLEAN && r2.Type == BOOLEAN {
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1.BoolData && r2.BoolData})
		}
	case OpOr:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		if r1.Type == BOOLEAN && r2.Type == BOOLEAN {
			vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1.BoolData || r2.BoolData})
		}
	case OpNot:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		if r1.Type == BOOLEAN {
			vm.Reg.InsertRegister(int(instr.Oprand2.IntData), VMDataObject{Type: BOOLEAN, BoolData: !r1.BoolData})
		}
	case OpCmpEq:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1 == r2})
	case OpCmpNeq:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), VMDataObject{Type: BOOLEAN, BoolData: r1 != r2})
	case OpCmpGt:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performComparison(r1, r2, func(a, b float64) bool { return a > b }, func(a, b int64) bool { return a > b })
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpCmpLt:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performComparison(r1, r2, func(a, b float64) bool { return a < b }, func(a, b int64) bool { return a < b })
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpCmpGte:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performComparison(r1, r2, func(a, b float64) bool { return a >= b }, func(a, b int64) bool { return a >= b })
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpCmpLte:
		r1 := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		r2 := vm.Reg.GetRegister(int(instr.Oprand2.IntData))
		result := performComparison(r1, r2, func(a, b float64) bool { return a <= b }, func(a, b int64) bool { return a <= b })
		vm.Reg.InsertRegister(int(instr.Oprand3.IntData), result)
	case OpBrch:
		condition := vm.Reg.GetRegister(int(instr.Oprand1.IntData))
		if condition.Type != BOOLEAN {
			panic("Branch condition must be BOOLEAN type")
		}

		if condition.BoolData {
			vm.Reg.InsertResult(vm.Reg.GetRegister(int(instr.Oprand2.IntData)))
		} else {
			vm.Reg.InsertResult(vm.Reg.GetRegister(int(instr.Oprand3.IntData)))
		}

	case OpJmp:
		// This is a control flow instruction and should only be handled by the main Run loop.
		panic("OpJmp should not be called from executeInstruction")

	case OpJmpIfFalse:
		// This is a control flow instruction and should only be handled by the main Run loop.
		panic("OpJmpIfFalse should not be called from executeInstruction")

	case OpClearReg:
		vm.Reg.ClearRegisters()

	case OpHlt:
		// Stop execution
		return
	}
}
