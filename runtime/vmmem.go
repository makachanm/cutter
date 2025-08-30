package runtime

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

type VMArgumentRegisters struct {
	ArgumentRegisterMemory []VMDataObject
	ArgumentRegisterMap    map[int]int
	ReturnValueRegister    VMDataObject

	last_allocated_area int
}

func NewRegister() VMArgumentRegisters {
	return VMArgumentRegisters{
		ArgumentRegisterMemory: make([]VMDataObject, 0),
		ArgumentRegisterMap:    make(map[int]int, 0),
		ReturnValueRegister:    VMDataObject{},

		last_allocated_area: 0,
	}
}

func (rg *VMArgumentRegisters) ClearRegisters() {
	rg.ArgumentRegisterMemory = rg.ArgumentRegisterMemory[:0]
	for k := range rg.ArgumentRegisterMap {
		delete(rg.ArgumentRegisterMap, k)
	}

	rg.last_allocated_area = 0
}

func (rg *VMArgumentRegisters) InsertRegister(idx int, val VMDataObject) {
	rg.ArgumentRegisterMap[idx] = rg.last_allocated_area
	rg.ArgumentRegisterMemory = append(rg.ArgumentRegisterMemory, val)
	rg.last_allocated_area++
}

func (rg *VMArgumentRegisters) GetRegister(idx int) VMDataObject {
	pos, exist := rg.ArgumentRegisterMap[idx]
	if !exist {
		panic("VM Register Error - cannot find register")
	}
	return rg.ArgumentRegisterMemory[pos]
}

func (rg *VMArgumentRegisters) InsertResult(val VMDataObject) {
	rg.ReturnValueRegister = val
}

func (rg *VMArgumentRegisters) GetResult() VMDataObject {
	return rg.ReturnValueRegister
}

type VMMEMObjectTable struct {
	DataTable      map[string]int
	FunctionTable  map[string]int
	ArrayTable     map[string][]VMDataObject
	DataMemory     []VMDataObject
	FunctionMemory []VMFunctionObject

	currunt_free_dm_pointer int
	currunt_free_fm_pointer int
}

func NewVMMEMObjTable() VMMEMObjectTable {
	return VMMEMObjectTable{
		DataTable:      make(map[string]int),
		FunctionTable:  make(map[string]int),
		ArrayTable:     make(map[string][]VMDataObject),
		DataMemory:     make([]VMDataObject, 0),
		FunctionMemory: make([]VMFunctionObject, 0),

		currunt_free_dm_pointer: 0,
		currunt_free_fm_pointer: 0,
	}
}

func (v *VMMEMObjectTable) MakeObj(name string) {
	v.DataMemory = append(v.DataMemory, VMDataObject{})
	v.DataTable[name] = v.currunt_free_dm_pointer

	v.currunt_free_dm_pointer++
}

func (v *VMMEMObjectTable) GetObj(name string) *VMDataObject {
	idx, ok := v.DataTable[name]
	if !ok {
		panic("VMDataObject not found: " + name)
	}
	return &v.DataMemory[idx]
}

func (v *VMMEMObjectTable) SetObj(name string, data VMDataObject) {
	idx, ok := v.DataTable[name]
	if !ok {
		panic("VMDataObject not found: " + name)
	}
	v.DataMemory[idx] = data
}

func (v *VMMEMObjectTable) HasObj(name string) bool {
	idx, ok := v.DataTable[name]
	if !ok || idx >= len(v.DataMemory) {
		return false
	}
	return true
}

func (v *VMMEMObjectTable) MakeFunc(name string) {
	v.FunctionMemory = append(v.FunctionMemory, VMFunctionObject{})
	v.FunctionTable[name] = v.currunt_free_fm_pointer
	v.currunt_free_fm_pointer++
}

func (v *VMMEMObjectTable) GetFunc(name string) *VMFunctionObject {
	idx, ok := v.FunctionTable[name]
	if !ok || idx >= len(v.FunctionMemory) {
		panic("VMFunctionObject not found: " + name)
	}
	return &v.FunctionMemory[idx]
}

func (v *VMMEMObjectTable) SetFunc(name string, fn VMFunctionObject) {
	idx, ok := v.FunctionTable[name]
	if !ok || idx >= len(v.FunctionMemory) {
		panic("VMFunctionObject not found: " + name)
	}
	v.FunctionMemory[idx] = fn
}

func (v *VMMEMObjectTable) MakeArray(name string) {
	v.ArrayTable[name] = make([]VMDataObject, 0)
}

func (t *VMMEMObjectTable) GetArray(name string) []VMDataObject {
	arr, ok := t.ArrayTable[name]
	if !ok {
		panic("Array not found: " + name)
	}
	return arr
}

func (t *VMMEMObjectTable) PushArrayItem(name string, item VMDataObject) {
	t.ArrayTable[name] = append(t.ArrayTable[name], item)
}

func (t *VMMEMObjectTable) SetArrayItem(name string, idx int, item VMDataObject) {
	if idx >= len(t.ArrayTable[name]) {
		panic("Array index out of bounds: " + name)
	}
	t.ArrayTable[name][idx] = item

}

func (t *VMMEMObjectTable) HasArray(name string) bool {
	_, ok := t.ArrayTable[name]
	return ok
}
