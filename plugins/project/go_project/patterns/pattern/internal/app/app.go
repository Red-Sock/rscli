// Code generated by RedSock CLI. DO NOT EDIT.

package app

import (
	"github.com/Red-Sock/toolbox"
	"github.com/Red-Sock/toolbox/closer"
	errors "github.com/Red-Sock/trace-errors"
	"proj_name/internal/config"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type App struct {
	Ctx  context.Context
	Stop func()
	Cfg  config.Config

	Custom Custom
}

func New() (app App, err error) {
	logrus.Println("starting app")

	err = app.InitConfig()
	if err != nil {
		return App{}, errors.Wrap(err, "error initializing config")
	}

	err = app.Custom.Init(&app)
	if err != nil {
		return App{}, errors.Wrap(err, "error initializing custom app properties")
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
