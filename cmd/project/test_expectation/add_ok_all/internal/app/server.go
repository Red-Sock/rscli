package app

import (
	"github.com/Red-Sock/add_ok_all/internal/transport"
	errors "github.com/Red-Sock/trace-errors"
)

func (a *App) InitServers() (err error) {
	a.Server, err = transport.NewServerManager(a.Ctx, ":8080")
	if err != nil {
		return errors.Wrap(err, "error during server initialization on port: 8080")
	}

	return nil
}
