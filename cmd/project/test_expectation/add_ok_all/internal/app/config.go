package app

import (
	"context"

	"github.com/Red-Sock/toolbox/closer"
	errors "github.com/Red-Sock/trace-errors"

	"github.com/Red-Sock/add_ok_all/internal/config"
)

func (a *App) InitConfig() (err error) {
	a.Ctx, a.Stop = context.WithCancel(context.Background())
	closer.Add(func() error { a.Stop(); return nil })

	a.Cfg, err = config.Load()
	if err != nil {
		return errors.Wrap(err, "error reading config")
	}

	return nil
}
