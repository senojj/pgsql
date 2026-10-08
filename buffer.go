package pgsql

import (
	"bytes"
	"slices"
)

type buffer struct {
	data []byte
}

func (buf *buffer) Bytes() []byte {
	return buf.data
}

func (buf *buffer) Len() int {
	return len(buf.data)
}

func (buf *buffer) Cap() int {
	return cap(buf.data)
}

func (buf *buffer) Count(b []byte) int {
	return bytes.Count(buf.data, b)
}

func (buf *buffer) Grow(n int) {
	buf.data = slices.Grow(buf.data, n)
}

func (buf *buffer) ShiftByte() (value byte, err error) {
	value, buf.data, err = shiftByte(buf.data)
	return
}

func (buf *buffer) ShiftBytes(length int) (value []byte, err error) {
	value, buf.data, err = shiftBytes(buf.data, length)
	return
}

func (buf *buffer) ShiftInt8() (value int8, err error) {
	value, buf.data, err = shiftInt8(buf.data)
	return
}

func (buf *buffer) ShiftInt16() (value int16, err error) {
	value, buf.data, err = shiftInt16(buf.data)
	return
}

func (buf *buffer) ShiftInt32() (value int32, err error) {
	value, buf.data, err = shiftInt32(buf.data)
	return
}

func (buf *buffer) ShiftString() (value string, err error) {
	value, buf.data, err = shiftString(buf.data)
	return
}

func (buf *buffer) AppendByte(i ...byte) {
	buf.data = appendByte(buf.data, i...)
}

func (buf *buffer) AppendInt8(i ...int8) {
	buf.data = appendInt8(buf.data, i...)
}

func (buf *buffer) AppendInt16(i ...int16) {
	buf.data = appendInt16(buf.data, i...)
}

func (buf *buffer) AppendInt32(i ...int32) {
	buf.data = appendInt32(buf.data, i...)
}

func (buf *buffer) AppendInt64(i ...int64) {
	buf.data = appendInt64(buf.data, i...)
}

func (buf *buffer) AppendString(s ...string) {
	buf.data = appendString(buf.data, s...)
}
