package tester

import (
	"database/sql"
	"embed"
	"os"

	adapterdb "github.com/ARUMANDESU/go-revise/internal/adapters/db"
	"github.com/ARUMANDESU/go-revise/pkg/teardowns"
)

//go:embed migrations/*
var MockDataMigrationsFS embed.FS

var sqlite *Sqlite

func GetSqlite() *Sqlite {
	return sqlite
}

func SetupSqlite() error {
	tds := teardowns.New()

	dbFilePath, err := adapterdb.GetFile("go-revise")
	if err != nil {
		return err
	}
	tds.Append(func() error { return os.Remove(dbFilePath) })

	db, err := adapterdb.NewSqlite(dbFilePath)
	if err != nil {
		return err
	}
	tds.Append(db.Close)

	adapterdb.MigrateSchema(db, adapterdb.MigrationsFS, "", nil)
	adapterdb.MigrateSchema(db, MockDataMigrationsFS, "test-mock", nil)

	sqlite = &Sqlite{
		db:        db,
		teardowns: tds,
	}

	return nil
}

type Sqlite struct {
	db        *sql.DB
	teardowns teardowns.Funcs
}

func newSqlite(db *sql.DB) *Sqlite {
	return &Sqlite{
		db: db,
	}
}

func (s *Sqlite) DB() *sql.DB {
	return s.db
}

func (s *Sqlite) NewDB(dbName string) *sql.DB {
	// @@TODO: create new database or scheme based on default one
	return nil
}
