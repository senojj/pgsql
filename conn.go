package pgsql

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/senojj/pgsql/pgmsg"
)

const (
	DefaultHost               string = "localhost"
	DefaultOption             string = ""
	DefaultChannelBinding     string = "disable"
	DefaultTargetSessionAttrs string = "any"
	DefaultLoadBalanceHosts   string = "disable"
	DefaultSSLMode            string = "disable"
	DefaultSSLCertMode        string = "disable"
	DefaultSSLNegotiation     string = "postgres"
	DefaultGSSMode            string = "disable"
)

type ConnInfo struct {
	Service                 string
	User                    string
	Password                string
	Passfile                string
	ChannelBinding          string
	ConnectTimeout          string
	DBName                  string
	Host                    string
	HostAddr                string
	Port                    string
	ClientEncoding          string
	Options                 string
	ApplicationName         string
	FallbackApplicationName string
	Keepalives              string
	KeepalivesIdle          string
	KeepalivesInterval      string
	KeepalivesCount         string
	TCPUserTimeout          string
	SSLMode                 string
	SSLNegotiation          string
	SSLCompression          string
	SSLCert                 string
	SSLKey                  string
	SSLCertMode             string
	SSLPassword             string
	SSLRootCert             string
	SSLCrl                  string
	SSLCrlDir               string
	SSLSNI                  string
	RequirePeer             string
	RequireAuth             string
	MinProtocolVersion      string
	MaxProtocolVersion      string
	SSLMinProtocolVersion   string
	SSLMaxProtocolVersion   string
	GSSEncMode              string
	KrbSrvName              string
	GSSLib                  string
	GSSDelegation           string
	Replication             string
	TargetSessionAttrs      string
	LoadBalanceHosts        string
	ScramClientKey          string
	ScramServerKey          string
	OAuthIssuer             string
	OAuthClientID           string
	OAuthClientSecret       string
	OAuthScope              string
	OAuthCAFile             string
	SSLKeyLogFile           string
}

type ConnInfoOption struct {
	Keyword string
	Envvar  string
	Default string
	Display byte // '' as-is, '*' password, 'D' debug
	Value   string
	offset  uintptr
}

var connInfoOptions = [51]ConnInfoOption{
	{"service", "PGSERVICE", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).Service)},

	{"user", "PGUSER", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).User)},

	{"password", "PGPASSWORD", "", '*', "",
		unsafe.Offsetof((*ConnInfo)(nil).Password)},

	{"passfile", "PGPASSFILE", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).Passfile)},

	{"channel_binding", "PGCHANNELBINDING", DefaultChannelBinding, 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).ChannelBinding)},

	{"connect_timeout", "PGCONNECT_TIMEOUT", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).ConnectTimeout)},

	{"dbname", "PGDATABASE", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).DBName)},

	{"host", "PGHOST", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).Host)},

	{"hostaddr", "PGHOSTADDR", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).HostAddr)},

	{"port", "PGPORT", "5432", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).Port)},

	{"client_encoding", "PGCLIENTENCODING", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).ClientEncoding)},

	{"options", "PGOPTIONS", DefaultOption, 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).Options)},

	{"application_name", "PGAPPNAME", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).ApplicationName)},

	{"fallback_application_name", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).FallbackApplicationName)},

	{"keepalives", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).Keepalives)},

	{"keepalives_idle", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).KeepalivesIdle)},

	{"keepalives_interval", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).KeepalivesInterval)},

	{"keepalives_count", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).KeepalivesCount)},

	{"tcp_user_timeout", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).TCPUserTimeout)},

	{"sslmode", "PGSSLMODE", DefaultSSLMode, 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLMode)},

	{"sslnegotiation", "PGSSLNEGOTIATION", DefaultSSLNegotiation, 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLNegotiation)},

	{"sslcompression", "PGSSLCOMPRESSION", "0", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLCompression)},

	{"sslcert", "PGSSLCERT", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLCert)},

	{"sslkey", "PGSSLKEY", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLKey)},

	{"sslcertmode", "PGSSLCERTMODE", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLCertMode)},

	{"sslpassword", "", "", '*', "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLPassword)},

	{"sslrootcert", "PGSSLROOTCERT", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLRootCert)},

	{"sslcrl", "PGSSLCRL", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLCrl)},

	{"sslcrldir", "PGSSLCRLDIR", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLCrlDir)},

	{"sslsni", "PGSSLSNI", "1", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLSNI)},

	{"requirepeer", "PGREQUIREPEER", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).RequirePeer)},

	{"require_auth", "PGREQUIREAUTH", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).RequireAuth)},

	{"min_protocol_version", "PGMINPROTOCOLVERSION", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).MinProtocolVersion)},

	{"max_protocol_version", "PGMAXPROTOCOLVERSION", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).MaxProtocolVersion)},

	{"ssl_min_protocol_version", "PGSSLMINPROTOCOLVERSION", "TLSv1.2", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLMinProtocolVersion)},

	{"ssl_max_protocol_version", "PGSSLMAXPROTOCOLVERSION", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLMaxProtocolVersion)},

	{"gssencmode", "PGGSSENCMODE", DefaultGSSMode, 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).GSSEncMode)},

	{"krbsrvname", "PGKRBSRVNAME", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).KrbSrvName)},

	{"gsslib", "PGGSSLIB", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).GSSLib)},

	{"gssdelegation", "PGGSSDELEGATION", "0", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).GSSDelegation)},

	{"replication", "", "", 'D', "",
		unsafe.Offsetof((*ConnInfo)(nil).Replication)},

	{
		"target_session_attrs",
		"PGTARGETSESSIONATTRS",
		DefaultTargetSessionAttrs,
		0,
		"",
		unsafe.Offsetof((*ConnInfo)(nil).TargetSessionAttrs),
	},

	{
		"load_balance_hosts",
		"PGLOADBALANCEHOSTS",
		DefaultLoadBalanceHosts,
		0,
		"",
		unsafe.Offsetof((*ConnInfo)(nil).LoadBalanceHosts),
	},

	{"scram_client_key", "", "", 'D', "",
		unsafe.Offsetof((*ConnInfo)(nil).ScramClientKey)},

	{"scram_server_key", "", "", 'D', "",
		unsafe.Offsetof((*ConnInfo)(nil).ScramServerKey)},

	{"oauth_issuer", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).OAuthIssuer)},

	{"oauth_client_id", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).OAuthClientID)},

	{"oauth_client_secret", "", "", '*', "",
		unsafe.Offsetof((*ConnInfo)(nil).OAuthClientSecret)},

	{"oauth_scope", "", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).OAuthScope)},

	{"oauth_ca_file", "PGOAUTHCAFILE", "", 0, "",
		unsafe.Offsetof((*ConnInfo)(nil).OAuthCAFile)},

	{"sslkeylogfile", "", "", 'D', "",
		unsafe.Offsetof((*ConnInfo)(nil).SSLKeyLogFile)},
}

