package runtime

import (
	"cutter/parser"
	"fmt"
)

type Compiler struct {
	reg           *regAlloc
	funcInfo      map[string]parser.FunctionObject
	standardFuncs map[string][]VMInstr
}

func NewCompiler() *Compiler {
	return &Compiler{
		reg:           &regAlloc{},
		funcInfo:      make(map[string]parser.FunctionObject),
		standardFuncs: StandardFuncs,
	}
}

type regAlloc struct {
	next int
}

func (r *regAlloc) alloc() int {
	idx := r.next
	r.next++
	return idx
}

func (r *regAlloc) reset() {
	r.next = 0
}

func (c *Compiler) CompileASTToVMInstr(input parser.HeadNode) []VMInstr {
	instructions := make([]VMInstr, 0)
	c.reg.reset()

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
			// Function definitions are handled by the pre-scan, but we need to generate code for them.
			// However, the VM will handle skipping over this code during execution.
			// We just need to make sure the instructions are there.
			instructions = append(instructions, c.CompileFunctionDefToVMInstr(items.Func)...)

		case parser.FUNCTION_CALL:
			callInstructions := c.CompileFunctionCallToVMInstr(items.Call, nil)
			instructions = append(instructions, callInstructions...)
			// After a top-level call, store the result in stdout
			instructions = append(instructions, VMInstr{Op: OpRslStr, Oprand1: makeStrValueObj("stdout")})
			instructions = append(instructions, VMInstr{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)})

			c.reg.reset()
			instructions = append(instructions, VMInstr{Op: OpClearReg})
		case parser.NORM_STRINGS:
			tmpReg := c.reg.alloc()
			instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(tmpReg)), Oprand2: makeStrValueObj(items.Norm.Data)})
			instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj("stdout"), Oprand2: makeIntValueObj(int64(tmpReg))})
			instructions = append(instructions, VMInstr{Op: OpSyscall, Oprand1: makeIntValueObj(SYS_IO_FLUSH)})

			c.reg.reset()
			instructions = append(instructions, VMInstr{Op: OpClearReg})
		}
	}
	return instructions
}

func (c *Compiler) CompileFunctionDefToVMInstr(fnc parser.FunctionObject) []VMInstr {
	if _, exists := c.standardFuncs[fnc.Name]; exists {
		panic("Standard function already exists: " + fnc.Name)
	}

	instructions := make([]VMInstr, 0)
	c.reg.reset()

	// Define the function entry point
	instructions = append(instructions, VMInstr{Op: OpDefFunc, Oprand1: makeStrValueObj(fnc.Name)})

	// If there are no arguments, it's a simple value definition
	if len(fnc.Args) == 0 {
		tempReg := c.reg.alloc()
		instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(tempReg)), Oprand2: transformToVMDataObject(fnc.StaticData)})
		instructions = append(instructions, VMInstr{Op: OpRslSet, Oprand1: makeIntValueObj(int64(tempReg))})
	} else {
		// The last argument is the function body
		body := fnc.Args[len(fnc.Args)-1]
		argNames := make([]string, len(fnc.Args)-1)
		for i, arg := range fnc.Args[:len(fnc.Args)-1] {
			argNames[i] = arg.Name
		}

		// Load arguments from memory into registers
		for i, name := range argNames {
			instructions = append(instructions, VMInstr{Op: OpLdr, Oprand1: makeIntValueObj(int64(i)), Oprand2: makeStrValueObj(name)})
		}

		// Compile the body of the function
		bodyInstr := c.CompileFunctionCallToVMInstr(body, &argNames)
		instructions = append(instructions, bodyInstr...)
	}

	instructions = append(instructions, VMInstr{Op: OpReturn})
	return instructions
}

func (c *Compiler) CompileFunctionCallToVMInstr(call parser.CallObject, localVars *[]string) []VMInstr {
	instructions := make([]VMInstr, 0)

	// Check if it's a standard function
	if _, isStandard := c.standardFuncs[call.Name]; isStandard {
		// Handle callable arguments first
		for i, carg := range call.CallableArgs {
			nestedCallInstructions := c.CompileFunctionCallToVMInstr(carg, localVars)
			instructions = append(instructions, nestedCallInstructions...)
			// Move result to the correct argument register
			instructions = append(instructions, VMInstr{Op: OpRslMov, Oprand1: makeIntValueObj(int64(len(call.Args) + i))})
		}
		// Handle literal arguments
		for i, arg := range call.Args {
			instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(i)), Oprand2: transformToVMDataObject(arg)})
		}
	} else { // User-defined function
		funcInfo, isUserDefined := c.funcInfo[call.Name]
		if !isUserDefined {
			// It's not a function, so it must be a variable.
			// Check if it's a local variable
			isLocal := false
			if localVars != nil {
				for i, name := range *localVars {
					if name == call.Name {
						// It's a local variable, which is already in a register.
						// We just need to move it to the result register.
						instructions = append(instructions, VMInstr{Op: OpRslSet, Oprand1: makeIntValueObj(int64(i))})
						isLocal = true
						break
					}
				}
			}
			if !isLocal {
				// It's a global variable (a function with no args)
				// We just call it.
				instructions = append(instructions, VMInstr{Op: OpCall, Oprand1: makeStrValueObj(call.Name)})
			}
			return instructions
		}

		// Handle arguments for user-defined functions
		for i, arg := range call.Args {
			if i < len(funcInfo.Args)-1 {
				paramName := funcInfo.Args[i].Name
				instructions = append(instructions, VMInstr{Op: OpMemSet, Oprand1: makeStrValueObj(paramName)})
				tmpReg := c.reg.alloc()
				instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(tmpReg)), Oprand2: transformToVMDataObject(arg)})
				instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(paramName), Oprand2: makeIntValueObj(int64(tmpReg))})
			}
		}
		for i, carg := range call.CallableArgs {
			if i+len(call.Args) < len(funcInfo.Args)-1 {
				paramName := funcInfo.Args[i+len(call.Args)].Name
				instructions = append(instructions, VMInstr{Op: OpMemSet, Oprand1: makeStrValueObj(paramName)})
				nestedCallInstructions := c.CompileFunctionCallToVMInstr(carg, localVars)
				instructions = append(instructions, nestedCallInstructions...)
				tmpReg := c.reg.alloc()
				instructions = append(instructions, VMInstr{Op: OpRslMov, Oprand1: makeIntValueObj(int64(tmpReg))})
				instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(paramName), Oprand2: makeIntValueObj(int64(tmpReg))})
			}
		}
	}

	// Perform the call
	instructions = append(instructions, VMInstr{Op: OpCall, Oprand1: makeStrValueObj(call.Name)})

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
		panic(fmt.Sprintf("Unknown ValueType: %d", val.Type))
	}
}
