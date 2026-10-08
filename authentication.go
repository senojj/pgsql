package pgsql

import (
	"math"
)

type msgAuthenticationOk struct{}

func (x *msgAuthenticationOk) message() {}

func (x *msgAuthenticationOk) backend() {}

func (x *msgAuthenticationOk) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authOk)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationOk) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authOk != auth {
		return unexpectedAuth(auth, authOk)
	}

	if len(b) > 0 {
		return ErrValueOverflow
	}
	return nil
}

type msgAuthenticationKerberosV5 struct{}

func (x *msgAuthenticationKerberosV5) message() {}

func (x *msgAuthenticationKerberosV5) backend() {}

func (x *msgAuthenticationKerberosV5) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authKerberosV5)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationKerberosV5) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authKerberosV5 != auth {
		return unexpectedAuth(auth, authKerberosV5)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgAuthenticationCleartextPassword struct{}

func (x *msgAuthenticationCleartextPassword) message() {}

func (x *msgAuthenticationCleartextPassword) backend() {}

func (x *msgAuthenticationCleartextPassword) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authClearTextPassword)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationCleartextPassword) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authClearTextPassword != auth {
		return unexpectedAuth(auth, authClearTextPassword)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgAuthenticationMD5Password struct {
	Salt [4]byte
}

func (x *msgAuthenticationMD5Password) message() {}

func (x *msgAuthenticationMD5Password) backend() {}

func (x *msgAuthenticationMD5Password) AppendBinary(b []byte) ([]byte, error) {
	const sizeSalt = 4
	const length = sizeMessageLength + sizeAuth + sizeSalt
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authMD5Password)
	buf.AppendByte(x.Salt[:]...)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationMD5Password) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authMD5Password != auth {
		return unexpectedAuth(auth, authMD5Password)
	}
	copy(x.Salt[:], b)
	return nil
}

type msgAuthenticationGSS struct{}

func (x *msgAuthenticationGSS) message() {}

func (x *msgAuthenticationGSS) backend() {}

func (x *msgAuthenticationGSS) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authGSS)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationGSS) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authGSS != auth {
		return unexpectedAuth(auth, authGSS)
	}
	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgAuthenticationGSSContinue struct {
	Data []byte
}

func (x *msgAuthenticationGSSContinue) message() {}

func (x *msgAuthenticationGSSContinue) backend() {}

func (x *msgAuthenticationGSSContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)
	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authGSSContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationGSSContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authGSSContinue != auth {
		return unexpectedAuth(auth, authGSSContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type msgAuthenticationSSPI struct{}

func (x *msgAuthenticationSSPI) message() {}

func (x *msgAuthenticationSSPI) backend() {}

func (x *msgAuthenticationSSPI) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authSSPI)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationSSPI) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authSSPI != auth {
		return unexpectedAuth(auth, authSSPI)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgAuthenticationSASL struct {
	Mechanisms []string
}

func (x *msgAuthenticationSASL) message() {}

func (x *msgAuthenticationSASL) backend() {}

func (x *msgAuthenticationSASL) AppendBinary(b []byte) ([]byte, error) {
	countMechanisms := len(x.Mechanisms)
	sizeMechanisms := 0

	for i := range countMechanisms {
		sizeMechanisms += len(x.Mechanisms[i]) + 1 // null terminated string
	}
	sizeMechanisms += 1 // null terminated list

	length := sizeMessageLength + sizeAuth + sizeMechanisms

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authSASL)
	buf.AppendString(x.Mechanisms...)
	buf.AppendByte(0)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationSASL) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	auth, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if authSASL != auth {
		return unexpectedAuth(auth, authSASL)
	}
	x.Mechanisms = make([]string, 0, buf.Count(zero))

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

type msgAuthenticationSASLContinue struct {
	Data []byte
}

func (x *msgAuthenticationSASLContinue) message() {}

func (x *msgAuthenticationSASLContinue) backend() {}

func (x *msgAuthenticationSASLContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authSASLContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationSASLContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authSASLContinue != auth {
		return unexpectedAuth(auth, authSASLContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type msgAuthenticationSASLFinal struct {
	Data []byte
}

func (x *msgAuthenticationSASLFinal) message() {}

func (x *msgAuthenticationSASLFinal) backend() {}

func (x *msgAuthenticationSASLFinal) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(authSASLFinal)
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *msgAuthenticationSASLFinal) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if authSASLFinal != auth {
		return unexpectedAuth(auth, authSASLFinal)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}
