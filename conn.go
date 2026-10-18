package pgsql

import (
	"context"
	"net"
	"net/url"
	"strconv"

	"github.com/senojj/pgsql/pgproto"
)

func connect(ctx context.Context, s string) error {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return err
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
	if database == "" {
		database = "postgres"
	}

	address := hostname + ":" + strconv.FormatUint(uint64(port), 10)

	user := u.User.Username()

	c, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	frontend := pgproto.NewFrontend(c)

	err = frontend.Startup(ctx, pgproto.ProtocolVersion3_2, map[string]string{
		pgproto.ParamUser:     user,
		pgproto.ParamDatabase: database,
	})
	if err != nil {
		return err
	}
	return nil
}
