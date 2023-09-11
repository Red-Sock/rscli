package resources

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

// TODO
func (p *Redis) GetEnv() map[string]string {
	return map[string]string{}
}

func (p *Redis) GetType() DataSourceName {
	return DataSourceRedis
}
