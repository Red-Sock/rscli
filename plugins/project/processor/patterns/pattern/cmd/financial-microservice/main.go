package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"financial-microservice/internal/config"
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
}

// rscli comment: an obligatory function for tool to work properly.
// must be called in the main function above
// also this is a LP song name reference, so no rules can be applied to the function name
func waitingForTheEnd() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
