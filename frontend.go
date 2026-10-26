package pgsql

import (
	"bytes"
	"fmt"
	"net"

	"github.com/senojj/pgsql/pgmsg"
)

const readBufferSize int = 1 << 13 // 8kb

func rtoerr(msg *pgmsg.ErrorResponse) error {
	idxSeverity := bytes.IndexByte(msg.Fields, FieldSeverity)
	idxCode := bytes.IndexByte(msg.Fields, FieldCode)
	idxMessage := bytes.IndexByte(msg.Fields, FieldMessage)
	severity := msg.Values[idxSeverity]
	code := msg.Values[idxCode]
	message := msg.Values[idxMessage]

	return fmt.Errorf(
		"severity: %s, code: %s, message: %s",
		severity,
		code,
		message,
	)
}

type Frontend struct {
	raw     net.Conn
	encoder *pgmsg.Encoder
	decoder *pgmsg.Decoder
}

func NewFrontend(conn net.Conn) *Frontend {
	return &Frontend{
		raw:     conn,
		encoder: pgmsg.NewEncoder(conn),
		decoder: pgmsg.NewDecoder(conn),
	}
}

func (f *Frontend) Close() error {
	return f.raw.Close()
}
