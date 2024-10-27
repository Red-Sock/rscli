package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Red-Sock/toolbox/closer"
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/init_project_ok/internal/config"
	//_transport_imports
)

func main() {
	logrus.Println("starting app")

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("error reading config %s", err.Error())
	}

	if cfg.AppInfo.StartupDuration == 0 {
		logrus.Fatalf("no startup duration in config")
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.AppInfo.StartupDuration)
	closer.Add(
		func() error {
			cancel()
			return nil
		})

	waitingForTheEnd()

	logrus.Println("shutting down the app")

	if err = closer.Close(); err != nil {
		logrus.Fatalf("errors while shutting down application %s", err.Error())
	}
}

// rscli comment: an obligatory function for tool to work properly.
// must be called in the main function above
// also this is the LP's song name reference, so no linting rules can be applied to the function name
func waitingForTheEnd() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
