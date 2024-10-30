package db

import (
	"database/sql"
	"embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

//go:embed migrations/*
var MigrationsFS embed.FS

func NewSqlite(filePath string) (*sql.DB, error) {
	const op = "db.sqlite.new_slite"

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, errs.NewMsgError(op, err, "failed to connect to sqlite database")
	}

	return db, nil
}

// MigrateSchema performs migrations on the database
//
//	filePath_ is the path to the database file
//	migrationsFs is the filesystem containing the migrations
//	migrationTable is the name of the table to store migration information, by default it is "migrations"
//	nSteps is the number of steps to migrate, if nil, all migrations will be applied
func MigrateSchema(
	filePath string,
	migrationsFs embed.FS,
	migrationTable string,
	nSteps *int,
) error {
	const op = "adapters.db.sqlite"

	if migrationTable == "" {
		migrationTable = "migrations"
	}
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return errs.NewMsgError(op, err, "failed to open sqlite connection")
	}

	migrateDriver, err := sqlite.WithInstance(db, &sqlite.Config{
		MigrationsTable: migrationTable,
	})
	if err != nil {
		return errs.NewMsgError(op, err, "failed to create migration driver")
	}

	srcDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return errs.NewMsgError(op, err, "failed to create migration source")
	}

	preparedMigrations, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"",
		migrateDriver,
	)
	if err != nil {
		return errs.NewMsgError(op, err, "failed to create prepared migrations")
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
		return errs.NewMsgError(op, err, "failed to apply migrations")
	}

	return nil
}

// GetFile function returns filePath of created in temp directory file
func GetFile(fileName string) (string, error) {
	fileDir := filepath.Join(os.TempDir(), "go-revise")
	filePath := filepath.Join(fileDir, fileName)
	os.MkdirAll(fileDir, os.FileMode(0755))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
		file.Close()
	}

	return filePath, nil
}
