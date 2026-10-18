package pgio

import (
	"bytes"
	"slices"
)

type Buffer struct {
	Bytes []byte
}

func (buf *Buffer) Len() int {
	return len(buf.Bytes)
}

func (buf *Buffer) Cap() int {
	return cap(buf.Bytes)
}

func (buf *Buffer) Count(b []byte) int {
	return bytes.Count(buf.Bytes, b)
}

func (buf *Buffer) Grow(n int) {
	buf.Bytes = slices.Grow(buf.Bytes, n)
}

func (buf *Buffer) ShiftByte() (value byte, err error) {
	value, buf.Bytes, err = ShiftByte(buf.Bytes)
	return
}

func (buf *Buffer) ShiftBytes(length int) (value []byte, err error) {
	value, buf.Bytes, err = ShiftBytes(buf.Bytes, length)
	return
}

func (buf *Buffer) ShiftInt8() (value int8, err error) {
	value, buf.Bytes, err = ShiftInt8(buf.Bytes)
	return
}

func (buf *Buffer) ShiftInt16() (value int16, err error) {
	value, buf.Bytes, err = ShiftInt16(buf.Bytes)
	return
}

func (buf *Buffer) ShiftInt32() (value int32, err error) {
	value, buf.Bytes, err = ShiftInt32(buf.Bytes)
	return
}

func (buf *Buffer) ShiftString() (value string, err error) {
	value, buf.Bytes, err = ShiftString(buf.Bytes)
	return
}

func (buf *Buffer) AppendByte(i ...byte) {
	buf.Bytes = AppendByte(buf.Bytes, i...)
}

func (buf *Buffer) AppendInt8(i ...int8) {
	buf.Bytes = AppendInt8(buf.Bytes, i...)
}

func (buf *Buffer) AppendInt16(i ...int16) {
	buf.Bytes = AppendInt16(buf.Bytes, i...)
}

func (buf *Buffer) AppendInt32(i ...int32) {
	buf.Bytes = AppendInt32(buf.Bytes, i...)
}

func (buf *Buffer) AppendInt64(i ...int64) {
	buf.Bytes = AppendInt64(buf.Bytes, i...)
}

func (buf *Buffer) AppendString(s ...string) {
	buf.Bytes = AppendString(buf.Bytes, s...)
}