func connInfoFind(options []ConnInfoOption, keyword string) *ConnInfoOption {
	for i := 0; i < len(options); i++ {
		option := &options[i]

		if option.Keyword == keyword {
			return option
		}
	}
	return nil
}

func connInfoAddDefaults(options []ConnInfoOption) {
	var sslModeDefault, sslRootCert *ConnInfoOption

	for i := 0; i < len(options); i++ {
		option := &options[i]

		if option.Keyword == "sslrootcert" {
			sslRootCert = option
		}

		if option.Value != "" {
			continue
		}

		if option.Envvar != "" {
			option.Value = os.Getenv(option.Envvar)
			if option.Value != "" {
				continue
			}
		}

		if option.Keyword == "sslmode" {
			requireSSL := os.Getenv("PGREQUIRESSL")

			if requireSSL != "" && requireSSL[0] == '1' {
				option.Value = "require"
				continue
			}

			sslModeDefault = option
		}

		if option.Default != "" {
			option.Value = option.Default
			continue
		}

		if option.Keyword == "user" {
			usr, err := user.Current()
			if err == nil {
				option.Value = usr.Username
			}
			continue
		}
	}

	if sslModeDefault != nil && sslRootCert != nil {
		if sslRootCert.Value == "system" {
			sslModeDefault.Value = "verify-full"
		}
	}
}

const (
	longURIDesignator  = "postgresql://"
	shortURIDesignator = "postgres://"
)

func hasURIPrefix(connstr string) bool {
	if strings.HasPrefix(connstr, longURIDesignator) {
		return true
	}
	return strings.HasPrefix(connstr, shortURIDesignator)
}

func connInfoWrite(options []ConnInfoOption, keyword, value string) error {
	if keyword == "requiressl" {
		keyword = "sslmode"
		if value[0] == '1' {
			value = "require"
		} else {
			value = "prefer"
		}
	}

	option := connInfoFind(options, keyword)
	if option == nil {
		return fmt.Errorf("invalid connection option %s", keyword)
	}

	option.Value = value
	return nil
}

