package bootstrap

import (
	"context"
	"log"

	"financial-microservice/internal/config"
	"financial-microservice/internal/transport"
)

func ApiEntryPoint(ctx context.Context, cfg *config.Config) func(context.Context) error {
	mngr := transport.NewManager()

	go func() {
		err := mngr.Start(ctx)
		if err != nil {
			log.Fatalf("error starting server %s", err.Error())
		}
	}()
	return mngr.Stop
}
