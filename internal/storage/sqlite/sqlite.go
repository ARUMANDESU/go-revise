package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var migrationsFs embed.FS

type Storage struct {
	db *sql.DB
}

func NewStorage(fileName string) (*Storage, error) {
	db, err := sql.Open("sqlite", getDataSource(fileName))
	if err != nil {
		return nil, fmt.Errorf("open sqlite connection: %w", err)
	}

	if err := migrateSchema(db, nil); err != nil {
		return nil, fmt.Errorf("failed to perform migrations: %w", err)
	}

	return &Storage{db: db}, nil
}

func getDataSource(fileName string) string {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "go-revise")
	os.MkdirAll(dataDir, os.FileMode(0755))

	// if file is not found, it will be created automatically
	if _, err := os.Stat(filepath.Join(dataDir, fileName)); os.IsNotExist(err) {
		file, err := os.Create(filepath.Join(dataDir, fileName))
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	return filepath.Join(dataDir, fileName)
}

func migrateSchema(db *sql.DB, nSteps *int) error {
	migrateDriver, err := sqlite.WithInstance(db, &sqlite.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	srcDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source driver: %w", err)
	}

	preparedMigrations, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"",
		migrateDriver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration tooling instance: %w", err)
	}

	defer func() {
		preparedMigrations.Close()
		db.Close()
	}()
	if nSteps != nil {
		fmt.Printf("stepping migrations %d...\n", *nSteps)
		err = preparedMigrations.Steps(*nSteps)
	} else {
		err = preparedMigrations.Up()
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Successfully applied db migrations")
	return nil
}
