package pgproto

import (
	"context"
	"fmt"

	"github.com/senojj/pgsql/pgproto/pgmsg"
)

func (f *Frontend) Startup(
	ctx context.Context,
	protocolVersion int32,
	parameters map[string]string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	req := pgmsg.StartupMessage{
		ProtocolVersion: protocolVersion,
		Parameters:      parameters,
	}

	err := f.encoder.Encode(&req)
	if err != nil {
		return err
	}

	err = f.encoder.Flush()
	if err != nil {
		return err
	}

	for {
		if err = ctx.Err(); err != nil {
			return err
		}

		resp, err := f.decoder.Decode()
		if err != nil {
			return err
		}

		switch msg := resp.(type) {
		case *pgmsg.ErrorResponse:
			return rtoerr(msg)
		case *pgmsg.ReadyForQuery:
			println("ready for query")
			return nil
		case *pgmsg.AuthenticationOk:
			println("authentication OK")
		case *pgmsg.ParameterStatus:
			println("parameter status:", msg.Name, msg.Value)
		case *pgmsg.BackendKeyData:
			fmt.Printf("backend key data: %d\n", msg.ProcessID)
		default:
			println("unhandled message")
			return nil
		}
	}
}
