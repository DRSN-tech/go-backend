package main

import (
	"os"

	"github.com/DRSN-tech/go-backend/internal/app"
	config "github.com/DRSN-tech/go-backend/internal/cfg"
	"github.com/DRSN-tech/go-backend/pkg/logger"
)

func main() {
	log := logger.NewSlogLogger()

	cfg, err := config.Load(log)
	if err != nil {
		log.Errorf(err, "failed to load config")
		os.Exit(1)
	}

	application, err := app.NewApp(cfg, log)
	if err != nil {
		log.Errorf(err, "failed to initialize app")
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		os.Exit(1)
	}
}
