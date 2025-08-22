package runtime

import (
	"cutter/etc"
	"cutter/lexer"
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
		standardFuncs: GetStandardFuncs(),
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
	r.next = 10
}

func (c *Compiler) CompileASTToVMInstr(input parser.HeadNode) []VMInstr {
	instructions := make([]VMInstr, 0)
	c.reg.reset()

	// Pre-pass: handle includes
	newBodys := make([]parser.BodyObject, 0)
	for _, item := range input.Bodys {
		if item.Type == parser.FUNCTION_CALL && item.Call.Name == "include" {
			if len(item.Call.Arguments) != 1 {
				panic("'include' function requires 1 argument: a file path")
			}
			filePathArg := item.Call.Arguments[0]
			if filePathArg.Type != parser.ARG_LITERAL || filePathArg.Literal.Type != parser.STRING {
				panic("'include' function argument must be a string literal")
			}
			filePath := filePathArg.Literal.StringData
			content, err := etc.ReadFile(filePath)
			if err != nil {
				panic(fmt.Sprintf("failed to read file: %s", err))
			}
			lex := lexer.NewLexer()
			tokens := lex.DoLex(content)
			p := parser.NewParser()
			ast := p.DoParse(tokens)

			newBodys = append(newBodys, ast.Bodys...)
		} else {
			newBodys = append(newBodys, item)
		}
	}
	input.Bodys = newBodys

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
			if items.Call.Name == "include" {
				continue
			}
			callInstructions := c.CompileFunctionCallToVMInstr(items.Call, []string{}, len(instructions))
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
	instructions = append(instructions, VMInstr{Op: OpDefFunc, Oprand1: makeStrValueObj(fnc.Name)})

	if fnc.StaticData.Type != 0 {
		tempReg := c.reg.alloc()
		instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(tempReg)), Oprand2: transformToVMDataObject(fnc.StaticData)})
		instructions = append(instructions, VMInstr{Op: OpRslSet, Oprand1: makeIntValueObj(int64(tempReg))})
	} else {
		argNames := fnc.Parameters

		for i, name := range argNames {
			instructions = append(instructions, VMInstr{Op: OpLdr, Oprand1: makeIntValueObj(int64(i)), Oprand2: makeStrValueObj(name)})
		}

		body := fnc.Body
		_, isStandard := c.standardFuncs[body.Name]
		_, isUserFunc := c.funcInfo[body.Name]

		// If the body is not a known function and has no arguments, treat it as a variable lookup.
		if len(body.Arguments) == 0 && !isStandard && !isUserFunc {
			// It's a variable lookup. The value should be loaded and set as the result.
			tempReg := c.reg.alloc()

			// OpLdr loads from memory (where params and global vars are) into a register.
			instructions = append(instructions, VMInstr{Op: OpLdr, Oprand1: makeIntValueObj(int64(tempReg)), Oprand2: makeStrValueObj(body.Name)})
			// OpRslSet sets the result register from another register.
			instructions = append(instructions, VMInstr{Op: OpRslSet, Oprand1: makeIntValueObj(int64(tempReg))})
		} else {
			// It's a real function call.
			bodyInstructions := c.CompileFunctionCallToVMInstr(body, argNames, len(instructions))
			instructions = append(instructions, bodyInstructions...)
		}
	}

	instructions = append(instructions, VMInstr{Op: OpReturn})
	return instructions
}

