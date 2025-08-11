package runtime

import "cutter/parser"

type Function struct {
	Body parser.FunctionObject
}

type StdFunction struct {
	Action func() parser.ValueObject
}

type FunctionLookupTable struct {
	FunctionTables       map[string]Function
	NonRedefineableTable map[string]StdFunction
}

func NewFunctionLookupTable() FunctionLookupTable {
	return FunctionLookupTable{}
}

func (t *FunctionLookupTable) InsertToTable(fnc parser.FunctionObject) {
	t.FunctionTables[fnc.Name] = Function{Body: fnc}
}

func (t *FunctionLookupTable) LookupTable(name string) Function {
	fnc, exist := t.NonRedefineableTable[name]
	if exist {
		data := fnc.Action()
		return Function{Body: parser.FunctionObject{Type: parser.VALUE_FUNCTION, ValueBodys: data}}
	}

	nfnc, nexist := t.FunctionTables[name]
	if nexist {
		return nfnc
	} else {
		panic("function not defined: " + name)
	}
}
