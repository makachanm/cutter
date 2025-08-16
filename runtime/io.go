package runtime

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TODO: File IO

type RuntimeIO struct {
	buffer []string
	writer io.Writer // New field to write to
}

func NewIO() RuntimeIO {
	return RuntimeIO{
		buffer: make([]string, 0),
		writer: nil, // Default to nil, meaning no direct output
	}
}

// NewIOWithWriter creates a RuntimeIO that writes to the given writer.
func NewIOWithWriter(w io.Writer) RuntimeIO {
	return RuntimeIO{
		buffer: make([]string, 0),
		writer: w,
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
	// Removed io.FlushIO() from here. FlushIO will be called explicitly by syscall.
}

// FlushIO writes the buffered content to the configured writer and clears the buffer.
func (io *RuntimeIO) FlushIO() {
	if io.writer != nil {
		for _, elem := range io.buffer {
			fmt.Fprint(io.writer, elem)
		}
	}
	io.buffer = io.buffer[:0] // Clear the buffer
}

// ReadBuffer returns the current content of the buffer without clearing it.
// This is primarily for testing purposes to inspect the buffer before it's flushed.
func (io *RuntimeIO) ReadBuffer() string {
	return strings.Join(io.buffer, "")
}
