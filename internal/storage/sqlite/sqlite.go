package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var MigrationsFs embed.FS

type Storage struct {
	DB *sql.DB
}

func NewStorage(fileName string) (*Storage, error) {
	const op = "storage.sqlite.new_storage"

	filePath, err := getDataSource(fileName)
	if err != nil {
		return nil, domain.WrapErrorWithOp(err, op, "failed to get data source")
	}

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, domain.WrapErrorWithOp(err, op, "failed to open sqlite connection")
	}

	if err := MigrateSchema(filePath, MigrationsFs, "", nil); err != nil {
		return nil, domain.WrapErrorWithOp(err, op, "failed to migrate schema")
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) Close() error {
	if s == nil || s.DB == nil {
		return nil
	}

	return s.DB.Close()
}

// MigrateSchema performs migrations on the database
//
//	filePath_ is the path to the database file
//	migrationsFs is the filesystem containing the migrations
//	migrationTable is the name of the table to store migration information, by default it is "migrations"
//	nSteps is the number of steps to migrate, if nil, all migrations will be applied
func MigrateSchema(filePath string, migrationsFs embed.FS, migrationTable string, nSteps *int) error {
	const op = "storage.sqlite.migrate_schema"

	if migrationTable == "" {
		migrationTable = "migrations"
	}
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to open sqlite connection")
	}

	migrateDriver, err := sqlite.WithInstance(db, &sqlite.Config{
		MigrationsTable: migrationTable,
	})
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to create migration driver")
	}

	srcDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to create migration source")
	}

	preparedMigrations, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"",
		migrateDriver,
	)
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to create prepared migrations")
	}

	defer func() {
		preparedMigrations.Close()
	}()
	if nSteps != nil {
		err = preparedMigrations.Steps(*nSteps)
	} else {
		err = preparedMigrations.Up()
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return domain.WrapErrorWithOp(err, op, "failed to apply migrations")
	}

	return nil
}

func getDataSource(fileName string) (string, error) {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "go-revise")
	os.MkdirAll(dataDir, os.FileMode(0755))

	filePath := filepath.Join(dataDir, fileName)

	// if file is not found, it will be created automatically
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return "", domain.WrapError(err, "failed to create database file")
		}
		file.Close()
	}

	return filePath, nil
}
