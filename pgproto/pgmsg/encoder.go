package pgmsg

import (
	"io"
)

type Encoder struct {
	writer io.Writer
	buf    []byte
}

func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{
		writer: writer,
		buf:    make([]byte, 0),
	}
}

func (e *Encoder) Encode(msg Frontend) error {
	var err error
	e.buf, err = msg.AppendBinary(e.buf)
	if err != nil {
		return err
	}
	return nil
}

func (e *Encoder) Flush() error {
	if len(e.buf) == 0 {
		return nil
	}

	defer func() {
		e.buf = e.buf[:0]
	}()

	_, err := e.writer.Write(e.buf)
	if err != nil {
		return err
	}
	return nil
}
