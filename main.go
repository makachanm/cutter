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
	debugFlag := flag.Bool("d", false, "Debug Mode")
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

	if *debugFlag {
		fmt.Println(" ----- INSTRUCTIONS -----")

		for i, instr := range vmInstr {
			fmt.Print(i, " ")
			fmt.Println(runtime.ResolveVMInstruction(instr))
		}
	}

	vm := runtime.NewVM(vmInstr)
	vm.Run()

	if *debugFlag {
		fmt.Println(" ----- DATA TABLE -----")
		for i, memdata := range vm.Mem.DataTable {
			fmt.Println(i, ":", memdata)
		}

		fmt.Println(" ----- FUNCTION TABLE -----")
		for i, memdata := range vm.Mem.FunctionTable {
			fmt.Println(i, ":", memdata)
		}

		fmt.Println(" ----- DATA MEMORY -----")
		for i, memdata := range vm.Mem.DataMemory {
			fmt.Println(i, ":", memdata)
		}

		fmt.Println(" ----- FUNCTION MEMORY -----")
		for i, memdata := range vm.Mem.FunctionMemory {
			fmt.Println(i, ":", memdata)
		}

	}

	result := vm.IO.ReadBuffer()

	if *writeToFileFlag != "" {
		etc.WriteFile(*writeToFileFlag, result)
	} else {
		fmt.Println(result)
	}

}
