package runtime

import (
	"fmt"
)

func DumpRegisters(vm *VM) {
	fmt.Println(" ----- REGISTERS -----")
	if len(vm.Reg.ArgumentRegisterMap) <= 0 {
		fmt.Println("REGISTER IS EMPTY")
	}

	for k, reg := range vm.Reg.ArgumentRegisterMap {
		fmt.Printf("R%d: %s ", k, formatVMDataObject(vm.Reg.GetRegister(reg)))
	}
	fmt.Println("\nValue Register: ", formatVMDataObject(vm.Reg.ReturnValueRegister))
	fmt.Println("Register Clear Count: ", vm.Reg.register_cleared_count)
}

func DumpMemory(vm *VM) {
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

func ResolveVMInstruction(instr VMInstr) string {
	opCode := ""
	switch instr.Op {
	case OpRegSet:
		opCode = "OpRegSet"
	case OpMemSet:
		opCode = "OpMemSet"
	case OpRslSet:
		opCode = "OpRslSet"
	case OpRegMov:
		opCode = "OpRegMov"
	case OpMemMov:
		opCode = "OpMemMov"
	case OpRslMov:
		opCode = "OpRslMov"
	case OpLdr:
		opCode = "OpLdr"
	case OpStr:
		opCode = "OpStr"
	case OpRslStr:
		opCode = "OpRslStr"
	case OpDefFunc:
		opCode = "OpDefFunc"
	case OpCall:
		opCode = "OpCall"
	case OpReturn:
		opCode = "OpReturn"
	case OpSyscall:
		opCode = "OpSyscall"
	case OpAdd:
		opCode = "OpAdd"
	case OpSub:
		opCode = "OpSub"
	case OpMul:
		opCode = "OpMul"
	case OpDiv:
		opCode = "OpDiv"
	case OpMod:
		opCode = "OpMod"
	case OpAnd:
		opCode = "OpAnd"
	case OpOr:
		opCode = "OpOr"
	case OpNot:
		opCode = "OpNot"
	case OpCmpEq:
		opCode = "OpCmpEq"
	case OpCmpNeq:
		opCode = "OpCmpNeq"
	case OpBrch:
		opCode = "OpBrch"
	case OpClearReg:
		opCode = "OpClearReg"
	case OpHlt:
		opCode = "OpHlt"
	default:
		opCode = fmt.Sprintf("UnknownOp(%d)", instr.Op)
	}

	// Format operands
	oprand1Str := formatVMDataObject(instr.Oprand1)
	oprand2Str := formatVMDataObject(instr.Oprand2)
	oprand3Str := formatVMDataObject(instr.Oprand3)

	return fmt.Sprintf("%s %s %s %s", opCode, oprand1Str, oprand2Str, oprand3Str)
}

func formatVMDataObject(obj VMDataObject) string {
	switch obj.Type {
	case INTGER:
		return fmt.Sprintf("INT(%d)", obj.IntData)
	case REAL:
		return fmt.Sprintf("REAL(%f)", obj.FloatData)
	case STRING:
		if obj.StringData == "\n" {
			return fmt.Sprintf("STR(%s)", "newline")
		}
		return fmt.Sprintf("STR(%s)", obj.StringData)
	case BOOLEAN:
		return fmt.Sprintf("BOOL(%t)", obj.BoolData)
	default:
		return "EMPTY"
	}
}
