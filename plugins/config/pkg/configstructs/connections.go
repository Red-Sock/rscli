package configstructs

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

type Telegram struct {
	ApiKey string
}