func (c *Compiler) compileArgument(arg parser.Argument, argNames []string, targetReg int, currentOffset int) []VMInstr {
	instructions := make([]VMInstr, 0)
	switch arg.Type {
	case parser.ARG_LITERAL:
		instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(targetReg)), Oprand2: transformToVMDataObject(arg.Literal)})
	case parser.ARG_VARIABLE:
		isParam := false
		paramIndex := -1
		for j, name := range argNames {
			if name == arg.VarName {
				isParam = true
				paramIndex = j
				break
			}
		}

		if isParam {
			instructions = append(instructions, VMInstr{Op: OpRegMov, Oprand1: makeIntValueObj(int64(paramIndex)), Oprand2: makeIntValueObj(int64(targetReg))})
		} else if _, isUserFunc := c.funcInfo[arg.VarName]; isUserFunc {
			callObj := parser.CallObject{Name: arg.VarName, Arguments: []parser.Argument{}}
			nestedCallInstructions := c.CompileFunctionCallToVMInstr(callObj, argNames, currentOffset+len(instructions))
			instructions = append(instructions, nestedCallInstructions...)
			instructions = append(instructions, VMInstr{Op: OpRslMov, Oprand1: makeIntValueObj(int64(targetReg))})
		} else if _, isStdFunc := c.standardFuncs[arg.VarName]; isStdFunc {
			instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(targetReg)), Oprand2: makeStrValueObj(arg.VarName)})
		} else {
			instructions = append(instructions, VMInstr{Op: OpLdr, Oprand1: makeIntValueObj(int64(targetReg)), Oprand2: makeStrValueObj(arg.VarName)})
		}
	case parser.ARG_CALLABLE:
		nestedCallInstructions := c.CompileFunctionCallToVMInstr(arg.Callable, argNames, currentOffset+len(instructions))
		instructions = append(instructions, nestedCallInstructions...)
		instructions = append(instructions, VMInstr{Op: OpRslMov, Oprand1: makeIntValueObj(int64(targetReg))})
	}
	return instructions
}

