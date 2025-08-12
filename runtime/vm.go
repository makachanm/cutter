package runtime

import "fmt"

type VMOp int
type ValueType int

const (
	OpCall VMOp = iota + 1
	OpSet
	OpGet
	OpLdr
	OpStr
	OpDefFunc
	OpReturn
	OpStdout
	OpSysFlush
)

type VMInstr struct {
	Op   VMOp
	Args VMArgumentRegisters
}

type CallStack struct {
	stack []int
}

func NewCallStack() *CallStack {
	return &CallStack{stack: make([]int, 0)}
}

func (cs *CallStack) Push(pc int) {
	cs.stack = append(cs.stack, pc)
}

func (cs *CallStack) Pop() int {
	if len(cs.stack) == 0 {
		panic("CallStack underflow")
	}
	val := cs.stack[len(cs.stack)-1]
	cs.stack = cs.stack[:len(cs.stack)-1]
	return val
}

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

	// 1. 함수 정의(DefFunc/OpReturn)만 먼저 등록
	pc := 0
	for pc < len(vm.Program) {
		instr := vm.Program[pc]
		if instr.Op == OpDefFunc {
			funcName := instr.Args.GetRegister(0).StringData
			funcObj := VMFunctionObject{
				JumpPc: int(instr.Args.GetRegister(1).IntData),
			}
			vm.Mem.MakeFunc(funcName)
			vm.Mem.SetFunc(funcName, funcObj)
		}
		if instr.Op == OpReturn && vm.Program[pc+1].Op != OpDefFunc {
			pc++
			break // main 코드 시작점
		}
		pc++
	}

	// 2. main 코드 실행 (함수 정의 이후부터)
	vm.PC = pc
	for vm.PC < len(vm.Program) {
		instr := vm.Program[vm.PC]
		//fmt.Println("Executing instruction:", instr)
		switch instr.Op {
		case OpCall:
			funcObj := vm.Mem.GetFunc(instr.Args.GetRegister(0).StringData)
			targetpc := funcObj.JumpPc
			vm.Stack.Push(vm.PC + 1) // OpCall 다음 명령어로 복귀
			vm.PC = targetpc

			fmt.Println("Calling function:", instr.Args.GetRegister(0).StringData, "at PC:", targetpc)
			continue
		case OpReturn:
			vm.PC = vm.Stack.Pop()
			vm.Reg.ClearRegisters()
			continue
		case OpSet:
			if !vm.Mem.HasObj(instr.Args.GetRegister(0).StringData) {
				vm.Mem.MakeObj(instr.Args.GetRegister(0).StringData)
				fmt.Println("Created new VMDataObject:", instr.Args.GetRegister(0).StringData)
			}
			vm.Mem.SetObj(instr.Args.GetRegister(0).StringData, instr.Args.GetRegister(1))
		case OpGet:
			obj := vm.Mem.GetObj(instr.Args.GetRegister(0).StringData)
			vm.Reg.InsertRegister(int(instr.Args.GetRegister(1).IntData), *obj)
		case OpLdr:
			vm.Reg.InsertRegister(int(instr.Args.GetRegister(0).IntData), *vm.Mem.GetObj(instr.Args.GetRegister(1).StringData))
		case OpStr:
			vm.Mem.SetObj(instr.Args.GetRegister(0).StringData, vm.Reg.GetRegister(int(instr.Args.GetRegister(1).IntData)))
		case OpStdout:
			vm.IO.WriteObjectToStream(*vm.Mem.GetObj("stdout"))
		case OpSysFlush:
			vm.IO.FlushIO()
		default:
			fmt.Println("Unknown VM instruction:", instr.Op)
			return
		}
		vm.PC++
	}
}
