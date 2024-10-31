package tester

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	adapterdb "github.com/ARUMANDESU/go-revise/internal/adapters/db"
	"github.com/ARUMANDESU/go-revise/pkg/teardowns"
)

//go:embed migrations/*
var MockDataMigrationsFS embed.FS

func NewSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()

	tds := teardowns.New()
	handleErr := func(err error) {
		terr := tds.Execute()
		if terr != nil {
			errors.Join(err, terr)
		}
		t.Fatal(err)
	}

	dbFileName := strings.ReplaceAll(t.Name(), "/", "-") + ".db"
	dbFilePath, err := adapterdb.GetFile(dbFileName)
	if err != nil {
		handleErr(fmt.Errorf("failed to get/create db file path/file: %w", err))
		return nil
	}
	tds.Append(func() error { return os.Remove(dbFilePath) })

	db, err := adapterdb.NewSqlite(dbFilePath)
	if err != nil {
		handleErr(fmt.Errorf("failed to create new new sqlite db: %w", err))
		return nil
	}
	tds.Append(db.Close)

	err = db.Ping()
	if err != nil {
		handleErr(fmt.Errorf("failed to ping db: %w", err))
		return nil
	}

	err = adapterdb.MigrateSchema(dbFilePath, adapterdb.MigrationsFS, "", nil)
	if err != nil {
		handleErr(fmt.Errorf("failed to migrate init: %w", err))
		return nil
	}
	err = adapterdb.MigrateSchema(dbFilePath, MockDataMigrationsFS, "test_mock", nil)
	if err != nil {
		handleErr(fmt.Errorf("failed to migrate mock data: %w", err))
		return nil
	}
	t.Cleanup(func() {
		err := tds.Execute()
		if err != nil {
			t.Error(err)
		}
	})
	return db
}
