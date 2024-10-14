package user

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/db/sqlc"
)

type SQLiteRepo struct {
	db *sql.DB
}

func UserToModel(u User) sqlc.User {
	return sqlc.User{
		ID:     u.ID().String(),
		ChatID: int64(u.ChatID()),
	}
}

// SaveUser saves a user.
func (r *SQLiteRepo) SaveUser(ctx context.Context, u User) (_ error) {
	params := sqlc.SaveUserParams{
		ID:     u.ID().String(),
		ChatID: int64(u.ChatID()),
	}

	err := sqlc.New(r.db).SaveUser(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates user data.
// The updateFn function is called with the user data to be updated.
func (SQLiteRepo) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	updateFn func(*User) (*User, error),
) (_ error) {
	panic("not implemented") // TODO: Implement
}
