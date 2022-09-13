package config

import (
	"github.com/Red-Sock/rscli/internal/service/config/model"
)

const (
	dbFlag = "db"
)

type partition interface {
	GetParts(nestingLevel int) []model.Part
}

var patterns = map[string]func(argument []string) partition{
	dbFlag: newDbPattern,
}
