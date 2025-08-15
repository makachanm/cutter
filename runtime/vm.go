package runtime

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
		//fmt.Println("Executing instruction:", instr)
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
			vm.Mem.SetObj(instr.Oprand1.StringData, instr.Oprand2)

		case OpMemMov:
			value := vm.Mem.GetObj(instr.Oprand1.StringData)
			vm.Mem.SetObj(instr.Oprand2.StringData, *value)

		case OpLdr:
			vm.Reg.InsertRegister(int(instr.Oprand1.IntData), *vm.Mem.GetObj(instr.Oprand2.StringData))

		case OpStr:
			vm.Mem.SetObj(instr.Oprand1.StringData, vm.Reg.GetRegister(int(instr.Oprand2.IntData)))

		case OpSyscall:
		}
		vm.PC++
	}
}
