package pgmsg_test

import (
	"testing"

	"github.com/senojj/pgsql/pgmsg"
	"github.com/senojj/pgsql/pgio"

	"github.com/stretchr/testify/require"
)

func TestCopyData(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeCopyData)
	buf.AppendInt32(15)
	buf.AppendByte([]byte("hello world")...)

	var m pgmsg.CopyData

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, []byte("hello world"), m.Data)
	})
}

func TestCopyDone(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(pgmsg.TypeCopyDone)
	buf.AppendInt32(4)

	var m pgmsg.CopyDone

	testMessage(t, buf.Bytes, &m, nil)
}
