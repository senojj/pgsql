package pgmsg

import (
	"math"

	"github.com/senojj/pgsql/pgproto/pgio"
)

type CopyData struct {
	Data []byte
}

func (_ *CopyData) frontend() {}

func (_ *CopyData) backend() {}

func (x *CopyData) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)
	length := sizeMessageLength + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeCopyData)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *CopyData) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyData, b)
	if err != nil {
		return invalidFormat(err)
	}

	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type CopyDone struct{}

func (_ *CopyDone) frontend() {}

func (_ *CopyDone) backend() {}

func (x *CopyDone) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeCopyDone)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *CopyDone) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyDone, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}
