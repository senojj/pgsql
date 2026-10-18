package pgmsg_test

import (
	"testing"

	"github.com/senojj/pgsql/pgmsg"
	"github.com/senojj/pgsql/pgio"

	"github.com/stretchr/testify/require"
)

func TestAuthenticationOk(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(pgmsg.AuthTypeOk)

	var m pgmsg.AuthenticationOk

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationKerberosV5(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(pgmsg.AuthTypeKerberosV5)

	var m pgmsg.AuthenticationKerberosV5

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationClearTextPassword(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(pgmsg.AuthTypeClearTextPassword)

	var m pgmsg.AuthenticationCleartextPassword

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationMD5Password(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(12)
	buf.AppendInt32(pgmsg.AuthTypeMD5Password)
	buf.AppendByte([]byte("4321")...)

	var m pgmsg.AuthenticationMD5Password

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "4321", string(m.Salt[:]))
	})
}

func TestAuthenticationGSS(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(pgmsg.AuthTypeGSS)

	var m pgmsg.AuthenticationGSS

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationGSSContinue(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(pgmsg.AuthTypeGSSContinue)
	buf.AppendByte([]byte("hello world")...)

	var m pgmsg.AuthenticationGSSContinue

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestAuthenticationSSPI(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(pgmsg.AuthTypeSSPI)

	var m pgmsg.AuthenticationSSPI

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationSASL(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(21)
	buf.AppendInt32(pgmsg.AuthTypeSASL)
	buf.AppendString("hello", "world")
	buf.AppendByte(0)

	var m pgmsg.AuthenticationSASL

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, []string{"hello", "world"}, m.Mechanisms)
	})
}

func TestAuthenticationSASLContinue(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(pgmsg.AuthTypeSASLContinue)
	buf.AppendByte([]byte("hello world")...)

	var m pgmsg.AuthenticationSASLContinue

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestAuthenticationSASLFinal(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(pgmsg.AuthTypeSASLFinal)
	buf.AppendByte([]byte("hello world")...)

	var m pgmsg.AuthenticationSASLFinal

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}
