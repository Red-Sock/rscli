package project

import (
	_ "embed"
)

//go:embed patterns/main.go.pattern
var mainFile []byte

//go:embed patterns/redis.go.pattern
var connectionToRedis []byte

//go:embed patterns/config.go.pattern
var configurator string
