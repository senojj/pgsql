package pgsql

import (
	"math"
)

type msgCopyData struct {
	Data []byte
}

func (x *msgCopyData) message() {}

func (x *msgCopyData) frontend() {}

func (x *msgCopyData) backend() {}

func (x *msgCopyData) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)
	length := sizeMessageLength + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCopyData)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *msgCopyData) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyData, b)
	if err != nil {
		return invalidFormat(err)
	}

	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type msgCopyDone struct{}

func (x *msgCopyDone) message() {}

func (x *msgCopyDone) frontend() {}

func (x *msgCopyDone) backend() {}

func (x *msgCopyDone) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCopyDone)
	buf.AppendInt32(int32(length))
	return buf.Bytes(), nil
}

func (x *msgCopyDone) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyDone, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}
