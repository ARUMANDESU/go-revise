package main

import (
	"log/slog"

	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log, close := logger.Setup(cfg.Env)
	defer close()

	log.Info("starting the app", slog.Attr{Key: "env", Value: slog.StringValue(cfg.Env)})

	// TODO: initialize the app

	// TODO: run the app

	log.Info("shutting down the app")

	// TODO: graceful shutdown the app

	log.Info("the app is down")
}
