package sources

import (
	"github.com/Red-Sock/rscli/internal/service/config/model"
)

type rds struct {
	connectionName string

	host string
	port string

	user string
	pwd  string
	db   int
}

func DefaultRdsPattern() *rds {
	r := &rds{}

	r.user = ""
	r.pwd = ""

	r.host = "0.0.0.0"
	r.port = "6379"

	return r
}

func (r *rds) GetParts(nl int) []model.Part {

	out := make([]model.Part, 1, 4)

	if r.connectionName
	return nil
}
