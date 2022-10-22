package help

import "strings"

const Command = "help"

const (
	Header = `
========================
    RedSock CLI tool
========================
`
	Options = `
configure, cfg, c - creates configuration file (by default - API application with PostgresSQL connection)
ui - opens visual menu
`
)

func Run() string {
	return FormMessage()
}

func FormMessage(additionalInfo ...string) string {
	messageSB := &strings.Builder{}
	messageSB.WriteString(Header)

	messageSB.WriteString(Options)

	for _, item := range additionalInfo {
		messageSB.WriteString(item)
	}

	return messageSB.String()
}
