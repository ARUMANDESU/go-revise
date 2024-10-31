package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mattn/go-sqlite3"

	"github.com/ARUMANDESU/go-revise/internal/adapters/db/sqlc"
	"github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
	"github.com/ARUMANDESU/go-revise/pkg/pointers"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(db *sql.DB) SQLiteRepo {
	return SQLiteRepo{db: db}
}

// CreateUser creates a new user.
func (r *SQLiteRepo) CreateUser(ctx context.Context, u user.User) (_ error) {
	const op = "domain.user.sqlite.create_user"

	reminderTime := fmt.Sprintf(
		"%d:%d",
		u.Settings().ReminderTime.Hour,
		u.Settings().ReminderTime.Minute,
	)
	params := sqlc.CreateUserParams{
		ID:           u.ID().String(),
		ChatID:       int64(u.ChatID()),
		CreatedAt:    u.CreatedAt(),
		UpdatedAt:    u.UpdatedAt(),
		Language:     sql.NullString{String: u.Settings().Language.String(), Valid: true},
		ReminderTime: reminderTime,
	}

	err := sqlc.New(r.db).CreateUser(ctx, params)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch {
			case
				errors.Is(sqliteErr.Code, sqlite3.ErrConstraint),
				errors.Is(sqliteErr.Code, sqlite3.ErrConstraintUnique):
				return errs.NewConflictError(op, err, "user already exists")
			case errors.Is(sqliteErr.Code, sqlite3.ErrConstraintNotNull):
				return errs.NewIncorrectInputError(op, err, "missing required fields")
			}
		}
		return errs.NewMsgError(op, err, "failed to create new user")
	}

	return nil
}

// UpdateUser updates user data.
// The updateFn function is called with the user data to be updated.
func (r *SQLiteRepo) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	updateFn func(*user.User) (*user.User, error),
) (_ error) {
	const op = "domain.user.sqlite.update_user"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.NewMsgError(op, err, "failed to begin transaction")
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
		return errs.NewMsgError(op, err, "failed to get user by id")
	}

	domainUser, err := modelToUser(userModel)
	if err != nil {
		return errs.NewMsgError(op, err, "failed to convert user model to domain")
	}

	domainUser, err = updateFn(domainUser)
	if err != nil {
		return err
	}

	userModel = userToModel(domainUser)
	err = qtx.UpdateUser(ctx, sqlc.UpdateUserParams{
		UpdatedAt:    userModel.UpdatedAt,
		Language:     userModel.Language,
		ReminderTime: userModel.ReminderTime,
		ID:           userModel.ID,
	})
	if err != nil {
		return errs.NewMsgError(op, err, "failed to update user")
	}

	return tx.Commit()
}

func (r *SQLiteRepo) GetUsersForNotification(ctx context.Context) ([]user.User, error) {
	const op = "domain.user.sqlite.get_users_for_notification"

	q := sqlc.New(r.db)

	reminderTimeModel := reminderTimeToModel(user.ReminderTime{
		Hour:   uint8(time.Now().Hour()),
		Minute: uint8(time.Now().Minute()),
	})
	userModels, err := q.GetUsersByReminderTime(ctx, reminderTimeModel)
	if err != nil {
		return nil, errs.NewMsgError(op, err, "failed to get users by reminder time")
	}

	return modelsToUsers(userModels)
}

func (r *SQLiteRepo) GetUserByID(ctx context.Context, id uuid.UUID) (query.User, error) {
	const op = "domain.user.sqlite.get_user_by_id"
	q := sqlc.New(r.db)

	userModel, err := q.GetUserByID(ctx, id.String())
	if err != nil {
		return query.User{}, errs.NewMsgError(op, err, "failed to get user by id")
	}

	return modelToQueryUser(userModel)
}

func (r *SQLiteRepo) GetUserByChatID(
	ctx context.Context,
	chatID user.TelegramID,
) (query.User, error) {
	const op = "domain.user.sqlite.get_user_by_chat_id"
	q := sqlc.New(r.db)

	userModel, err := q.GetUserByChatID(ctx, int64(chatID))
	if err != nil {
		return query.User{}, errs.NewMsgError(op, err, "failed to get user by chat id")
	}

	queryUser, err := modelToQueryUser(userModel)
	if err != nil {
		return query.User{}, errs.NewMsgError(op, err, "failed to convert user model into query")
	}
	return queryUser, nil
}

