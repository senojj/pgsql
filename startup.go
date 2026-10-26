package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/senojj/pgsql/pgmsg"
)

func startup(
	ctx context.Context,
	conn *Conn,
) {
	if err := ctx.Err(); err != nil {
		conn.err = err
		return
	}

	req := pgmsg.StartupMessage{
		ProtocolVersion: conn.protocol,
		Keywords:        conn.paramKeywords,
		Values:          conn.paramValues,
	}

	err := conn.encoder.Encode(&req)
	if err != nil {
		conn.err = err
		return
	}

	err = conn.encoder.Flush()
	if err != nil {
		conn.err = err
		return
	}

	for {
		if err = ctx.Err(); err != nil {
			conn.err = err
			return
		}

		resp, err := conn.decoder.Decode()
		if err != nil {
			conn.err = err
			return
		}

		switch msg := resp.(type) {
		case *pgmsg.ErrorResponse:
			conn.err = rtoerr(msg)
			return
		case *pgmsg.ReadyForQuery:
			conn.status = ConnectionOK
			return
		case *pgmsg.AuthenticationOk:
		case *pgmsg.ParameterStatus:
		case *pgmsg.BackendKeyData:
			conn.processID = msg.ProcessID
		default:
			conn.err = errors.New("unhandled message")
			return
		}
	}
}
