package pgwire

import (
	"gopsql/pgio"
	"math"
)

var _ Message = &MsgAuthenticationOk{}
var _ Backend = &MsgAuthenticationOk{}

type MsgAuthenticationOk struct{}

func (x *MsgAuthenticationOk) message() {}

func (x *MsgAuthenticationOk) backend() {}

func (x *MsgAuthenticationOk) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthOk)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationOk) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthOk != auth {
		return unexpectedAuth(auth, AuthOk)
	}

	if len(b) > 0 {
		return pgio.ErrValueOverflow
	}
	return nil
}

var _ Message = &MsgAuthenticationKerberosV5{}
var _ Backend = &MsgAuthenticationKerberosV5{}

type MsgAuthenticationKerberosV5 struct{}

func (x *MsgAuthenticationKerberosV5) message() {}

func (x *MsgAuthenticationKerberosV5) backend() {}

func (x *MsgAuthenticationKerberosV5) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthKerberosV5)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationKerberosV5) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthKerberosV5 != auth {
		return unexpectedAuth(auth, AuthKerberosV5)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

var _ Message = &MsgAuthenticationCleartextPassword{}
var _ Backend = &MsgAuthenticationCleartextPassword{}

type MsgAuthenticationCleartextPassword struct{}

func (x *MsgAuthenticationCleartextPassword) message() {}

func (x *MsgAuthenticationCleartextPassword) backend() {}

func (x *MsgAuthenticationCleartextPassword) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthClearTextPassword)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationCleartextPassword) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthClearTextPassword != auth {
		return unexpectedAuth(auth, AuthClearTextPassword)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

var _ Message = &MsgAuthenticationMD5Password{}
var _ Backend = &MsgAuthenticationMD5Password{}

type MsgAuthenticationMD5Password struct {
	Salt [4]byte
}

func (x *MsgAuthenticationMD5Password) message() {}

func (x *MsgAuthenticationMD5Password) backend() {}

func (x *MsgAuthenticationMD5Password) AppendBinary(b []byte) ([]byte, error) {
	const sizeSalt = 4
	const length = sizeMessageLength + sizeAuth + sizeSalt
	const size = sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthMD5Password)
	buf.AppendByte(x.Salt[:]...)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationMD5Password) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthMD5Password != auth {
		return unexpectedAuth(auth, AuthMD5Password)
	}
	copy(x.Salt[:], b)
	return nil
}

var _ Message = &MsgAuthenticationGSS{}
var _ Backend = &MsgAuthenticationGSS{}

type MsgAuthenticationGSS struct{}

func (x *MsgAuthenticationGSS) message() {}

func (x *MsgAuthenticationGSS) backend() {}

func (x *MsgAuthenticationGSS) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthGSS)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationGSS) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthGSS != auth {
		return unexpectedAuth(auth, AuthGSS)
	}
	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

var _ Message = &MsgAuthenticationGSSContinue{}
var _ Backend = &MsgAuthenticationGSSContinue{}

type MsgAuthenticationGSSContinue struct {
	Data []byte
}

func (x *MsgAuthenticationGSSContinue) message() {}

func (x *MsgAuthenticationGSSContinue) backend() {}

func (x *MsgAuthenticationGSSContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)
	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthGSSContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationGSSContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthGSSContinue != auth {
		return unexpectedAuth(auth, AuthGSSContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

var _ Message = &MsgAuthenticationSSPI{}
var _ Backend = &MsgAuthenticationSSPI{}

type MsgAuthenticationSSPI struct{}

func (x *MsgAuthenticationSSPI) message() {}

func (x *MsgAuthenticationSSPI) backend() {}

func (x *MsgAuthenticationSSPI) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength + sizeAuth
	const size = sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthSSPI)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationSSPI) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthSSPI != auth {
		return unexpectedAuth(auth, AuthSSPI)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

var _ Message = &MsgAuthenticationSASL{}
var _ Backend = &MsgAuthenticationSASL{}

type MsgAuthenticationSASL struct {
	Mechanisms []string
}

func (x *MsgAuthenticationSASL) message() {}

func (x *MsgAuthenticationSASL) backend() {}

func (x *MsgAuthenticationSASL) AppendBinary(b []byte) ([]byte, error) {
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

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthSASL)
	buf.AppendString(x.Mechanisms...)
	buf.AppendByte(0x0)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationSASL) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.NewBuffer(b)

	auth, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if AuthSASL != auth {
		return unexpectedAuth(auth, AuthSASL)
	}
	x.Mechanisms = make([]string, 0, buf.Count(nullByte))

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

var _ Message = &MsgAuthenticationSASLContinue{}
var _ Backend = &MsgAuthenticationSASLContinue{}

type MsgAuthenticationSASLContinue struct {
	Data []byte
}

func (x *MsgAuthenticationSASLContinue) message() {}

func (x *MsgAuthenticationSASLContinue) backend() {}

func (x *MsgAuthenticationSASLContinue) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthSASLContinue)
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationSASLContinue) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthSASLContinue != auth {
		return unexpectedAuth(auth, AuthSASLContinue)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

var _ Message = &MsgAuthenticationSASLFinal{}
var _ Backend = &MsgAuthenticationSASLFinal{}

type MsgAuthenticationSASLFinal struct {
	Data []byte
}

func (x *MsgAuthenticationSASLFinal) message() {}

func (x *MsgAuthenticationSASLFinal) backend() {}

func (x *MsgAuthenticationSASLFinal) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeAuth + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.NewBuffer(b)
	buf.Grow(size)
	buf.AppendByte(TypeAuthentication)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(AuthSASLFinal)
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *MsgAuthenticationSASLFinal) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeAuthentication, b)
	if err != nil {
		return invalidFormat(err)
	}

	auth, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if AuthSASLFinal != auth {
		return unexpectedAuth(auth, AuthSASLFinal)
	}
	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}
