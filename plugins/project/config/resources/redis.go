package resources

import (
	"strconv"

	errors "github.com/Red-Sock/trace-errors"
)

var (
	EnvVarRedisHost     = "REDIS_HOST"
	EnvVarRedisPort     = "REDIS_PORT"
	EnvVarRedisUser     = "REDIS_USER"
	EnvVarRedisPassword = "REDIS_PWD"
	EnvVarRedisDbNum    = "REDIS_DB"
)

type Redis struct {
	ResourceName string `yaml:"-"`

	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`

	User string `yaml:"user"`
	Pwd  string `yaml:"pwd"`
	Db   int    `yaml:"db"`
}

func (p *Redis) GetName() string {
	return p.ResourceName
}

func (p *Redis) GetEnv() map[string]string {
	return map[string]string{
		EnvVarRedisHost:     p.Host,
		EnvVarRedisPort:     strconv.FormatUint(uint64(p.Port), 10),
		EnvVarRedisUser:     p.User,
		EnvVarRedisPassword: p.Pwd,
		EnvVarRedisDbNum:    strconv.FormatUint(uint64(p.Db), 10),
	}
}

func (p *Redis) GetType() DataSourceName {
	return DataSourceRedis
}

func (p *Redis) FillFromEnv(env map[string]string) error {
	p.Host = env[EnvVarRedisHost]

	p.Pwd = env[EnvVarRedisPassword]
	p.User = env[EnvVarRedisUser]

	dbNumStr, ok := env[EnvVarRedisDbNum]
	if ok {
		dbNum, err := strconv.ParseInt(dbNumStr, 10, 64)
		if err != nil {
			return errors.Wrap(err, "error parsing redis db number")
		}
		p.Db = int(dbNum)
	}

	dbPortStr, ok := env[EnvVarRedisPort]
	if ok {
		redisPort, err := strconv.ParseUint(dbPortStr, 10, 64)
		if err != nil {
			return errors.Wrap(err, "error parsing redis port")
		}

		p.Port = uint16(redisPort)
	}

	return nil
}
