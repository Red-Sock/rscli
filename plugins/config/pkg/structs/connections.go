package structs

type Postgres struct {
	Host string
	Port uint16

	Name    string
	User    string
	Pwd     string
	SSLMode string
}

type Redis struct {
	Host string
	Port uint16

	User string
	Pwd  string
	Db   int
}

type RestApi struct {
	Port uint16
}

type Telegram struct {
	ApiKey string
}
