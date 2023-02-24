package main

import (
	"context"
	"log"

	"financial-microservice/internal/config"
	"financial-microservice/internal/transport"
)

func apiEntryPoint(ctx context.Context, cfg *config.Config) transport.Server {
	mngr := transport.NewManager()

	go func() {
		err := mngr.Start(ctx)
		if err != nil {
			log.Fatalf("error starting server %s", err.Error())
		}
	}()
	return mngr
}
