package pgsql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgCopyData(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCopyData)
	buf.AppendInt32(15)
	buf.AppendByte([]byte("hello world")...)

	var m msgCopyData

	testMessage(t, buf.Bytes(), &m, func(t *testing.T) {
		require.Equal(t, []byte("hello world"), m.Data)
	})
}

func TestMsgCopyDone(t *testing.T) {
	t.Parallel()

	buf := buffer{}
	buf.AppendByte(typeCopyDone)
	buf.AppendInt32(4)

	var m msgCopyDone

	testMessage(t, buf.Bytes(), &m, nil)
}
