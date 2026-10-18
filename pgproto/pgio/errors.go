package pgio

import "errors"

var (
	ErrValueUnderflow = errors.New("partial value")
	ErrValueOverflow  = errors.New("value too large")
)
