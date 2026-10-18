package pgmsg

import (
	"testing"

	"github.com/senojj/pgsql/pgproto/pgio"

	"github.com/stretchr/testify/require"
)

func TestBackendKeyData(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeBackendKeyData)
	buf.AppendInt32(19)
	buf.AppendInt32(4321)
	buf.AppendByte([]byte("hello world")...)

	var m BackendKeyData

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int32(4321), m.ProcessID)
		require.Equal(t, "hello world", string(m.SecretKey))
	})
}

func TestBindComplete(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeBindComplete)
	buf.AppendInt32(4)

	var m BindComplete

	testMessage(t, buf.Bytes, &m, nil)
}

func TestCloseComplete(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCloseComplete)
	buf.AppendInt32(4)

	var m CloseComplete

	testMessage(t, buf.Bytes, &m, nil)
}

func TestCommandComplete(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCommandComplete)
	buf.AppendInt32(16)
	buf.AppendString("hello world")

	var m CommandComplete

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, "hello world", m.Tag)
	})
}

func TestCopyInResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCopyInResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(0))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(0))
	buf.AppendInt16(int16(0))
	buf.AppendInt16(int16(0))

	var m CopyInResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int8(0), m.Format)
		require.Equal(t, []int16{
			int16(0),
			int16(0),
			int16(0),
		}, m.Columns)
	})
}

func TestCopyOutResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCopyOutResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(0))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(0))
	buf.AppendInt16(int16(0))
	buf.AppendInt16(int16(0))

	var m CopyOutResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int8(0), m.Format)
		require.Equal(t, []int16{
			int16(0),
			int16(0),
			int16(0),
		}, m.Columns)
	})
}

func TestCopyBothResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCopyBothResponse)
	buf.AppendInt32(13)
	buf.AppendInt8(int8(0))
	buf.AppendInt16(3)
	buf.AppendInt16(int16(0))
	buf.AppendInt16(int16(0))
	buf.AppendInt16(int16(0))

	var m CopyBothResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int8(0), m.Format)
		require.Equal(t, []int16{
			int16(0),
			int16(0),
			int16(0),
		}, m.Columns)
	})
}

func TestDataRow(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeDataRow)
	buf.AppendInt32(28)
	buf.AppendInt16(3)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("hello")...)
	buf.AppendInt32(-1)
	buf.AppendInt32(5)
	buf.AppendByte([]byte("world")...)

	var m DataRow

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Len(t, m.Columns, 3)
		require.Equal(t, []byte("hello"), m.Columns[0])
		require.Nil(t, m.Columns[1])
		require.Equal(t, []byte("world"), m.Columns[2])
	})
}

func TestEmptyQueryResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeEmptyQueryResponse)
	buf.AppendInt32(4)

	var m EmptyQueryResponse

	testMessage(t, buf.Bytes, &m, nil)
}

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeErrorResponse)
	buf.AppendInt32(19)
	buf.AppendByte('S')
	buf.AppendString("ERROR")
	buf.AppendByte('C')
	buf.AppendString("00000")
	buf.AppendByte(0)

	var m ErrorResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Len(t, m.Fields, 2)
		require.Len(t, m.Values, 2)
		require.Equal(t, byte('S'), m.Fields[0])
		require.Equal(t, "ERROR", m.Values[0])
		require.Equal(t, byte('C'), m.Fields[1])
		require.Equal(t, "00000", m.Values[1])
	})
}

func TestFunctionCallResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeFunctionCallResponse)
	buf.AppendInt32(8)
	buf.AppendInt32(-1)

	var m FunctionCallResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Len(t, m.Result, 0)
		require.Equal(t, []byte(nil), m.Result)
	})
}

func TestNegotiateProtocolVersion(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeNegotiateProtocolVersion)
	buf.AppendInt32(24)
	buf.AppendInt32(1)
	buf.AppendInt32(2)
	buf.AppendString("hello")
	buf.AppendString("world")

	var m NegotiateProtocolVersion

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int32(1), m.MinorVersionSupported)
		require.Len(t, m.UnrecognizedOptions, 2)
		require.Equal(t, "hello", m.UnrecognizedOptions[0])
		require.Equal(t, "world", m.UnrecognizedOptions[1])
	})
}

func TestNoData(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeNoData)
	buf.AppendInt32(4)

	var m NoData

	testMessage(t, buf.Bytes, &m, nil)
}

func TestNoticeResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeNoticeResponse)
	buf.AppendInt32(32)
	buf.AppendByte('A')
	buf.AppendString("ERROR")
	buf.AppendByte('B')
	buf.AppendString("00000")
	buf.AppendByte('C')
	buf.AppendString("hello world")
	buf.AppendByte(0)

	var m NoticeResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Len(t, m.Fields, 3)
		require.Len(t, m.Values, 3)
		require.Equal(t, uint8('A'), m.Fields[0])
		require.Equal(t, "ERROR", m.Values[0])
		require.Equal(t, uint8('B'), m.Fields[1])
		require.Equal(t, "00000", m.Values[1])
		require.Equal(t, uint8('C'), m.Fields[2])
		require.Equal(t, "hello world", m.Values[2])
	})
}

func TestNotificationResponse(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeNotificationResponse)
	buf.AppendInt32(26)
	buf.AppendInt32(123)
	buf.AppendString("channel_a")
	buf.AppendString("payload")

	var m NotificationResponse

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, int32(123), m.ProcessID)
		require.Equal(t, "channel_a", m.Channel)
		require.Equal(t, "payload", m.Payload)
	})
}

func TestParameterDescription(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeParameterDescription)
	buf.AppendInt32(18)
	buf.AppendInt16(3)
	buf.AppendInt32(3)
	buf.AppendInt32(2)
	buf.AppendInt32(1)

	var m ParameterDescription

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Len(t, m.Parameters, 3)
		require.Equal(t, int32(3), m.Parameters[0])
		require.Equal(t, int32(2), m.Parameters[1])
		require.Equal(t, int32(1), m.Parameters[2])
	})
}
