package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	envLocal = "local"
	envTest  = "test"
	envDev   = "dev"
	envProd  = "prod"
)

func Setup(env string) (*slog.Logger, func()) {
	var log *slog.Logger

	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "go-revise")

	err := os.MkdirAll(dataDir, os.FileMode(0755))
	if err != nil {
		panic(err)
	}

	// Open or create the log file
	logFile, err := os.OpenFile(filepath.Join(dataDir, "log.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// Create a multi-writer to write to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	switch env {
	case envLocal, envTest:
		log = slog.New(
			slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log, func() {
		logFile.Close()
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Plug() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
