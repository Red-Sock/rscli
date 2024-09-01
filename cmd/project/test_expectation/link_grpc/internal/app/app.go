package app

import (
	"github.com/Red-Sock/toolbox"
	"github.com/Red-Sock/toolbox/closer"
	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/Red-Sock/link_grpc/internal/config"
)

type App struct {
	Ctx  context.Context
	Stop func()
	Cfg  config.Config
	/* Data source connection */
}

func New() (app App, err error) {
	logrus.Println("starting app")

	err = app.InitConfig()
	if err != nil {
		return App{}, errors.Wrap(err, "error initializing config")
	}

	err = app.InitDataSources()
	if err != nil {
		return App{}, errors.Wrap(err, "error during data sources initialization")
	}

	return app, nil
}

func (a *App) Start() (err error) {
	toolbox.WaitForInterrupt()

	logrus.Println("shutting down the app")

	err = closer.Close()
	if err != nil {
		return errors.Wrap(err, "error while shutting down application")
	}

	return nil
}
