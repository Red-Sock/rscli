package help

import "strings"

const (
	header = `
    RedSock CLI tool
========================
`
	options = `
configure, cfg, c - creates configuration file (by default - API application with PostgresSQL connection)  
`
)

type helpTool struct {
}

func NewHelpTool() *helpTool {
	return &helpTool{}
}

func (h *helpTool) Run(_ []string) string {
	return FormMessage()
}

func (h *helpTool) HelpMessage() string {
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
