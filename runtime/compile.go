package runtime

import (
	"cutter/parser"
	"fmt"
)

type Compiler struct {
	reg      *regAlloc
	last_pos int
	funcInfo map[string]parser.FunctionObject
}

func NewCompiler() *Compiler {
	return &Compiler{reg: &regAlloc{}, last_pos: 0, funcInfo: make(map[string]parser.FunctionObject)}
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

	// First pass: gather all function definitions
	for _, items := range input.Bodys {
		if items.Type == parser.FUCNTION_DEFINITION {
			c.funcInfo[items.Func.Name] = items.Func
		}
	}

	// Second pass: compile instructions
	for _, items := range input.Bodys {
		switch items.Type {
		case parser.FUCNTION_DEFINITION:
			instructions = append(instructions, c.CompileFunctionDefToVMInstr(items.Func)...)
		case parser.FUNCTION_CALL:
			callInstructions := c.CompileFunctionCallToVMInstr(items.Call)
			instructions = append(instructions, callInstructions...)
			// After a top-level call, store the result in stdout
			resultReg := c.reg.next - 1 // The result is in the last allocated register
			instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj("stdout"), Oprand2: makeIntValueObj(int64(resultReg))})
			instructions = append(instructions, VMInstr{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)})
		case parser.NORM_STRINGS:
			tmpReg := c.reg.alloc()
			instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(tmpReg)), Oprand2: makeStrValueObj(items.Norm.Data)})
			instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj("stdout"), Oprand2: makeIntValueObj(int64(tmpReg))})
			instructions = append(instructions, VMInstr{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)})
		}
	}
	return instructions
}

func (c *Compiler) CompileFunctionDefToVMInstr(fnc parser.FunctionObject) []VMInstr {
	instructions := make([]VMInstr, 0)

	instructions = append(instructions, VMInstr{Op: OpDefFunc, Oprand1: makeStrValueObj(fnc.Name), Oprand2: makeIntValueObj(int64(c.last_pos))})

	// The function body only contains the return value expression
	if len(fnc.Args) > 0 {
		returnValue := fnc.Args[len(fnc.Args)-1]
		if returnValue.Name != "" { // It's a variable
			tmpReg := c.reg.alloc()
			instructions = append(instructions, VMInstr{Op: OpLdr, Oprand1: makeIntValueObj(int64(tmpReg)), Oprand2: makeStrValueObj(returnValue.Name)})
			instructions = append(instructions, VMInstr{Op: OpRslSet, Oprand1: makeIntValueObj(int64(tmpReg))})
		} else { // It's a nested function call
			instructions = append(instructions, c.CompileFunctionCallToVMInstr(returnValue)...)
		}
	}

	instructions = append(instructions, VMInstr{Op: OpReturn})
	c.last_pos += len(instructions)
	return instructions
}

func (c *Compiler) CompileFunctionCallToVMInstr(call parser.CallObject) []VMInstr {
	instructions := make([]VMInstr, 0)
	c.reg.next = 0 // Reset registers for each call

	funcInfo, ok := c.funcInfo[call.Name]
	if !ok {
		// Handle standard library functions or error
	} else {
		// Create and populate memory slots for parameters
		for i, arg := range call.Args {
			if i < len(funcInfo.Args)-1 {
				paramName := funcInfo.Args[i].Name
				// Create memory slot
				instructions = append(instructions, VMInstr{Op: OpMemSet, Oprand1: makeStrValueObj(paramName), Oprand2: VMDataObject{}})
				// Load argument into register and then store it in the memory slot
				tmpReg := c.reg.alloc()
				instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(tmpReg)), Oprand2: transformToVMDataObject(arg)})
				instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(paramName), Oprand2: makeIntValueObj(int64(tmpReg))})
			}
		}
		for i, carg := range call.CallableArgs {
			if i+len(call.Args) < len(funcInfo.Args)-1 {
				paramName := funcInfo.Args[i+len(call.Args)].Name
				// Create memory slot
				instructions = append(instructions, VMInstr{Op: OpMemSet, Oprand1: makeStrValueObj(paramName), Oprand2: VMDataObject{}})
				// Load argument into register and then store it in the memory slot
				tmpReg := c.reg.alloc()
				instructions = append(instructions, VMInstr{Op: OpLdr, Oprand1: makeIntValueObj(int64(tmpReg)), Oprand2: makeStrValueObj(carg.Name)})
				instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(paramName), Oprand2: makeIntValueObj(int64(tmpReg))})
			}
		}
	}

	// Perform the call
	instructions = append(instructions, VMInstr{Op: OpCall, Oprand1: makeStrValueObj(call.Name)})

	// Move result from result register to a general-purpose register
	tmpReg := c.reg.alloc()
	instructions = append(instructions, VMInstr{Op: OpRelMov, Oprand1: makeIntValueObj(int64(tmpReg))})

	c.last_pos += len(instructions)
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

func makeRealValueObj(f float64) VMDataObject {
	return VMDataObject{
		Type:      REAL,
		FloatData: f,
	}
}

func makeBoolValueObj(b bool) VMDataObject {
	return VMDataObject{
		Type:     BOOLEAN,
		BoolData: b,
	}
}

func transformToVMDataObject(val parser.ValueObject) VMDataObject {
	switch val.Type {
	case parser.INTGER:
		return makeIntValueObj(val.IntData)
	case parser.REAL:
		return makeRealValueObj(val.FloatData)
	case parser.STRING:
		return makeStrValueObj(val.StringData)
	case parser.BOOLEAN:
		return makeBoolValueObj(val.BoolData)
	default:
		panic("Unknown ValueType")
	}
}