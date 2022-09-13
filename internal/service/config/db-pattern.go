package config

import (
	"github.com/Red-Sock/rscli/internal/service/config/model"
	"github.com/Red-Sock/rscli/internal/service/config/sources"
	"strings"
)

const (
	sourceNamePostgres = "postgres"
	sourceNamePg       = "pg"

	sourceNameRedis = "redis"
	sourceNameRds   = "rds"
)

type databases struct {
	configs []partition
}

func (d *databases) getParts(nl int) []model.Part {
	parts := make([]model.Part, 1, len(d.configs)+1)

	parts[0] = model.Part{NestingLevel: nl, Key: "postgres"}
	for _, db := range d.configs {
		parts = append(parts, db.GetParts(nl+1)...)
	}

	return parts
}

func newDbPattern(args []string) partition {
	db := &databases{}

	if len(args) == 0 {
		db.configs = []partition{sources.DefaultPgPattern()}
	}

	for _, arg := range args {
		idx := strings.Index(arg, "_")
		var source, name string

		source = arg[:idx]
		if idx != -1 {
			name = arg[idx+1:]
		} else {
			name = source
		}

		switch source {
		case sourceNamePg, sourceNamePostgres:
			pg := sources.DefaultPgPattern()
			pg.ConnectionName = name
			db.configs = append(db.configs)
		case sourceNameRds, sourceNameRedis:
			rds := sources.DefaultRdsPattern()
			rds.ConnectionName = name
			db.configs = append(db.configs)
		}
	}

	return db
}
