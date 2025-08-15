package runtime

import (
	"cutter/parser"
	"fmt"
)

type Compiler struct {
	reg *regAlloc
}

func NewCompiler() *Compiler {
	return &Compiler{reg: &regAlloc{}}
}

type regAlloc struct {
	next int
}

func (r *regAlloc) alloc() int {
	idx := r.next
	r.next++
	return idx
}

func (r *regAlloc) tmpVar() string {
	name := fmt.Sprintf("_tmp%d", r.next)
	r.next++
	return name
}

func (c *Compiler) CompileASTToVMInstr(input parser.HeadNode) []VMInstr {
	instructions := make([]VMInstr, 0)
	for _, items := range input.Bodys {
		switch items.Type {
		case parser.FUCNTION_DEFINITION:
		case parser.FUNCTION_CALL:
		case parser.NORM_STRINGS:
		}
	}

	return instructions
}

func (c *Compiler) CompileFunctionDefToVMInstr(fnc parser.FunctionObject, last_pos int) []VMInstr {
	instructions := make([]VMInstr, 0)

	defregs := NewRegister()
	defregs.InsertRegister(0, makeStrValueObj(fnc.Name))
	defregs.InsertRegister(1, makeIntValueObj(int64(last_pos)))
	instructions = append(instructions, VMInstr{Op: OpDefFunc, Args: defregs})

	// Allocate instr for function arguments
	for i, arg := range fnc.Args {
		if i >= len(fnc.Args)-1 {
			break
		}

		regs := NewRegister()
		regs.InsertRegister(0, makeStrValueObj(arg.Name))

		instructions = append(instructions, VMInstr{Op: OpSet, Args: regs})
	}

	// Compile function body
	for _, body := range fnc.FuncBodys {

	}

	instructions = append(instructions, VMInstr{Op: OpReturn, Args: NewRegister()})

	return instructions
}

func makeIntValueObj(i int64) VMDataObject {
	return VMDataObject{
		Type:    INTGER,
		IntData: i,
	}
}

func makeStrValueObj(s string) VMDataObject {
	return VMDataObject{
		Type:       STRING,
		StringData: s,
	}
}

// NewCompiler: Compiler 구조체 생성자
