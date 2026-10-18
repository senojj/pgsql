package pgsql

import "encoding"

type message interface {
	encoding.BinaryAppender
	encoding.BinaryUnmarshaler
}
