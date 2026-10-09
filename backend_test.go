package pgsql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgBackendKeyData(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeBackendKeyData)
	buf.AppendInt32(19)
	buf.AppendInt32(4321)
	buf.AppendByte([]byte("hello world")...)

	var m msgBackendKeyData

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ProcessID)
		require.Equal(t, "hello world", string(m.SecretKey))
	})
}

func TestMsgBindComplete(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeBindComplete)
	buf.AppendInt32(4)

	var m msgBindComplete

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgCloseComplete(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCloseComplete)
	buf.AppendInt32(4)

	var m msgCloseComplete

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgCommandComplete(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCommandComplete)
	buf.AppendInt32(16)
	buf.AppendString("hello world")

	var m msgCommandComplete

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.Tag)
	})
}

func TestMsgCopyInResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCopyInResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(formatBinary))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(formatBinary))
	buf.AppendInt16(int16(formatBinary))
	buf.AppendInt16(int16(formatBinary))

	var m msgCopyInResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int8(formatBinary), m.Format)
		require.Equal(t, []int16{
			int16(formatBinary),
			int16(formatBinary),
			int16(formatBinary),
		}, m.Columns)
	})
}

func TestMsgCopyOutResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCopyOutResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(formatBinary))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(formatBinary))
	buf.AppendInt16(int16(formatBinary))
	buf.AppendInt16(int16(formatBinary))

	var m msgCopyOutResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int8(formatBinary), m.Format)
		require.Equal(t, []int16{
			int16(formatBinary),
			int16(formatBinary),
			int16(formatBinary),
		}, m.Columns)
	})
}

func TestMsgCopyBothResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCopyBothResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(formatBinary))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(formatBinary))
	buf.AppendInt16(int16(formatBinary))
	buf.AppendInt16(int16(formatBinary))

	var m msgCopyBothResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int8(formatBinary), m.Format)
		require.Equal(t, []int16{
			int16(formatBinary),
			int16(formatBinary),
			int16(formatBinary),
		}, m.Columns)
	})
}

func TestMsgDataRow(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeDataRow)
	buf.AppendInt32(28)
	buf.AppendInt16(3)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("hello")...)
	buf.AppendInt32(-1)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("world")...)

	var m msgDataRow

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Columns, 3)
		require.Equal(t, []byte("hello"), m.Columns[0])
		require.Nil(t, m.Columns[1])
		require.Equal(t, []byte("world"), m.Columns[2])
	})
}

func TestMsgEmptyQueryResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeEmptyQueryResponse)
	buf.AppendInt32(4)

	var m msgEmptyQueryResponse

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgErrorResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeErrorResponse)
	buf.AppendInt32(19)
	buf.AppendByte('S')
	buf.AppendString("ERROR")
	buf.AppendByte('C')
	buf.AppendString("00000")
	buf.AppendByte(0)

	var m msgErrorResponse

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

	buf := buffer{}
	buf.AppendByte(typeFunctionCallResponse)
	buf.AppendInt32(8)
	buf.AppendInt32(-1)

	var m msgFunctionCallResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Result, 0)
		require.Equal(t, []byte(nil), m.Result)
	})
}

func TestMsgNegotiateProtocolVersion(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeNegotiateProtocolVersion)
	buf.AppendInt32(24)
	buf.AppendInt32(1)
	buf.AppendInt32(2)
	buf.AppendString("hello")
	buf.AppendString("world")

	var m msgNegotiateProtocolVersion

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(1), m.MinorVersionSupported)
		require.Len(t, m.UnrecognizedOptions, 2)
		require.Equal(t, "hello", m.UnrecognizedOptions[0])
		require.Equal(t, "world", m.UnrecognizedOptions[1])
	})
}

func TestMsgNoData(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeNoData)
	buf.AppendInt32(4)

	var m msgNoData

	testMessage(t, buf.Bytes(), &m, nil)
}

func TestMsgNoticeResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeNoticeResponse)
	buf.AppendInt32(32)
	buf.AppendByte(fieldSeverity)
	buf.AppendString("ERROR")
	buf.AppendByte(fieldCode)
	buf.AppendString("00000")
	buf.AppendByte(fieldMessage)
	buf.AppendString("hello world")
	buf.AppendByte(0)

	var m msgNoticeResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Fields, 3)
		require.Len(t, m.Values, 3)
		require.Equal(t, fieldSeverity, m.Fields[0])
		require.Equal(t, "ERROR", m.Values[0])
		require.Equal(t, fieldCode, m.Fields[1])
		require.Equal(t, "00000", m.Values[1])
		require.Equal(t, fieldMessage, m.Fields[2])
		require.Equal(t, "hello world", m.Values[2])
	})
}

func TestMsgNotificationResponse(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeNotificationResponse)
	buf.AppendInt32(26)
	buf.AppendInt32(123)
	buf.AppendString("channel_a")
	buf.AppendString("payload")

	var m msgNotificationResponse

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, int32(123), m.ProcessID)
		require.Equal(t, "channel_a", m.Channel)
		require.Equal(t, "payload", m.Payload)
	})
}

func TestMsgParameterDescription(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeParameterDescription)
	buf.AppendInt32(18)
	buf.AppendInt16(3)
	buf.AppendInt32(3)
	buf.AppendInt32(2)
	buf.AppendInt32(1)

	var m msgParameterDescription

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Len(t, m.Parameters, 3)
		require.Equal(t, int32(3), m.Parameters[0])
		require.Equal(t, int32(2), m.Parameters[1])
		require.Equal(t, int32(1), m.Parameters[2])
	})
}
