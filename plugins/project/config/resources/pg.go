package resources

var (
	EnvVarPostgresUser     = "POSTGRES_USER"
	EnvVarPostgresPassword = "POSTGRES_PASSWORD"
	EnvVarPostgresDB       = "POSTGRES_DB"
)

type Postgres struct {
	ResourceName string `yaml:"-"`
	Host         string `yaml:"host"`
	Port         uint16 `yaml:"port"`

	Name    string `yaml:"name"`
	User    string `yaml:"user"`
	Pwd     string `yaml:"pwd"`
	SSLMode string `yaml:"ssl_mode"`
}

func (p *Postgres) GetName() string {
	return p.ResourceName
}

func (p *Postgres) GetType() DataSourceName {
	return DataSourcePostgres
}

func (p *Postgres) GetEnv() map[string]string {
	return map[string]string{
		EnvVarPostgresUser:     p.User,
		EnvVarPostgresPassword: p.Pwd,
		EnvVarPostgresDB:       p.Name,
	}
}
