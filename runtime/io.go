package runtime

import (
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

func (io *RuntimeIO) WriteObjectToStream(data VMDataObject) {
	switch data.Type {
	case STRING:
		io.buffer = append(io.buffer, data.StringData)
	case INTGER:
		io.buffer = append(io.buffer, strconv.FormatInt(data.IntData, 10))
	case REAL:
		io.buffer = append(io.buffer, strconv.FormatFloat(data.FloatData, 'f', -1, 64))
	case BOOLEAN:
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

func (io *RuntimeIO) FlushIO() {
	for _, elem := range io.buffer {
		fmt.Print(elem)
	}
	io.buffer = io.buffer[:0]
}