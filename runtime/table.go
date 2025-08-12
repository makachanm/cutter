package runtime

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

func (v *VMMEMObjectTable) SetObj(name string, data VMDataObject) bool {
	idx, ok := v.ValueTable[name]
	if !ok || idx >= len(v.DataMemory) {
		panic("VMDataObject not found: " + name)
	}
	v.DataMemory[idx] = data
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

func (v *VMMEMObjectTable) SetFunc(name string, fn VMFunctionObject) bool {
	idx, ok := v.ValueTable[name]
	if !ok || idx >= len(v.FunctionMemory) {
		panic("VMFunctionObject not found: " + name)
	}
	v.FunctionMemory[idx] = fn
	return true
}
