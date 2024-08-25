package tgbot

import (
	"context"

	"github.com/ARUMANDESU/go-revise/internal/domain"
)

type Bot struct {
}

type ReviseService interface {
	Get(ctx context.Context, id string) (domain.ReviseItem, error)
	List(ctx context.Context, userID string) ([]domain.ReviseItem, domain.PaginationMetadata, error)
	Create(ctx context.Context, revise domain.ReviseItem) (string, error)
	Update(ctx context.Context, revise domain.ReviseItem) error
	Delete(ctx context.Context, id string) error
}
