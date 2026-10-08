package pgsql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgAuthenticationOk(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authOk)

	var m msgAuthenticationOk

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationKerberosV5(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authKerberosV5)

	var m msgAuthenticationKerberosV5

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationClearTextPassword(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authClearTextPassword)

	var m msgAuthenticationCleartextPassword

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationMD5Password(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(12)
	buf.AppendInt32(authMD5Password)
	buf.AppendByte([]byte("4321")...)

	var m msgAuthenticationMD5Password

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "4321", string(m.Salt[:]))
	})
}

func TestMsgAuthenticationGSS(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authGSS)

	var m msgAuthenticationGSS

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationGSSContinue(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(authGSSContinue)
	buf.AppendByte([]byte("hello world")...)

	var m msgAuthenticationGSSContinue

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestMsgAuthenticationSSPI(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authSSPI)

	var m msgAuthenticationSSPI

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationSASL(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(21)
	buf.AppendInt32(authSASL)
	buf.AppendString("hello", "world")
	buf.AppendByte(0)

	var m msgAuthenticationSASL

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, []string{"hello", "world"}, m.Mechanisms)
	})
}

func TestMsgAuthenticationSASLContinue(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(authSASLContinue)
	buf.AppendByte([]byte("hello world")...)

	var m msgAuthenticationSASLContinue

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestMsgAuthenticationSASLFinal(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(authSASLFinal)
	buf.AppendByte([]byte("hello world")...)

	var m msgAuthenticationSASLFinal

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}
