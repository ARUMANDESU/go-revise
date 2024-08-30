package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ARUMANDESU/go-revise/internal/app"
	"github.com/ARUMANDESU/go-revise/internal/config"
	"github.com/ARUMANDESU/go-revise/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}

	cfg := config.MustLoad()

	log, close := logger.Setup(cfg.Env)
	defer close()

	log.Info("starting the app", slog.Attr{Key: "env", Value: slog.StringValue(cfg.Env)})

	application := app.NewApp(*cfg, log)

	go application.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	defer log.Info("application stopped", slog.String("signal", sign.String()))
	log.Info("stopping application", slog.String("signal", sign.String()))
}
