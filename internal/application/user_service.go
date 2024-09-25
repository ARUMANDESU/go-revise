package application

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"

	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/i18n"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

// UserProvider handles the retrieval of user data.
type UserProvider interface {
	// GetUserByID returns a user by ID(UUID).
	GetUserByID(ctx context.Context, id uuid.UUID) (domainUser.User, error)
	// GetUserByTelegramID returns a user by Telegram ID(int64).
	GetUserByTelegramID(ctx context.Context, telegramID domainUser.TelegramID) (domainUser.User, error)
}

// UserRepository handles the persistence of user data.
type UserRepository interface {
	// Save saves a user.
	Save(ctx context.Context, u domainUser.User) error
	// UpdateSettings updates user settings.
	UpdateSettings(ctx context.Context, userID uuid.UUID, settings domainUser.Settings) error
}

type UserService struct {
	log            slog.Logger
	userRepository UserRepository
	userProvider   UserProvider
}

func (s UserService) GetUserByID(ctx context.Context, id domainUser.Identifier) (domainUser.User, error) {
	const op = "UserService.GetUserByID"
	log := s.log.With("op", op)

	if !id.IsValid() {
		return domainUser.User{}, domainUser.ErrInvalidIdentifier
	}

	var (
		user domainUser.User
		err  error
	)
	switch id := id.(type) {
	case domainUser.UUIDIdentifier:
		user, err = s.userProvider.GetUserByID(ctx, id.UUID())
	case domainUser.TelegramIDWrapper:
		user, err = s.userProvider.GetUserByTelegramID(ctx, id.TelegramID())
	default:
		return domainUser.User{}, domainUser.ErrInvalidIdentifier
	}
	if err != nil {
		log.Error("failed to get domainUser", logutil.Err(err))
		return domainUser.User{}, err
	}

	return user, nil
}

type NewUserServiceParams struct {
	ChatID       domainUser.TelegramID    `json:"chat_id"`
	Language     i18n.Language            `json:"language"`
	ReminderTime *domainUser.ReminderTime `json:"reminder_time"`
}

func (s UserService) SaveUser(ctx context.Context, u NewUserServiceParams) error {
	const op = "UserService.SaveUser"
	log := s.log.With("op", op)

	settings := domainUser.Settings{
		ID:           uuid.Must(uuid.NewV7()),
		Language:     u.Language,
		ReminderTime: domainUser.DefaultReminderTime(),
	}
	if u.ReminderTime != nil {
		settings.ReminderTime = *u.ReminderTime
	}

	user := domainUser.NewUser(u.ChatID, domainUser.WithSettings(settings))

	if err := s.userRepository.Save(ctx, user); err != nil {
		log.Error("failed to save domainUser", logutil.Err(err))
		return err
	}

	return nil
}

func (s UserService) UpdateUserSettings(ctx context.Context, id domainUser.Identifier, settings domainUser.Settings) error {
	const op = "UserService.UpdateUserSettings"
	log := s.log.With("op", op)

	if !id.IsValid() {
		return domainUser.ErrInvalidIdentifier
	}

	var userID uuid.UUID
	switch id := id.(type) {
	case domainUser.UUIDIdentifier:
		userID = id.UUID()
	case domainUser.TelegramIDWrapper:
		user, err := s.userProvider.GetUserByTelegramID(ctx, id.TelegramID())
		if err != nil {
			log.Error("failed to get domainUser", logutil.Err(err))
			return err
		}
		userID = user.ID()
	default:
		return domainUser.ErrInvalidIdentifier
	}

	err := s.userRepository.UpdateSettings(ctx, userID, settings)
	if err != nil {
		log.Error("failed to update domainUser settings", logutil.Err(err))
		return err
	}

	return nil
}
