package pgwire_test

import (
	"gopsql/pgio"
	"gopsql/pgwire"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgBind(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeBind)
	buf.AppendInt32(73)
	buf.AppendString("hello world")
	buf.AppendString("lorem ipsum")
	buf.AppendInt16(3)
	buf.AppendInt16(
		int16(pgwire.FormatBinary),
		int16(pgwire.FormatBinary),
		int16(pgwire.FormatBinary),
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
		int16(pgwire.FormatBinary),
		int16(pgwire.FormatBinary),
		int16(pgwire.FormatBinary),
	)

	var m pgwire.MsgBind

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.DestinationName)

		require.Equal(t, "lorem ipsum", m.SourceName)

		require.Equal(t, []byte{
			pgwire.FormatBinary,
			pgwire.FormatBinary,
			pgwire.FormatBinary,
		}, m.ParameterFormatCodes)

		require.Equal(t, [][]byte{
			[]byte("lorem"),
			[]byte("ipsum"),
			[]byte("dolor"),
		}, m.ParameterData)

		require.Equal(t, []byte{
			pgwire.FormatBinary,
			pgwire.FormatBinary,
			pgwire.FormatBinary,
		}, m.ColumnFormatCodes)
	})
}

func TestMsgCancelRequest(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendInt32(23)
	buf.AppendInt32(pgwire.CodeCancelRequest)
	buf.AppendInt32(4321)
	buf.AppendByte([]byte("hello world")...)

	var m pgwire.MsgCancelRequest

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ProcessID)
		require.Equal(t, "hello world", string(m.SecretKey))
	})
}

func TestMsgClose(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeClose)
	buf.AppendInt32(17)
	buf.AppendInt8(int8(pgwire.ObjectPortal))
	buf.AppendString("hello world")

	var m pgwire.MsgClose

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, pgwire.ObjectPortal, m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestMsgCopyFail(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeCopyFail)
	buf.AppendInt32(16)
	buf.AppendString("hello world")

	var m pgwire.MsgCopyFail

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.Message)
	})
}

func TestMsgDescribe(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeDescribe)
	buf.AppendInt32(17)
	buf.AppendByte(pgwire.ObjectPortal)
	buf.AppendString("hello world")

	var m pgwire.MsgDescribe

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, pgwire.ObjectPortal, m.ObjectKind)
		require.Equal(t, "hello world", m.ObjectName)
	})
}

func TestMsgExecute(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeExecute)
	buf.AppendInt32(20)
	buf.AppendString("hello world")
	buf.AppendInt32(5)

	var m pgwire.MsgExecute

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.PortalName)
		require.Equal(t, int32(5), m.RowLimit)
	})
}

func TestMsgFlush(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeFlush)
	buf.AppendInt32(4)

	var m pgwire.MsgFlush

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgFunctionCall(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeFunctionCall)
	buf.AppendInt32(39)
	buf.AppendInt32(4321)
	buf.AppendInt16(1)
	buf.AppendInt16(int16(pgwire.FormatBinary))
	buf.AppendInt16(3)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("hello")...)
	buf.AppendInt32(1)
	buf.AppendByte([]byte(" ")...)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("world")...)
	buf.AppendInt16(int16(pgwire.FormatBinary))

	var m pgwire.MsgFunctionCall

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ObjectID)
		require.Len(t, m.ArgumentFormats, 1)
		require.Equal(t, []byte{pgwire.FormatBinary}, m.ArgumentFormats)
		require.Len(t, m.ArgumentValues, 3)
		require.Equal(t, [][]byte{
			[]byte("hello"),
			[]byte(" "),
			[]byte("world"),
		}, m.ArgumentValues)
		require.Equal(t, pgwire.FormatBinary, m.ResultFormat)
	})
}

func TestMsgGSSENCRequest(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendInt32(8)
	buf.AppendInt32(pgwire.CodeEncryptionRequest)

	var m pgwire.MsgGSSENCRequest

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgGSSResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeGSSResponse)
	buf.AppendInt32(15)
	buf.AppendByte([]byte("hello world")...)

	var m pgwire.MsgGSSResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, []byte("hello world"), m.Data)
	})
}

func TestMsgParse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeParse)
	buf.AppendInt32(36)
	buf.AppendString("destination")
	buf.AppendString("query")
	buf.AppendInt16(3)
	buf.AppendInt32(3)
	buf.AppendInt32(2)
	buf.AppendInt32(1)

	var m pgwire.MsgParse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "destination", m.DestinationStatementName)
		require.Equal(t, "query", m.Query)
		require.Len(t, m.ParameterDataTypes, 3)
		require.Equal(t, int32(3), m.ParameterDataTypes[0])
		require.Equal(t, int32(2), m.ParameterDataTypes[1])
		require.Equal(t, int32(1), m.ParameterDataTypes[2])
	})
}
