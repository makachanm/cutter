package runtime

import "cutter/parser"

func compileCall(call parser.CallObject) []VMInstr {
	var instrs []VMInstr
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
	instrs = append(instrs, VMInstr{Op: OpCall, Arg: call.Name})
	return instrs
}

func CompileASTToVMInstr(ast parser.HeadNode) []VMInstr {
	var instrs []VMInstr
	for _, body := range ast.Bodys {
		switch body.Type {
		case parser.FUCNTION_DEFINITION:
			// 함수 정의는 환경에 등록 (여기선 생략)
		case parser.FUNCTION_CALL:
			instrs = append(instrs, compileCall(body.Call)...)
		case parser.NORM_STRINGS:
			instrs = append(instrs, VMInstr{Op: OpPushString, Arg: body.Norm.Data})
		}
	}
	return instrs
}
