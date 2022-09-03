package config

type Format []string

var (
	YmlFormat  Format = []string{"yml", "yaml"}
	JSONFormat Format = []string{"json"}
)
