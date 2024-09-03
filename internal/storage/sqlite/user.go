package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/storage"
	"github.com/mattn/go-sqlite3"
)

func (s *Storage) GetUser(ctx context.Context, id string) (domain.User, error) {
	const op = "storage.sqlite.user.getUser"

	query := `
		SELECT id, telegram_id 
		FROM users
		WHERE id = ?
		`

	var user domain.User
	err := s.DB.QueryRowContext(ctx, query, id).Scan(user.ID, user.TelegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.WrapErrorWithOp(
				storage.ErrNotFound,
				op,
				"failed to get user",
			)
		}
		return domain.User{}, domain.WrapErrorWithOp(err, op, "failed to get user")
	}

	return user, nil
}

func (s *Storage) GetUserByChatID(ctx context.Context, chatID int64) (domain.User, error) {
	const op = "storage.sqlite.user.getUserByTelegramID"

	query := `
		SELECT id, telegram_id
		FROM users
		WHERE telegram_id = ?
		`

	var user domain.User
	err := s.DB.QueryRowContext(ctx, query, chatID).Scan(user.ID, user.TelegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.WrapErrorWithOp(
				storage.ErrNotFound,
				op,
				"failed to get user",
			)
		}
		return domain.User{}, domain.WrapErrorWithOp(err, op, "failed to get user")
	}

	return user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user domain.User) error {
	const op = "storage.sqlite.user.createUser"

	query := `
		INSERT INTO users (id, telegram_id)
		VALUES (?, ?)
	`

	result, err := s.DB.Exec(query, user.ID, user.TelegramID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
				return domain.WrapErrorWithOp(storage.ErrAlreadyExists, op, "user already exists")
			}
		}
		return domain.WrapErrorWithOp(err, op, "failed to insert user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.WrapErrorWithOp(err, op, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return domain.WrapErrorWithOp(storage.ErrInsert, op, "failed to insert user (no rows affected)")
	}

	return nil
}
