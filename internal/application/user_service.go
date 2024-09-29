package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid"
	"go.uber.org/multierr"

	"golang.org/x/text/language"

	domainUser "github.com/ARUMANDESU/go-revise/internal/domain/user"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

var (
	ErrInvalidArguments = errors.New("invalid arguments")
)

type UserService struct {
	log            *slog.Logger
	userRepository domainUser.Repository
	userProvider   domainUser.Provider
}

func NewUserService(log *slog.Logger, userRepository domainUser.Repository, userProvider domainUser.Provider) UserService {
	return UserService{
		log:            log,
		userRepository: userRepository,
		userProvider:   userProvider,
	}
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
	case domainUser.UUID:
		user, err = s.userProvider.GetUserByID(ctx, id.ID().(uuid.UUID))
	case domainUser.TelegramID:
		user, err = s.userProvider.GetUserByTelegramID(ctx, id.ID().(domainUser.TelegramID))
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
	Language     language.Tag             `json:"language"`
	ReminderTime *domainUser.ReminderTime `json:"reminder_time"`
}

func (s UserService) RegisterUser(ctx context.Context, u NewUserServiceParams) error {
	const op = "UserService.RegisterUser"
	log := s.log.With("op", op)

	if !u.ChatID.IsValid() {
		return fmt.Errorf("%w: chat_id", ErrInvalidArguments)
	}
	if u.Language == language.Und {
		return fmt.Errorf("%w: language", ErrInvalidArguments)
	}
	if !u.ReminderTime.IsValid() {
		return fmt.Errorf("%w: reminder_time", ErrInvalidArguments)
	}

	settings := domainUser.Settings{
		Language:     u.Language,
		ReminderTime: domainUser.DefaultReminderTime(),
	}
	if u.ReminderTime != nil {
		settings.ReminderTime = *u.ReminderTime
	}

	user, err := domainUser.NewUser(u.ChatID, domainUser.WithSettings(settings))
	if err != nil {
		log.Error("failed to create domainUser", logutil.Err(err))
		return err
	}

	if err := s.userRepository.SaveUser(ctx, user); err != nil {
		log.Error("failed to save domainUser", logutil.Err(err))
		return err
	}

	return nil
}

func (s UserService) UpdateUserSettings(ctx context.Context, id domainUser.Identifier, settings domainUser.Settings) error {
	const op = "UserService.UpdateUserSettings"
	log := s.log.With("op", op)

	switch {
	case !id.IsValid():
		return multierr.Combine(ErrInvalidArguments, domainUser.ErrInvalidIdentifier)
	}
	if !settings.IsValid() {
		return multierr.Combine(ErrInvalidArguments, domainUser.ErrInvalidSettings)
	}

	var userID uuid.UUID
	switch id := id.(type) {
	case domainUser.UUID:
		userID = id.ID().(uuid.UUID)
	case domainUser.TelegramID:
		user, err := s.userProvider.GetUserByTelegramID(ctx, id.ID().(domainUser.TelegramID))
		if err != nil {
			log.Error("failed to get domainUser", logutil.Err(err))
			return err
		}
		userID = user.ID()
	default:
		return domainUser.ErrInvalidIdentifier
	}

	err := s.userRepository.UpdateUser(ctx, userID, func(user *domainUser.User) (*domainUser.User, error) {
		user.UpdateSettings(settings)
		return user, nil
	})
	if err != nil {
		log.Error("failed to update domainUser settings", logutil.Err(err))
		return err
	}

	return nil
}
