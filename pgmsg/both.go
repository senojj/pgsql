package pgmsg

import (
	"math"

	"github.com/senojj/pgsql/pgio"
)

type CopyData struct {
	Data []byte
}

func (x *CopyData) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)
	length := sizeMessageLength + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCopyData)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *CopyData) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCopyData, b)
	if err != nil {
		return invalidFormat(err)
	}

	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type CopyDone struct{}

func (x *CopyDone) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCopyDone)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *CopyDone) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCopyDone, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}
