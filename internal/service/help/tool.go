package help

import "strings"

const Command = "help"

const (
	header = `
    RedSock CLI tool
========================
`
	options = `
configure, cfg, c - creates configuration file (by default - API application with PostgresSQL connection)  
`
)

func Run() string {
	return FormMessage()
}

func FormMessage(additionalInfo ...string) string {
	messageSB := &strings.Builder{}
	messageSB.WriteString(header)

	messageSB.WriteString(options)

	for _, item := range additionalInfo {
		messageSB.WriteString(item)
	}

	return messageSB.String()
}
