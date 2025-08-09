package parser

type ValueType int

const (
	INTGER ValueType = iota + 1
	REAL
	STRING
	BOOLEAN
)

type BodyObject struct {
	Name string
	Args []ValueObject
	//when we da- da- da- dance, when we da- a- a- ance
	Bodys []BodyObject
}

type ValueObject struct {
	Type ValueType

	IntData    int64
	FloatData  float64
	BoolData   bool
	StringData string
}

type HeadNode struct {
	Bodys []BodyObject
}
