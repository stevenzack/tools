package ioToolkit

import (
	"io"
)

type Pipe struct {
	writer io.Writer
}

func NewPipe(writer io.Writer) *Pipe {
	return &Pipe{
		writer: writer,
	}
}

func (p *Pipe) Read(b []byte) (int, error) {
	return p.writer.Write(b)
}
