package tester

import (
	"database/sql"
	"embed"
	"os"
	"testing"

	adapterdb "github.com/ARUMANDESU/go-revise/internal/adapters/db"
	"github.com/ARUMANDESU/go-revise/pkg/teardowns"
)

//go:embed migrations/*
var MockDataMigrationsFS embed.FS

func NewSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()

	tds := teardowns.New()
	defer t.Cleanup(func() {
		err := tds.Execute()
		if err != nil {
			t.Error(err)
		}
	})
	handleErr := func(err error) {
		// we don't call teardowns execute because it's already in t.Cleanup
		t.Fatal(err)
	}

	dbFileName := t.Name()
	dbFilePath, err := adapterdb.GetFile(dbFileName)
	if err != nil {
		handleErr(err)
		return nil
	}
	tds.Append(func() error { return os.Remove(dbFilePath) })

	db, err := adapterdb.NewSqlite(dbFilePath)
	if err != nil {
		handleErr(err)
		return nil
	}
	tds.Append(db.Close)

	err = adapterdb.MigrateSchema(db, adapterdb.MigrationsFS, "", nil)
	if err != nil {
		handleErr(err)
		return nil
	}
	err = adapterdb.MigrateSchema(db, MockDataMigrationsFS, "test-mock", nil)
	if err != nil {
		handleErr(err)
		return nil
	}
	return db
}
