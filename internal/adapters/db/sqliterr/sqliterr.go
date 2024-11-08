package sqliterr

import (
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func Handle(op errs.Op, err error, msg string) *errs.Error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return handleNotFound(op, err, msg)
	}

	var sqliteErr *sqlite.Error
	if !errors.As(err, &sqliteErr) {
		return handleUnknown(op, err, msg)
	}

	switch sqliteErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
		return handleUniqueConstraint(op, err, msg)
	case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
		return handleForeignKeyConstraint(op, err, msg)
	case sqlite3.SQLITE_CONSTRAINT:
		return handleConstraint(op, err, msg)
	case sqlite3.SQLITE_BUSY:
		return handleBusy(op, err, msg)
	case sqlite3.SQLITE_LOCKED:
		return handleLocked(op, err, msg)
	case sqlite3.SQLITE_READONLY:
		return handleReadOnly(op, err, msg)
	case sqlite3.SQLITE_NOTFOUND:
		return handleNotFound(op, err, msg)
	case sqlite3.SQLITE_FULL:
		return handleDatabaseFull(op, err, msg)
	case sqlite3.SQLITE_CORRUPT:
		return handleCorrupt(op, err, msg)
	case sqlite3.SQLITE_IOERR:
		return handleIO(op, err, msg)
	default:
		return handleUnknown(op, err, msg)
	}
}

// HandleTx wraps transaction-specific error handling
func HandleTx(op errs.Op, err error, msg string) *errs.Error {
	if err == nil {
		return nil
	}

	var sqliteErr *sqlite.Error
	if !errors.As(err, &sqliteErr) {
		return handleUnknown(op, err, msg)
	}

	switch sqliteErr.Code() {
	case sqlite3.SQLITE_BUSY:
		return errs.NewConflictError(op, err, withMsg(msg, "database is busy")).
			WithMessages([]errs.Message{{Key: "error", Value: "please retry the operation"}})
	case sqlite3.SQLITE_LOCKED:
		return errs.NewConflictError(op, err, withMsg(msg, "database is locked")).
			WithMessages([]errs.Message{{Key: "error", Value: "please retry the operation"}})
	default:
		return Handle(op, err, msg)
	}
}

// Individual error handlers
func handleUniqueConstraint(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewAlreadyExistsError(op, err, withMsg(msg, "record already exists")).
		WithMessages([]errs.Message{{Key: "error", Value: "record already exists"}})
}

func handleForeignKeyConstraint(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewIncorrectInputError(op, err, withMsg(msg, "referenced record does not exist")).
		WithMessages([]errs.Message{{Key: "error", Value: "referenced record does not exist"}})
}

func handleConstraint(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewIncorrectInputError(op, err, withMsg(msg, "data constraints were violated")).
		WithMessages([]errs.Message{{Key: "error", Value: "data constraints were violated"}})
}

func handleBusy(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewConflictError(op, err, withMsg(msg, "database is busy")).
		WithMessages([]errs.Message{{Key: "error", Value: "please retry the operation"}})
}

func handleLocked(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewConflictError(op, err, withMsg(msg, "database is locked")).
		WithMessages([]errs.Message{{Key: "error", Value: "please retry the operation"}})
}

func handleReadOnly(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewIncorrectInputError(op, err, withMsg(msg, "database is read-only")).
		WithMessages([]errs.Message{{Key: "error", Value: "cannot modify read-only database"}})
}

func handleNotFound(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewNotFound(op, err, withMsg(msg, "requested record was not found")).
		WithMessages([]errs.Message{{Key: "error", Value: "requested record was not found"}})
}

func handleDatabaseFull(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewUnknownError(op, err, withMsg(msg, "database is full"))
}

func handleCorrupt(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewUnknownError(op, err, withMsg(msg, "database is corrupt"))
}

func handleIO(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewUnknownError(op, err, withMsg(msg, "database I/O error"))
}

func handleUnknown(op errs.Op, err error, msg string) *errs.Error {
	return errs.NewUnknownError(op, err, fmt.Sprintf("%s: unexpected error: %v", msg, err))
}

func withMsg(msg, err string) string {
	if msg == "" {
		return err
	}
	return fmt.Sprintf("%s: %s", msg, err)
}
