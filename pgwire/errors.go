package pgwire

import "errors"

var (
	ErrInvalidFormat  = errors.New("invalid format")
	ErrUnexpectedType = errors.New("unexpected type")
)
