package help

import "strings"

const (
	Header = `
========================
    RedSock CLI tool
========================
`
	Options = `
ui - opens visual menu
get {link to plugin} - downloads plugin from git
fix - sets environment variable for plugins dir
del {link to plugin} - deletes plugin from plugin dir. "del github.com/SomeUser/* - to delete all packages started with given path"
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
