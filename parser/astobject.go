package parser

type ValueType int
type FunctionType int

const (
	VALUE_FUNCTION FunctionType = iota + 1
	EXCUTABLE_FUNCTION
)

const (
	INTGER ValueType = iota + 1
	REAL
	STRING
	BOOLEAN
)

type FunctionObject struct {
	Name string
	Type FunctionType

	Args       []CallObject
	FuncBodys  CallObject
	ValueBodys ValueObject
}

type CallObject struct {
	Name         string
	Args         []ValueObject
	CallableArgs []CallObject
}

type ValueObject struct {
	Type ValueType

	IntData    int64
	FloatData  float64
	BoolData   bool
	StringData string
}

func makeIntValueObj(input int64) ValueObject {
	return ValueObject{Type: INTGER, IntData: input}
}

func makeRealValueObj(input float64) ValueObject {
	return ValueObject{Type: REAL, FloatData: input}
}

func makeBoolValueObj(input bool) ValueObject {
	return ValueObject{Type: BOOLEAN, BoolData: input}
}

func makeStrValueObj(input string) ValueObject {
	return ValueObject{Type: STRING, StringData: input}
}

type BodyType int

const (
	FUCNTION_DEFINITION BodyType = iota + 1
	FUNCTION_CALL
)

type HeadNode struct {
	Bodys []BodyObject
}

type BodyObject struct {
	Type BodyType

	Func FunctionObject
	Call CallObject
}

func NewFunctionBodyObject(funs FunctionObject) BodyObject {
	return BodyObject{Type: FUCNTION_DEFINITION, Func: funs}
}

func NewCallBodyObject(calls CallObject) BodyObject {
	return BodyObject{Type: FUNCTION_CALL, Call: calls}
}
