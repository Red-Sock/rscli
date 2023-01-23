package scripts

type redisEnvs struct{}

func (r *redisEnvs) GetEnvs() string {
	return `
PROJECT_NAME_CAPS_REDIS_PORT=%d
`
}
func (r *redisEnvs) GetStartingPort() int {
	return 6379
}

type pgEnvs struct{}

func (r *pgEnvs) GetEnvs() string {
	return `
PROJECT_NAME_CAPS_POSTGRES_PORT=%d
`
}
func (r *pgEnvs) GetStartingPort() int {
	return 5432
}
