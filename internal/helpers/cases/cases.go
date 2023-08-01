package cases

import (
	"strings"

	"github.com/Red-Sock/rscli/internal/helpers/slices"
)

func SnakeToCamel(v string) string {
	parts := strings.Split(v, "_")
	for i := range parts {
		if slices.Contains(initialisms, parts[i]) {
			parts[i] = strings.ToUpper(parts[i])
		} else {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
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
