package bootstrap

import (
	"context"

	"github.com/sirupsen/logrus"

	"financial-microservice/internal/config"
	"financial-microservice/internal/transport"
)

func ApiEntryPoint(ctx context.Context, cfg *config.Config) func(context.Context) error {
	mngr := transport.NewManager()

	go func() {
		err := mngr.Start(ctx)
		if err != nil {
			logrus.Fatalf("error starting server %s", err.Error())
		}
	}()
	return mngr.Stop
}
