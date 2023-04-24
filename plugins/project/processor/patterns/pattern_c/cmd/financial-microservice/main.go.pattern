package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"financial-microservice/internal/config"
	"financial-microservice/internal/utils/closer"
	//_transport_imports
)

func main() {
	log.Println("starting app")

	ctx := context.Background()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config %s", err.Error())
	}

	startupDuration, err := cfg.GetDuration(config.AppInfoStartupDuration)
	if err != nil {
		log.Fatalf("error extracting startup duration %s", err)
	}
	context.WithTimeout(ctx, startupDuration)

	waitingForTheEnd()

	log.Println("shutting down the app")

	if err = closer.Close(); err != nil {
		log.Fatalf("errors while shutting down application %s", err.Error())
	}
}

// rscli comment: an obligatory function for tool to work properly.
// must be called in the main function above
// also this is a LP song name reference, so no rules can be applied to the function name
func waitingForTheEnd() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