func connInfoURIParseOptions(
	options []ConnInfoOption,
	uri string,
) error {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	username := u.User.Username()

	err = connInfoWrite(options[:], "user", username)
	if err != nil {
		return err
	}

	if password, ok := u.User.Password(); ok {
		err = connInfoWrite(options[:], "password", password)
		if err != nil {
			return err
		}
	}

	var head, tail string
	tail = u.Hostname()

	hosts := make([]string, 0, bytes.Count([]byte(tail), []byte{','}))
	ports := make([]string, 0, bytes.Count([]byte(tail), []byte{':'}))

	for ; tail != ""; head, tail, _ = strings.Cut(tail, ",") {
		host, port, found := strings.Cut(head, ":")
		hosts = append(hosts, host)
		if found {
			ports = append(ports, port)
		}
	}

	if len(hosts) > 0 {
		err = connInfoWrite(options[:], "host", strings.Join(hosts, ","))
		if err != nil {
			return err
		}
	}

	if len(ports) > 0 {
		err = connInfoWrite(options, "port", strings.Join(ports, ","))
		if err != nil {
			return err
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
	if database != "" {
		err = connInfoWrite(options, "dbname", database)
		if err != nil {
			return err
		}
	}

	return nil
}

func connInfoURIParse(uri string, useDefaults bool) ([]ConnInfoOption, error) {
	options := connInfoOptions

	err := connInfoURIParseOptions(options[:], uri)
	if err != nil {
		return nil, err
	}

	if useDefaults {
		connInfoAddDefaults(options[:])
	}

	return options[:], nil
}

func parseConnectionString(
	connstr string,
	useDefaults bool,
) ([]ConnInfoOption, error) {
	if hasURIPrefix(connstr) {
		return connInfoURIParse(connstr, useDefaults)
	}
	return nil, nil
}

func ConnDefaults() []ConnInfoOption {
	options := connInfoOptions
	connInfoAddDefaults(options[:])
	return options[:]
}

type ConnStatus int32

const (
	ConnectionOK ConnStatus = iota
	ConnectionBad
	ConnectionStarted
	ConnectionMade
	ConnectionAwaitingResponse
	ConnectionAuthOK
	ConnectionSSLStartup
	ConnectionGSSSTartup
	ConnectionCheckWritable
	ConnectionCheckStandby
	ConnectionConsume
)

type ExecStatus int

const (
	StatusCmdOK ExecStatus = iota
	StatusTuplesOK
	StatusFatalError
)

type Conn struct {
	info      ConnInfo
	raw       net.Conn
	err       error
	decoder   *pgmsg.Decoder
	encoder   *pgmsg.Encoder
	protocol  int32
	status    ConnStatus
	processID int32
}

func (c *Conn) fill(options []ConnInfoOption) {
	for i := 0; i < len(connInfoOptions); i++ {
		option := &connInfoOptions[i]

		if option.offset < 0 {
			continue
		}

		tmp := connInfoFind(options, option.Keyword)
		if tmp == nil || tmp.Value == "" {
			continue
		}

		*(*string)(
			unsafe.Pointer(uintptr(unsafe.Pointer(&c.info)) + option.offset),
		) = tmp.Value
	}
}

func Connect(ctx context.Context, connstr string) *Conn {
	var conn Conn
	conn.status = ConnectionBad
	conn.protocol = ProtocolVersion3_2

	options, err := parseConnectionString(connstr, true)
	if err != nil {
		conn.err = err
		return &conn
	}

	conn.fill(options)

	/*
		c, err := net.Dial("tcp", address)
		if err != nil {
			conn.err = err
			return &conn
		}

		conn.raw = c
		conn.encoder = pgmsg.NewEncoder(c)
		conn.decoder = pgmsg.NewDecoder(c)

		startup(ctx, &conn)
	*/
	return &conn
}

func (c *Conn) Status() ConnStatus { return ConnectionOK }

func (c *Conn) Err() error { return nil }

// SocketPoll polls a connection's underlying socket descriptor, optionally
// with a timeout. The timeout is specified by `deadline`, which is the time
// to stop waiting. Timeout is infinite if `deadline` is the zero value of a
// `time.Time`.
//
// This function returns `true` if the socket descriptor has data available
// for reading, `false` if a timeout occurred, or `false` and a non-nil error
// if an error occurred.
//
// SocketPoll is implemented using a 0-byte read on the connection's underlying
// `net.Conn`.
func (c *Conn) SocketPoll(deadline time.Time) (bool, error) {
	if err := c.raw.SetReadDeadline(deadline); err != nil {
		return false, err
	}
	defer c.raw.SetReadDeadline(time.Time{})

	_, err := c.raw.Read([]byte{})
	if err == nil {
		return true, nil
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return false, nil
	}

	return false, err
}

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
