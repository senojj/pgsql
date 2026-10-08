package pgsql

import "encoding"

type message interface {
	encoding.BinaryAppender
	encoding.BinaryUnmarshaler

	message()
}

type backend interface {
	message

	backend()
}

type frontend interface {
	message

	frontend()
}
