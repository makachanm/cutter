package runtime

import (
	"cutter/parser"
	"fmt"
	"strconv"
)

// TODO: File IO

type RuntimeIO struct {
	buffer []string
}

func NewIO() RuntimeIO {
	return RuntimeIO{
		buffer: make([]string, 0),
	}
}

func (io *RuntimeIO) WriteObjectToStream(data parser.ValueObject) {
	switch data.Type {
	case parser.STRING:
		io.buffer = append(io.buffer, data.StringData)
	case parser.INTGER:
		io.buffer = append(io.buffer, strconv.FormatInt(data.IntData, 10))
	case parser.REAL:
		io.buffer = append(io.buffer, strconv.FormatFloat(data.FloatData, 'f', -1, 64))
	case parser.BOOLEAN:
		var d string
		if data.BoolData {
			d = "!t"
		} else {
			d = "!f"
		}

		io.buffer = append(io.buffer, d)
	}
	io.FlushIO()
}

func (io *RuntimeIO) WriteNorm(data parser.NormStringObject) {
	io.buffer = append(io.buffer, data.Data)
	io.FlushIO()
}

func (io *RuntimeIO) FlushIO() {
	for elem := range io.buffer {
		fmt.Print(elem)
	}
	io.buffer = io.buffer[:0]
}
