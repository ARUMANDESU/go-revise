package usersvc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ARUMANDESU/go-revise/internal/domain"
	"github.com/ARUMANDESU/go-revise/internal/service"
	"github.com/ARUMANDESU/go-revise/internal/storage"
	"github.com/ARUMANDESU/go-revise/pkg/logger"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid"
)

//go:generate mockery --name UserProvider --output ./mocks
type UserProvider interface {
	GetUser(ctx context.Context, id string) (domain.User, error)
	GetUserByChatID(ctx context.Context, chatID int64) (domain.User, error)
}

//go:generate mockery --name UserCreator --output ./mocks
type UserCreator interface {
	CreateUser(ctx context.Context, user domain.User) error
}

type Service struct {
	log      *slog.Logger
	provider UserProvider
	creator  UserCreator
}

func NewService(log *slog.Logger, provider UserProvider, creator UserCreator) Service {
	return Service{
		log:      log,
		provider: provider,
		creator:  creator,
	}
}

func (s Service) Get(ctx context.Context, id string) (domain.User, error) {
	const op = "service.user.get"
	log := s.log.With("op", op)

	err := validation.Validate(id, validation.Required, is.UUID)
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %w", service.ErrInvalidArgument, err)
	}

	user, err := s.provider.GetUser(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.User{}, service.ErrNotFound
		default:
			log.Error("failed to get user", logger.Err(err))
			return domain.User{}, service.ErrInternal
		}
	}

	return user, nil
}

func (s Service) GetByChatID(ctx context.Context, chatID int64) (domain.User, error) {
	const op = "service.user.getByChatID"
	log := s.log.With("op", op)

	user, err := s.provider.GetUserByChatID(ctx, chatID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return domain.User{}, service.ErrNotFound
		default:
			log.Error("failed to get user", logger.Err(err))
			return domain.User{}, service.ErrInternal
		}
	}

	return user, nil
}

func (s Service) Create(ctx context.Context, chatID int64) (domain.User, error) {
	const op = "service.user.create"
	log := s.log.With("op", op)

	uid, err := uuid.NewV7()
	if err != nil {
		return domain.User{}, domain.WrapErrorWithOp(err, op, "failed to generate UUID")
	}

	user, err := s.provider.GetUserByChatID(ctx, chatID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		switch {
		case errors.Is(err, storage.ErrAlreadyExists):
			return domain.User{}, service.ErrAlreadyExists
		default:
			log.Error("failed to get user", logger.Err(err))
			return domain.User{}, service.ErrInternal
		}
	}

	if user.ID != uuid.Nil {
		return domain.User{}, service.ErrAlreadyExists
	}

	user = domain.User{
		ID:         uid,
		TelegramID: chatID,
	}

	err = s.creator.CreateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrAlreadyExists):
			return domain.User{}, service.ErrAlreadyExists
		default:
			log.Error("failed to create user", logger.Err(err))
			return domain.User{}, service.ErrInternal
		}
	}

	return user, nil
}
