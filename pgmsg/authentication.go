package pgmsg

import (
	"math"

	"github.com/senojj/pgsql/pgio"
)

type AuthenticationOk struct{}

func (_ *AuthenticationOk) backend() {}

func (x *AuthenticationOk) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeOk)
	return buf.Bytes, nil
}

func (x *AuthenticationOk) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeOk != auth {
		return unexpectedAuth(auth, authTypeOk)
	}

	if len(b) > 0 {
		return pgio.ErrValueOverflow
	}
	return nil
}

type AuthenticationKerberosV5 struct{}

func (_ *AuthenticationKerberosV5) backend() {}

func (x *AuthenticationKerberosV5) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeKerberosV5)
	return buf.Bytes, nil
}

func (x *AuthenticationKerberosV5) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeKerberosV5 != auth {
		return unexpectedAuth(auth, authTypeKerberosV5)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationCleartextPassword struct{}

func (_ *AuthenticationCleartextPassword) backend() {}

func (x *AuthenticationCleartextPassword) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeClearTextPassword)
	return buf.Bytes, nil
}

func (x *AuthenticationCleartextPassword) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeClearTextPassword != auth {
		return unexpectedAuth(auth, authTypeClearTextPassword)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationMD5Password struct {
	Salt [4]byte
}

func (_ *AuthenticationMD5Password) backend() {}

func (x *AuthenticationMD5Password) AppendBinary(b []byte) ([]byte, error) {
	const sizeSalt = 4
	const length = sizeMessageLength + sizeAuth + sizeSalt
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeMD5Password)
	buf.AppendByte(x.Salt[:]...)
	return buf.Bytes, nil
}

func (x *AuthenticationMD5Password) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeMD5Password != auth {
		return unexpectedAuth(auth, authTypeMD5Password)
	}
	copy(x.Salt[:], b)
	return nil
}

type AuthenticationGSS struct{}

func (_ *AuthenticationGSS) backend() {}

func (x *AuthenticationGSS) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeGSS)
	return buf.Bytes, nil
}

func (x *AuthenticationGSS) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeGSS != auth {
		return unexpectedAuth(auth, authTypeGSS)
	}
	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationGSSContinue struct {
	Data []byte
}

func (_ *AuthenticationGSSContinue) backend() {}

func (x *AuthenticationGSSContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)
	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeGSSContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *AuthenticationGSSContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeGSSContinue != auth {
		return unexpectedAuth(auth, authTypeGSSContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type AuthenticationSSPI struct{}

func (_ *AuthenticationSSPI) backend() {}

func (x *AuthenticationSSPI) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeSSPI)
	return buf.Bytes, nil
}

func (x *AuthenticationSSPI) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeSSPI != auth {
		return unexpectedAuth(auth, authTypeSSPI)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationSASL struct {
	Mechanisms []string
}

func (_ *AuthenticationSASL) backend() {}

func (x *AuthenticationSASL) AppendBinary(b []byte) ([]byte, error) {
	countMechanisms := len(x.Mechanisms)
	sizeMechanisms := 0

	for i := range countMechanisms {
		sizeMechanisms += len(x.Mechanisms[i]) + 1 // null terminated string
	}
	sizeMechanisms += 1 // null terminated list

	length := sizeMessageLength + sizeAuth + sizeMechanisms

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeSASL)
	buf.AppendString(x.Mechanisms...)
	buf.AppendByte(0)
	return buf.Bytes, nil
}

func (x *AuthenticationSASL) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	auth, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeSASL != auth {
		return unexpectedAuth(auth, authTypeSASL)
	}
	x.Mechanisms = make([]string, 0, buf.Count(pgio.Zero))

	for {
		var mechanism string
		var err error

		mechanism, err = buf.ShiftString()
		if err != nil {
			return invalidFormat(err)
		}

		if len(mechanism) == 0 {
			break
		}
		x.Mechanisms = append(x.Mechanisms, mechanism)
	}
	return nil
}

type AuthenticationSASLContinue struct {
	Data []byte
}

func (_ *AuthenticationSASLContinue) backend() {}

func (x *AuthenticationSASLContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeSASLContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *AuthenticationSASLContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeSASLContinue != auth {
		return unexpectedAuth(auth, authTypeSASLContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type AuthenticationSASLFinal struct {
	Data []byte
}

func (_ *AuthenticationSASLFinal) backend() {}

func (x *AuthenticationSASLFinal) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authTypeSASLFinal)
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *AuthenticationSASLFinal) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authTypeSASLFinal != auth {
		return unexpectedAuth(auth, authTypeSASLFinal)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}
