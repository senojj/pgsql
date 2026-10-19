package pgproto

import (
	"context"
	"errors"
	"fmt"

	"github.com/senojj/pgsql/pgproto/pgmsg"
)

func startup(
	ctx context.Context,
	conn *Conn,
) {
	if err := ctx.Err(); err != nil {
		return err
	}

	req := pgmsg.StartupMessage{
		ProtocolVersion: conn.version,
		Parameters:      conn.params,
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
			conn.params[msg.Name] = msg.Value
		case *pgmsg.BackendKeyData:
			conn.processID = msg.ProcessID
		default:
			conn.err = errors.New("unhandled message")	
			return
		}
	}
}
