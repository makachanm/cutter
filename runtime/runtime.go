package runtime

import "cutter/parser"

type CutterRuntime struct {
	fntable FunctionLookupTable
	io      RuntimeIO
}

func NewRuntime() CutterRuntime {
	return CutterRuntime{
		fntable: NewFunctionLookupTable(),
		io:      NewIO(),
	}
}

func (r *CutterRuntime) RunAST(input parser.HeadNode) {

}