func (c *Compiler) CompileFunctionCallToVMInstr(call parser.CallObject, argNames []string, currentOffset int) []VMInstr {
	instructions := make([]VMInstr, 0)

	if call.Name == "for" {
		if len(call.Arguments) != 2 {
			panic("'for' function requires 2 arguments: a condition and a body")
		}

		condArg := call.Arguments[0]
		bodyArg := call.Arguments[1]

		loopStartOffset := currentOffset

		// Compile condition
		condReg := c.reg.alloc()
		condInstructions := c.compileArgument(condArg, argNames, condReg, currentOffset)
		instructions = append(instructions, condInstructions...)

		// Conditional jump to end of loop
		jmpIfFalseInstruction := VMInstr{Op: OpJmpIfFalse, Oprand1: makeIntValueObj(int64(condReg))}
		instructions = append(instructions, jmpIfFalseInstruction)

		// Compile body
		bodyReg := c.reg.alloc()
		bodyInstructions := c.compileArgument(bodyArg, argNames, bodyReg, currentOffset+len(instructions))
		instructions = append(instructions, bodyInstructions...)

		// Unconditional jump back to the start
		instructions = append(instructions, VMInstr{Op: OpJmp, Oprand1: makeIntValueObj(int64(loopStartOffset))})

		// Patch the conditional jump
		loopEndOffset := currentOffset + len(instructions)
		instructions[len(condInstructions)].Oprand2 = makeIntValueObj(int64(loopEndOffset))

		return instructions
	} else if call.Name == "chain" {
		if len(call.Arguments) == 0 {
			return instructions
		}

		intermediateReg := c.reg.alloc()

		// Handle the first argument
		firstArg := call.Arguments[0]
		argInstr := c.compileArgument(firstArg, argNames, intermediateReg, currentOffset+len(instructions))
		instructions = append(instructions, argInstr...)

		// Chain the rest of the arguments
		for i := 1; i < len(call.Arguments); i++ {
			regStateBeforeSubCall := c.reg.next

			nextArg := call.Arguments[i]
			var nextCall parser.CallObject

			if nextArg.Type == parser.ARG_CALLABLE {
				nextCall = nextArg.Callable
			} else if nextArg.Type == parser.ARG_VARIABLE {
				nextCall = parser.CallObject{Name: nextArg.VarName, Arguments: []parser.Argument{}}
			} else {
				panic("chain arguments from the second one must be callable or a function name")
			}

			_, isStandard := c.standardFuncs[nextCall.Name]
			userFunc, isUserFunc := c.funcInfo[nextCall.Name]

			// Compile arguments for the next call
			argRegs := make([]int, len(nextCall.Arguments))
			for j := range nextCall.Arguments {
				argRegs[j] = c.reg.alloc()
			}
			for j, arg := range nextCall.Arguments {
				argCompileInstr := c.compileArgument(arg, argNames, argRegs[j], currentOffset+len(instructions))
				instructions = append(instructions, argCompileInstr...)
			}

			// Pass arguments
			if isStandard {
				instructions = append(instructions, VMInstr{Op: OpRegMov, Oprand1: makeIntValueObj(int64(intermediateReg)), Oprand2: makeIntValueObj(0)})
				for j, reg := range argRegs {
					instructions = append(instructions, VMInstr{Op: OpRegMov, Oprand1: makeIntValueObj(int64(reg)), Oprand2: makeIntValueObj(int64(j + 1))})
				}
			} else if isUserFunc {
				if len(userFunc.Parameters) > 0 {
					paramName := userFunc.Parameters[0]
					instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(paramName), Oprand2: makeIntValueObj(int64(intermediateReg))})
				}
				if len(userFunc.Parameters) < len(argRegs)+1 {
					panic("too many arguments in chain call to " + nextCall.Name)
				}
				for j, reg := range argRegs {
					paramName := userFunc.Parameters[j+1]
					instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(paramName), Oprand2: makeIntValueObj(int64(reg))})
				}
			} else {
				panic("chained function not found: " + nextCall.Name)
			}

			// Perform the call and store the result for the next iteration
			instructions = append(instructions, VMInstr{Op: OpCall, Oprand1: makeStrValueObj(nextCall.Name)})
			instructions = append(instructions, VMInstr{Op: OpRslMov, Oprand1: makeIntValueObj(int64(intermediateReg))})

			c.reg.next = regStateBeforeSubCall
		}

		instructions = append(instructions, VMInstr{Op: OpRslSet, Oprand1: makeIntValueObj(int64(intermediateReg))})

		return instructions
	}

	argRegs := make([]int, len(call.Arguments))
	for i := range call.Arguments {
		argRegs[i] = c.reg.alloc()
	}

	if _, isStandard := c.standardFuncs[call.Name]; isStandard {
		for i := len(call.Arguments) - 1; i >= 0; i-- {
			arg := call.Arguments[i]

			if call.Name == "set" && i == 0 {
				if arg.Type == parser.ARG_VARIABLE {
					instructions = append(instructions, VMInstr{Op: OpRegSet, Oprand1: makeIntValueObj(int64(argRegs[0])), Oprand2: makeStrValueObj(arg.VarName)})
					continue
				} else {
					panic("Compiler error: First argument to 'set' must be a function name.")
				}
			}

			argInstructions := c.compileArgument(arg, argNames, argRegs[i], currentOffset+len(instructions))
			instructions = append(instructions, argInstructions...)
		}
		for i := range call.Arguments {
			instructions = append(instructions, VMInstr{Op: OpRegMov, Oprand1: makeIntValueObj(int64(argRegs[i])), Oprand2: makeIntValueObj(int64(i))})
		}
	} else { // User-defined function
		userFunc, _ := c.funcInfo[call.Name]
		for i, arg := range call.Arguments {
			sname := userFunc.Parameters[i]
			switch arg.Type {
			case parser.ARG_LITERAL:
				instructions = append(instructions, VMInstr{Op: OpMemSet, Oprand1: makeStrValueObj(sname), Oprand2: transformToVMDataObject(arg.Literal)})
			case parser.ARG_VARIABLE:
				tempReg := c.reg.alloc()
				instructions = append(instructions, VMInstr{Op: OpLdr, Oprand1: makeIntValueObj(int64(tempReg)), Oprand2: makeStrValueObj(arg.VarName)})
				instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(sname), Oprand2: makeIntValueObj(int64(tempReg))})
			case parser.ARG_CALLABLE:
				tempReg := c.reg.alloc()
				nestedCallInstructions := c.CompileFunctionCallToVMInstr(arg.Callable, argNames, currentOffset+len(instructions))
				instructions = append(instructions, nestedCallInstructions...)
				instructions = append(instructions, VMInstr{Op: OpRslMov, Oprand1: makeIntValueObj(int64(tempReg))})
				instructions = append(instructions, VMInstr{Op: OpStr, Oprand1: makeStrValueObj(sname), Oprand2: makeIntValueObj(int64(tempReg))})
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
