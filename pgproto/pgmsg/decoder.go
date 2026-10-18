package pgmsg

import (
	"encoding/binary"
	"io"

	"github.com/senojj/pgsql/pgproto/pgio"
)

const readBufferSize int = 1 << 13 // 8kb

type Decoder struct {
	reader io.Reader
	buf    []byte
	rp, wp int

	authenticationOk                AuthenticationOk
	authenticationKerberosV5        AuthenticationKerberosV5
	authenticationCleartextPassword AuthenticationCleartextPassword
	authenticationMD5Password       AuthenticationMD5Password
	authenticationGSS               AuthenticationGSS
	authenticationGSSContinue       AuthenticationGSSContinue
	authenticationSSPI              AuthenticationSSPI
	authenticationSASL              AuthenticationSASL
	authenticationSASLContinue      AuthenticationSASLContinue
	authenticationSASLFinal         AuthenticationSASLFinal
	errorResponse                   ErrorResponse
	readyForQuery                   ReadyForQuery
	parameterStatus                 ParameterStatus
	backendKeyData                  BackendKeyData
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		reader: reader,
		buf:    make([]byte, readBufferSize),
	}
}

func (d *Decoder) resolve(data []byte) (Backend, error) {
	var m Backend

	switch data[0] {
	case typeAuthentication:
		if len(data) < 9 {
			return nil, pgio.ErrValueUnderflow
		}
		authType := int32(binary.BigEndian.Uint32(data[5:9]))

		switch authType {
		case authTypeOk:
			m = &d.authenticationOk
		case authTypeKerberosV5:
			m = &d.authenticationKerberosV5
		case authTypeClearTextPassword:
			m = &d.authenticationCleartextPassword
		case authTypeMD5Password:
			m = &d.authenticationMD5Password
		case authTypeGSS:
			m = &d.authenticationGSS
		case authTypeGSSContinue:
			m = &d.authenticationGSSContinue
		case authTypeSSPI:
			m = &d.authenticationSSPI
		case authTypeSASL:
			m = &d.authenticationSASL
		case authTypeSASLContinue:
			m = &d.authenticationSASLContinue
		case authTypeSASLFinal:
			m = &d.authenticationSASLFinal
		default:
			println("none of those authentication methods")
		}
		println("received authentication message")
	case typeErrorResponse:
		m = &d.errorResponse
	case typeReadyForQuery:
		m = &d.readyForQuery
	case typeParameterStatus:
		m = &d.parameterStatus
	case typeBackendKeyData:
		m = &d.backendKeyData
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

func (d *Decoder) Decode() (Backend, error) {
	var size uint32

	for {
		for {
			available := d.wp - d.rp
			if available < 5 {
				break
			}

			if available < int(size) {
				break
			}

			if size > 0 {
				data := d.buf[d.rp : d.rp+int(size)]
				msg, err := d.resolve(data)
				d.rp += int(size)
				return msg, err
			}

			size = binary.BigEndian.Uint32(d.buf[d.rp+1 : d.rp+5])
			size += 1
		}

		if d.rp > 0 {
			if d.rp == d.wp {
				d.rp, d.wp = 0, 0
			} else {
				d.wp = copy(d.buf[0:], d.buf[d.rp:d.wp])
				d.rp = 0
			}
		}

		if len(d.buf)-d.wp < 1024 {
			tmp := make([]byte, len(d.buf)*2)
			copy(tmp, d.buf[:d.wp])
			d.buf = tmp
		}

		n, err := d.reader.Read(d.buf[d.wp:])
		d.wp += n
		if err != nil {
			return nil, err
		}
	}
}
