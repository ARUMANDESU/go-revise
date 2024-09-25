package logutil

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Env represents the environment in which the application is running
type Env string

const (
	envLocal Env = "local"
	envTest  Env = "test"
	envDev   Env = "dev"
	envProd  Env = "prod"
)

func (e Env) Validate() bool {
	switch e {
	case envLocal, envTest, envDev, envProd:
		return true
	default:
		return false
	}
}

// Setup creates a new logger based on the environment
//
//	Note: the log file is stored in the user's cache directory,
//	don't forget to call the cleanup function to close the log file
//	Warning: panics if the log file cannot be created
func Setup(env Env) (*slog.Logger, func()) {
	ok := env.Validate()
	if !ok {
		panic(fmt.Errorf(
			"invalid environment: %s is not a valid environment, must be one of: %s, %s, %s, %s",
			env,
			envLocal,
			envTest,
			envDev,
			envProd,
		))
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(fmt.Errorf("failed to get user cache directory: %w", err))
	}

	// Create the data directory if it does not exist
	dataDir := filepath.Join(cacheDir, "go-revise")
	err = os.MkdirAll(dataDir, os.FileMode(0755))
	if err != nil {
		panic(fmt.Errorf("failed to create data directory: %w", err))
	}

	// Open or create the log file
	logFile, err := os.OpenFile(filepath.Join(dataDir, "log.txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("failed to open log file: %w", err))
	}

	// Create a multi-writer to write to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var log *slog.Logger
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
