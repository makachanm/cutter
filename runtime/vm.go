package runtime

import (
	"fmt"
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
	return &VM{
		Stack:   NewCallStack(),
		Program: input,
		Reg:     NewRegister(),
		Mem:     NewVMMEMObjTable(),
		IO:      NewIO(),
		PC:      0,

		isFuncDefineState: false,
	}
}

func (vm *VM) Run() {
	vm.Mem.MakeObj("stdout")

	pc := 0
	for pc < len(vm.Program) {
		instr := vm.Program[pc]
		if instr.Op == OpDefFunc {
			funcName := instr.Oprand1.StringData
			funcObj := VMFunctionObject{
				JumpPc: int(instr.Oprand2.IntData),
			}
			vm.Mem.MakeFunc(funcName)
			vm.Mem.SetFunc(funcName, funcObj)
		}
		if instr.Op == OpReturn && vm.Program[pc+1].Op != OpDefFunc {
			pc++
			break
		}
		pc++
	}

	vm.PC = pc
	for vm.PC < len(vm.Program) {
		instr := vm.Program[vm.PC]
		fmt.Println("Executing instruction:", instr)
		switch instr.Op {
		case OpCall:
			funcName := instr.Oprand1.StringData
			funcObj := vm.Mem.GetFunc(funcName)

			vm.Stack.Push(vm.PC + 1)
			vm.PC = funcObj.JumpPc

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

		case OpRelMov:
			value := vm.Reg.GetResult()
			vm.Reg.InsertRegister(int(instr.Oprand1.IntData), value)

		case OpLdr:
			vm.Reg.InsertRegister(int(instr.Oprand1.IntData), *vm.Mem.GetObj(instr.Oprand2.StringData))

		case OpStr:
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

		case OpClearReg:
			vm.Reg.ClearRegisters()
		}
		vm.PC++
	}
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
			i, _ := strconv.ParseInt(r2.StringData, 10, 64)
			return VMDataObject{Type: INTGER, IntData: intOp(r1.IntData, i)}
		}
	case REAL:
		switch r2.Type {
		case INTGER:
			return VMDataObject{Type: REAL, FloatData: floatOp(r1.FloatData, float64(r2.IntData))}
		case REAL:
			return VMDataObject{Type: REAL, FloatData: floatOp(r1.FloatData, r2.FloatData)}
		case STRING:
			f, _ := strconv.ParseFloat(r2.StringData, 64)
			return VMDataObject{Type: REAL, FloatData: floatOp(r1.FloatData, f)}
		}
	case STRING:
		switch r2.Type {
		case INTGER:
			i, _ := strconv.ParseInt(r1.StringData, 10, 64)
			return VMDataObject{Type: INTGER, IntData: intOp(i, r2.IntData)}
		case REAL:
			f, _ := strconv.ParseFloat(r1.StringData, 64)
			return VMDataObject{Type: REAL, FloatData: floatOp(f, r2.FloatData)}
		case STRING:
			if strOp != nil {
				return VMDataObject{Type: STRING, StringData: strOp(r1.StringData, r2.StringData)}
			}
		}
	}
	return VMDataObject{}
}
