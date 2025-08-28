package runtime

import (
	"strconv"
)

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

func (d1 VMDataObject) IsEqualTo(d2 VMDataObject) bool {
	if d1.Type != d2.Type {
		return false
	}
	switch d1.Type {
	case INTGER:
		return d1.IntData == d2.IntData
	case REAL:
		return d1.FloatData == d2.FloatData
	case STRING:
		return d1.StringData == d2.StringData
	case BOOLEAN:
		return d1.BoolData == d2.BoolData
	}
	return false
}

func (d1 VMDataObject) IsNotEqualTo(d2 VMDataObject) bool {
	if d1.Type != d2.Type {
		return false
	}
	switch d1.Type {
	case INTGER:
		return d1.IntData != d2.IntData
	case REAL:
		return d1.FloatData != d2.FloatData
	case STRING:
		return d1.StringData != d2.StringData
	case BOOLEAN:
		return d1.BoolData != d2.BoolData
	}
	return false
}

func (r1 VMDataObject) Compare(r2 VMDataObject, floatOp func(float64, float64) bool, intOp func(int64, int64) bool) VMDataObject {
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

func (r1 VMDataObject) Operate(r2 VMDataObject, floatOp func(float64, float64) float64, intOp func(int64, int64) int64, strOp func(string, string) string) VMDataObject {
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

func (obj *VMDataObject) CastTo(d_type ValueType) VMDataObject {
	switch d_type {
	case INTGER:
		switch obj.Type {
		case REAL:
			val := int64(obj.FloatData)
			return makeIntValueObj(val)
		case STRING:
			val, err := strconv.ParseInt(obj.StringData, 10, 64)
			if err != nil {
				panic("Error Occured in Converting Object - " + err.Error())
			}
			return makeIntValueObj(val)

		default:
			panic("Object cannot be converted to " + string(d_type))

		}

	case REAL:
		switch obj.Type {
		case INTGER:
			val := float64(obj.IntData)
			return makeRealValueObj(val)
		case STRING:
			val, err := strconv.ParseFloat(obj.StringData, 64)
			if err != nil {
				panic("Error Occured in Converting Object - " + err.Error())
			}
			return makeRealValueObj(val)

		default:
			panic("Object cannot be converted to " + string(d_type))

		}

	case STRING:
		switch obj.Type {
		case INTGER:
			return makeStrValueObj(string(obj.IntData))
		case REAL:
			return makeStrValueObj(strconv.FormatFloat(obj.FloatData, 'f', -1, 64))

		case BOOLEAN:
			if obj.BoolData {
				return makeStrValueObj("!t")
			} else {
				return makeStrValueObj("!f")
			}

		default:
			panic("Object cannot be converted to " + string(d_type))

		}

	default:
		panic("Object cannot be converted to " + string(d_type))
	}
}

type VMFunctionObject struct {
	JumpPc       int
	IsStandard   bool
	Instructions []VMInstr
}

type VMArgumentRegisters struct {
	ArgumentRegisters   []VMDataObject
	ReturnValueRegister VMDataObject
}

func NewRegister() VMArgumentRegisters {
	return VMArgumentRegisters{
		ArgumentRegisters:   make([]VMDataObject, 1024),
		ReturnValueRegister: VMDataObject{},
	}
}

func (rg *VMArgumentRegisters) ClearRegisters() {
	rg.ArgumentRegisters = rg.ArgumentRegisters[:0]
	rg.ArgumentRegisters = make([]VMDataObject, 1024)
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

type VMOp int
type ValueType int

const (
	OpRegSet VMOp = iota + 1
	OpMemSet
	OpRslSet
	OpRegMov
	OpMemMov
	OpRslMov
	OpLdr
	OpStr
	OpRslStr
	OpStrReg

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
	OpCmpGt
	OpCmpLt
	OpCmpGte
	OpCmpLte

	OpBrch
	OpJmp
	OpJmpIfFalse

	OpCstInt
	OpCstReal
	OpCstStr

	OpClearReg
	OpHlt
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
