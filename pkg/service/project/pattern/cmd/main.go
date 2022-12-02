package main

import (
	"context"
	"financial-microservice/internal/config"
	"financial-microservice/internal/transport/rest"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("starting app")

	ctx := context.Background()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config %s", err.Error())
	}
	server := rest.NewServer(*cfg)
	go func() {
		err = server.Start()
		if err != nil {
			log.Fatalf("error starting server %s", err.Error())
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	err = server.Stop(ctx)
	if err != nil {
		log.Printf("error stopping web server %v", err.Error())
	}

	log.Println("shutting down the app")
}
