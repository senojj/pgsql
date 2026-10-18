package pgmsg

import (
	"testing"

	"github.com/senojj/pgsql/pgproto/pgio"

	"github.com/stretchr/testify/require"
)

func TestAuthenticationOk(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authTypeOk)

	var m AuthenticationOk

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationKerberosV5(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authTypeKerberosV5)

	var m AuthenticationKerberosV5

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationClearTextPassword(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authTypeClearTextPassword)

	var m AuthenticationCleartextPassword

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationMD5Password(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(12)
	buf.AppendInt32(authTypeMD5Password)
	buf.AppendByte([]byte("4321")...)

	var m AuthenticationMD5Password

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "4321", string(m.Salt[:]))
	})
}

func TestAuthenticationGSS(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authTypeGSS)

	var m AuthenticationGSS

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationGSSContinue(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(authTypeGSSContinue)
	buf.AppendByte([]byte("hello world")...)

	var m AuthenticationGSSContinue

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestAuthenticationSSPI(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(8)
	buf.AppendInt32(authTypeSSPI)

	var m AuthenticationSSPI

	testMessage(t, buf.Bytes, &m, nil)
}

func TestAuthenticationSASL(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(21)
	buf.AppendInt32(authTypeSASL)
	buf.AppendString("hello", "world")
	buf.AppendByte(0)

	var m AuthenticationSASL

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, []string{"hello", "world"}, m.Mechanisms)
	})
}

func TestAuthenticationSASLContinue(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(authTypeSASLContinue)
	buf.AppendByte([]byte("hello world")...)

	var m AuthenticationSASLContinue

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}

func TestAuthenticationSASLFinal(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeAuthentication)
	buf.AppendInt32(19)
	buf.AppendInt32(authTypeSASLFinal)
	buf.AppendByte([]byte("hello world")...)

	var m AuthenticationSASLFinal

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", string(m.Data))
	})
}