func (r *SQLiteRepo) GetUserByTelegramID(
	ctx context.Context,
	id user.TelegramID,
) (*user.User, error) {
	const op = "domain.user.sqlite.get_user_by_telegram_id"
	q := sqlc.New(r.db)

	userModel, err := q.GetUserByChatID(ctx, int64(id))
	if err != nil {
		return nil, errs.NewMsgError(op, err, "failed to get user by telegram id")
	}

	return modelToUser(userModel)
}

func userToModel(u *user.User) sqlc.User {
	return sqlc.User{
		ID:           u.ID().String(),
		ChatID:       int64(u.ChatID()),
		CreatedAt:    u.CreatedAt(),
		UpdatedAt:    u.UpdatedAt(),
		Language:     sql.NullString{String: u.Settings().Language.String(), Valid: true},
		ReminderTime: reminderTimeToModel(u.Settings().ReminderTime),
	}
}

func modelToUser(u sqlc.User) (*user.User, error) {
	reminderTime, err := modelToReminderTime(u.ReminderTime)
	if err != nil {
		return nil, err
	}

	// @@TODO: change language from default one to the one from the database, i did not figure out how to do it
	settings, err := user.NewSettings(pointers.New(user.DefaultLanguage()), reminderTime)
	if err != nil {
		return nil, err
	}

	return user.NewUser(
		uuid.FromStringOrNil(u.ID),
		user.TelegramID(u.ChatID),
		user.WithCreatedAt(u.CreatedAt),
		user.WithUpdatedAt(u.UpdatedAt),
		user.WithSettings(settings),
	)
}

func modelsToUsers(models []sqlc.User) ([]user.User, error) {
	users := make([]user.User, 0, len(models))
	for _, model := range models {
		domainUser, err := modelToUser(model)
		if err != nil {
			return nil, err
		}
		users = append(users, *domainUser)
	}

	return users, nil
}

func modelToQueryUser(u sqlc.User) (query.User, error) {
	const op = "domain.user.sqlite.model_to_query_user"

	reminderTime, err := modelToReminderTime(u.ReminderTime)
	if err != nil {
		return query.User{}, errs.NewMsgError(
			op,
			err,
			"failed to convert reminder time model into domain",
		)
	}

	var language string
	if u.Language.Valid {
		language = u.Language.String
	} else {
		language = user.DefaultLanguage().String()
	}

	return query.User{
		ID:     u.ID,
		ChatID: u.ChatID,
		Settings: query.Settings{
			Language: language,
			ReminderTime: query.ReminderTime{
				Hour:   reminderTime.Hour,
				Minute: reminderTime.Minute,
			},
		},
	}, nil
}

func reminderTimeToModel(rt user.ReminderTime) string {
	return fmt.Sprintf("%d:%d", rt.Hour, rt.Minute)
}

func modelToReminderTime(rt string) (user.ReminderTime, error) {
	const op = "domain.user.sqlite.model_to_reminder_time"
	rt = strings.TrimSpace(rt)
	if rt == "" {
		return user.ReminderTime{}, errs.NewMsgError(
			op,
			fmt.Errorf("reminder time is empty, have to be in this format: HOUR:MINUTE"),
			"reminder time is empty",
		)
	}
	var hour, minute uint8
	parsed, err := fmt.Sscanf(rt, "%d:%d", &hour, &minute)
	if err != nil {
		return user.ReminderTime{}, errs.NewMsgError(
			op,
			err,
			fmt.Sprintf("failed to parse model reminder time, model: %s", rt),
		)
	}
	if parsed != 2 {
		return user.ReminderTime{}, errs.NewMsgError(
			op,
			err,
			"number of parsed items are not equal to 2",
		)
	}

	return user.ReminderTime{Hour: hour, Minute: minute}, nil
}
