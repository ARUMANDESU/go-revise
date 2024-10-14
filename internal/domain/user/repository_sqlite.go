package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/db/sqlc"
	"github.com/ARUMANDESU/go-revise/pkg/pointers"
)

type SQLiteRepo struct {
	db *sql.DB
}

func userToModel(u User) sqlc.User {
	return sqlc.User{
		ID:           u.ID().String(),
		ChatID:       int64(u.ChatID()),
		CreatedAt:    u.createdAt,
		UpdatedAt:    u.updatedAt,
		Language:     sql.NullString{String: u.Settings().Language.String(), Valid: true},
		ReminderTime: reminderTimeToModel(u.Settings().ReminderTime),
	}
}
func modelToUser(u sqlc.User) (*User, error) {
	// @@TODO: change language from default one to the one from the database, i did not figure out how to do it
	settings, err := NewSettings(pointers.New(DefaultLanguage()), modelToReminderTime(u.ReminderTime))
	if err != nil {
		return nil, err
	}
	return &User{
		id:        uuid.FromStringOrNil(u.ID),
		chatID:    TelegramID(u.ChatID),
		createdAt: u.CreatedAt,
		updatedAt: u.UpdatedAt,
		settings:  settings,
	}, nil
}
func reminderTimeToModel(rt ReminderTime) string {
	return fmt.Sprintf("%d:%d", rt.Hour, rt.Minute)
}

func modelToReminderTime(rt string) ReminderTime {
	var hour, minute uint8
	fmt.Sscanf(rt, "%d:%d", &hour, &minute)
	return ReminderTime{Hour: hour, Minute: minute}
}

// CreateUser creates a new user.
func (r *SQLiteRepo) CreateUser(ctx context.Context, u User) (_ error) {
	params := sqlc.CreateUserParams{
		ID:        u.ID().String(),
		ChatID:    int64(u.ChatID()),
		CreatedAt: u.createdAt,
		UpdatedAt: u.updatedAt,
		Language:  sql.NullString{String: u.Settings().Language.String(), Valid: true},
	}

	err := sqlc.New(r.db).CreateUser(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates user data.
// The updateFn function is called with the user data to be updated.
func (r *SQLiteRepo) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	updateFn func(*User) (*User, error),
) (_ error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	qtx := sqlc.New(tx)
	userModel, err := qtx.GetUserByID(ctx, userID.String())
	if err != nil {
		return err
	}

	user, err := modelToUser(userModel)
	if err != nil {
		return err
	}

	user, err = updateFn(user)
	if err != nil {
		return err
	}

	userModel = userToModel(*user)
	err = qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		UpdatedAt:    userModel.UpdatedAt,
		Language:     userModel.Language,
		ReminderTime: userModel.ReminderTime,
		ID:           userModel.ID,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}
