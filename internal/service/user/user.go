package usersvc

import (
	"context"
	"log/slog"

	"github.com/ARUMANDESU/go-revise/internal/domain"
)

type UserProvider interface {
	Get(ctx context.Context, id string) (domain.User, error)
	GetByChatID(ctx context.Context, chatID int64) (domain.User, error)
}

type UserCreator interface {
	Create(ctx context.Context, chatID int64) (domain.User, error)
}

type Service struct {
	log      *slog.Logger
	provider UserProvider
	creator  UserCreator
}

func NewService(log *slog.Logger, provider UserProvider, creator UserCreator) *Service {
	return &Service{
		log:      log,
		provider: provider,
		creator:  creator,
	}
}

func (s Service) Get(ctx context.Context, id string) (domain.User, error) {
	// TODO: implement this method
	panic("not implemented")
}

func (s Service) GetByChatID(ctx context.Context, chatID int64) (domain.User, error) {
	// TODO: implement this method
	panic("not implemented")
}

func (s Service) Create(ctx context.Context, chatID int64) (domain.User, error) {
	// TODO: implement this method
	panic("not implemented")
}
