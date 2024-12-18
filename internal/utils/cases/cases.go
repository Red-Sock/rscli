package cases

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/Red-Sock/rscli/internal/utils/slices"
)

func SnakeToPascal(v string) string {
	parts := strings.Split(v, "_")
	for i := range parts {
		if slices.Contains(initialisms, parts[i]) {
			parts[i] = strings.ToUpper(parts[i])
		} else {
			if parts[i] == "" {
				continue
			}
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
	}

	return strings.Join(parts, "")
}

func KebabToSnake(v string) string {
	parts := strings.Split(v, "-")

	return strings.ToLower(strings.Join(parts, "_"))
}

var initialisms = []string{"acl", "api", "ascii", "cpu", "css", "dns",
	"eof", "guid", "html", "http", "https", "id",
	"ip", "json", "qps", "ram", "rpc", "sla",
	"smtp", "sql", "ssh", "tcp", "tls", "ttl",
	"udp", "ui", "gid", "uid", "uuid", "uri",
	"url", "utf8", "vm", "xml", "xmpp", "xsrf",
	"xss", "sip", "rtp", "amqp", "db", "ts"}

func ToPascal(newName string) string {
	pascalNameSB := bytes.Buffer{}
	nameRuned := []rune(newName)

	pascalNameSB.WriteRune(unicode.ToUpper(nameRuned[0]))

	nextUpper := false
	for _, n := range nameRuned[1:] {
		if n == '-' || n == '_' {
			nextUpper = true
			continue
		}

		if nextUpper {
			nextUpper = false
			n = unicode.ToUpper(n)
		}

		pascalNameSB.WriteRune(n)
	}

	return pascalNameSB.String()
}
