package pgproto

type ConnStatus int

const (
	ConnectionOK ConnStatus = iota
	ConnectionBad
)

type ExecStatus int

const (
	StatusCmdOK ExecStatus = iota
	StatusTuplesOK
	StatusFatalError
)

type Conn struct {
}

func Connect(info string) *Conn {

}

func (c *Conn) Status() ConnStatus { return ConnectionOK }

func (c *Conn) Err() error { return nil }

func (c *Conn) Finish() {}

func (c *Conn) Exec(query string) *Result { return nil }

func (c *Conn) SendQuery(query string) bool { return false }

func (c *Conn) GetResult() *Result { return nil }

type Result struct {
}

func (r *Result) Status() ExecStatus { return StatusCmdOK }

func (r *Result) NTuples() int { return 0 }

func (r *Result) NFields() int { return 0 }

func (r *Result) FName(idx int) string { return "" }

func (r *Result) GetValue(row, field int) string { return "" }

func (r *Result) IsNull(row, field int) bool { return false }

func (r *Result) CmdStatus() string { return "" }

func (r *Result) Clear() {}
