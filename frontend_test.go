package pgsql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgBind(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeBind)
	buf.AppendInt32(73)
	buf.AppendString("hello world")
	buf.AppendString("lorem ipsum")
	buf.AppendInt16(3)
	buf.AppendInt16(
		int16(formatBinary),
		int16(formatBinary),
		int16(formatBinary),
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
		int16(formatBinary),
		int16(formatBinary),
		int16(formatBinary),
	)

	var m msgBind

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.DestinationName)

		require.Equal(t, "lorem ipsum", m.SourceName)

		require.Equal(t, []byte{
			formatBinary,
			formatBinary,
			formatBinary,
		}, m.ParameterFormatCodes)

		require.Equal(t, [][]byte{
			[]byte("lorem"),
			[]byte("ipsum"),
			[]byte("dolor"),
		}, m.ParameterData)

		require.Equal(t, []byte{
			formatBinary,
			formatBinary,
			formatBinary,
		}, m.ColumnFormatCodes)
	})
}

func TestMsgCancelRequest(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendInt32(23)
	buf.AppendInt32(codeCancelRequest)
	buf.AppendInt32(4321)
	buf.AppendByte([]byte("hello world")...)

	var m msgCancelRequest

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ProcessID)
		require.Equal(t, "hello world", string(m.SecretKey))
	})
}

func TestMsgClose(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeClose)
	buf.AppendInt32(17)
	buf.AppendInt8(int8(objectPortal))
	buf.AppendString("hello world")

	var m msgClose

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, objectPortal, m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestMsgCopyFail(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCopyFail)
	buf.AppendInt32(16)
	buf.AppendString("hello world")

	var m msgCopyFail

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.Message)
	})
}

func TestMsgDescribe(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeDescribe)
	buf.AppendInt32(17)
	buf.AppendByte(objectPortal)
	buf.AppendString("hello world")

	var m msgDescribe

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, objectPortal, m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestMsgExecute(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeExecute)
	buf.AppendInt32(20)
	buf.AppendString("hello world")
	buf.AppendInt32(5)

	var m msgExecute

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.PortalName)
		require.Equal(t, int32(5), m.RowLimit)
	})
}

func TestMsgFlush(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeFlush)
	buf.AppendInt32(4)

	var m msgFlush

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgFunctionCall(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeFunctionCall)
	buf.AppendInt32(39)
	buf.AppendInt32(4321)
	buf.AppendInt16(1)
	buf.AppendInt16(int16(formatBinary))
	buf.AppendInt16(3)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("hello")...)
	buf.AppendInt32(1)
	buf.AppendByte([]byte(" ")...)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("world")...)
	buf.AppendInt16(int16(formatBinary))

	var m msgFunctionCall

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ObjectID)
		require.Len(t, m.ArgumentFormats, 1)
		require.Equal(t, []byte{formatBinary}, m.ArgumentFormats)
		require.Len(t, m.ArgumentValues, 3)
		require.Equal(t, [][]byte{
			[]byte("hello"),
			[]byte(" "),
			[]byte("world"),
		}, m.ArgumentValues)
		require.Equal(t, formatBinary, m.ResultFormat)
	})
}

func TestMsgGSSENCRequest(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendInt32(8)
	buf.AppendInt32(codeEncryptionRequest)

	var m msgGSSENCRequest

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgGSSResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeGSSResponse)
	buf.AppendInt32(15)
	buf.AppendByte([]byte("hello world")...)

	var m msgGSSResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, []byte("hello world"), m.Data)
	})
}

func TestMsgParse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeParse)
	buf.AppendInt32(36)
	buf.AppendString("destination")
	buf.AppendString("query")
	buf.AppendInt16(3)
	buf.AppendInt32(3)
	buf.AppendInt32(2)
	buf.AppendInt32(1)

	var m msgParse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "destination", m.DestinationStatementName)
		require.Equal(t, "query", m.Query)
		require.Len(t, m.ParameterDataTypes, 3)
		require.Equal(t, int32(3), m.ParameterDataTypes[0])
		require.Equal(t, int32(2), m.ParameterDataTypes[1])
		require.Equal(t, int32(1), m.ParameterDataTypes[2])
	})
}
