package resources

import (
	"strconv"

	errors "github.com/Red-Sock/trace-errors"
)

var (
	EnvVarPostgresHost     = "POSTGRES_HOST"
	EnvVarPostgresPort     = "POSTGRES_PORT"
	EnvVarPostgresUser     = "POSTGRES_USER"
	EnvVarPostgresPassword = "POSTGRES_PWD"
	EnvVarPostgresDbName   = "POSTGRES_NAME"
)

type Postgres struct {
	ResourceName string `yaml:"-"`
	Host         string `yaml:"host"`
	Port         uint64 `yaml:"port"`

	Name string `yaml:"name"`
	User string `yaml:"user"`

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
		EnvVarPostgresDbName:   p.Name,

		EnvVarPostgresHost: p.Host,
		EnvVarPostgresPort: strconv.FormatUint(p.Port, 10),
	}
}

func (p *Postgres) FillFromEnv(in map[string]string) (err error) {
	p.User = in[EnvVarPostgresUser]
	p.Pwd = in[EnvVarPostgresPassword]
	p.Name = in[EnvVarPostgresDbName]

	p.Host = in[EnvVarPostgresHost]
	p.Port, err = strconv.ParseUint(in[EnvVarPostgresPort], 10, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing port value")
	}
	return nil
}
