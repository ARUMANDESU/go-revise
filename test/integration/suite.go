package integration

import (
	"database/sql"
	"embed"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/ARUMANDESU/go-revise/internal/storage/sqlite"
	"github.com/thejerf/slogassert"
)

//go:embed migrations/*
var migrationsFs embed.FS

type Suite struct {
	*testing.T
	LogHandler *slogassert.Handler
	Service    service.Revise
}

func NewSuite(t *testing.T) (*Suite, func()) {
	t.Helper()

	handler := slogassert.New(t, slog.LevelWarn, nil)
	log := slog.New(handler)

	storage, cleanup := setupSqlite(t)

	return &Suite{
			T:          t,
			LogHandler: handler,
			Service: service.NewRevise(log,
				service.ReviseStorages{
					ReviseProvider: storage,
					ReviseManager:  storage,
					UserProvider:   storage,
				},
			)},
		func() { cleanup() }
}

func setupSqlite(t *testing.T) (*sqlite.Storage, func()) {
	t.Helper()
	var storage *sqlite.Storage
	dbFilePatn := getDataSource(t, "test.db")

	cleanUp := func() {
		// cleanup
		storage.Close()

		// remove test database
		if err := os.Remove(dbFilePatn); err != nil {
			t.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite", dbFilePatn)
	if err != nil {
		cleanUp()
		t.Fatal(err)
	}
	// create test database base schema
	if err := sqlite.MigrateSchema(dbFilePatn, sqlite.MigrationsFs, "", nil); err != nil {
		cleanUp()
		t.Fatal(err)
	}

	// add mock data
	if err := sqlite.MigrateSchema(dbFilePatn, migrationsFs, "test", nil); err != nil {
		cleanUp()
		t.Fatal(err)
	}

	storage = &sqlite.Storage{DB: db}

	return storage, cleanUp
}

func getDataSource(t *testing.T, fileName string) string {
	t.Helper()

	dataDir := filepath.Join(os.TempDir(), "go-revise")
	filePath := filepath.Join(dataDir, fileName)
	os.MkdirAll(dataDir, os.FileMode(0755))

	// if file is not found, it will be created automatically
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	return filePath
}
