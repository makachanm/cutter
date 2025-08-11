package runtime

import (
	"cutter/parser"
	"fmt"
	"strconv"
)

// TODO: File IO

type RuntimeIO struct {
	buffer chan string
}

func NewIO() RuntimeIO {
	return RuntimeIO{
		buffer: make(chan string),
	}
}

func (io *RuntimeIO) WriteObjectToStream(data Function) {
	switch data.Body.ValueBodys.Type {
	case parser.STRING:
		io.buffer <- data.Body.ValueBodys.StringData
	case parser.INTGER:
		io.buffer <- strconv.FormatInt(data.Body.ValueBodys.IntData, 10)
	case parser.REAL:
		io.buffer <- strconv.FormatFloat(data.Body.ValueBodys.FloatData, 'f', -1, 64)
	case parser.BOOLEAN:
		var d string
		if data.Body.ValueBodys.BoolData {
			d = "!t"
		} else {
			d = "!f"
		}

		io.buffer <- d
	}
}

func (io *RuntimeIO) FlushIO() {
	for elem := range io.buffer {
		fmt.Print(elem)
	}
}
