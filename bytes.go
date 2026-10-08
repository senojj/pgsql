package pgsql

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"
)

func shiftByte(b []byte) (byte, []byte, error) {
	if len(b) == 0 {
		return 0, b, io.EOF
	}
	return b[0], b[1:], nil
}

func shiftBytes(b []byte, length int) ([]byte, []byte, error) {
	if len(b) < length {
		return nil, b, ErrValueUnderflow
	}
	output := make([]byte, length)
	copy(output, b[:length])
	return output, b[length:], nil
}

func shiftInt8(b []byte) (int8, []byte, error) {
	v, b, err := shiftByte(b)
	if err != nil {
		return 0, b, err
	}
	return int8(v), b, nil
}

func shiftInt16(b []byte) (int16, []byte, error) {
	if len(b) < 2 {
		return 0, b, ErrValueUnderflow
	}
	i := int16(binary.BigEndian.Uint16(b))
	return i, b[2:], nil
}

func shiftInt32(b []byte) (int32, []byte, error) {
	if len(b) < 4 {
		return 0, b, ErrValueUnderflow
	}
	i := int32(binary.BigEndian.Uint32(b))
	return i, b[4:], nil
}

var zero []byte = []byte{0}

func shiftString(b []byte) (string, []byte, error) {
	s, b, found := bytes.Cut(b, zero)
	if found {
		return unsafe.String(unsafe.SliceData(s), len(s)), b, nil
	}
	return "", b, ErrValueUnderflow
}

func appendByte(b []byte, i ...byte) []byte {
	return append(b, i...)
}

func appendInt8(b []byte, i ...int8) []byte {
	for _, v := range i {
		b = append(b, byte(v))
	}
	return b
}

func appendInt16(b []byte, i ...int16) []byte {
	for _, v := range i {
		b = binary.BigEndian.AppendUint16(b, uint16(v))
	}
	return b
}

func appendInt32(b []byte, i ...int32) []byte {
	for _, v := range i {
		b = binary.BigEndian.AppendUint32(b, uint32(v))
	}
	return b
}

func appendInt64(b []byte, i ...int64) []byte {
	for _, v := range i {
		b = binary.BigEndian.AppendUint64(b, uint64(v))
	}
	return b
}

func appendString(b []byte, s ...string) []byte {
	for _, v := range s {
		b = append(b, v...)
		b = append(b, 0)
	}
	return b
}
