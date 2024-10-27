package app

import (
	"github.com/Red-Sock/add_ok_all/internal/clients/redis"
	"github.com/Red-Sock/add_ok_all/internal/clients/sqldb"
	"github.com/Red-Sock/add_ok_all/internal/clients/telegram"
	errors "github.com/Red-Sock/trace-errors"
)

func (a *App) InitDataSources() (err error) {
	a.Postgres, err = sqldb.New(a.Cfg.DataSources.Postgres)
	if err != nil {
		return errors.Wrap(err, "error during sql connection initialization")
	}

	a.Redis, err = redis.New(a.Cfg.DataSources.Redis)
	if err != nil {
		return errors.Wrap(err, "error during redis connection initialization")
	}

	a.Telegram, err = telegram.New(a.Cfg.DataSources.Telegram)
	if err != nil {
		return errors.Wrap(err, "error during telegram bot initialization")
	}

	a.Sqlite, err = sqldb.New(a.Cfg.DataSources.Sqlite)
	if err != nil {
		return errors.Wrap(err, "error during sql connection initialization")
	}

	return nil
}
