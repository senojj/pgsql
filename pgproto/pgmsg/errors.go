package pgmsg

import (
	"errors"
)

var (
	ErrUnknownMessageType = errors.New("unknown message type")
	ErrUnknownAuthType    = errors.New("unknown authentication type")
	ErrUnknownCode        = errors.New("unknown code")
	ErrInvalidFormat      = errors.New("invalid format")
	ErrUnexpectedType     = errors.New("unexpected type")
)
