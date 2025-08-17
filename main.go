package main

import (
	"cutter/etc"
	"cutter/lexer"
	"cutter/parser"
	"cutter/runtime"
	"flag"
	"fmt"
)

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
	vmInstr := com.CompileASTToVMInstr(ast)

	vm := runtime.NewVM(vmInstr)
	vm.Run()

	result := vm.IO.ReadBuffer()

	if *writeToFileFlag != "" {
		etc.WriteFile(*writeToFileFlag, result)
	} else {
		fmt.Println(result)
	}

}
