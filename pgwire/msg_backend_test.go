package pgwire_test

import (
	"gopsql/pgio"
	"gopsql/pgwire"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgBackendKeyData(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeBackendKeyData)
	buf.AppendInt32(19)
	buf.AppendInt32(4321)
	buf.AppendByte([]byte("hello world")...)

	var m pgwire.MsgBackendKeyData

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ProcessID)
		require.Equal(t, "hello world", string(m.SecretKey))
	})
}

func TestMsgBindComplete(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeBindComplete)
	buf.AppendInt32(4)

	var m pgwire.MsgBindComplete

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgCloseComplete(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeCloseComplete)
	buf.AppendInt32(4)

	var m pgwire.MsgCloseComplete

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgCommandComplete(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeCommandComplete)
	buf.AppendInt32(16)
	buf.AppendString("hello world")

	var m pgwire.MsgCommandComplete

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.Tag)
	})
}

func TestMsgCopyInResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeCopyInResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(pgwire.FormatBinary))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(pgwire.FormatBinary))
	buf.AppendInt16(int16(pgwire.FormatBinary))
	buf.AppendInt16(int16(pgwire.FormatBinary))

	var m pgwire.MsgCopyInResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int8(pgwire.FormatBinary), m.Format)
		require.Equal(t, []int16{
			int16(pgwire.FormatBinary),
			int16(pgwire.FormatBinary),
			int16(pgwire.FormatBinary),
		}, m.Columns)
	})
}

func TestMsgCopyOutResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeCopyOutResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(pgwire.FormatBinary))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(pgwire.FormatBinary))
	buf.AppendInt16(int16(pgwire.FormatBinary))
	buf.AppendInt16(int16(pgwire.FormatBinary))

	var m pgwire.MsgCopyOutResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int8(pgwire.FormatBinary), m.Format)
		require.Equal(t, []int16{
			int16(pgwire.FormatBinary),
			int16(pgwire.FormatBinary),
			int16(pgwire.FormatBinary),
		}, m.Columns)
	})
}

func TestMsgCopyBothResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeCopyBothResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(pgwire.FormatBinary))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(pgwire.FormatBinary))
	buf.AppendInt16(int16(pgwire.FormatBinary))
	buf.AppendInt16(int16(pgwire.FormatBinary))

	var m pgwire.MsgCopyBothResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int8(pgwire.FormatBinary), m.Format)
		require.Equal(t, []int16{
			int16(pgwire.FormatBinary),
			int16(pgwire.FormatBinary),
			int16(pgwire.FormatBinary),
		}, m.Columns)
	})
}

func TestMsgDataRow(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeDataRow)
	buf.AppendInt32(28)
	buf.AppendInt16(3)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("hello")...)
	buf.AppendInt32(-1)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("world")...)

	var m pgwire.MsgDataRow

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Columns, 3)
		require.Equal(t, []byte("hello"), m.Columns[0])
		require.Nil(t, m.Columns[1])
		require.Equal(t, []byte("world"), m.Columns[2])
	})
}

func TestMsgEmptyQueryResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeEmptyQueryResponse)
	buf.AppendInt32(4)

	var m pgwire.MsgEmptyQueryResponse

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgErrorResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeErrorResponse)
	buf.AppendInt32(19)
	buf.AppendByte('S')
	buf.AppendString("ERROR")
	buf.AppendByte('C')
	buf.AppendString("00000")
	buf.AppendByte(0)

	var m pgwire.MsgErrorResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Fields, 2)
		require.Len(t, m.Values, 2)
		require.Equal(t, byte('S'), m.Fields[0])
		require.Equal(t, "ERROR", m.Values[0])
		require.Equal(t, byte('C'), m.Fields[1])
		require.Equal(t, "00000", m.Values[1])
	})
}

func TestMsgFunctionCallResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeFunctionCallResponse)
	buf.AppendInt32(8)
	buf.AppendInt32(-1)

	var m pgwire.MsgFunctionCallResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Result, 0)
		require.Equal(t, []byte(nil), m.Result)
	})
}

func TestMsgNegotiateProtocolVersion(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeNegotiateProtocolVersion)
	buf.AppendInt32(24)
	buf.AppendInt32(1)
	buf.AppendInt32(2)
	buf.AppendString("hello")
	buf.AppendString("world")

	var m pgwire.MsgNegotiateProtocolVersion

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(1), m.MinorVersionSupported)
		require.Len(t, m.UnrecognizedOptions, 2)
		require.Equal(t, "hello", m.UnrecognizedOptions[0])
		require.Equal(t, "world", m.UnrecognizedOptions[1])
	})
}

func TestMsgNoData(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeNoData)
	buf.AppendInt32(4)

	var m pgwire.MsgNoData

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgNoticeResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.NewBuffer(nil)
	buf.AppendByte(pgwire.TypeNoticeResponse)
	buf.AppendInt32(32)
	buf.AppendByte(pgwire.FieldSeverity)
	buf.AppendString("ERROR")
	buf.AppendByte(pgwire.FieldCode)
	buf.AppendString("00000")
	buf.AppendByte(pgwire.FieldMessage)
	buf.AppendString("hello world")	
	buf.AppendByte(0)

	var m pgwire.MsgNoticeResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Fields, 3) 
		require.Len(t, m.Values, 3)
		require.Equal(t, pgwire.FieldSeverity, m.Fields[0])
		require.Equal(t, "ERROR", m.Values[0])
		require.Equal(t, pgwire.FieldCode, m.Fields[1])
		require.Equal(t, "00000", m.Values[1])
		require.Equal(t, pgwire.FieldMessage, m.Fields[2])
		require.Equal(t, "hello world", m.Values[2])
	})
}
