package bootstrap

import (
	"context"

	"github.com/sirupsen/logrus"

	"proj_name/internal/config"
	"proj_name/internal/transport"
)

func ApiEntryPoint(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
	mngr := transport.NewManager()

	go func() {
		err := mngr.Start(ctx)
		if err != nil {
			logrus.Fatalf("error starting server %s", err.Error())
		}
	}()

	return mngr.Stop, nil
}
