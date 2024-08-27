package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/storage"
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
			return domain.User{}, domain.WrapErrorWithOp(storage.ErrNotFound, op, "failed to get user")
		}
		return domain.User{}, domain.WrapErrorWithOp(err, op, "failed to get user")
	}

	return user, nil
}

func (s *Storage) GetUserByTelegramID(ctx context.Context, telegramID int64) (domain.User, error) {
	const op = "storage.sqlite.user.getUserByTelegramID"

	query := `
		SELECT id, telegram_id
		FROM users
		WHERE telegram_id = ?
		`

	var user domain.User
	err := s.DB.QueryRowContext(ctx, query, telegramID).Scan(user.ID, user.TelegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.WrapErrorWithOp(storage.ErrNotFound, op, "failed to get user")
		}
		return domain.User{}, domain.WrapErrorWithOp(err, op, "failed to get user")
	}

	return user, nil
}
