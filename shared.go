package pgsql

import (
	"fmt"
)

func invalidFormat(cause error) error {
	return fmt.Errorf("%w: %w", ErrInvalidFormat, cause)
}

func unexpectedType(got, want byte) error {
	return fmt.Errorf("%w: got '%d', want '%d'", ErrUnexpectedType, got, want)
}

func unexpectedAuth(got, want int32) error {
	return fmt.Errorf("%w: got '%d', want '%d'", ErrUnexpectedType, got, want)
}

func shiftLength(in []byte) ([]byte, error) {
	length, b, err := shiftInt32(in)
	if err != nil {
		return in, err
	}

	if length < 4 {
		return in, ErrValueUnderflow
	}

	size := int(length) - 4

	if size > len(b) {
		return in, ErrValueUnderflow
	}
	return b[:size], nil
}

func shiftHeader(msg byte, in []byte) ([]byte, error) {
	kind, b, err := shiftByte(in)
	if err != nil {
		return in, err
	}

	if msg != kind {
		return in, unexpectedType(kind, msg)
	}

	return shiftLength(b)
}
