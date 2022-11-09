package config

type Postgres struct {
	Host string
	Port string

	Name string
	User string
	Pwd  string
}

type Redis struct {
	Host string
	Port string

	User string
	Pwd  string
	Db   int
}

type RestApi struct {
	Port uint16
}
