package runtime

import (
	"cutter/parser"
	"fmt"
)

func compileCall(call parser.CallObject) []VMInstr {
	var instrs []VMInstr
	// 값 인자 push
	for _, arg := range call.Args {
		switch arg.Type {
		case parser.INTGER:
			instrs = append(instrs, VMInstr{Op: OpPushInt, Arg: arg.IntData})
		case parser.REAL:
			instrs = append(instrs, VMInstr{Op: OpPushReal, Arg: arg.FloatData})
		case parser.STRING:
			instrs = append(instrs, VMInstr{Op: OpPushString, Arg: arg.StringData})
		case parser.BOOLEAN:
			instrs = append(instrs, VMInstr{Op: OpPushBool, Arg: arg.BoolData})
		}
	}
	// CallableArgs(중첩 함수 호출)도 재귀적으로 변환
	for _, subcall := range call.CallableArgs {
		instrs = append(instrs, compileCall(subcall)...)
	}
	// 함수 호출
	instrs = append(instrs, VMInstr{Op: OpCall, Arg: call.Name})
	return instrs
}

func compileFuncDef(fun parser.FunctionObject) []VMInstr {
	// 함수 정의를 OpSetFunc으로 등록
	// 인자 이름 추출
	argNames := make([]string, 0, len(fun.Args))
	for _, arg := range fun.Args {
		argNames = append(argNames, arg.Name)
	}
	// 함수 본문(여기선 FuncBodys) 컴파일
	bodyInstrs := compileCall(fun.FuncBodys)
	return []VMInstr{
		{
			Op:   OpSetFunc,
			Arg:  fun.Name,
			Args: []interface{}{argNames, bodyInstrs},
		},
	}
}

func CompileASTToVMInstr(ast parser.HeadNode) []VMInstr {
	var instrs []VMInstr
	for _, body := range ast.Bodys {
		switch body.Type {
		case parser.FUCNTION_DEFINITION:
			instrs = append(instrs, compileFuncDef(body.Func)...)
		case parser.FUNCTION_CALL:
			instrs = append(instrs, compileCall(body.Call)...)
		case parser.NORM_STRINGS:
			instrs = append(instrs, VMInstr{Op: OpPushString, Arg: body.Norm.Data})
		default:
			fmt.Println("Unknown BodyObject type in CompileASTToVMInstr")
		}
	}
	return instrs
}
