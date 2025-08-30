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
			if intOp != nil {
				return VMDataObject{Type: INTGER, IntData: intOp(r1.IntData, r2.IntData)}
			}
		case REAL:
			if floatOp != nil {
				return VMDataObject{Type: REAL, FloatData: floatOp(float64(r1.IntData), r2.FloatData)}
			}
		case STRING:
			if strOp != nil {
				return VMDataObject{Type: STRING, StringData: strOp(strconv.FormatInt(r1.IntData, 10), r2.StringData)}
			}
		}
	case REAL:
		switch r2.Type {
		case INTGER:
			if floatOp != nil {
				return VMDataObject{Type: REAL, FloatData: floatOp(r1.FloatData, float64(r2.IntData))}
			}
		case REAL:
			if floatOp != nil {
				return VMDataObject{Type: REAL, FloatData: floatOp(r1.FloatData, r2.FloatData)}
			}
		case STRING:
			if strOp != nil {
				return VMDataObject{Type: STRING, StringData: strOp(strconv.FormatFloat(r1.FloatData, 'f', -1, 64), r2.StringData)}
			}
		}
	case STRING:
		switch r2.Type {
		case INTGER:
			if strOp != nil {
				return VMDataObject{Type: STRING, StringData: strOp(r1.StringData, strconv.FormatInt(r2.IntData, 10))}
			}
		case REAL:
			if strOp != nil {
				return VMDataObject{Type: STRING, StringData: strOp(r1.StringData, strconv.FormatFloat(r2.FloatData, 'f', -1, 64))}
			}
		case STRING:
			if strOp != nil {
				return VMDataObject{Type: STRING, StringData: strOp(r1.StringData, r2.StringData)}
			}
		}
	}
	panic("Unsupported operation")
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
			return makeStrValueObj(strconv.FormatInt(obj.IntData, 10))
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
