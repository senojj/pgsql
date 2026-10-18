package pgmsg

import (
	"math"

	"github.com/senojj/pgsql/pgio"
)

type AuthenticationOk struct{}

func (x *AuthenticationOk) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeOk)
	return buf.Bytes, nil
}

func (x *AuthenticationOk) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeOk != auth {
		return unexpectedAuth(auth, AuthTypeOk)
	}

	if len(b) > 0 {
		return pgio.ErrValueOverflow
	}
	return nil
}

type AuthenticationKerberosV5 struct{}

func (x *AuthenticationKerberosV5) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeKerberosV5)
	return buf.Bytes, nil
}

func (x *AuthenticationKerberosV5) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeKerberosV5 != auth {
		return unexpectedAuth(auth, AuthTypeKerberosV5)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationCleartextPassword struct{}

func (x *AuthenticationCleartextPassword) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeClearTextPassword)
	return buf.Bytes, nil
}

func (x *AuthenticationCleartextPassword) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeClearTextPassword != auth {
		return unexpectedAuth(auth, AuthTypeClearTextPassword)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationMD5Password struct {
	Salt [4]byte
}

func (x *AuthenticationMD5Password) AppendBinary(b []byte) ([]byte, error) {
	const sizeSalt = 4
	const length = sizeMessageLength + sizeAuth + sizeSalt
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeMD5Password)
	buf.AppendByte(x.Salt[:]...)
	return buf.Bytes, nil
}

func (x *AuthenticationMD5Password) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeMD5Password != auth {
		return unexpectedAuth(auth, AuthTypeMD5Password)
	}
	copy(x.Salt[:], b)
	return nil
}

type AuthenticationGSS struct{}

func (x *AuthenticationGSS) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeGSS)
	return buf.Bytes, nil
}

func (x *AuthenticationGSS) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeGSS != auth {
		return unexpectedAuth(auth, AuthTypeGSS)
	}
	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationGSSContinue struct {
	Data []byte
}

func (x *AuthenticationGSSContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)
	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeGSSContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *AuthenticationGSSContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeGSSContinue != auth {
		return unexpectedAuth(auth, AuthTypeGSSContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type AuthenticationSSPI struct{}

func (x *AuthenticationSSPI) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeSSPI)
	return buf.Bytes, nil
}

func (x *AuthenticationSSPI) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeSSPI != auth {
		return unexpectedAuth(auth, AuthTypeSSPI)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type AuthenticationSASL struct {
	Mechanisms []string
}

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
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeSASL)
	buf.AppendString(x.Mechanisms...)
	buf.AppendByte(0)
	return buf.Bytes, nil
}

func (x *AuthenticationSASL) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	auth, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeSASL != auth {
		return unexpectedAuth(auth, AuthTypeSASL)
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

func (x *AuthenticationSASLContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeSASLContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *AuthenticationSASLContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeSASLContinue != auth {
		return unexpectedAuth(auth, AuthTypeSASLContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type AuthenticationSASLFinal struct {
	Data []byte
}

func (x *AuthenticationSASLFinal) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthTypeSASLFinal)
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *AuthenticationSASLFinal) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthTypeSASLFinal != auth {
		return unexpectedAuth(auth, AuthTypeSASLFinal)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}
