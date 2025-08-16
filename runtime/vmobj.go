package runtime

const (
	INTGER ValueType = iota + 1
	REAL
	STRING
	BOOLEAN
)

type VMDataObject struct {
	Type ValueType

	IntData    int64
	FloatData  float64
	BoolData   bool
	StringData string
}

type VMFunctionObject struct {
	JumpPc      int
	IsStandard  bool
	Instructions []VMInstr
}

type VMArgumentRegisters struct {
	ArgumentRegisters   []VMDataObject
	ReturnValueRegister VMDataObject
}

func NewRegister() VMArgumentRegisters {
	return VMArgumentRegisters{
		ArgumentRegisters:   make([]VMDataObject, 64),
		ReturnValueRegister: VMDataObject{},
	}
}

func (rg *VMArgumentRegisters) ClearRegisters() {
	rg.ArgumentRegisters = rg.ArgumentRegisters[:0]
	rg.ArgumentRegisters = make([]VMDataObject, 64)
}

func (rg *VMArgumentRegisters) InsertRegister(idx int, val VMDataObject) {
	rg.ArgumentRegisters[idx] = val
}

func (rg *VMArgumentRegisters) GetRegister(idx int) VMDataObject {
	return rg.ArgumentRegisters[idx]
}

func (rg *VMArgumentRegisters) InsertResult(val VMDataObject) {
	rg.ReturnValueRegister = val
}

func (rg *VMArgumentRegisters) GetResult() VMDataObject {
	return rg.ReturnValueRegister
}

type VMMEMObjectTable struct {
	ValueTable     map[string]int
	DataMemory     []VMDataObject
	FunctionMemory []VMFunctionObject

	currunt_free_dm_pointer int
	currunt_free_fm_pointer int
}

func NewVMMEMObjTable() VMMEMObjectTable {
	return VMMEMObjectTable{
		ValueTable:     make(map[string]int),
		DataMemory:     make([]VMDataObject, 0),
		FunctionMemory: make([]VMFunctionObject, 0),

		currunt_free_dm_pointer: 0,
		currunt_free_fm_pointer: 0,
	}
}

func (v *VMMEMObjectTable) MakeObj(name string) {
	v.DataMemory = append(v.DataMemory, VMDataObject{})
	v.ValueTable[name] = v.currunt_free_dm_pointer

	v.currunt_free_dm_pointer++
}

func (v *VMMEMObjectTable) GetObj(name string) *VMDataObject {
	idx, ok := v.ValueTable[name]
	if !ok || idx >= len(v.DataMemory) {
		panic("VMDataObject not found: " + name)
	}
	return &v.DataMemory[idx]
}

func (v *VMMEMObjectTable) SetObj(name string, data VMDataObject) {
	idx, ok := v.ValueTable[name]
	if !ok || idx >= len(v.DataMemory) {
		panic("VMDataObject not found: " + name)
	}
	v.DataMemory[idx] = data
}

func (v *VMMEMObjectTable) HasObj(name string) bool {
	idx, ok := v.ValueTable[name]
	if !ok || idx >= len(v.DataMemory) {
		return false
	}
	return true
}

func (v *VMMEMObjectTable) MakeFunc(name string) {
	v.FunctionMemory = append(v.FunctionMemory, VMFunctionObject{})
	v.ValueTable[name] = v.currunt_free_fm_pointer
	v.currunt_free_fm_pointer++
}

func (v *VMMEMObjectTable) GetFunc(name string) *VMFunctionObject {
	idx, ok := v.ValueTable[name]
	if !ok || idx >= len(v.FunctionMemory) {
		panic("VMFunctionObject not found: " + name)
	}
	return &v.FunctionMemory[idx]
}

func (v *VMMEMObjectTable) SetFunc(name string, fn VMFunctionObject) {
	idx, ok := v.ValueTable[name]
	if !ok || idx >= len(v.FunctionMemory) {
		panic("VMFunctionObject not found: " + name)
	}
	v.FunctionMemory[idx] = fn
}

type VMOp int
type ValueType int

const (
	OpRegSet VMOp = iota + 1
	OpMemSet
	OpRslSet
	OpRegMov
	OpMemMov
	OpRelMov
	OpLdr
	OpStr

	OpDefFunc
	OpCall
	OpReturn

	OpSyscall

	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpAnd
	OpOr
	OpNot
	OpCmpEq
	OpCmpNeq

	OpClearReg
)

type VMInstr struct {
	Op      VMOp
	Oprand1 VMDataObject
	Oprand2 VMDataObject
	Oprand3 VMDataObject
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
