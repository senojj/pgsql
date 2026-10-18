package pgmsg

import (
	"testing"

	"github.com/senojj/pgsql/pgproto/pgio"

	"github.com/stretchr/testify/require"
)

func TestBind(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeBind)
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

	var m Bind

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
	buf.AppendInt32(codeCancelRequest)
	buf.AppendInt32(4321)
	buf.AppendByte([]byte("hello world")...)

	var m CancelRequest

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ProcessID)
		require.Equal(t, "hello world", string(m.SecretKey))
	})
}

func TestClose(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeClose)
	buf.AppendInt32(17)
	buf.AppendInt8(int8(1))
	buf.AppendString("hello world")

	var m Close

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, uint8(1), m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestCopyFail(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCopyFail)
	buf.AppendInt32(16)
	buf.AppendString("hello world")

	var m CopyFail

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.Message)
	})
}

func TestDescribe(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeDescribe)
	buf.AppendInt32(17)
	buf.AppendByte(1)
	buf.AppendString("hello world")

	var m Describe

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, uint8(1), m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestExecute(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeExecute)
	buf.AppendInt32(20)
	buf.AppendString("hello world")
	buf.AppendInt32(5)

	var m Execute

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.PortalName)
		require.Equal(t, int32(5), m.RowLimit)
	})
}

func TestFlush(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeFlush)
	buf.AppendInt32(4)

	var m Flush

	testMessage(t, buf.Bytes, &m, nil)
}

func TestFunctionCall(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeFunctionCall)
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

	var m FunctionCall

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
	buf.AppendInt32(codeEncryptionRequest)

	var m GSSENCRequest

	testMessage(t, buf.Bytes, &m, nil)
}

func TestGSSResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeGSSResponse)
	buf.AppendInt32(15)
	buf.AppendByte([]byte("hello world")...)

	var m GSSResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, []byte("hello world"), m.Data)
	})
}

func TestParse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeParse)
	buf.AppendInt32(36)
	buf.AppendString("destination")
	buf.AppendString("query")
	buf.AppendInt16(3)
	buf.AppendInt32(3)
	buf.AppendInt32(2)
	buf.AppendInt32(1)

	var m Parse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "destination", m.DestinationStatementName)
		require.Equal(t, "query", m.Query)
		require.Len(t, m.ParameterDataTypes, 3)
		require.Equal(t, int32(3), m.ParameterDataTypes[0])
		require.Equal(t, int32(2), m.ParameterDataTypes[1])
		require.Equal(t, int32(1), m.ParameterDataTypes[2])
	})
}
