package runtime

const (
	INTGER ValueType = iota + 1
	REAL
	STRING
	BOOLEAN
)

type VMDataObject struct {
	Type ValueType

	IntData    int64
	FloatData  float64
	BoolData   bool
	StringData string
}

type VMFunctionObject struct {
	JumpPc int
}

type VMArgumentRegisters struct {
	ArgumentRegisters []VMDataObject
}

func NewRegister() VMArgumentRegisters {
	return VMArgumentRegisters{
		ArgumentRegisters: make([]VMDataObject, 32),
	}
}

func (rg *VMArgumentRegisters) ClearRegisters() {
	rg.ArgumentRegisters = rg.ArgumentRegisters[:0]
	rg.ArgumentRegisters = make([]VMDataObject, 32)
}

func (rg *VMArgumentRegisters) InsertRegister(idx int, val VMDataObject) {
	rg.ArgumentRegisters[idx] = val
}

func (rg *VMArgumentRegisters) GetRegister(idx int) VMDataObject {
	return rg.ArgumentRegisters[idx]
}
