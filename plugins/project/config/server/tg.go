package server

type Telegram struct {
	name string

	ApiKey string `yaml:"api_key"`
}

func (t *Telegram) GetName() string {
	return t.name
}

func (t *Telegram) GetPort() uint16 {
	return 0
}
