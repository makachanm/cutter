package runtime

type VMOp int
type ValueType int

const (
	OpPushInt VMOp = iota + 1
	OpPushReal
	OpPushString
	OpPushBool
	OpPushSymbol
	OpCall
	OpSet
	OpSetFunc
	OpReturn
)

type VMInstr struct {
	Op   VMOp
	Arg  interface{}
	Args []interface{}
}

type VM struct {
	Stack   []interface{}
	Program []VMInstr
	Reg     VMArgumentRegisters
	PC      int
}

func NewVM() *VM {
	return &VM{
		Stack:      make([]interface{}, 0),
		ValueTable: make(map[string]int64),
		Memory:     make([]VMDataObject, 0),
		PC:         0,
	}
}

func (vm *VM) Run() {
	for vm.PC < len(vm.Program) {
		instr := vm.Program[vm.PC]
		switch instr.Op {
		case OpPushInt, OpPushReal, OpPushString, OpPushBool:
			vm.Stack = append(vm.Stack, instr.Arg)
		case OpCall:

		}
		vm.PC++
	}
}
