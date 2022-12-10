package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"financial-microservice/internal/config"
	"financial-microservice/internal/transport"
	//_transport_imports
)

func main() {
	log.Println("starting app")

	ctx := context.Background()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config %s", err.Error())
	}

	server := apiEntryPoint(ctx, cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	err = server.Stop(ctx)
	if err != nil {
		log.Printf("error stopping web server %v", err.Error())
	}

	log.Println("shutting down the app")
}

func apiEntryPoint(ctx context.Context, cfg *config.Config) transport.Server {
	mngr := transport.NewManager()

	//_initiation_of_servers

	go func() {
		err := mngr.Start(ctx)
		if err != nil {
			log.Fatalf("error starting server %s", err.Error())
		}
	}()
	return mngr
}
