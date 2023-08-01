package _const

import (
	"strings"
)

const (
	PostgresConnectionString = pgConnectionStringPrefix + "%s:%s@%s:%d/%s"

	pgConnectionStringPrefix = "postgresql://"
)

func ParsePgConnectionString(cs string) (user, pwd, host, port, name string) {
	cs = cs[len(pgConnectionStringPrefix):]

	creds := strings.Split(cs, "@")
	{
		up := strings.Split(creds[0], ":")
		user, pwd = up[0], up[1]
	}

	{
		hpn := strings.Split(creds[1], "/")
		name = hpn[1]

		hp := strings.Split(hpn[0], ":")
		host, port = hp[0], hp[1]
	}

	return
}
