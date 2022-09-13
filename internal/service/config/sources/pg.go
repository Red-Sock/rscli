package sources

import (
	"github.com/Red-Sock/rscli/internal/service/config/model"
)

func DefaultPgPattern() *dbConf {
	d := &dbConf{}
	d.name = "postgres"
	d.user = "postgres"
	d.pwd = "postgres"
	d.host = "0.0.0.0"
	d.port = "5432"

	return d
}

type dbConf struct {
	ConnectionName string

	host string
	port string

	name string
	user string
	pwd  string
}

func (d *dbConf) GetParts(nl int) []model.Part {
	out := make([]model.Part, 6)

	if d.ConnectionName == "" {
		d.ConnectionName = "postgres"
	}

	out[0] = model.Part{NestingLevel: nl, Key: d.ConnectionName}
	out[1] = model.Part{NestingLevel: nl, Key: "name", Value: d.name}
	out[2] = model.Part{NestingLevel: nl, Key: "user", Value: d.user}
	out[3] = model.Part{NestingLevel: nl, Key: "pwd", Value: d.pwd}
	out[4] = model.Part{NestingLevel: nl, Key: "host", Value: d.host}
	out[5] = model.Part{NestingLevel: nl, Key: "port", Value: d.port}
	return out
}
