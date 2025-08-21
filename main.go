package main

import (
	"cutter/etc"
	"cutter/lexer"
	"cutter/parser"
	"cutter/runtime"
	"flag"
	"fmt"
	"path/filepath"
)

func preprocessIncludes(ast parser.HeadNode, currentPath string) (parser.HeadNode, error) {
	newBodys := make([]parser.BodyObject, 0)
	for _, body := range ast.Bodys {
		if body.Type == parser.FUNCTION_CALL && body.Call.Name == "include" {
			if len(body.Call.Arguments) != 1 {
				return ast, fmt.Errorf("include function requires exactly one argument")
			}
			arg := body.Call.Arguments[0]
			if arg.Type != parser.ARG_LITERAL || arg.Literal.Type != parser.STRING {
				return ast, fmt.Errorf("include function argument must be a string literal")
			}

			includePath := arg.Literal.StringData
			if !filepath.IsAbs(includePath) {
				includePath = filepath.Join(filepath.Dir(currentPath), includePath)
			}

			source, err := etc.ReadFile(includePath)
			if err != nil {
				return ast, fmt.Errorf("failed to read include file: %w", err)
			}

			lex := lexer.NewLexer()
			pas := parser.NewParser()

			tokens := lex.DoLex(source)
			includedAst := pas.DoParse(tokens)

			processedIncludedAst, err := preprocessIncludes(includedAst, includePath)
			if err != nil {
				return ast, err
			}

			newBodys = append(newBodys, processedIncludedAst.Bodys...)
		} else {
			newBodys = append(newBodys, body)
		}
	}
	return parser.HeadNode{Bodys: newBodys}, nil
}


func main() {
	versionFlag := flag.Bool("v", false, "Show Version")
	writeToFileFlag := flag.String("w", "", "Write excution result to file")
	input := flag.String("i", "", "Input file")

	flag.Parse()

	if *versionFlag {
		fmt.Println("Cutter Runtime Version: " + etc.RUNTIMEVERSION)
		return
	}

	if *input == "" {
		fmt.Println("No input file specified")
		return
	}

	source, err := etc.ReadFile(*input)
	if err != nil {
		panic(err)
	}

	lex := lexer.NewLexer()
	pas := parser.NewParser()
	com := runtime.NewCompiler()

	tokens := lex.DoLex(source)
	ast := pas.DoParse(tokens)

	absPath, err := filepath.Abs(*input)
	if err != nil {
		panic(err)
	}

	processedAst, err := preprocessIncludes(ast, absPath)
	if err != nil {
		panic(err)
	}

	vmInstr := com.CompileASTToVMInstr(processedAst)

	vm := runtime.NewVM(vmInstr)
	vm.Run()

	result := vm.IO.ReadBuffer()

	if *writeToFileFlag != "" {
		etc.WriteFile(*writeToFileFlag, result)
	} else {
		fmt.Println(result)
	}

}
