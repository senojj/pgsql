package pgwire_test

import (
	"gopsql/pgio"
	"gopsql/pgwire"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgAuthenticationOk(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(8)
	buf.AppendInt32(int32(pgwire.AuthOk))

	var m pgwire.MsgAuthenticationOk

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationKerberosV5(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(8)
	buf.AppendInt32(int32(pgwire.AuthKerberosV5))

	var m pgwire.MsgAuthenticationKerberosV5

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationClearTextPassword(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(8)
	buf.AppendInt32(int32(pgwire.AuthClearTextPassword))

	var m pgwire.MsgAuthenticationCleartextPassword

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationMD5Password(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(12)
	buf.AppendInt32(int32(pgwire.AuthMD5Password))
	buf.AppendByte([]byte("4321")...)

	var m pgwire.MsgAuthenticationMD5Password

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "4321", string(m.Salt[:]))
	})
}

func TestMsgAuthenticationGSS(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(8)
	buf.AppendInt32(int32(pgwire.AuthGSS))

	var m pgwire.MsgAuthenticationGSS

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationGSSContinue(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(19)
	buf.AppendInt32(int32(pgwire.AuthGSSContinue))
	buf.AppendByte([]byte("hello world")...)

	var m pgwire.MsgAuthenticationGSSContinue

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestMsgAuthenticationSSPI(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(8)
	buf.AppendInt32(int32(pgwire.AuthSSPI))

	var m pgwire.MsgAuthenticationSSPI

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgAuthenticationSASL(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(21)
	buf.AppendInt32(int32(pgwire.AuthSASL))
	buf.AppendString("hello", "world")
	buf.AppendByte(0)

	var m pgwire.MsgAuthenticationSASL

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, []string{"hello", "world"}, m.Mechanisms)
	})
}

func TestMsgAuthenticationSASLContinue(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(19)
	buf.AppendInt32(int32(pgwire.AuthSASLContinue))
	buf.AppendByte([]byte("hello world")...)

	var m pgwire.MsgAuthenticationSASLContinue

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestMsgAuthenticationSASLFinal(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(byte(pgwire.TypeAuthentication))
	buf.AppendInt32(19)
	buf.AppendInt32(int32(pgwire.AuthSASLFinal))
	buf.AppendByte([]byte("hello world")...)

	var m pgwire.MsgAuthenticationSASLFinal

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}
