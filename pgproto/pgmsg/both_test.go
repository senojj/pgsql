package pgmsg

import (
	"testing"

	"github.com/senojj/pgsql/pgproto/pgio"

	"github.com/stretchr/testify/require"
)

func TestCopyData(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCopyData)
	buf.AppendInt32(15)
	buf.AppendByte([]byte("hello world")...)

	var m CopyData

	testMessage(t, buf.Bytes, &m, func(t *testing.T) {
		require.Equal(t, []byte("hello world"), m.Data)
	})
}

func TestCopyDone(t *testing.T) {
	t.Parallel()

	buf := pgio.Buffer{}
	buf.AppendByte(typeCopyDone)
	buf.AppendInt32(4)

	var m CopyDone

	testMessage(t, buf.Bytes, &m, nil)
}
