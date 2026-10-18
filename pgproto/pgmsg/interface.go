package pgmsg

import (
	"encoding"
)

type Backend interface {
	encoding.BinaryUnmarshaler

	backend()
}

type Frontend interface {
	encoding.BinaryAppender

	frontend()
}
