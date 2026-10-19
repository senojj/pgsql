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
	raw net.Conn
	decoder *pgmsg.Decoder
	encoder *pgmsg.Encoder
	params map[string]string
	protocol int32
	status ConnStatus
	err error
	processID int32
}

func Connect(ctx context.Context, info string) *Conn {
	var conn Conn
	conn.status = ConnectionBad
	conn.version = ProtocolVersion3_2

	u, err := url.ParseRequestURI(s)
	if err != nil {
		conn.err = err
		return &conn
	}

	hostname := u.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}

	port := uint64(5432)
	portStr := u.Port()
	if portStr != "" {
		port, err = strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			conn.err = err
			return &conn
		}
	}

	database := u.Path
	length := len(database)
	if length > 0 && database[0] == '/' {
		database = database[1:]
		length--
	}
	if length > 0 && database[length-1] == '/' {
		database = database[:length-1]
	}
	if database == "" {
		database = "postgres"
	}

	address := hostname + ":" + strconv.FormatUint(uint64(port), 10)

	user := u.User.Username()

	c, err := net.Dial("tcp", address)
	if err != nil {
		conn.err = err
		return &conn
	}

	conn.raw = c
	conn.encoder = pgmsg.NewEncoder(c)
	conn.decoder = pgmsg.NewDecoder(c)
	conn.params = map[string]string{
		ParamUser: user,
		ParamDatabase: database,
	}

	startup(ctx, &conn)
	return &conn
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
