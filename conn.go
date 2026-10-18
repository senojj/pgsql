package pgsql

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/senojj/pgsql/pgmsg"
	"github.com/senojj/pgsql/pgio"
)

const readBufferSize int = 1 << 13 // 8kb

func rtoerr(msg *pgmsg.ErrorResponse) error {
	idxSeverity := bytes.IndexByte(msg.Fields, fieldSeverity)
	idxCode := bytes.IndexByte(msg.Fields, fieldCode)
	idxMessage := bytes.IndexByte(msg.Fields, fieldMessage)
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

type frontend struct {
	raw    net.Conn
	rbuf   []byte
	rp, wp int

	wbuf []byte

	authenticationOk                pgmsg.AuthenticationOk
	authenticationKerberosV5        pgmsg.AuthenticationKerberosV5
	authenticationCleartextPassword pgmsg.AuthenticationCleartextPassword
	authenticationMD5Password       pgmsg.AuthenticationMD5Password
	authenticationGSS               pgmsg.AuthenticationGSS
	authenticationGSSContinue       pgmsg.AuthenticationGSSContinue
	authenticationSSPI              pgmsg.AuthenticationSSPI
	authenticationSASL              pgmsg.AuthenticationSASL
	authenticationSASLContinue      pgmsg.AuthenticationSASLContinue
	authenticationSASLFinal         pgmsg.AuthenticationSASLFinal
	errorResponse                   pgmsg.ErrorResponse
	readyForQuery                   pgmsg.ReadyForQuery
	parameterStatus                 pgmsg.ParameterStatus
	backendKeyData                  pgmsg.BackendKeyData
}

func newFrontend(conn net.Conn) *frontend {
	return &frontend{
		raw:  conn,
		rbuf: make([]byte, readBufferSize),
		wbuf: make([]byte, 0),
	}
}

func (f *frontend) resolve(data []byte) (message, error) {
	var m message

	switch data[0] {
	case pgmsg.TypeAuthentication:
		if len(data) < 9 {
			return nil, pgio.ErrValueUnderflow
		}
		authType := int32(binary.BigEndian.Uint32(data[5:9]))

		switch authType {
		case pgmsg.AuthTypeOk:
			m = &f.authenticationOk
		case pgmsg.AuthTypeKerberosV5:
			m = &f.authenticationKerberosV5
		case pgmsg.AuthTypeClearTextPassword:
			m = &f.authenticationCleartextPassword
		case pgmsg.AuthTypeMD5Password:
			m = &f.authenticationMD5Password
		case pgmsg.AuthTypeGSS:
			m = &f.authenticationGSS
		case pgmsg.AuthTypeGSSContinue:
			m = &f.authenticationGSSContinue
		case pgmsg.AuthTypeSSPI:
			m = &f.authenticationSSPI
		case pgmsg.AuthTypeSASL:
			m = &f.authenticationSASL
		case pgmsg.AuthTypeSASLContinue:
			m = &f.authenticationSASLContinue
		case pgmsg.AuthTypeSASLFinal:
			m = &f.authenticationSASLFinal
		default:
			println("none of those authentication methods")
		}
		println("received authentication message")
	case pgmsg.TypeErrorResponse:
		m = &f.errorResponse
	case pgmsg.TypeReadyForQuery:
		m = &f.readyForQuery
	case pgmsg.TypeParameterStatus:
		m = &f.parameterStatus
	case pgmsg.TypeBackendKeyData:
		m = &f.backendKeyData
	default:
		println("received kind " + string(data[0:1]))
	}

	if m != nil {
		err := m.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (f *frontend) recv() (message, error) {
	var size uint32

	for {
		for {
			available := f.wp - f.rp
			if available < 5 {
				break
			}

			if available < int(size) {
				break
			}

			if size > 0 {
				data := f.rbuf[f.rp : f.rp+int(size)]
				msg, err := f.resolve(data)
				f.rp += int(size)
				return msg, err
			}

			size = binary.BigEndian.Uint32(f.rbuf[f.rp+1 : f.rp+5])
			size += 1
		}

		if f.rp > 0 {
			if f.rp == f.wp {
				f.rp, f.wp = 0, 0
			} else {
				f.wp = copy(f.rbuf[0:], f.rbuf[f.rp:f.wp])
				f.rp = 0
			}
		}

		if len(f.rbuf)-f.wp < 1024 {
			tmp := make([]byte, len(f.rbuf)*2)
			copy(tmp, f.rbuf[:f.wp])
			f.rbuf = tmp
		}

		n, err := f.raw.Read(f.rbuf[f.wp:])
		f.wp += n
		if err != nil {
			return nil, err
		}
	}
}

func (f *frontend) send(msg message) error {
	var err error
	f.wbuf, err = msg.AppendBinary(f.wbuf)
	if err != nil {
		return err
	}
	return nil
}

func (f *frontend) flush() error {
	if len(f.wbuf) == 0 {
		return nil
	}

	defer func() {
		f.wbuf = f.wbuf[:0]
	}()

	_, err := f.raw.Write(f.wbuf)
	if err != nil {
		return err
	}
	return nil
}

func (f *frontend) close() error {
	return f.raw.Close()
}

func (f *frontend) startup(
	ctx context.Context,
	req *pgmsg.StartupMessage,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := f.send(req)
	if err != nil {
		return err
	}

	err = f.flush()
	if err != nil {
		return err
	}

	for {
		if err = ctx.Err(); err != nil {
			return err
		}

		resp, err := f.recv()
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
			fmt.Printf("backend key data: %d", msg.ProcessID)
		default:
			println("unhandled message")
			return nil
		}
	}
}

func connect(ctx context.Context, s string) error {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}

	hostname := u.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}

	port := uint64(5432)
	portStr := u.Port()
	if portStr != "" {
		port, err = strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			return err
		}
	}

	database := u.Path
	length := len(database)
	if length > 0 && database[0] == '/' {
		database = database[1:]
		length--
	}
	if length > 0 && database[length-1] == '/' {
		database = database[:length-1]
	}
	if database == "" {
		database = "postgres"
	}

	address := hostname + ":" + strconv.FormatUint(uint64(port), 10)

	user := u.User.Username()

	c, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	front := newFrontend(c)

	start := pgmsg.StartupMessage{
		ProtocolVersion: protocolVersion3_2,
		Parameters: map[string]string{
			paramUser:     user,
			paramDatabase: database,
		},
	}

	err = front.startup(ctx, &start)
	if err != nil {
		return err
	}
	return nil
}
