package logutil

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ARUMANDESU/go-revise/pkg/env"
)

// Setup creates a new logger based on the environment
//
// WARNING: panics if any errors occur during setup
//
// REMINDER: don't forget to call the cleanup function to close the log file
func Setup(mode env.Mode) (*slog.Logger, func()) {
	ok := mode.Validate()
	if !ok {
		panic(fmt.Errorf(
			"invalid environment mode: %s is not a valid environment, must be one of: %s, %s, %s, %s",
			mode,
			env.Local,
			env.Test,
			env.Dev,
			env.Prod,
		))
	}

	logFile := OpenFile("log.txt")

	// Create a multi-writer to write to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	var log *slog.Logger
	switch mode {
	case env.Local, env.Test:
		log = slog.New(
			slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.Dev:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env.Prod:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	slog.SetDefault(log)

	return log, func() {
		err := logFile.Close()
		if err != nil {
			log.Error("failed to close log file", Err(err))
		}
	}
}

// OpenFile opens or creates a log file in the user's cache directory and returns the file
//
// WARNING: panics if the log file cannot be created
func OpenFile(fileName string) *os.File {
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
	file, err := os.OpenFile(filepath.Join(dataDir, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("failed to open log file: %w", err))
	}

	return file
}

// Err returns an error attribute for logging
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Plug() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
