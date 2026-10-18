package pgmsg_test

import (
	"testing"

	"github.com/senojj/pgsql/pgmsg"
	"github.com/senojj/pgsql/pgio"

	"github.com/stretchr/testify/require"
)

func TestBind(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeBind)
	buf.AppendInt32(73)
	buf.AppendString("hello world")
	buf.AppendString("lorem ipsum")
	buf.AppendInt16(3)
	buf.AppendInt16(
		int16(0),
		int16(0),
		int16(0),
	)
	buf.AppendInt16(3)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("lorem")...)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("ipsum")...)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("dolor")...)
	buf.AppendInt16(3)
	buf.AppendInt16(
		int16(0),
		int16(0),
		int16(0),
	)

	var m pgmsg.Bind

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.DestinationName)

		require.Equal(t, "lorem ipsum", m.SourceName)

		require.Equal(t, []byte{
			0,
			0,
			0,
		}, m.ParameterFormatCodes)

		require.Equal(t, [][]byte{
			[]byte("lorem"),
			[]byte("ipsum"),
			[]byte("dolor"),
		}, m.ParameterData)

		require.Equal(t, []byte{
			0,
			0,
			0,
		}, m.ColumnFormatCodes)
	})
}

func TestCancelRequest(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendInt32(23)
	buf.AppendInt32(pgmsg.CodeCancelRequest)
	buf.AppendInt32(4321)
	buf.AppendByte([]byte("hello world")...)

	var m pgmsg.CancelRequest

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ProcessID)
		require.Equal(t, "hello world", string(m.SecretKey))
	})
}

func TestClose(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeClose)
	buf.AppendInt32(17)
	buf.AppendInt8(int8(1))
	buf.AppendString("hello world")

	var m pgmsg.Close

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, uint8(1), m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestCopyFail(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeCopyFail)
	buf.AppendInt32(16)
	buf.AppendString("hello world")

	var m pgmsg.CopyFail

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.Message)
	})
}

func TestDescribe(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeDescribe)
	buf.AppendInt32(17)
	buf.AppendByte(1)
	buf.AppendString("hello world")

	var m pgmsg.Describe

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, uint8(1), m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestExecute(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeExecute)
	buf.AppendInt32(20)
	buf.AppendString("hello world")
	buf.AppendInt32(5)

	var m pgmsg.Execute

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.PortalName)
		require.Equal(t, int32(5), m.RowLimit)
	})
}

func TestFlush(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeFlush)
	buf.AppendInt32(4)

	var m pgmsg.Flush

	testMessage(t, buf.Bytes, &m, nil)
}

func TestFunctionCall(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeFunctionCall)
	buf.AppendInt32(39)
	buf.AppendInt32(4321)
	buf.AppendInt16(1)
	buf.AppendInt16(int16(0))
	buf.AppendInt16(3)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("hello")...)
	buf.AppendInt32(1)
	buf.AppendByte([]byte(" ")...)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("world")...)
	buf.AppendInt16(int16(0))

	var m pgmsg.FunctionCall

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ObjectID)
		require.Len(t, m.ArgumentFormats, 1)
		require.Equal(t, []byte{0}, m.ArgumentFormats)
		require.Len(t, m.ArgumentValues, 3)
		require.Equal(t, [][]byte{
			[]byte("hello"),
			[]byte(" "),
			[]byte("world"),
		}, m.ArgumentValues)
		require.Equal(t, uint8(0), m.ResultFormat)
	})
}

func TestGSSENCRequest(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendInt32(8)
	buf.AppendInt32(pgmsg.CodeEncryptionRequest)

	var m pgmsg.GSSENCRequest

	testMessage(t, buf.Bytes, &m, nil)
}

func TestGSSResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeGSSResponse)
	buf.AppendInt32(15)
	buf.AppendByte([]byte("hello world")...)

	var m pgmsg.GSSResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, []byte("hello world"), m.Data)
	})
}

func TestParse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeParse)
	buf.AppendInt32(36)
	buf.AppendString("destination")
	buf.AppendString("query")
	buf.AppendInt16(3)
	buf.AppendInt32(3)
	buf.AppendInt32(2)
	buf.AppendInt32(1)

	var m pgmsg.Parse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "destination", m.DestinationStatementName)
		require.Equal(t, "query", m.Query)
		require.Len(t, m.ParameterDataTypes, 3)
		require.Equal(t, int32(3), m.ParameterDataTypes[0])
		require.Equal(t, int32(2), m.ParameterDataTypes[1])
		require.Equal(t, int32(1), m.ParameterDataTypes[2])
	})
}
