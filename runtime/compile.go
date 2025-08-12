package runtime

import (
	"fmt"
)

/* Co-generated with GPT-4.1 */

// Compiler: Cutter 컴파일러 상태 및 메서드 집합
// regAlloc 등 상태를 포함

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

// compileValue: ValueObject를 임시 변수에 OpSet으로 저장

func (c *Compiler) CompileASTToVMInstr() []VMInstr {

}

// NewCompiler: Compiler 구조체 생성자
