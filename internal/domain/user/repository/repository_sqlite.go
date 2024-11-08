package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ARUMANDESU/go-revise/internal/adapters/db/sqlc"
	"github.com/ARUMANDESU/go-revise/internal/adapters/db/sqliterr"
	"github.com/ARUMANDESU/go-revise/internal/application/user/query"
	"github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
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
	op := errs.Op("domain.user.sqlite.create_user")

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
		return sqliterr.Handle(op, err, "failed to create user").WithContext("user", u)
	}

	return nil
}

func (r *SQLiteRepo) withTx(ctx context.Context, op errs.Op, fn func(*sqlc.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return sqliterr.HandleTx(op, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.
					With(slog.String("op", string(op))).
					Error("failed to rollback transaction", logutil.Err(rollbackErr), "original_error", err)
			}
		}
	}()

	qtx := sqlc.New(tx)
	if err = fn(qtx); err != nil {
		return err // Already wrapped with operation
	}

	if err = tx.Commit(); err != nil {
		return sqliterr.HandleTx(op, err, "failed to commit transaction")
	}

	return nil
}

func (r *SQLiteRepo) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	updateFn func(*user.User) (*user.User, error),
) error {
	op := errs.Op("domain.user.sqlite.update_user")

	return r.withTx(ctx, op, func(q *sqlc.Queries) error {
		userModel, err := q.GetUserByID(ctx, userID.String())
		if err != nil {
			return sqliterr.Handle(op, err, "failed to get user by id").WithContext("id", userID)
		}

		domainUser, err := modelToUser(userModel)
		if err != nil {
			return errs.WithOp(op, err, "failed to convert model to user")
		}

		domainUser, err = updateFn(domainUser)
		if err != nil {
			return errs.WithOp(op, err, "failed to update user")
		}

		userModel = userToModel(domainUser)
		err = q.UpdateUser(ctx, sqlc.UpdateUserParams{
			UpdatedAt:    userModel.UpdatedAt,
			Language:     userModel.Language,
			ReminderTime: userModel.ReminderTime,
			ID:           userModel.ID,
		})
		if err != nil {
			return sqliterr.Handle(op, err, "failed to update user")
		}

		return nil
	})
}

func (r *SQLiteRepo) GetUsersForNotification(ctx context.Context) ([]user.User, error) {
	op := errs.Op("domain.user.sqlite.get_users_for_notification")
	q := sqlc.New(r.db)

	reminderTimeModel := reminderTimeToModel(user.ReminderTime{
		Hour:   uint8(time.Now().Hour()),
		Minute: uint8(time.Now().Minute()),
	})
	userModels, err := q.GetUsersByReminderTime(ctx, reminderTimeModel)
	if err != nil {
		return nil, sqliterr.
			Handle(op, err, "failed to get users by reminder time").
			WithContext("reminder_time", reminderTimeModel)
	}

	users, err := modelsToUsers(userModels)
	if err != nil {
		return nil, errs.WithOp(op, err, "failed to convert models to users")
	}
	return users, nil
}

func (r *SQLiteRepo) GetUserByID(ctx context.Context, id uuid.UUID) (query.User, error) {
	op := errs.Op("domain.user.sqlite.get_user_by_id")
	q := sqlc.New(r.db)

	userModel, err := q.GetUserByID(ctx, id.String())
	if err != nil {
		return query.User{}, sqliterr.
			Handle(op, err, "failed to get user by id").
			WithContext("id", id)
	}

	queryUser, err := modelToQueryUser(userModel)
	if err != nil {
		return query.User{}, errs.WithOp(op, err, "failed to convert model to query user")
	}
	return queryUser, nil
}

func (r *SQLiteRepo) GetUserByChatID(
	ctx context.Context,
	chatID user.TelegramID,
) (query.User, error) {
	op := errs.Op("domain.user.sqlite.get_user_by_chat_id")
	q := sqlc.New(r.db)

	userModel, err := q.GetUserByChatID(ctx, int64(chatID))
	if err != nil {
		return query.User{}, sqliterr.
			Handle(op, err, "failed to get user by chat id").
			WithContext("chat_id", chatID)
	}

	queryUser, err := modelToQueryUser(userModel)
	if err != nil {
		return query.User{}, errs.WithOp(op, err, "failed to convert model to query user")
	}
	return queryUser, nil
}

func (r *SQLiteRepo) GetUserByTelegramID(
	ctx context.Context,
	id user.TelegramID,
) (*user.User, error) {
	op := errs.Op("domain.user.sqlite.get_user_by_telegram_id")
	q := sqlc.New(r.db)

	userModel, err := q.GetUserByChatID(ctx, int64(id))
	if err != nil {
		return nil, sqliterr.
			Handle(op, err, "failed to get user by chat id").
			WithContext("chat_id", id)
	}

	domainUser, err := modelToUser(userModel)
	if err != nil {
		return nil, errs.WithOp(op, err, "failed to convert model to user")
	}
	return domainUser, nil
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
	op := errs.Op("domain.user.sqlite.model_to_user")

	reminderTime, err := modelToReminderTime(u.ReminderTime)
	if err != nil {
		return nil, errs.WithOp(op, err, "failed to convert reminder time")
	}

	// @@TODO: change language from default one to the one from the database, i did not figure out how to do it
	settings, err := user.NewSettings(pointers.New(user.DefaultLanguage()), reminderTime)
	if err != nil {
		return nil, errs.WithOp(op, err, "failed to create settings")
	}

	domainUser, err := user.NewUser(
		uuid.FromStringOrNil(u.ID),
		user.TelegramID(u.ChatID),
		user.WithCreatedAt(u.CreatedAt),
		user.WithUpdatedAt(u.UpdatedAt),
		user.WithSettings(settings),
	)
	if err != nil {
		return nil, errs.WithOp(op, err, "failed to create user")
	}
	return domainUser, nil
}

func modelsToUsers(models []sqlc.User) ([]user.User, error) {
	op := errs.Op("domain.user.sqlite.models_to_users")

	users := make([]user.User, 0, len(models))
	for _, model := range models {
		domainUser, err := modelToUser(model)
		if err != nil {
			return nil, errs.WithOp(op, err, "failed to convert model to user")
		}
		users = append(users, *domainUser)
	}

	return users, nil
}

func modelToQueryUser(u sqlc.User) (query.User, error) {
	const op = "domain.user.sqlite.model_to_query_user"

	reminderTime, err := modelToReminderTime(u.ReminderTime)
	if err != nil {
		return query.User{}, errs.WithOp(op, err, "failed to convert reminder time")
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
	if isValidTimeFormat(rt) {
		return user.ReminderTime{}, errs.
			NewIncorrectInputError(op, nil, "invalid model reminder time").
			WithContext("value", rt)
	}
	var hour, minute uint8
	parsed, err := fmt.Sscanf(rt, "%d:%d", &hour, &minute)
	if err != nil {
		return user.ReminderTime{}, errs.
			NewUnknownError(op, err, "failed to parse reminder time").
			WithContext("reminder_time", rt)
	}
	if parsed != 2 {
		return user.ReminderTime{}, errs.
			NewUnknownError(op, err, "number of parsed items are not equal to 2").
			WithContext("reminder_time", rt).
			WithContext("parsed", parsed)
	}

	return user.ReminderTime{Hour: hour, Minute: minute}, nil
}

func isValidTimeFormat(s string) bool {
	_, err := time.Parse("15:04", s)
	return err == nil
}
